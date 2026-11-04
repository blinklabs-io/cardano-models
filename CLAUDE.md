# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go library (`package models`) for serializing/deserializing Cardano application-layer data structures not part of the blockchain specification—primarily CIP (Cardano Improvement Proposal) standards. Part of the Blink Labs ecosystem.

## Commands

```bash
# Run tests (includes mod-tidy, uses -race flag)
make test

# Run a single test
go test -v -race -run TestFunctionName ./...

# Lint (uses extensive .golangci.yml config with 40+ linters)
golangci-lint run

# Format code
make format        # gofmt
make golines       # 80-char lines with chain-split-dots

# Update dependencies
make mod-tidy
```

## Architecture

**Single package structure**: All code lives in the root `models` package.

**Models implemented**:
- CIP-20 (`cip20.go`): Transaction messages (label 674)
- CIP-25 v1/v2 (`cip25.go`): NFT metadata (label 721)
- CIP-27 (`cip27.go`): Royalty metadata (label 777)
- CDNS (`cdns.go`): Cardano DNS
- TUNA v1/v2 (`tuna.go`): Mining smart contract datum

**Key dependencies**:
- `github.com/blinklabs-io/gouroboros` - CBOR utilities and Plutus constructors
- `github.com/fxamacker/cbor/v2` - CBOR encoding
- `github.com/go-playground/validator/v10` - Struct validation

## Code Patterns

**Flexible field types** handle both single values and arrays in JSON/CBOR:
- `HexBytes` - String that may be hex-encoded bytes
- `UriField` - Single URI or array of URIs
- `DescField` - Single description or array
- `AddrField` - Single address or array

**CBOR with Plutus constructors** (for smart contract datums like TUNA):
```go
func (t *Type) MarshalCBOR() ([]byte, error) {
    tmp := cbor.NewConstructor(0, cbor.IndefLengthList{...fields...})
    return cbor.Encode(&tmp)
}

func (t *Type) UnmarshalCBOR(data []byte) error {
    var tmpConstr cbor.Constructor
    if _, err := cbor.Decode(data, &tmpConstr); err != nil {
        return err
    }
    type tType Type
    var tmp tType
    if _, err := cbor.Decode(tmpConstr.FieldsCbor(), &tmp); err != nil {
        return err
    }
    *t = Type(tmp)
    return nil
}
```

**Forward compatibility**: Use `Extra map[string]any` fields to preserve unknown JSON properties during round-trip.

**Validation**: Use `validator` struct tags and `NewXxx()` constructors that validate before returning.

## Testing Conventions

- Table-driven tests with CBOR hex strings and JSON data
- Round-trip serialization tests (CBOR → struct → CBOR, JSON → struct → JSON)
- Use `testify/require` for assertions
- Tests run with `-race` flag

## Commit Convention

Enforced conventional commits: `type(scope): description` (e.g., `fix:`, `feat:`, `chore(deps):`, `build:`)
