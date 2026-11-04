# CDNS Models Improvement Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers improvements to the existing Cardano DNS (CDNS) domain models including JSON serialization support and decoding improvements.

**Related Issues**:
- [#108](https://github.com/blinklabs-io/cardano-models/issues/108) - Support JSON serde in CDNS
- [#13](https://github.com/blinklabs-io/cardano-models/issues/13) - Improve decoding for CardanoDnsDomain

---

## Current State

### Existing Implementation

The current CDNS implementation in `cdns.go` provides:
- `CardanoDnsDomain` struct for domain datum parsing
- Custom CBOR unmarshaling for Plutus data format
- Constructor format handling (tag 121+)

### Known Issues

1. **No JSON Support**: CBOR-only, can't be used for API responses
2. **Decoding Limitations**: Error handling could be improved
3. **Limited Validation**: Minimal field validation

---

## Scope

### JSON Serialization (#108)

Add full JSON marshaling/unmarshaling support for:
- `CardanoDnsDomain`
- All nested types
- Round-trip compatibility

### Improved Decoding (#13)

- Better error messages
- Handle edge cases in Plutus data
- Graceful handling of malformed data

---

## API Design

### JSON Support

```go
package cdns

// Add JSON tags to all fields
type CardanoDnsDomain struct {
    Origin       string            `json:"origin" cbor:"origin"`
    Records      map[string]string `json:"records" cbor:"records"`
    Owner        []byte            `json:"owner" cbor:"owner"`
    // ... other fields
}

// MarshalJSON implements json.Marshaler
func (d *CardanoDnsDomain) MarshalJSON() ([]byte, error)

// UnmarshalJSON implements json.Unmarshaler
func (d *CardanoDnsDomain) UnmarshalJSON(data []byte) error
```

### Improved Decoding

```go
// DecodeError provides detailed error information
type DecodeError struct {
    Field   string
    Message string
    Cause   error
}

func (e *DecodeError) Error() string
func (e *DecodeError) Unwrap() error

// DecodeDomainDatum decodes with detailed error reporting
func DecodeDomainDatum(data []byte) (*CardanoDnsDomain, error)

// DecodeDomainDatumLenient decodes with lenient parsing
// Returns partial result even if some fields fail
func DecodeDomainDatumLenient(data []byte) (*CardanoDnsDomain, []error)

// ValidateDomain validates domain structure
func ValidateDomain(domain *CardanoDnsDomain) error
```

### Builder Pattern

```go
// NewCardanoDnsDomain creates a domain builder
func NewCardanoDnsDomain() *CardanoDnsDomainBuilder

type CardanoDnsDomainBuilder struct {
    domain *CardanoDnsDomain
}

func (b *CardanoDnsDomainBuilder) WithOrigin(origin string) *CardanoDnsDomainBuilder
func (b *CardanoDnsDomainBuilder) WithRecord(key, value string) *CardanoDnsDomainBuilder
func (b *CardanoDnsDomainBuilder) WithOwner(owner []byte) *CardanoDnsDomainBuilder
func (b *CardanoDnsDomainBuilder) Build() (*CardanoDnsDomain, error)
```

---

## Implementation Phases

### Phase 1: JSON Support

- [ ] Add JSON tags to `CardanoDnsDomain`
- [ ] Implement custom `MarshalJSON()` for byte fields
- [ ] Implement custom `UnmarshalJSON()`
- [ ] Handle hex encoding for byte arrays

### Phase 2: Error Handling

- [ ] Define `DecodeError` type
- [ ] Add field-specific error messages
- [ ] Implement `DecodeDomainDatum()` with detailed errors
- [ ] Implement `DecodeDomainDatumLenient()`

### Phase 3: Validation

- [ ] Implement `ValidateDomain()`
- [ ] Validate origin format
- [ ] Validate record types
- [ ] Validate owner credential

### Phase 4: Builder

- [ ] Implement `CardanoDnsDomainBuilder`
- [ ] Add validation on build
- [ ] Generate valid Plutus datum output

### Phase 5: Testing

- [ ] Unit tests for JSON round-trip
- [ ] Tests for error cases
- [ ] Tests with real CDNS datums from mainnet
- [ ] Edge case handling tests

---

## Success Criteria

- [ ] JSON serialization works correctly
- [ ] Error messages are clear and actionable
- [ ] Lenient parsing recovers partial data
- [ ] Backward compatible with existing code
- [ ] Comprehensive test coverage

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| encoding/json | JSON marshaling |
| encoding/hex | Byte array encoding |

---

## Related Plans

None specific - this is an improvement to existing functionality.
