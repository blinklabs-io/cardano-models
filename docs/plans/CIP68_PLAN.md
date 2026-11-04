# CIP-68 Datum Metadata Standard Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers implementation of CIP-68 (Datum Metadata Standard) which defines a standardized way to store NFT/token metadata in datums rather than transaction metadata, enabling updates without re-minting.

---

## Scope

### CIP-68 Specification

CIP-68 defines a reference token pattern where:
- A **reference NFT** (label 100) holds the metadata in a datum
- A **user token** (label 222/333/444) is what users hold
- Metadata can be updated by spending the reference UTxO

### Token Labels (Asset Name Prefixes)

| Label | Hex Prefix | Token Type | Description |
|-------|------------|------------|-------------|
| 100 | `000643b0` | Reference NFT | Holds metadata datum |
| 222 | `000de140` | NFT (Non-Fungible) | User-held NFT |
| 333 | `0014df10` | FT (Fungible) | Fungible token |
| 444 | `001bc280` | RFT (Rich Fungible) | FT with extra metadata |

### Datum Structure

```cddl
; Reference datum structure
cip68_datum = {
  metadata: metadata_map,    ; Key-value metadata
  version: uint,             ; Version number (typically 1)
  ? extra: plutus_data       ; Optional extra data
}

metadata_map = { * metadata_key => metadata_value }
metadata_key = bytes / uint
metadata_value = bytes / uint / [* metadata_value] / { * metadata_key => metadata_value }
```

### Standard Metadata Keys

| Key (bytes) | Key (uint) | Description |
|-------------|------------|-------------|
| "name" | - | Token name |
| "image" | - | Image URI |
| "mediaType" | - | MIME type |
| "description" | - | Description |
| "files" | - | Additional files array |

---

## API Design

### Core Structures

```go
package cip68

// TokenLabel represents CIP-68 token label types
type TokenLabel uint32

const (
    TokenLabelReference TokenLabel = 100
    TokenLabelNFT       TokenLabel = 222
    TokenLabelFT        TokenLabel = 333
    TokenLabelRFT       TokenLabel = 444
)

// TokenLabelPrefix returns the 4-byte prefix for the label
func (t TokenLabel) Prefix() []byte

// CIP68Datum represents a CIP-68 reference datum
type CIP68Datum struct {
    Metadata MetadataMap `cbor:"0,keyasint"`
    Version  uint        `cbor:"1,keyasint"`
    Extra    any         `cbor:"2,keyasint,omitempty"`
}

// MetadataMap is the key-value metadata store
type MetadataMap map[any]any

// CIP68Metadata provides typed access to standard fields
type CIP68Metadata struct {
    Name        string
    Image       string
    MediaType   string
    Description string
    Files       []CIP68File
    Extra       map[string]any
}

// CIP68File represents a file entry in metadata
type CIP68File struct {
    Name      string `cbor:"name"`
    MediaType string `cbor:"mediaType"`
    Src       string `cbor:"src"`
}
```

### Token Name Helpers

```go
// ParseAssetName parses a CIP-68 prefixed asset name
func ParseAssetName(assetName []byte) (TokenLabel, []byte, error)

// BuildAssetName creates a CIP-68 prefixed asset name
func BuildAssetName(label TokenLabel, name []byte) []byte

// GetReferenceAssetName converts user token name to reference token name
func GetReferenceAssetName(userAssetName []byte) ([]byte, error)

// GetUserAssetName converts reference token name to user token name
func GetUserAssetName(refAssetName []byte, label TokenLabel) ([]byte, error)

// IsCIP68Token checks if asset name has a CIP-68 prefix
func IsCIP68Token(assetName []byte) bool

// GetTokenLabel extracts the label from a CIP-68 asset name
func GetTokenLabel(assetName []byte) (TokenLabel, bool)
```

### Datum Functions

```go
// ParseCIP68Datum parses a CIP-68 datum from Plutus data
func ParseCIP68Datum(data any) (*CIP68Datum, error)

// ParseCIP68DatumCBOR parses a CIP-68 datum from CBOR bytes
func ParseCIP68DatumCBOR(cbor []byte) (*CIP68Datum, error)

// GetMetadata extracts typed metadata from datum
func (d *CIP68Datum) GetMetadata() (*CIP68Metadata, error)

// GetString gets a string value from metadata
func (d *CIP68Datum) GetString(key string) (string, bool)

// GetUint gets a uint value from metadata
func (d *CIP68Datum) GetUint(key string) (uint64, bool)

// GetBytes gets a bytes value from metadata
func (d *CIP68Datum) GetBytes(key string) ([]byte, bool)

// GetArray gets an array value from metadata
func (d *CIP68Datum) GetArray(key string) ([]any, bool)

// GetMap gets a map value from metadata
func (d *CIP68Datum) GetMap(key string) (map[any]any, bool)
```

### Builder Pattern

