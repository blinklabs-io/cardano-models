# DeFi Datums Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers common DeFi protocol datum patterns used across Cardano decentralized finance applications including DEXes, lending protocols, and staking contracts.

---

## Scope

### Protocol Categories

| Category | Examples | Use Cases |
|----------|----------|-----------|
| DEX | Minswap, SundaeSwap, Wingriders | Liquidity pools, swaps |
| Lending | Liqwid, Lenfi | Lending/borrowing positions |
| Staking | ISPO, NFT staking | Delegated staking |
| Orderbook | Axo, GeniusYield | Limit orders |

### Common Patterns

1. **Pool State Datums**: Track liquidity pool state
2. **Order Datums**: Represent user orders/positions
3. **Config Datums**: Protocol configuration
4. **Receipt Datums**: Proof of action/position

---

## API Design

### DEX Pool Patterns

```go
package defi

// LiquidityPoolDatum represents common LP datum fields
type LiquidityPoolDatum struct {
    TokenA       AssetClass `cbor:"0,keyasint"`
    TokenB       AssetClass `cbor:"1,keyasint"`
    ReserveA     uint64     `cbor:"2,keyasint"`
    ReserveB     uint64     `cbor:"3,keyasint"`
    TotalLP      uint64     `cbor:"4,keyasint"`
    Fee          Rational   `cbor:"5,keyasint"`
    ExtraFields  any        `cbor:"6,keyasint,omitempty"`
}

// AssetClass identifies a native token
type AssetClass struct {
    PolicyId  []byte `cbor:"0,keyasint"`
    AssetName []byte `cbor:"1,keyasint"`
}

// Rational represents a fraction
type Rational struct {
    Numerator   int64 `cbor:"0,keyasint"`
    Denominator int64 `cbor:"1,keyasint"`
}

// SwapOrderDatum represents a swap request
type SwapOrderDatum struct {
    Owner         []byte     `cbor:"0,keyasint"` // Payment credential
    OfferAsset    AssetClass `cbor:"1,keyasint"`
    OfferAmount   uint64     `cbor:"2,keyasint"`
    RequestAsset  AssetClass `cbor:"3,keyasint"`
    MinReceive    uint64     `cbor:"4,keyasint"`
    Deadline      uint64     `cbor:"5,keyasint,omitempty"`
}
```

### Lending Patterns

```go
// LendingPositionDatum represents a lending position
type LendingPositionDatum struct {
    Owner           []byte     `cbor:"0,keyasint"`
    CollateralAsset AssetClass `cbor:"1,keyasint"`
    CollateralAmt   uint64     `cbor:"2,keyasint"`
    DebtAsset       AssetClass `cbor:"3,keyasint"`
    DebtAmount      uint64     `cbor:"4,keyasint"`
    InterestRate    Rational   `cbor:"5,keyasint"`
    StartTime       uint64     `cbor:"6,keyasint"`
}

// LendingPoolDatum represents a lending pool state
type LendingPoolDatum struct {
    Asset            AssetClass `cbor:"0,keyasint"`
    TotalDeposits    uint64     `cbor:"1,keyasint"`
    TotalBorrowed    uint64     `cbor:"2,keyasint"`
    UtilizationRate  Rational   `cbor:"3,keyasint"`
    InterestModel    []byte     `cbor:"4,keyasint"` // Script hash
}
```

### Staking Patterns

```go
// StakingPositionDatum represents a staked position
type StakingPositionDatum struct {
    Owner         []byte     `cbor:"0,keyasint"`
    StakedAsset   AssetClass `cbor:"1,keyasint"`
    StakedAmount  uint64     `cbor:"2,keyasint"`
    StakeTime     uint64     `cbor:"3,keyasint"`
    LockPeriod    uint64     `cbor:"4,keyasint,omitempty"` // 0 = no lock
    RewardsClaimed uint64    `cbor:"5,keyasint"`
}

// StakingPoolDatum represents a staking pool
type StakingPoolDatum struct {
    StakeAsset        AssetClass `cbor:"0,keyasint"`
    RewardAsset       AssetClass `cbor:"1,keyasint"`
    TotalStaked       uint64     `cbor:"2,keyasint"`
    RewardsPerToken   Rational   `cbor:"3,keyasint"`
    LastUpdateTime    uint64     `cbor:"4,keyasint"`
}
```

### Orderbook Patterns

