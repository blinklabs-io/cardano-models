# Cardano Models Plans

**Last Updated**: 2026-11-04

This directory contains detailed implementation plans for cardano-models.

## Plan Index

### Foundation & Improvements

| Plan | Description | Status | Notes |
|------|-------------|--------|-------|
| [TESTING_PLAN.md](TESTING_PLAN.md) | Testing patterns and utilities | Planned | Per-model tests exist but no shared test infrastructure yet |
| [PATTERNS_PLAN.md](PATTERNS_PLAN.md) | Model design patterns documentation | Planned | |
| [TUNA_MODELS_PLAN.md](TUNA_MODELS_PLAN.md) | TUNA model improvements and tests | In Progress | `tuna_test.go` covers constructor validation only; round-trip tests outstanding (#30) |
| [CDNS_MODELS_PLAN.md](CDNS_MODELS_PLAN.md) | CDNS decoding and JSON serde | In Progress | CBOR tests and decode validation done (#13 closed); JSON serde (#108) outstanding |

### CIP Metadata Standards

| Plan | Description | Status | Notes |
|------|-------------|--------|-------|
| [CIP15_PLAN.md](CIP15_PLAN.md) | Catalyst registration metadata | Planned | #3 |
| [CIP36_PLAN.md](CIP36_PLAN.md) | Catalyst/Voltaire registration | Planned | #4 |
| [CIP68_PLAN.md](CIP68_PLAN.md) | Datum metadata standard | Planned | #210 |

### Application Datums & Redeemers

| Plan | Description | Status | Notes |
|------|-------------|--------|-------|
| [DEFI_DATUMS_PLAN.md](DEFI_DATUMS_PLAN.md) | DeFi protocol datum patterns | Planned | |
| [REDEEMERS_PLAN.md](REDEEMERS_PLAN.md) | Common redeemer patterns | Planned | |

## Plan Status Legend

- **Planned**: Plan created, not yet started
- **In Progress**: Active development underway
- **Complete**: Implementation finished (plan will be removed)

## Related Documentation

- [ROADMAP.md](../../ROADMAP.md) - Main roadmap and milestones
- [GitHub Issues](https://github.com/blinklabs-io/cardano-models/issues) - Open issues to address