```go
// NewCIP68Datum creates a new CIP-68 datum builder
func NewCIP68Datum() *CIP68DatumBuilder

type CIP68DatumBuilder struct {
    metadata MetadataMap
    version  uint
    extra    any
}

func (b *CIP68DatumBuilder) WithName(name string) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) WithImage(uri string) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) WithMediaType(mimeType string) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) WithDescription(desc string) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) WithFile(name, mediaType, src string) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) WithMetadata(key string, value any) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) WithMetadataBytes(key []byte, value any) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) WithVersion(version uint) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) WithExtra(extra any) *CIP68DatumBuilder
func (b *CIP68DatumBuilder) Build() (*CIP68Datum, error)
func (b *CIP68DatumBuilder) BuildCBOR() ([]byte, error)
```

---

## Usage Examples

### Parse CIP-68 Token

```go
import "github.com/blinklabs-io/cardano-models/cip68"

// Check if token is CIP-68
if cip68.IsCIP68Token(assetName) {
    label, name, err := cip68.ParseAssetName(assetName)
    if err != nil {
        // Invalid CIP-68 asset name
    }

    switch label {
    case cip68.TokenLabelReference:
        fmt.Println("Reference token:", hex.EncodeToString(name))
    case cip68.TokenLabelNFT:
        fmt.Println("NFT:", hex.EncodeToString(name))
    case cip68.TokenLabelFT:
        fmt.Println("FT:", hex.EncodeToString(name))
    case cip68.TokenLabelRFT:
        fmt.Println("Rich FT:", hex.EncodeToString(name))
    }
}
```

### Parse Reference Datum

```go
// Parse datum from UTxO
datum, err := cip68.ParseCIP68DatumCBOR(utxo.DatumCBOR())
if err != nil {
    // Not a valid CIP-68 datum
}

// Get typed metadata
meta, err := datum.GetMetadata()
if err != nil {
    // Error extracting metadata
}

fmt.Println("Name:", meta.Name)
fmt.Println("Image:", meta.Image)
fmt.Println("Files:", len(meta.Files))

// Get individual fields
if desc, ok := datum.GetString("description"); ok {
    fmt.Println("Description:", desc)
}
```

### Build CIP-68 Datum

```go
datum, err := cip68.NewCIP68Datum().
    WithName("My NFT").
    WithImage("ipfs://QmXyz...").
    WithMediaType("image/png").
    WithDescription("A unique NFT").
    WithFile("animation", "video/mp4", "ipfs://QmAbc...").
    WithMetadata("rarity", "legendary").
    WithVersion(1).
    Build()
```

### Find Reference Token

```go
// Given a user NFT, find its reference token
refAssetName, err := cip68.GetReferenceAssetName(userAssetName)
if err != nil {
    // Not a CIP-68 token
}

// Look up reference UTxO by policy ID and reference asset name
refUtxo := findUtxo(policyId, refAssetName)
datum, _ := cip68.ParseCIP68DatumCBOR(refUtxo.DatumCBOR())
```

---

## Implementation Phases

### Phase 1: Token Labels

- [ ] Define `TokenLabel` type and constants
- [ ] Implement `Prefix()` for each label
- [ ] Implement `ParseAssetName()`
- [ ] Implement `BuildAssetName()`
- [ ] Implement `IsCIP68Token()`
- [ ] Implement `GetTokenLabel()`

### Phase 2: Reference Token Helpers

- [ ] Implement `GetReferenceAssetName()`
- [ ] Implement `GetUserAssetName()`
- [ ] Add validation for label conversions

### Phase 3: Datum Parsing

- [ ] Define `CIP68Datum` struct
- [ ] Define `MetadataMap` type
- [ ] Implement `ParseCIP68Datum()`
- [ ] Implement `ParseCIP68DatumCBOR()`
- [ ] Handle Plutus data encoding

### Phase 4: Metadata Access

- [ ] Define `CIP68Metadata` struct
- [ ] Implement `GetMetadata()`
- [ ] Implement typed getters (GetString, GetUint, etc.)
- [ ] Handle standard metadata keys

### Phase 5: Builder & Serialization

- [ ] Implement `CIP68DatumBuilder`
- [ ] Implement CBOR marshaling
- [ ] Generate valid Plutus datum CBOR

### Phase 6: Testing

- [ ] Unit tests for asset name parsing
- [ ] Unit tests for datum parsing
- [ ] Tests with real CIP-68 tokens from mainnet
- [ ] Round-trip serialization tests

---

## Success Criteria

- [ ] Parse CIP-68 prefixed asset names correctly
- [ ] Parse CIP-68 reference datums from mainnet
- [ ] Convert between reference and user token names
- [ ] Extract typed metadata from datums
- [ ] Generate valid CIP-68 datums
- [ ] Comprehensive test coverage

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| gouroboros | Plutus data types |
| cbor/v2 | CBOR encoding |

---

## Related Plans

- [CIP25_PLAN.md](CIP25_PLAN.md) - Legacy NFT metadata (already implemented)
