# CIP-83 Encrypted Transaction Messages Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers implementation of CIP-83 (Encrypted Transaction Messages) which extends CIP-20 to support encrypted messages in transaction metadata.

**Related Issue**: [#6](https://github.com/blinklabs-io/cardano-models/issues/6)

---

## Scope

### CIP-83 Specification

CIP-83 defines a standard for encrypted messages in transaction metadata, building on CIP-20's message format. Messages are encrypted using recipient public keys, allowing private communication on-chain.

### Metadata Label

| Label | Purpose |
|-------|---------|
| 674 | Transaction messages (shared with CIP-20) |

### Message Structure

```cddl
; CIP-83 extends CIP-20 with encrypted messages
transaction_message = {
  "msg": [* text],                    ; CIP-20 plaintext
  / "enc": encrypted_message_type     ; CIP-83 encrypted
}

; Encryption type identifier
encrypted_message_type = "basic" / "other"

; Encrypted payload structure
encrypted_message = {
  "enc": encrypted_message_type,
  "payload": [* encrypted_chunk],
  ? "recipients": [* recipient_info]
}

encrypted_chunk = {
  "ciphertext": text,        ; Base64-encoded encrypted data
  "nonce": text,             ; Base64-encoded nonce
  ? "ephemeral_pk": text     ; Base64-encoded ephemeral public key
}

recipient_info = {
  "address": text,           ; Recipient address
  ? "enc_key": text          ; Encrypted symmetric key for recipient
}
```

### Encryption Scheme ("basic")

The "basic" encryption scheme uses:
- X25519 key exchange (Curve25519)
- ChaCha20-Poly1305 AEAD encryption
- HKDF-SHA256 key derivation

---

## API Design

### Core Structures

```go
package cip83

// EncryptionType identifies the encryption scheme
type EncryptionType string

const (
    EncryptionBasic EncryptionType = "basic"
)

// CIP83Message represents an encrypted transaction message
type CIP83Message struct {
    EncType    EncryptionType   `json:"enc" validate:"required"`
    Payload    []EncryptedChunk `json:"payload" validate:"required,min=1"`
    Recipients []RecipientInfo  `json:"recipients,omitempty"`
}

// EncryptedChunk represents a chunk of encrypted data
type EncryptedChunk struct {
    Ciphertext  string `json:"ciphertext" validate:"required"`
    Nonce       string `json:"nonce" validate:"required"`
    EphemeralPK string `json:"ephemeral_pk,omitempty"`
}

// RecipientInfo contains recipient-specific encryption data
type RecipientInfo struct {
    Address string `json:"address" validate:"required"`
    EncKey  string `json:"enc_key,omitempty"`
}
```

### Extraction Functions

```go
// ExtractCIP83Message extracts encrypted message from metadata
func ExtractCIP83Message(metadata map[uint64]any) (*CIP83Message, error)

// HasCIP83Message checks if metadata contains encrypted message
func HasCIP83Message(metadata map[uint64]any) bool

// IsCIP83 checks if 674 label contains "enc" field
func IsCIP83(metadata map[uint64]any) bool

// IsCIP20 checks if 674 label contains "msg" field (plaintext)
func IsCIP20(metadata map[uint64]any) bool
```

### Decryption Functions

```go
// DecryptMessage decrypts a CIP-83 message using recipient private key
func DecryptMessage(msg *CIP83Message, privateKey []byte) (string, error)

// DecryptChunk decrypts a single encrypted chunk
func DecryptChunk(chunk EncryptedChunk, sharedSecret []byte) ([]byte, error)

// DeriveSharedSecret derives shared secret from private and ephemeral keys
func DeriveSharedSecret(privateKey, ephemeralPublicKey []byte) ([]byte, error)
```

### Encryption Functions

```go
// EncryptMessage encrypts a message for recipients
func EncryptMessage(plaintext string, recipientPubKeys [][]byte) (*CIP83Message, error)

// EncryptChunk encrypts a chunk of data
func EncryptChunk(data []byte, sharedSecret []byte) (*EncryptedChunk, error)

// GenerateEphemeralKeyPair generates a new X25519 key pair
func GenerateEphemeralKeyPair() (publicKey, privateKey []byte, err error)
```

### Builder Pattern

```go
// NewCIP83Message creates a new encrypted message builder
func NewCIP83Message() *CIP83MessageBuilder

type CIP83MessageBuilder struct {
    encType    EncryptionType
    plaintext  string
    recipients [][]byte  // Public keys
}

func (b *CIP83MessageBuilder) WithBasicEncryption() *CIP83MessageBuilder
func (b *CIP83MessageBuilder) WithPlaintext(text string) *CIP83MessageBuilder
func (b *CIP83MessageBuilder) WithRecipient(publicKey []byte) *CIP83MessageBuilder
func (b *CIP83MessageBuilder) WithRecipientAddress(address string, publicKey []byte) *CIP83MessageBuilder
func (b *CIP83MessageBuilder) Build() (*CIP83Message, error)
```

### Key Derivation

```go
// DeriveEncryptionKey derives X25519 key from Cardano payment key
func DeriveEncryptionKey(paymentKey []byte) ([]byte, error)

// GetPaymentKeyFromAddress extracts payment key hash from address
func GetPaymentKeyFromAddress(address string) ([]byte, error)
```

---

## Usage Examples

### Check for Encrypted Message

```go
import "github.com/blinklabs-io/cardano-models/cip83"

// Check what type of message
if cip83.IsCIP83(tx.Metadata()) {
    fmt.Println("Contains encrypted message")
} else if cip83.IsCIP20(tx.Metadata()) {
    fmt.Println("Contains plaintext message")
}
```

### Decrypt Message

```go
// Extract encrypted message
msg, err := cip83.ExtractCIP83Message(tx.Metadata())
if err != nil {
    // No valid CIP-83 message
}

// Decrypt using recipient's private key
plaintext, err := cip83.DecryptMessage(msg, myPrivateKey)
if err != nil {
    // Decryption failed (wrong key or corrupted)
}

fmt.Println("Decrypted message:", plaintext)
```

### Encrypt Message for Recipients

```go
// Build encrypted message
msg, err := cip83.NewCIP83Message().
    WithBasicEncryption().
    WithPlaintext("Secret message for recipients").
    WithRecipient(recipientPubKey1).
    WithRecipient(recipientPubKey2).
    Build()

// Include in transaction metadata
metadata := map[uint64]any{
    674: msg,
}
```

### Key Derivation from Address

```go
// Derive encryption key from Cardano address
paymentKeyHash, err := cip83.GetPaymentKeyFromAddress(address)
if err != nil {
    // Invalid address
}

// Note: Full encryption requires the actual public key, not just hash
// This is typically obtained through off-chain key exchange
```

---

## Implementation Phases

### Phase 1: Core Types

- [ ] Define `EncryptionType` type
- [ ] Define `CIP83Message` struct
- [ ] Define `EncryptedChunk` struct
- [ ] Define `RecipientInfo` struct

### Phase 2: Detection & Extraction

- [ ] Implement `HasCIP83Message()`
- [ ] Implement `IsCIP83()`
- [ ] Implement `IsCIP20()`
- [ ] Implement `ExtractCIP83Message()`

### Phase 3: Cryptographic Primitives

- [ ] Implement X25519 key exchange
- [ ] Implement ChaCha20-Poly1305 encryption
- [ ] Implement HKDF key derivation
- [ ] Add nonce generation

### Phase 4: Decryption

- [ ] Implement `DeriveSharedSecret()`
- [ ] Implement `DecryptChunk()`
- [ ] Implement `DecryptMessage()`
- [ ] Handle multiple recipients

### Phase 5: Encryption

- [ ] Implement `GenerateEphemeralKeyPair()`
- [ ] Implement `EncryptChunk()`
- [ ] Implement `EncryptMessage()`
- [ ] Implement `CIP83MessageBuilder`

### Phase 6: Testing

- [ ] Unit tests for encryption/decryption
- [ ] Round-trip tests
- [ ] Tests with real CIP-83 messages from mainnet
- [ ] Cross-reference with reference implementations

---

## Success Criteria

- [ ] Detect CIP-83 vs CIP-20 messages
- [ ] Decrypt CIP-83 messages with correct key
- [ ] Encrypt messages for multiple recipients
- [ ] Compatible with reference CIP-83 implementations
- [ ] Comprehensive test coverage

---

## Security Considerations

- Private keys must be handled securely
- Nonces must be unique per message
- Ephemeral keys should be generated fresh
- Consider timing attacks in decryption

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| gouroboros | Transaction metadata types |
| golang.org/x/crypto/chacha20poly1305 | AEAD encryption |
| golang.org/x/crypto/curve25519 | X25519 key exchange |
| golang.org/x/crypto/hkdf | Key derivation |

---

## Related Plans

- CIP-20 (already implemented) - Plaintext message format
