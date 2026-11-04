# CIP-15 Catalyst Registration Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers implementation of CIP-15 (Catalyst Registration) metadata format used for Catalyst voting registration prior to CIP-36.

**Related Issue**: [#3](https://github.com/blinklabs-io/cardano-models/issues/3)

---

## Scope

### CIP-15 Specification

CIP-15 defines the original Catalyst voter registration format using transaction metadata labels 61284 and 61285.

### Metadata Labels

| Label | Purpose |
|-------|---------|
| 61284 | Registration payload |
| 61285 | Registration witness (signature) |

### Registration Payload (61284)

```cddl
registration_payload = {
  1: voting_key,          ; Ed25519 public key (32 bytes)
  2: staking_key,         ; Ed25519 public key (32 bytes)
  3: reward_address,      ; Shelley address for rewards
  4: nonce                ; Slot number for replay protection
}

voting_key = bytes .size 32
staking_key = bytes .size 32
reward_address = bytes
nonce = uint
```

### Registration Witness (61285)

```cddl
registration_witness = {
  1: signature            ; Ed25519 signature over payload
}

signature = bytes .size 64
```

---

## API Design

### Registration Structure

```go
package cip15

// CIP15Registration represents a Catalyst voter registration
type CIP15Registration struct {
    Payload CIP15RegistrationPayload `cbor:"61284,keyasint"`
    Witness CIP15RegistrationWitness `cbor:"61285,keyasint"`
}

// CIP15RegistrationPayload contains the registration data
type CIP15RegistrationPayload struct {
    VotingKey     []byte `cbor:"1,keyasint" validate:"required,len=32"`
    StakingKey    []byte `cbor:"2,keyasint" validate:"required,len=32"`
    RewardAddress []byte `cbor:"3,keyasint" validate:"required"`
    Nonce         uint64 `cbor:"4,keyasint" validate:"required"`
}

// CIP15RegistrationWitness contains the signature
type CIP15RegistrationWitness struct {
    Signature []byte `cbor:"1,keyasint" validate:"required,len=64"`
}
```

### Builder Pattern

```go
// NewCIP15Registration creates a new registration builder
func NewCIP15Registration() *CIP15RegistrationBuilder

type CIP15RegistrationBuilder struct {
    votingKey     []byte
    stakingKey    []byte
    rewardAddress []byte
    nonce         uint64
    signature     []byte
}

func (b *CIP15RegistrationBuilder) WithVotingKey(key []byte) *CIP15RegistrationBuilder
func (b *CIP15RegistrationBuilder) WithStakingKey(key []byte) *CIP15RegistrationBuilder
func (b *CIP15RegistrationBuilder) WithRewardAddress(addr []byte) *CIP15RegistrationBuilder
func (b *CIP15RegistrationBuilder) WithNonce(slot uint64) *CIP15RegistrationBuilder
func (b *CIP15RegistrationBuilder) WithSignature(sig []byte) *CIP15RegistrationBuilder
func (b *CIP15RegistrationBuilder) Build() (*CIP15Registration, error)
```

### Extraction Functions

```go
// ExtractCIP15Registration extracts registration from transaction metadata
func ExtractCIP15Registration(metadata map[uint64]any) (*CIP15Registration, error)

// HasCIP15Registration checks if metadata contains CIP-15 registration
func HasCIP15Registration(metadata map[uint64]any) bool
```

### Validation

```go
// Validate checks registration validity
func (r *CIP15Registration) Validate() error

// VerifySignature verifies the registration signature
func (r *CIP15Registration) VerifySignature() error

// GetVotingPower calculates voting power from staking key delegation
func (r *CIP15Registration) GetVotingPower(utxos []common.Utxo) uint64
```

---

## Usage Examples

### Parse Registration from Metadata

```go
import "github.com/blinklabs-io/cardano-models/cip15"

// Extract from transaction metadata
reg, err := cip15.ExtractCIP15Registration(tx.Metadata())
if err != nil {
    // Not a valid CIP-15 registration
}

// Validate the registration
if err := reg.Validate(); err != nil {
    // Invalid registration data
}

// Verify signature
if err := reg.VerifySignature(); err != nil {
    // Invalid signature
}
```

### Build Registration

```go
reg, err := cip15.NewCIP15Registration().
    WithVotingKey(votingPubKey).
    WithStakingKey(stakingPubKey).
    WithRewardAddress(rewardAddr).
    WithNonce(currentSlot).
    WithSignature(signature).
    Build()
```

---

## Implementation Phases

### Phase 1: Core Types

- [ ] Define `CIP15RegistrationPayload` struct
- [ ] Define `CIP15RegistrationWitness` struct
- [ ] Define `CIP15Registration` wrapper struct
- [ ] Add CBOR tags for metadata labels

### Phase 2: Serialization

- [ ] Implement CBOR marshaling
- [ ] Implement CBOR unmarshaling
- [ ] Add JSON support for off-chain use

### Phase 3: Extraction & Validation

- [ ] Implement `ExtractCIP15Registration()`
- [ ] Implement `HasCIP15Registration()`
- [ ] Implement `Validate()`
- [ ] Add validation error types

### Phase 4: Signature Verification

- [ ] Implement `VerifySignature()`
- [ ] Add Ed25519 signature verification
- [ ] Test with real mainnet registrations

### Phase 5: Testing

- [ ] Unit tests with synthetic data
- [ ] Round-trip serialization tests
- [ ] Tests with real Catalyst registrations from mainnet

---

## Success Criteria

- [ ] Parse real CIP-15 registrations from mainnet
- [ ] Validate registration signatures
- [ ] CBOR output matches cardano-cli format
- [ ] Comprehensive test coverage

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| gouroboros | Transaction metadata types |
| ed25519 | Signature verification |
| cbor/v2 | CBOR encoding |

---

## Related Plans

- [CIP36_PLAN.md](CIP36_PLAN.md) - Successor standard (CIP-36)
