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
	"encoding/json"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/require"
)

// cip83SpecVectorJSON and cip83SpecVectorMessages are taken directly from the
// CIP-83 specification's "basic" mode example, encrypted with the default
// passphrase "cardano".
const cip83SpecVectorJSON = `{
  "674": {
    "enc": "basic",
    "msg": [
      "U2FsdGVkX1/5Y0A7l8xK686rvLsmPviTlna2n3P/ADNm89Ynr1UPZ/Q6bynbe28Y",
      "/zWYOB9PAGt+bq1L0z/W2LNHe92HTN/Fwz16aHa98TOsgM3q8tAR4NSqrLZVu1H7"
    ]
  }
}`

var cip83SpecVectorMessages = []string{
	"Invoice-No: 123456789",
	"Order-No: 7654321",
	"Email: john@doe.com",
}

func TestCip83Metadata_JSONUnmarshalAndValidate(t *testing.T) {
	t.Parallel()

	var metadata Cip83Metadata
	require.NoError(
		t,
		json.Unmarshal([]byte(cip83SpecVectorJSON), &metadata),
	)
	require.NoError(t, metadata.Validate())
	require.Equal(t, Cip83EncBasic, metadata.Num674.Enc)
	require.Equal(t, 2, len(metadata.Num674.Msg))
}

func TestCip83Metadata_CBORRoundTrip(t *testing.T) {
	t.Parallel()

	metadata, err := NewCip83Metadata(Cip83EncBasic, []string{
		"U2FsdGVkX1/5Y0A7l8xK686rvLsmPviTlna2n3P/ADNm89Ynr1UPZ/Q6bynbe28Y",
		"/zWYOB9PAGt+bq1L0z/W2LNHe92HTN/Fwz16aHa98TOsgM3q8tAR4NSqrLZVu1H7",
	})
	require.NoError(t, err)

	data, err := cbor.Marshal(metadata)
	require.NoError(t, err)

	var decoded Cip83Metadata
	require.NoError(t, cbor.Unmarshal(data, &decoded))
	require.Equal(t, metadata.Num674, decoded.Num674)
}

func TestDecryptCip83Message_SpecVector(t *testing.T) {
	t.Parallel()

	var metadata Cip83Metadata
	require.NoError(
		t,
		json.Unmarshal([]byte(cip83SpecVectorJSON), &metadata),
	)

	messages, err := DecryptCip83Message(&metadata, "")
	require.NoError(t, err)
	require.Equal(t, cip83SpecVectorMessages, messages)

	// Explicit default passphrase should behave the same as empty.
	messages, err = DecryptCip83Message(&metadata, Cip83DefaultPassphrase)
	require.NoError(t, err)
	require.Equal(t, cip83SpecVectorMessages, messages)
}

func TestEncryptDecryptCip83Message_RoundTrip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		messages   []string
		passphrase string
	}{
		{
			name:       "default passphrase",
			messages:   []string{"hello world"},
			passphrase: "",
		},
		{
			name:       "custom passphrase",
			messages:   []string{"Invoice-No: 42", "Thanks for your order"},
			passphrase: "s3cr3t-passphrase",
		},
		{
			name: "long message requiring multiple base64 lines",
			messages: []string{
				strings.Repeat("a very long message indeed ", 10),
			},
			passphrase: "cardano",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			metadata, err := EncryptCip83Message(tc.messages, tc.passphrase)
			require.NoError(t, err)
			require.Equal(t, Cip83EncBasic, metadata.Num674.Enc)
			for _, line := range metadata.Num674.Msg {
				require.LessOrEqual(t, len(line), 64)
			}

			decrypted, err := DecryptCip83Message(metadata, tc.passphrase)
			require.NoError(t, err)
			require.Equal(t, tc.messages, decrypted)
		})
	}
}

func TestDecryptCip83Message_WrongPassphraseFails(t *testing.T) {
	t.Parallel()

	metadata, err := EncryptCip83Message([]string{"secret"}, "correct-horse")
	require.NoError(t, err)

	_, err = DecryptCip83Message(metadata, "wrong-passphrase")
	require.Error(t, err)
}

func TestEncryptCip83Message_EmptyMessagesRejected(t *testing.T) {
	t.Parallel()

	_, err := EncryptCip83Message(nil, "")
	require.Error(t, err)

	_, err = EncryptCip83Message([]string{}, "")
	require.Error(t, err)
}

func TestDecryptCip83Message_PlainPassthrough(t *testing.T) {
	t.Parallel()

	metadata, err := NewCip83Metadata(Cip83EncPlain, []string{"hello"})
	require.NoError(t, err)

	messages, err := DecryptCip83Message(metadata, "")
	require.NoError(t, err)
	require.Equal(t, []string{"hello"}, messages)

	// No enc value at all is also treated as plaintext.
	metadata, err = NewCip83Metadata("", []string{"hello"})
	require.NoError(t, err)
	messages, err = DecryptCip83Message(metadata, "")
	require.NoError(t, err)
	require.Equal(t, []string{"hello"}, messages)
}

func TestDecryptCip83Message_UnsupportedMethod(t *testing.T) {
	t.Parallel()

	metadata := &Cip83Metadata{
		Num674: Num674Enc{Enc: "chacha20", Msg: []string{"abc"}},
	}

	_, err := DecryptCip83Message(metadata, "")
	require.Error(t, err)
}

func TestNewCip83Metadata_ValidationErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		enc      string
		messages []string
	}{
		{
			name:     "invalid enc value",
			enc:      "rot13",
			messages: []string{"hello"},
		},
		{
			name:     "empty messages",
			enc:      Cip83EncPlain,
			messages: []string{},
		},
		{
			name:     "message line too long",
			enc:      Cip83EncPlain,
			messages: []string{strings.Repeat("a", 65)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewCip83Metadata(tc.enc, tc.messages)
			require.Error(t, err)
		})
	}
}
