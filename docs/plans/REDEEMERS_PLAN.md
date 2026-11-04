# Redeemer Patterns Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers common redeemer patterns used across Cardano smart contracts. Redeemers specify the action to take when spending a script UTxO.

---

## Scope

### Redeemer Categories

| Category | Description | Use Cases |
|----------|-------------|-----------|
| Action Enum | Simple action identifier | DEX operations, governance |
| Parameterized | Action with parameters | Orders, positions |
| Index-Based | References to inputs/outputs | Batched operations |
| Composite | Multiple actions combined | Complex transactions |

### Common Patterns

1. **Simple Actions**: Enum-style (Swap, Withdraw, Claim)
2. **Order Actions**: Order type + parameters
3. **Index References**: Input/output indices
4. **State Machine**: State transitions

---

## API Design

### Simple Action Redeemers

```go
package redeemer

// ActionRedeemer represents an enum-style action
type ActionRedeemer struct {
    Action uint8 `cbor:"0,keyasint"`
}

// Common action constants
const (
    ActionSwap     uint8 = 0
    ActionDeposit  uint8 = 1
    ActionWithdraw uint8 = 2
    ActionClaim    uint8 = 3
    ActionCancel   uint8 = 4
)

// NewActionRedeemer creates a simple action redeemer
func NewActionRedeemer(action uint8) *ActionRedeemer
```

### DEX Redeemers

```go
// SwapRedeemer for DEX swap operations
type SwapRedeemer struct {
    Action      uint8  `cbor:"0,keyasint"` // Usually 0
    Direction   uint8  `cbor:"1,keyasint"` // 0=A→B, 1=B→A
    MinReceive  uint64 `cbor:"2,keyasint,omitempty"`
}

// DepositRedeemer for adding liquidity
type DepositRedeemer struct {
    Action    uint8  `cbor:"0,keyasint"`
    AmountA   uint64 `cbor:"1,keyasint"`
    AmountB   uint64 `cbor:"2,keyasint"`
    MinLP     uint64 `cbor:"3,keyasint,omitempty"`
}

// WithdrawRedeemer for removing liquidity
type WithdrawRedeemer struct {
    Action   uint8  `cbor:"0,keyasint"`
    LPAmount uint64 `cbor:"1,keyasint"`
    MinA     uint64 `cbor:"2,keyasint,omitempty"`
    MinB     uint64 `cbor:"3,keyasint,omitempty"`
}
```

### Index-Based Redeemers

```go
// IndexRedeemer references transaction inputs/outputs by index
type IndexRedeemer struct {
    Action      uint8   `cbor:"0,keyasint"`
    InputIndex  uint32  `cbor:"1,keyasint,omitempty"`
    OutputIndex uint32  `cbor:"2,keyasint,omitempty"`
}

// BatchRedeemer for batched operations
type BatchRedeemer struct {
    Action  uint8    `cbor:"0,keyasint"`
    Indices []uint32 `cbor:"1,keyasint"` // Input indices to process
}
```

### Order Redeemers

```go
// OrderAction represents order-specific actions
type OrderAction uint8

const (
    OrderFill   OrderAction = 0
    OrderPartial OrderAction = 1
    OrderCancel OrderAction = 2
)

// OrderRedeemer for order-based protocols
type OrderRedeemer struct {
    Action     OrderAction `cbor:"0,keyasint"`
    FillAmount uint64      `cbor:"1,keyasint,omitempty"` // For partial fills
}

// ScriptPurpose identifies redeemer purpose
type ScriptPurpose uint8

const (
    PurposeMint   ScriptPurpose = 0
    PurposeSpend  ScriptPurpose = 1
    PurposeReward ScriptPurpose = 2
    PurposeCert   ScriptPurpose = 3
)
```

### Governance Redeemers

```go
// GovRedeemer for governance actions
type GovRedeemer struct {
    Action uint8 `cbor:"0,keyasint"`
}

const (
    GovVote     uint8 = 0
    GovPropose  uint8 = 1
    GovExecute  uint8 = 2
    GovDelegate uint8 = 3
)

// ProposalRedeemer for proposal execution
type ProposalRedeemer struct {
    Action     uint8  `cbor:"0,keyasint"`
    ProposalId uint64 `cbor:"1,keyasint"`
}
```

### Staking Redeemers

```go
// StakeRedeemer for staking operations
type StakeRedeemer struct {
    Action uint8 `cbor:"0,keyasint"`
}

const (
    StakeDeposit  uint8 = 0
    StakeWithdraw uint8 = 1
    StakeClaim    uint8 = 2
)

// ClaimRedeemer for reward claiming
type ClaimRedeemer struct {
    Action      uint8    `cbor:"0,keyasint"`
    EpochsToClaim []uint64 `cbor:"1,keyasint,omitempty"`
}
```

