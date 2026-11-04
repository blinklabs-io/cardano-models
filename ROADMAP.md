# Cardano Models Roadmap

**Last Updated**: 2026-02-22

## Vision

Cardano Models serves as the centralized library for Cardano application-layer data structures, specifically models that are not core ledger primitives. This includes metadata formats (CIP standards), application datums, redeemers, and other community-defined data shapes used across the Cardano ecosystem.

**Scope Boundary**: This library focuses on *application-layer* structures. Core ledger primitives (blocks, transactions, ledger state) belong in gouroboros. Mock generation belongs in ouroboros-mock.

---

## Current State (v0.7.0)

### Implemented Models

| Model | Standard | Description | Status |
|-------|----------|-------------|--------|
| CIP-20 | Transaction Messages | Label 674 metadata | Complete (with tests) |
| CIP-25 v1 | NFT Metadata | Label 721 basic | Complete (with tests) |
| CIP-25 v2 | NFT Metadata | Label 721 extended | Complete (with tests) |
| CIP-27 | Royalty Metadata | Label 777 | Complete (with tests) |
| CDNS | Cardano DNS | Domain resolution | Partial (CBOR serde + tests done, JSON serde outstanding) |
| TUNA v1 | Mining Contract | Datum format | Implementation complete, no tests |
| TUNA v2 | Mining Contract | Extended datum | Implementation complete, no tests |

### Open Issues

| Issue | Title | Priority |
|-------|-------|----------|
| #108 | Support JSON serde in CDNS | Medium |
| #30 | Create tests for Tuna models | High |
| #13 | Improve decoding for CardanoDnsDomain | Medium |
| #6 | CIP-83 Encrypted messages | High |
| #5 | CIP-60 Music token metadata | Medium |
| #4 | CIP-36 Catalyst/Voltaire registration | High |
| #3 | CIP-15 Catalyst registrations | Medium |

---

## Roadmap Phases

### Phase 1: Foundation & Existing Model Improvements

**Goal**: Stabilize existing models and establish patterns for new additions.

| Milestone | Description | Plan | Status |
|-----------|-------------|------|--------|
| M1.1 | Add comprehensive tests for TUNA models | [TUNA_MODELS_PLAN.md](docs/plans/TUNA_MODELS_PLAN.md) | Not started |
| M1.2 | Improve CDNS decoding and JSON serde | [CDNS_MODELS_PLAN.md](docs/plans/CDNS_MODELS_PLAN.md) | In progress (CBOR tests done, JSON serde + improved decoding remaining) |
| M1.3 | Establish testing patterns and utilities | [TESTING_PLAN.md](docs/plans/TESTING_PLAN.md) | Not started (tests exist for CIP-20/25/27/CDNS but no shared test infrastructure) |
| M1.4 | Document existing model patterns | [PATTERNS_PLAN.md](docs/plans/PATTERNS_PLAN.md) | Not started |

### Phase 2: Catalyst & Governance Metadata

**Goal**: Support Catalyst voting and Voltaire governance metadata.

| Milestone | Description | Plan | Status |
|-----------|-------------|------|--------|
| M2.1 | CIP-15 Catalyst registration metadata | [CIP15_PLAN.md](docs/plans/CIP15_PLAN.md) | Not started |
| M2.2 | CIP-36 Catalyst/Voltaire registration | [CIP36_PLAN.md](docs/plans/CIP36_PLAN.md) | Not started |

### Phase 3: Extended NFT & Token Standards

**Goal**: Support advanced NFT and token metadata formats.

| Milestone | Description | Plan | Status |
|-----------|-------------|------|--------|
| M3.1 | CIP-68 Datum Metadata Standard | [CIP68_PLAN.md](docs/plans/CIP68_PLAN.md) | Not started |
| M3.2 | CIP-60 Music Token Metadata | [CIP60_PLAN.md](docs/plans/CIP60_PLAN.md) | Not started |
| M3.3 | CIP-83 Encrypted Messages | [CIP83_PLAN.md](docs/plans/CIP83_PLAN.md) | Not started |

### Phase 4: DeFi & Application Datums

**Goal**: Support common DeFi protocol datum formats.

