# Testing Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers the testing strategy and infrastructure for cardano-models, ensuring high quality and reliability across all model implementations.

---

## Scope

### Testing Categories

| Category | Description | Purpose |
|----------|-------------|---------|
| Unit Tests | Individual function/struct tests | Correctness |
| Round-Trip Tests | Serialize → deserialize → compare | Data integrity |
| Real-World Tests | Parse actual blockchain data | Compatibility |
| Fuzz Tests | Random/malformed input | Robustness |
| Benchmark Tests | Performance measurement | Optimization |

---

## Testing Infrastructure

### Directory Structure

```
cardano-models/
├── testdata/
│   ├── cip20/           # CIP-20 test fixtures
│   ├── cip25/           # CIP-25 test fixtures
│   ├── cip68/           # CIP-68 test fixtures
│   ├── cdns/            # CDNS test fixtures
│   ├── tuna/            # TUNA test fixtures
│   └── defi/            # DeFi datum fixtures
├── *_test.go            # Test files
└── testutil/            # Test utilities
```

### Test Utilities Package

```go
package testutil

// CBOR Helpers
func MustDecodeHex(hex string) []byte
func MustEncodeCBOR(v any) []byte
func LoadTestFixture(path string) []byte
func SaveTestFixture(path string, data []byte)

// Comparison Helpers
func AssertCBOREqual(t *testing.T, expected, actual []byte)
func AssertJSONEqual(t *testing.T, expected, actual string)
func AssertStructEqual(t *testing.T, expected, actual any)

// Random Generation
func RandomBytes(n int) []byte
func RandomHex(n int) string
func RandomPolicyId() []byte
func RandomAssetName() []byte
func RandomAddress() []byte

// Test Data Fetching (optional)
func FetchMainnetDatum(txHash string, outputIdx int) ([]byte, error)
func FetchMainnetMetadata(txHash string) (map[uint64]any, error)
```

### Fixture Management

```go
// fixtures.go - manage test fixtures

type Fixture struct {
    Name        string
    Description string
    Source      string // e.g., "mainnet:txhash"
    CBOR        []byte
    Expected    any    // Expected parsed result
}

// LoadFixtures loads all fixtures for a model type
func LoadFixtures(modelType string) ([]Fixture, error)

// GenerateFixture creates a fixture from live data
func GenerateFixture(source string) (*Fixture, error)
```

---

## Test Patterns

### Unit Test Pattern

```go
func TestModelName_MethodName(t *testing.T) {
    tests := []struct {
        name     string
        input    any
        expected any
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedResult,
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    invalidInput,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := SomeMethod(tt.input)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Round-Trip Test Pattern

```go
func TestModelName_RoundTrip(t *testing.T) {
    original := &ModelName{
        Field1: "value1",
        Field2: 42,
    }

    // Marshal
    cbor, err := original.MarshalCBOR()
    require.NoError(t, err)

    // Unmarshal
    var decoded ModelName
    err = decoded.UnmarshalCBOR(cbor)
    require.NoError(t, err)

    // Compare
    assert.Equal(t, original, &decoded)

    // Re-marshal and compare CBOR
    remarshaled, err := decoded.MarshalCBOR()
    require.NoError(t, err)
    assert.Equal(t, cbor, remarshaled)
}
```

### Real-World Test Pattern

```go
func TestModelName_RealWorld(t *testing.T) {
    fixtures, err := testutil.LoadFixtures("modelname")
    require.NoError(t, err)

    for _, fixture := range fixtures {
        t.Run(fixture.Name, func(t *testing.T) {
            result, err := ParseModelName(fixture.CBOR)
            require.NoError(t, err)
            assert.Equal(t, fixture.Expected, result)
        })
    }
}
```

### Fuzz Test Pattern

```go
func FuzzModelName_UnmarshalCBOR(f *testing.F) {
    // Seed with valid examples
    f.Add(validCBOR1)
    f.Add(validCBOR2)

    f.Fuzz(func(t *testing.T, data []byte) {
        var m ModelName
        // Should not panic on any input
        _ = m.UnmarshalCBOR(data)
    })
}
```

### Benchmark Pattern

```go
func BenchmarkModelName_UnmarshalCBOR(b *testing.B) {
    data := testutil.MustDecodeHex("...")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var m ModelName
        _ = m.UnmarshalCBOR(data)
    }
}

func BenchmarkModelName_MarshalCBOR(b *testing.B) {
    m := &ModelName{...}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = m.MarshalCBOR()
    }
}
```

---

## Coverage Requirements

### Minimum Coverage Targets

| Package | Line Coverage | Branch Coverage |
|---------|--------------|-----------------|
| cip20 | 90% | 85% |
| cip25 | 90% | 85% |
| cip68 | 90% | 85% |
| cdns | 90% | 85% |
| tuna | 90% | 85% |
| defi | 80% | 75% |
| redeemer | 80% | 75% |

### Required Test Cases

Each model must have:
1. Valid input parsing
2. Invalid input rejection
3. Round-trip serialization
4. At least one real-world fixture
5. Edge cases (empty, max values, etc.)

---

## CI Integration

### Test Commands

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Run fuzz tests (limited time)
go test -fuzz=FuzzModelName -fuzztime=30s ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### GitHub Actions Workflow

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...

      - name: Check coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
          echo "Coverage: $coverage"

      - name: Run benchmarks
        run: go test -bench=. -benchmem ./... | tee benchmark.txt
```

---

## Implementation Phases

### Phase 1: Infrastructure

- [ ] Create `testutil/` package
- [ ] Create `testdata/` directory structure
- [ ] Add hex decode helpers
- [ ] Add CBOR comparison helpers

### Phase 2: Existing Model Tests

- [ ] Add CIP-20 tests (if missing)
- [ ] Add CIP-25 tests (if missing)
- [ ] Add CIP-27 tests (if missing)
- [ ] Add CDNS tests
- [ ] Add TUNA tests

### Phase 3: Fixtures

- [ ] Collect real-world CBOR fixtures
- [ ] Document fixture sources
- [ ] Add fixture validation

### Phase 4: Fuzz Testing

- [ ] Add fuzz tests for CBOR unmarshaling
- [ ] Add fuzz tests for JSON unmarshaling
- [ ] Verify no panics on malformed input

### Phase 5: Benchmarks

- [ ] Add marshal benchmarks
- [ ] Add unmarshal benchmarks
- [ ] Document baseline performance

### Phase 6: CI/CD

- [ ] Add GitHub Actions workflow
- [ ] Add coverage reporting
- [ ] Add benchmark tracking

---

## Success Criteria

- [ ] All models have unit tests
- [ ] All models have round-trip tests
- [ ] All models have real-world fixtures
- [ ] Coverage meets minimum targets
- [ ] No panics on fuzz testing
- [ ] CI runs on all PRs

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| testify | Assertions |
| go test | Standard testing |
| go fuzz | Fuzz testing |

---

## Related Plans

- [TUNA_MODELS_PLAN.md](TUNA_MODELS_PLAN.md) - TUNA-specific tests
- [PATTERNS_PLAN.md](PATTERNS_PLAN.md) - Model patterns to test