### Constructor Patterns

```go
// Many protocols use Plutus constructor tags
// Tag 121 = Constructor 0
// Tag 122 = Constructor 1
// etc.

// ConstructorRedeemer wraps data with constructor tag
type ConstructorRedeemer struct {
    Constructor uint64
    Fields      []any
}

// NewConstructorRedeemer creates a tagged constructor redeemer
func NewConstructorRedeemer(constructor uint64, fields ...any) *ConstructorRedeemer

// MarshalCBOR outputs proper constructor-tagged CBOR
func (r *ConstructorRedeemer) MarshalCBOR() ([]byte, error)
```

### Utilities

```go
// ParseRedeemer attempts to parse redeemer as known type
func ParseRedeemer(data []byte) (any, string, error)

// IsActionRedeemer checks if redeemer is simple action type
func IsActionRedeemer(data []byte) bool

// GetRedeemerAction extracts action byte from common patterns
func GetRedeemerAction(data []byte) (uint8, bool)

// BuildRedeemer creates CBOR for common redeemer types
func BuildRedeemer(action string, params map[string]any) ([]byte, error)
```

---

## Usage Examples

### Simple Action Redeemer

```go
import "github.com/blinklabs-io/cardano-models/redeemer"

// Create swap redeemer
swapRedeemer := redeemer.NewActionRedeemer(redeemer.ActionSwap)
cbor, err := swapRedeemer.MarshalCBOR()
```

### DEX Swap Redeemer

```go
// Create detailed swap redeemer
swap := &redeemer.SwapRedeemer{
    Action:     0,
    Direction:  0, // A → B
    MinReceive: 1000_000,
}
cbor, err := swap.MarshalCBOR()
```

### Constructor-Tagged Redeemer

```go
// Many protocols use constructor format
// Constructor 0 = Swap, Constructor 1 = Cancel, etc.
r := redeemer.NewConstructorRedeemer(0, minReceive)
cbor, err := r.MarshalCBOR()
```

### Parse Unknown Redeemer

```go
// Try to identify redeemer type
parsed, typeName, err := redeemer.ParseRedeemer(redeemerCbor)
if err != nil {
    // Unknown format
}

switch v := parsed.(type) {
case *redeemer.SwapRedeemer:
    fmt.Printf("Swap direction: %d\n", v.Direction)
case *redeemer.ActionRedeemer:
    fmt.Printf("Action: %d\n", v.Action)
}
```

---

## Implementation Phases

### Phase 1: Core Types

- [ ] Define `ActionRedeemer` type
- [ ] Define action constants
- [ ] Implement CBOR marshaling/unmarshaling

### Phase 2: DEX Redeemers

- [ ] Implement `SwapRedeemer`
- [ ] Implement `DepositRedeemer`
- [ ] Implement `WithdrawRedeemer`

### Phase 3: Index & Batch Redeemers

- [ ] Implement `IndexRedeemer`
- [ ] Implement `BatchRedeemer`

### Phase 4: Order & Governance Redeemers

- [ ] Implement `OrderRedeemer`
- [ ] Implement `GovRedeemer`
- [ ] Implement `ProposalRedeemer`

### Phase 5: Constructor Patterns

- [ ] Implement `ConstructorRedeemer`
- [ ] Handle Plutus constructor tags
- [ ] Add builder helpers

### Phase 6: Utilities

- [ ] Implement `ParseRedeemer()`
- [ ] Implement type detection
- [ ] Add convenience builders

### Phase 7: Testing

- [ ] Unit tests for each redeemer type
- [ ] Tests with real transaction redeemers
- [ ] Round-trip serialization tests

---

## Notes

### Protocol Specificity

Like datums, redeemers are protocol-specific. These patterns represent common abstractions found across many protocols.

### Constructor vs Inline

Some protocols use:
- **Constructor format**: `121([fields...])` - Plutus tagged
- **Inline format**: `{0: action, 1: param}` - Map-based

Both patterns should be supported.

---

## Success Criteria

- [ ] Common DEX redeemer patterns parseable
- [ ] Constructor-tagged redeemers handled correctly
- [ ] Detection heuristics identify common patterns
- [ ] CBOR output matches protocol expectations
- [ ] Comprehensive test coverage

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| gouroboros | Plutus data types |
| cbor/v2 | CBOR encoding |

---

## Related Plans

- [DEFI_DATUMS_PLAN.md](DEFI_DATUMS_PLAN.md) - Datum patterns paired with these redeemers
