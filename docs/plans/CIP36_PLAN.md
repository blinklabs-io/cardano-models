# CIP-36 Catalyst/Voltaire Registration Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers implementation of CIP-36 (Catalyst/Voltaire Registration) metadata format, the successor to CIP-15 that supports multiple delegations and Voltaire governance.

**Related Issue**: [#4](https://github.com/blinklabs-io/cardano-models/issues/4)

---

## Scope

### CIP-36 Specification

CIP-36 extends CIP-15 with:
- Multiple voting key delegations with weights
- Support for Voltaire governance voting
- Backward compatibility with CIP-15 single-key format

### Metadata Labels

| Label | Purpose |
|-------|---------|
| 61284 | Registration payload |
| 61285 | Registration witness (signature) |

### Registration Payload (61284)

```cddl
registration_payload = {
  1: delegations,         ; Voting key delegations
  2: staking_key,         ; Ed25519 public key (32 bytes)
  3: reward_address,      ; Address for rewards
  4: nonce,               ; Slot number for replay protection
  ? 5: voting_purpose     ; Optional: 0=Catalyst, 1=Other
}

; Single key (CIP-15 compatible)
delegations = voting_key / [ + delegation ]

; Multiple delegations with weights
delegation = [ voting_key, weight ]
voting_key = bytes .size 32
weight = uint

staking_key = bytes .size 32
reward_address = bytes
nonce = uint
voting_purpose = uint
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

### Core Structures

```go
package cip36

// VotingPurpose defines the registration purpose
type VotingPurpose uint

const (
    VotingPurposeCatalyst VotingPurpose = 0
    VotingPurposeOther    VotingPurpose = 1
)

// Delegation represents a voting key delegation with weight
type Delegation struct {
    VotingKey []byte `cbor:"0" validate:"required,len=32"`
    Weight    uint   `cbor:"1" validate:"required"`
}

// CIP36Registration represents a Catalyst/Voltaire voter registration
type CIP36Registration struct {
    Payload CIP36RegistrationPayload `cbor:"61284,keyasint"`
    Witness CIP36RegistrationWitness `cbor:"61285,keyasint"`
}

// CIP36RegistrationPayload contains the registration data
type CIP36RegistrationPayload struct {
    Delegations   []Delegation   `cbor:"1,keyasint" validate:"required,min=1"`
    StakingKey    []byte         `cbor:"2,keyasint" validate:"required,len=32"`
    RewardAddress []byte         `cbor:"3,keyasint" validate:"required"`
    Nonce         uint64         `cbor:"4,keyasint" validate:"required"`
    VotingPurpose *VotingPurpose `cbor:"5,keyasint,omitempty"`
}

// CIP36RegistrationWitness contains the signature
type CIP36RegistrationWitness struct {
    Signature []byte `cbor:"1,keyasint" validate:"required,len=64"`
}
```

### Builder Pattern

```go
// NewCIP36Registration creates a new registration builder
func NewCIP36Registration() *CIP36RegistrationBuilder

type CIP36RegistrationBuilder struct {
    delegations   []Delegation
    stakingKey    []byte
    rewardAddress []byte
    nonce         uint64
    votingPurpose *VotingPurpose
    signature     []byte
}

// Single delegation (CIP-15 compatible)
func (b *CIP36RegistrationBuilder) WithVotingKey(key []byte) *CIP36RegistrationBuilder

// Multiple delegations
func (b *CIP36RegistrationBuilder) WithDelegation(key []byte, weight uint) *CIP36RegistrationBuilder
func (b *CIP36RegistrationBuilder) WithDelegations(delegations []Delegation) *CIP36RegistrationBuilder

// Other fields
func (b *CIP36RegistrationBuilder) WithStakingKey(key []byte) *CIP36RegistrationBuilder
func (b *CIP36RegistrationBuilder) WithRewardAddress(addr []byte) *CIP36RegistrationBuilder
func (b *CIP36RegistrationBuilder) WithNonce(slot uint64) *CIP36RegistrationBuilder
func (b *CIP36RegistrationBuilder) WithVotingPurpose(purpose VotingPurpose) *CIP36RegistrationBuilder
func (b *CIP36RegistrationBuilder) ForCatalyst() *CIP36RegistrationBuilder
func (b *CIP36RegistrationBuilder) WithSignature(sig []byte) *CIP36RegistrationBuilder
func (b *CIP36RegistrationBuilder) Build() (*CIP36Registration, error)
```

### Extraction Functions

```go
// ExtractCIP36Registration extracts registration from transaction metadata
func ExtractCIP36Registration(metadata map[uint64]any) (*CIP36Registration, error)

// HasCIP36Registration checks if metadata contains CIP-36 registration
func HasCIP36Registration(metadata map[uint64]any) bool

// IsCIP36Format checks if registration uses CIP-36 multi-delegation format
// vs CIP-15 single-key format
func IsCIP36Format(metadata map[uint64]any) bool
```

### Validation & Analysis

```go
// Validate checks registration validity
func (r *CIP36Registration) Validate() error

// VerifySignature verifies the registration signature
func (r *CIP36Registration) VerifySignature() error

// IsSingleDelegation returns true if using CIP-15 compatible single key
func (r *CIP36Registration) IsSingleDelegation() bool

// TotalWeight returns the sum of all delegation weights
func (r *CIP36Registration) TotalWeight() uint

// NormalizedWeights returns weights as percentages (0.0 to 1.0)
func (r *CIP36Registration) NormalizedWeights() map[string]float64

// GetVotingPower calculates voting power from staking key delegation
func (r *CIP36Registration) GetVotingPower(utxos []common.Utxo) uint64
```

---

## Usage Examples

### Parse Multi-Delegation Registration

```go
import "github.com/blinklabs-io/cardano-models/cip36"

// Extract from transaction metadata
reg, err := cip36.ExtractCIP36Registration(tx.Metadata())
if err != nil {
    // Not a valid CIP-36 registration
}

// Check format
if reg.IsSingleDelegation() {
    fmt.Println("CIP-15 compatible single delegation")
} else {
    fmt.Printf("Multi-delegation with %d keys\n", len(reg.Payload.Delegations))
}

// Get normalized weights
weights := reg.NormalizedWeights()
for keyHex, pct := range weights {
    fmt.Printf("Key %s: %.2f%%\n", keyHex, pct*100)
}
```

### Build Multi-Delegation Registration

```go
reg, err := cip36.NewCIP36Registration().
    WithDelegation(votingKey1, 3).   // 60% weight
    WithDelegation(votingKey2, 2).   // 40% weight
    WithStakingKey(stakingPubKey).
    WithRewardAddress(rewardAddr).
    WithNonce(currentSlot).
    ForCatalyst().
    WithSignature(signature).
    Build()
```

### Build Single-Key Registration (CIP-15 Compatible)

```go
reg, err := cip36.NewCIP36Registration().
    WithVotingKey(votingPubKey).     // Single key = weight 1
    WithStakingKey(stakingPubKey).
    WithRewardAddress(rewardAddr).
    WithNonce(currentSlot).
    WithSignature(signature).
    Build()
```

---

## Implementation Phases

### Phase 1: Core Types

- [ ] Define `Delegation` struct
- [ ] Define `VotingPurpose` type
- [ ] Define `CIP36RegistrationPayload` struct
- [ ] Define `CIP36RegistrationWitness` struct
- [ ] Define `CIP36Registration` wrapper

### Phase 2: CBOR Serialization

- [ ] Handle single-key format (CIP-15 compatible)
- [ ] Handle multi-delegation array format
- [ ] Implement custom CBOR marshaling
- [ ] Implement custom CBOR unmarshaling

### Phase 3: Extraction & Detection

- [ ] Implement `ExtractCIP36Registration()`
- [ ] Implement `HasCIP36Registration()`
- [ ] Implement `IsCIP36Format()`
- [ ] Auto-detect format from payload structure

### Phase 4: Validation & Analysis

- [ ] Implement `Validate()`
- [ ] Implement `VerifySignature()`
- [ ] Implement `IsSingleDelegation()`
- [ ] Implement weight calculation methods

### Phase 5: Testing

- [ ] Unit tests with synthetic data
- [ ] Tests with CIP-15 format registrations
- [ ] Tests with CIP-36 multi-delegation registrations
- [ ] Tests with real Catalyst registrations from mainnet

---

## Success Criteria

- [ ] Parse both CIP-15 and CIP-36 format registrations
- [ ] Correctly detect single vs multi-delegation format
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

- [CIP15_PLAN.md](CIP15_PLAN.md) - Predecessor standard
