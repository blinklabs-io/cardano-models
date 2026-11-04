# Model Design Patterns Plan

**Last Updated**: 2026-01-17

## Overview

This plan documents the design patterns and conventions used across cardano-models to ensure consistency in new model implementations.

---

## Core Patterns

### 1. Struct Design

#### Field Tags

All exported structs should have appropriate tags for serialization:

```go
type ModelName struct {
    // CBOR tag with integer key (Plutus data format)
    Field1 string `cbor:"0,keyasint" json:"field1" validate:"required"`

    // Optional fields
    Field2 uint64 `cbor:"1,keyasint,omitempty" json:"field2,omitempty"`

    // Byte arrays (hex in JSON)
    Field3 []byte `cbor:"2,keyasint" json:"field3" validate:"len=32"`
}
```

#### Tag Conventions

| Tag | Usage | Example |
|-----|-------|---------|
| `cbor:"N,keyasint"` | Plutus datum fields | `cbor:"0,keyasint"` |
| `cbor:"N,keyasint,omitempty"` | Optional fields | `cbor:"1,keyasint,omitempty"` |
| `json:"name"` | JSON field name | `json:"fieldName"` |
| `json:"name,omitempty"` | Optional JSON | `json:"fieldName,omitempty"` |
| `validate:"required"` | Required validation | `validate:"required"` |
| `validate:"len=N"` | Length validation | `validate:"len=32"` |

### 2. Flexible Field Types

For fields that may have multiple valid representations, use custom types:

```go
// UriField handles URIs that may be string or array of strings
type UriField struct {
    value string
}

func (u *UriField) UnmarshalCBOR(data []byte) error {
    // Try string first
    var str string
    if err := cbor.Unmarshal(data, &str); err == nil {
        u.value = str
        return nil
    }

    // Try array of strings (chunked)
    var arr []string
    if err := cbor.Unmarshal(data, &arr); err == nil {
        u.value = strings.Join(arr, "")
        return nil
    }

    return errors.New("invalid URI field")
}

func (u *UriField) String() string {
    return u.value
}
```

### 3. CBOR Marshaling

#### Custom Marshal/Unmarshal

For complex types, implement custom CBOR methods:

```go
// MarshalCBOR implements cbor.Marshaler
func (m *ModelName) MarshalCBOR() ([]byte, error) {
    // Build the CBOR representation
    return cbor.Marshal(m.toCBOR())
}

// UnmarshalCBOR implements cbor.Unmarshaler
func (m *ModelName) UnmarshalCBOR(data []byte) error {
    // Parse the CBOR representation
    var raw cbor.RawMessage
    if err := cbor.Unmarshal(data, &raw); err != nil {
        return err
    }
    return m.fromCBOR(raw)
}
```

#### Plutus Constructor Format

For Plutus datums/redeemers, use constructor tags:

```go
const (
    constructorTag0 = 121 // d8799f
    constructorTag1 = 122 // d87a9f
    constructorTag2 = 123 // d87b9f
    // ... up to tag 127 for constructors 0-6
    // then 1280+ for constructors 7+
)

func (m *ModelName) MarshalCBOR() ([]byte, error) {
    // Wrap in constructor tag
    fields := []any{m.Field1, m.Field2}
    return cbor.Marshal(cbor.Tag{
        Number:  constructorTag0,
        Content: fields,
    })
}
```

### 4. JSON Handling

#### Byte Array Encoding

Byte arrays should be hex-encoded in JSON:

```go
func (m *ModelName) MarshalJSON() ([]byte, error) {
    type jsonModel struct {
        Field1 string `json:"field1"`
        Field2 string `json:"field2"` // hex-encoded bytes
    }

    return json.Marshal(jsonModel{
        Field1: m.Field1,
        Field2: hex.EncodeToString(m.Field2),
    })
}
```

### 5. Validation

Use `go-playground/validator` for field validation:

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func (m *ModelName) Validate() error {
    return validate.Struct(m)
}

// Custom validation
func init() {
    validate.RegisterValidation("policyid", validatePolicyId)
}

func validatePolicyId(fl validator.FieldLevel) bool {
    b := fl.Field().Bytes()
    return len(b) == 28 // Blake2b-224 hash
}
```

### 6. Builder Pattern

Provide fluent builders for complex types:

```go
// NewModelName creates a new builder
func NewModelName() *ModelNameBuilder {
    return &ModelNameBuilder{
        model: &ModelName{},
    }
}

type ModelNameBuilder struct {
    model *ModelName
    errs  []error
}

func (b *ModelNameBuilder) WithField1(value string) *ModelNameBuilder {
    b.model.Field1 = value
    return b
}