```go
// LimitOrderDatum represents a limit order
type LimitOrderDatum struct {
    Owner        []byte     `cbor:"0,keyasint"`
    OfferAsset   AssetClass `cbor:"1,keyasint"`
    OfferAmount  uint64     `cbor:"2,keyasint"`
    PriceAsset   AssetClass `cbor:"3,keyasint"`
    PriceAmount  uint64     `cbor:"4,keyasint"` // Price per unit
    Side         uint8      `cbor:"5,keyasint"` // 0=buy, 1=sell
    Expiration   uint64     `cbor:"6,keyasint,omitempty"`
}
```

### Generic Helpers

```go
// ParseDatum attempts to parse datum as common DeFi types
func ParseDatum(data []byte) (any, string, error)

// IsLiquidityPool checks if datum matches LP pattern
func IsLiquidityPool(data []byte) bool

// IsSwapOrder checks if datum matches swap order pattern
func IsSwapOrder(data []byte) bool

// AssetClassFromString parses "policyId.assetName" format
func AssetClassFromString(s string) (*AssetClass, error)

// AssetClassToString formats as "policyId.assetName"
func (a *AssetClass) String() string

// IsADA returns true if asset class represents ADA
func (a *AssetClass) IsADA() bool
```

---

## Usage Examples

### Parse Pool Datum

```go
import "github.com/blinklabs-io/cardano-models/defi"

// Parse liquidity pool datum
poolDatum, err := defi.ParseLiquidityPoolDatum(utxo.Datum())
if err != nil {
    // Not a standard pool datum
}

// Calculate price
priceAtoB := float64(poolDatum.ReserveB) / float64(poolDatum.ReserveA)
fmt.Printf("Price: %.6f %s per %s\n", priceAtoB,
    poolDatum.TokenB.String(), poolDatum.TokenA.String())
```

### Build Swap Order

```go
order := &defi.SwapOrderDatum{
    Owner:        myPaymentCredential,
    OfferAsset:   defi.ADAAssetClass(),
    OfferAmount:  100_000_000, // 100 ADA
    RequestAsset: *tokenAssetClass,
    MinReceive:   1000_000,    // Minimum tokens to receive
}

datumCbor, err := order.MarshalCBOR()
```

---

## Implementation Phases

### Phase 1: Core Types

- [ ] Define `AssetClass` type
- [ ] Define `Rational` type
- [ ] Add ADA helpers
- [ ] Add string conversion utilities

### Phase 2: DEX Patterns

- [ ] Implement `LiquidityPoolDatum`
- [ ] Implement `SwapOrderDatum`
- [ ] Add pool calculation helpers

### Phase 3: Lending Patterns

- [ ] Implement `LendingPositionDatum`
- [ ] Implement `LendingPoolDatum`
- [ ] Add interest calculation helpers

### Phase 4: Staking Patterns

- [ ] Implement `StakingPositionDatum`
- [ ] Implement `StakingPoolDatum`
- [ ] Add reward calculation helpers

### Phase 5: Orderbook Patterns

- [ ] Implement `LimitOrderDatum`
- [ ] Add order matching helpers

### Phase 6: Detection & Parsing

- [ ] Implement `ParseDatum()`
- [ ] Implement type detection functions
- [ ] Add heuristic matching

### Phase 7: Testing

- [ ] Unit tests with synthetic data
- [ ] Tests with real DeFi UTxOs from mainnet
- [ ] Cross-reference with protocol documentation

---

## Notes

### Protocol Specificity

These patterns represent **common abstractions**, not exact protocol implementations. Each DeFi protocol has its own specific datum format. This library provides:

1. **Generic patterns** that cover common fields
2. **Detection heuristics** for identifying datum types
3. **Utilities** for working with DeFi data

For exact protocol support, protocol-specific packages may be needed.

### Versioning

DeFi protocols frequently upgrade their contracts. Each pattern should:
- Document which protocols/versions it's compatible with
- Handle version detection where possible
- Fail gracefully on unknown formats

---

## Success Criteria

- [ ] Common DEX pool formats parseable
- [ ] Swap order patterns work with major DEXes
- [ ] Lending position patterns work with major lending protocols
- [ ] Detection heuristics have reasonable accuracy
- [ ] Comprehensive test coverage with real data

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| gouroboros | Plutus data types |
| cbor/v2 | CBOR encoding |

---

## Related Plans

- [REDEEMERS_PLAN.md](REDEEMERS_PLAN.md) - Redeemer patterns for DeFi actions