| Milestone | Description | Plan | Status |
|-----------|-------------|------|--------|
| M4.1 | Common DEX datum patterns | [DEFI_DATUMS_PLAN.md](docs/plans/DEFI_DATUMS_PLAN.md) | Not started |
| M4.2 | Lending protocol datums | [DEFI_DATUMS_PLAN.md](docs/plans/DEFI_DATUMS_PLAN.md) | Not started |
| M4.3 | Staking/delegation datums | [DEFI_DATUMS_PLAN.md](docs/plans/DEFI_DATUMS_PLAN.md) | Not started |

### Phase 5: Redeemer Patterns

**Goal**: Document and support common redeemer formats.

| Milestone | Description | Plan | Status |
|-----------|-------------|------|--------|
| M5.1 | Common redeemer patterns | [REDEEMERS_PLAN.md](docs/plans/REDEEMERS_PLAN.md) | Not started |
| M5.2 | Protocol-specific redeemers | [REDEEMERS_PLAN.md](docs/plans/REDEEMERS_PLAN.md) | Not started |

---

## Implementation Priority

### High Priority (Address Open Issues)
1. **CIP-36** (#4) - Catalyst/Voltaire critical for governance
2. **CIP-83** (#6) - Encrypted messaging for privacy
3. **TUNA tests** (#30) - Ensure existing code quality

### Medium Priority (Community Standards)
4. **CIP-68** - Widely adopted datum metadata standard
5. **CIP-15** (#3) - Legacy Catalyst support
6. **CIP-60** (#5) - Music industry adoption
7. **CDNS improvements** (#108, #13)

### Lower Priority (Future Standards)
8. **DeFi datums** - As protocols stabilize
9. **Redeemer patterns** - Documentation and utilities

---

## Design Principles

### 1. Separation of Concerns
- **cardano-models**: Application-layer structures (metadata, datums, redeemers)
- **gouroboros**: Core ledger primitives (blocks, transactions, certificates)
- **ouroboros-mock**: Mock generation for testing

### 2. Flexible Serialization
- Support both CBOR and JSON where applicable
- Use custom field types for flexible parsing (e.g., `UriField`, `DescField`)
- Handle on-chain (CBOR) and off-chain (JSON) representations

### 3. Validation
- Use `go-playground/validator` for field validation
- Provide clear error messages for invalid data
- Support partial parsing where possible

### 4. CBOR Compatibility
- Follow CDDL specifications exactly
- Use constructor patterns for Plutus datums
- Support both raw and tagged CBOR

### 5. Testing
- Comprehensive unit tests with real-world examples
- Round-trip serialization tests (CBOR → struct → CBOR)
- Cross-reference with cardano-cli/other implementations

---

### In-Progress Branches

The following local branches contain unmerged work:

| Branch | Description |
|--------|-------------|
| `build/makefile` | Make targets for formatting and tests |
| `fix/prealloc` | Prealloc fix in tuna for golangci-lint |
| `fix/remove-copy-var` | Remove variable copying |
| `test/field-alignment` | Optimize field alignment in test struct |

Note: The `feat/cip25` branch is stale (CIP-25 is already on main).

---

## Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| gouroboros | v0.153.1 | Ledger types, CBOR utilities |
| cbor/v2 | v2.9.0 | CBOR encoding/decoding |
| validator/v10 | v10.30.1 | Struct validation |
| testify | v1.11.1 | Testing assertions |

---

## Related Projects

| Project | Relationship |
|---------|--------------|
| [gouroboros](https://github.com/blinklabs-io/gouroboros) | Core ledger types, CBOR utilities |
| [ouroboros-mock](https://github.com/blinklabs-io/ouroboros-mock) | Mock generation using these models |
| [adder](https://github.com/blinklabs-io/adder) | Event processing using these models |
| [Dingo](https://github.com/blinklabs-io/dingo) | Node implementation using these models |
| [tx-submit-api](https://github.com/blinklabs-io/tx-submit-api) | Transaction submission |
| [bursa](https://github.com/blinklabs-io/bursa) | Wallet using these models |

---

## Contributing

When adding new models:

1. Check if a CIP exists for the format
2. Create a plan document in `plans/` first
3. Follow existing patterns (see [PATTERNS_PLAN.md](docs/plans/PATTERNS_PLAN.md))
4. Include comprehensive tests with real-world data
5. Support both CBOR and JSON where applicable
6. Update this roadmap with the new model