func (b *ModelNameBuilder) WithField2(value uint64) *ModelNameBuilder {
    b.model.Field2 = value
    return b
}

func (b *ModelNameBuilder) Build() (*ModelName, error) {
    if len(b.errs) > 0 {
        return nil, errors.Join(b.errs...)
    }
    if err := b.model.Validate(); err != nil {
        return nil, err
    }
    return b.model, nil
}
```

### 7. Extraction Functions

For metadata-based models, provide extraction helpers:

```go
// ExtractModelName extracts from transaction metadata
func ExtractModelName(metadata map[uint64]any) (*ModelName, error) {
    raw, ok := metadata[MetadataLabel]
    if !ok {
        return nil, ErrMetadataNotFound
    }

    // Re-encode and decode through CBOR for type safety
    cbor, err := cbor.Marshal(raw)
    if err != nil {
        return nil, fmt.Errorf("marshal metadata: %w", err)
    }

    var model ModelName
    if err := model.UnmarshalCBOR(cbor); err != nil {
        return nil, fmt.Errorf("unmarshal model: %w", err)
    }

    return &model, nil
}

// HasModelName checks if metadata contains this model type
func HasModelName(metadata map[uint64]any) bool {
    _, ok := metadata[MetadataLabel]
    return ok
}
```

### 8. Error Handling

Define specific error types:

```go
var (
    ErrMetadataNotFound = errors.New("metadata label not found")
    ErrInvalidFormat    = errors.New("invalid format")
    ErrMissingField     = errors.New("required field missing")
)

// DecodeError provides detailed error information
type DecodeError struct {
    Field   string
    Message string
    Cause   error
}

func (e *DecodeError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s: %v", e.Field, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (e *DecodeError) Unwrap() error {
    return e.Cause
}
```

---

## Package Organization

### File Naming

| Pattern | Example | Purpose |
|---------|---------|---------|
| `cipNN.go` | `cip25.go` | CIP implementation |
| `cipNN_test.go` | `cip25_test.go` | CIP tests |
| `modelname.go` | `tuna.go` | Application model |
| `types.go` | `types.go` | Shared types |
| `errors.go` | `errors.go` | Error definitions |

### Package Structure

```go
// Package documentation
// Package cip25 provides CIP-25 NFT metadata models.
package cip25

// Constants
const (
    MetadataLabel = 721
)

// Types (exported)
type CIP25Metadata struct { ... }

// Constructors
func NewCIP25Metadata() *CIP25MetadataBuilder { ... }

// Extraction functions
func ExtractCIP25Metadata(...) (*CIP25Metadata, error) { ... }

// Detection functions
func HasCIP25Metadata(...) bool { ... }

// Internal helpers (unexported)
func parseFields(...) { ... }
```

---

## Documentation

### Godoc Standards

```go
// CIP25Metadata represents NFT metadata as defined in CIP-25.
//
// CIP-25 defines a standard for NFT metadata stored in transaction
// metadata under label 721. This struct supports both v1 and v2 formats.
//
// Example:
//
//	meta := cip25.NewCIP25Metadata().
//	    WithName("My NFT").
//	    WithImage("ipfs://...").
//	    Build()
type CIP25Metadata struct {
    // Name is the display name of the NFT (required).
    Name string `json:"name" validate:"required"`

    // Image is the URI of the NFT image (required).
    // Can be IPFS, HTTPS, or other supported protocols.
    Image string `json:"image" validate:"required"`
}
```

### Example Tests

```go
func ExampleCIP25Metadata() {
    meta, _ := cip25.NewCIP25Metadata().
        WithName("My NFT").
        WithImage("ipfs://QmXyz...").
        WithDescription("A cool NFT").
        Build()

    fmt.Println(meta.Name)
    // Output: My NFT
}
```

---

## Implementation Checklist

When adding a new model:

- [ ] Define struct with appropriate tags
- [ ] Implement `MarshalCBOR()` / `UnmarshalCBOR()`
- [ ] Implement `MarshalJSON()` / `UnmarshalJSON()` if needed
- [ ] Implement `Validate()` method
- [ ] Create builder if struct is complex
- [ ] Add extraction function for metadata models
- [ ] Add detection function (`Has...()`)
- [ ] Write unit tests
- [ ] Write round-trip tests
- [ ] Add real-world test fixtures
- [ ] Add godoc with examples
- [ ] Update README if needed

---

## Success Criteria

- [ ] All models follow consistent patterns
- [ ] Documentation exists for all patterns
- [ ] New contributors can follow patterns easily
- [ ] Code review checklist based on patterns

---

## Related Plans

- [TESTING_PLAN.md](TESTING_PLAN.md) - Test patterns for models
