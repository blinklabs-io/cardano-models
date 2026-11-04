# TUNA Models Improvement Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers test coverage and improvements for the existing TUNA (Fortuna mining contract) datum models.

**Related Issue**: [#30](https://github.com/blinklabs-io/cardano-models/issues/30) - Create tests for Tuna models

---

## Current State

### Existing Implementation

The current TUNA implementation in `tuna.go` provides:
- `TunaDatumV1` - Original mining contract state
- `TunaDatumV2` - Extended mining contract state
- Custom CBOR marshaling with Plutus constructor format

### Coverage Gap

No test files exist for TUNA models. Need comprehensive test coverage.

---

## Scope

### Test Coverage

1. **Unit Tests**: Test individual struct serialization
2. **Round-Trip Tests**: CBOR → struct → CBOR
3. **Real-World Tests**: Parse actual TUNA UTxOs from mainnet
4. **Edge Cases**: Handle malformed data gracefully

### Potential Improvements

1. **Validation**: Add field validation
2. **Helpers**: Add convenience methods
3. **Documentation**: Add examples and usage docs

---

## API Design

### Current Structures (Review)

```go
// TunaDatumV1 represents the original TUNA mining state
type TunaDatumV1 struct {
    BlockNumber      int64  `cbor:"0,keyasint"`
    CurrentHash      []byte `cbor:"1,keyasint"`
    LeadingZeros     int64  `cbor:"2,keyasint"`
    DifficultyNumber int64  `cbor:"3,keyasint"`
    EpochTime        int64  `cbor:"4,keyasint"`
    CurrentPosixTime int64  `cbor:"5,keyasint"`
    // ... other fields
}

// TunaDatumV2 represents the extended TUNA mining state
type TunaDatumV2 struct {
    // Extended fields from V2
    // ...
}
```

### Proposed Additions

```go
// Validation
func (d *TunaDatumV1) Validate() error
func (d *TunaDatumV2) Validate() error

// Version detection
func DetectTunaVersion(data []byte) (int, error)

// Convenience methods
func (d *TunaDatumV1) CurrentDifficulty() uint64
func (d *TunaDatumV1) NextDifficulty() uint64
func (d *TunaDatumV1) TimeSinceLastBlock() time.Duration

// Parsing helpers
func ParseTunaDatumV1(data []byte) (*TunaDatumV1, error)
func ParseTunaDatumV2(data []byte) (*TunaDatumV2, error)
func ParseTunaDatum(data []byte) (any, int, error)  // Returns datum and version
```

---

## Test Plan

### Unit Tests

```go
func TestTunaDatumV1_MarshalCBOR(t *testing.T)
func TestTunaDatumV1_UnmarshalCBOR(t *testing.T)
func TestTunaDatumV2_MarshalCBOR(t *testing.T)
func TestTunaDatumV2_UnmarshalCBOR(t *testing.T)
```

### Round-Trip Tests

```go
func TestTunaDatumV1_RoundTrip(t *testing.T) {
    original := &TunaDatumV1{
        BlockNumber:      12345,
        CurrentHash:      []byte{...},
        LeadingZeros:     5,
        // ...
    }

    // Marshal
    cbor, err := original.MarshalCBOR()
    require.NoError(t, err)

    // Unmarshal
    var decoded TunaDatumV1
    err = decoded.UnmarshalCBOR(cbor)
    require.NoError(t, err)

    // Compare
    assert.Equal(t, original, &decoded)
}
```

### Real-World Tests

```go
// Test with actual TUNA UTxO data from mainnet
func TestTunaDatumV1_RealWorld(t *testing.T) {
    // CBOR from actual TUNA contract UTxO
    realCbor := mustDecodeHex("d8799f...")

    datum, err := ParseTunaDatumV1(realCbor)
    require.NoError(t, err)

    // Verify known values
    assert.Equal(t, int64(expectedBlock), datum.BlockNumber)
}

func TestTunaDatumV2_RealWorld(t *testing.T) {
    // Similar for V2
}
```

### Edge Case Tests

```go
func TestTunaDatum_EmptyData(t *testing.T)
func TestTunaDatum_InvalidConstructor(t *testing.T)
func TestTunaDatum_MissingFields(t *testing.T)
func TestTunaDatum_WrongVersion(t *testing.T)
```

---

## Implementation Phases

### Phase 1: Basic Tests

- [ ] Create `tuna_test.go`
- [ ] Add marshal/unmarshal unit tests for V1
- [ ] Add marshal/unmarshal unit tests for V2
- [ ] Add round-trip tests

### Phase 2: Real-World Tests

- [ ] Collect sample TUNA UTxOs from mainnet
- [ ] Add test fixtures with real CBOR data
- [ ] Test parsing against known values
- [ ] Document test data source

### Phase 3: Edge Cases

- [ ] Test empty/nil data handling
- [ ] Test invalid constructor tags
- [ ] Test missing required fields
- [ ] Test wrong version detection

### Phase 4: Validation & Helpers

- [ ] Implement `Validate()` methods
- [ ] Implement `DetectTunaVersion()`
- [ ] Add convenience methods
- [ ] Add parsing helpers

### Phase 5: Documentation

- [ ] Add godoc examples
- [ ] Document TUNA contract format
- [ ] Add usage examples in tests

---

## Test Data Sources

### Mainnet TUNA UTxOs

Collect datum CBOR from:
- TUNA v1 mining contract
- TUNA v2 mining contract
- Various block heights for diversity

### Synthetic Test Data

Create test fixtures for:
- Minimum valid datum
- Maximum values
- Edge cases

---

## Success Criteria

- [ ] 100% function coverage for TUNA models
- [ ] Real-world test data validates correctly
- [ ] Round-trip tests pass
- [ ] Edge cases handled gracefully
- [ ] Clear error messages for invalid data

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| testify | Test assertions |
| cbor/v2 | CBOR encoding |

---

## Related Plans

- [TESTING_PLAN.md](TESTING_PLAN.md) - Overall testing strategy
