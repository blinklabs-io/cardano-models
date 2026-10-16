// Copyright 2026 Blink Labs Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/pbkdf2"
)

// Encryption methods defined by CIP-83.
const (
	Cip83EncPlain = "plain"
	Cip83EncBasic = "basic"
)

// Cip83DefaultPassphrase is the well-known default passphrase used to encrypt
// or decrypt messages when no other passphrase is provided or known.
const Cip83DefaultPassphrase = "cardano"

const (
	cip83Pbkdf2Iterations = 10000
	cip83SaltLen          = 8
	cip83KeyLen           = 32
	cip83IvLen            = 16
	cip83MaxLineLength    = 64
)

var cip83SaltedPrefix = []byte("Salted__")

// validateCip83MsgLine implements the "cip83msgbytes" validator tag. Cardano
// limits each metadata text value to 64 *encoded bytes*, but validator's
// built-in "max" tag counts Unicode code points for strings, so a line with
// multibyte characters could pass max=64 while exceeding the 64-byte limit
// on the wire. len() on a Go string already counts UTF-8 encoded bytes, so
// comparing against cip83MaxLineLength directly enforces the real limit.
func validateCip83MsgLine(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) <= cip83MaxLineLength
}

func init() {
	if err := validate.RegisterValidation("cip83msgbytes", validateCip83MsgLine); err != nil {
		panic(err)
	}
}

// Cip83Metadata is the top-level container for CIP-83 encrypted transaction
// message metadata under the "674" tag. It is an addendum to CIP-20, adding
// an optional "enc" field alongside the existing "msg" field.
type Cip83Metadata struct {
	Num674 Num674Enc `cbor:"674,keyasint" json:"674" validate:"required"`
}

// Num674Enc represents the metadata structure for CIP-83.
type Num674Enc struct {
	Enc string   `cbor:"enc,omitempty" json:"enc,omitempty" validate:"omitempty,oneof=plain basic"`
	Msg []string `cbor:"msg"           json:"msg"           validate:"required,gt=0,dive,cip83msgbytes"`
}

func (c *Cip83Metadata) UnmarshalCBOR(data []byte) error {
	var raw map[any]cbor.RawMessage
	if err := cbor.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		if !isCip20Label(key) {
			continue
		}
		var num674 Num674Enc
		if err := cbor.Unmarshal(value, &num674); err != nil {
			return err
		}
		c.Num674 = num674
		return nil
	}
	return nil
}

// NewCip83Metadata builds a Cip83Metadata from an encryption method and
// message lines, validating the result.
func NewCip83Metadata(enc string, messages []string) (*Cip83Metadata, error) {
	metadata := &Cip83Metadata{Num674: Num674Enc{Enc: enc, Msg: messages}}
	if err := validate.Struct(metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func (c *Cip83Metadata) Validate() error {
	return validate.Struct(c)
}

// EncryptCip83Message encrypts the given message lines using CIP-83's
// "basic" method (AES-256-CBC, salted, key/IV derived via PBKDF2-HMAC-SHA256
// with 10000 iterations) and returns the resulting metadata. If passphrase
// is empty, Cip83DefaultPassphrase is used.
func EncryptCip83Message(
	messages []string,
	passphrase string,
) (*Cip83Metadata, error) {
	if len(messages) == 0 {
		return nil, errors.New("cip83: at least one message is required")
	}
	if passphrase == "" {
		passphrase = Cip83DefaultPassphrase
	}
	plaintext, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}
	salt := make([]byte, cip83SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	key, iv := cip83DeriveKeyIv(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padded := cip83Pkcs7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)
	out := append([]byte{}, cip83SaltedPrefix...)
	out = append(out, salt...)
	out = append(out, ciphertext...)
	encoded := base64.StdEncoding.EncodeToString(out)
	return NewCip83Metadata(Cip83EncBasic, cip83SplitLines(encoded))
}

// DecryptCip83Message returns the decrypted message lines from metadata. If
// the metadata is unencrypted ("plain" or no "enc" value), the messages are
// returned as-is. If passphrase is empty, Cip83DefaultPassphrase is used.
func DecryptCip83Message(
	metadata *Cip83Metadata,
	passphrase string,
) ([]string, error) {
	switch metadata.Num674.Enc {
	case "", Cip83EncPlain:
		return metadata.Num674.Msg, nil
	case Cip83EncBasic:
		// handled below
	default:
		return nil, fmt.Errorf(
			"unsupported encryption method: %s",
			metadata.Num674.Enc,
		)
	}
	if passphrase == "" {
		passphrase = Cip83DefaultPassphrase
	}
	data, err := base64.StdEncoding.DecodeString(
		strings.Join(metadata.Num674.Msg, ""),
	)
	if err != nil {
		return nil, err
	}
	if len(data) < len(cip83SaltedPrefix)+cip83SaltLen ||
		!bytes.Equal(data[:len(cip83SaltedPrefix)], cip83SaltedPrefix) {
		return nil, errors.New("cip83: missing or invalid Salted__ header")
	}
	salt := data[len(cip83SaltedPrefix) : len(cip83SaltedPrefix)+cip83SaltLen]
	ciphertext := data[len(cip83SaltedPrefix)+cip83SaltLen:]
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New(
			"cip83: ciphertext is not a multiple of the block size",
		)
	}
	key, iv := cip83DeriveKeyIv(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)
	unpadded, err := cip83Pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	var messages []string
	if err := json.Unmarshal(unpadded, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func cip83DeriveKeyIv(passphrase string, salt []byte) (key, iv []byte) {
	keyIv := pbkdf2.Key(
		[]byte(passphrase),
		salt,
		cip83Pbkdf2Iterations,
		cip83KeyLen+cip83IvLen,
		sha256.New,
	)
	return keyIv[:cip83KeyLen], keyIv[cip83KeyLen:]
}

func cip83SplitLines(s string) []string {
	var lines []string
	for i := 0; i < len(s); i += cip83MaxLineLength {
		end := min(i+cip83MaxLineLength, len(s))
		lines = append(lines, s[i:end])
	}
	return lines
}

func cip83Pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padByte := byte(padLen) //nolint:gosec // padLen is bounded by blockSize (16)
	padded := append([]byte{}, data...)
	padded = append(padded, bytes.Repeat([]byte{padByte}, padLen)...)
	return padded
}

func cip83Pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, errors.New("cip83: invalid padded data length")
	}
	padLen := int(data[length-1])
	if padLen == 0 || padLen > blockSize || padLen > length {
		return nil, errors.New("cip83: invalid padding")
	}
	for _, b := range data[length-padLen:] {
		if int(b) != padLen {
			return nil, errors.New("cip83: invalid padding bytes")
		}
	}
	return data[:length-padLen], nil
}
