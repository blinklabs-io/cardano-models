# CIP-60 Music Token Metadata Plan

**Last Updated**: 2026-01-17

## Overview

This plan covers implementation of CIP-60 (Music Token Metadata) which extends CIP-25 with music industry-specific metadata fields for audio/music NFTs.

**Related Issue**: [#5](https://github.com/blinklabs-io/cardano-models/issues/5)

---

## Scope

### CIP-60 Specification

CIP-60 defines additional metadata fields for music tokens under label 721, compatible with CIP-25 but adding music-specific properties.

### Metadata Label

| Label | Purpose |
|-------|---------|
| 721 | NFT metadata (shared with CIP-25) |

### Version Identifier

CIP-60 uses `music_metadata_version` field to identify the version:
- Version 1: Initial spec (deprecated)
- Version 2: Extended spec
- Version 3: Current spec with full feature set

### Music-Specific Fields

```cddl
music_metadata = {
  ; CIP-25 base fields
  "name": text,
  "image": uri,

  ; Music version identifier
  "music_metadata_version": uint,  ; 1, 2, or 3

  ; Music-specific fields
  ? "release": release_info,
  ? "files": [* music_file],
  ? "artists": [* artist],
  ? "genres": [* text],
  ? "copyright": text,
  ? "lyrics": text / uri,
  ? "contributing_artists": [* artist],
  ? "series": text,
  ? "collection": text,
  ? "set": text,
  ? "mood": text,
  ? "explicit": bool,
  ? "isrc": text,              ; International Standard Recording Code
  ? "iswc": text,              ; International Standard Musical Work Code
  ? "links": { * text => uri }
}

release_info = {
  ? "release_type": text,      ; "Single", "EP", "Album"
  ? "release_title": text,
  ? "release_date": text,      ; ISO 8601 date
  ? "publication_date": text,
  ? "catalog_number": text,
  ? "barcode": text,
  ? "country": text
}

artist = {
  "name": text,
  ? "id": text,                ; Platform-specific ID
  ? "links": { * text => uri }
}

music_file = {
  "name": text,
  "mediaType": text,           ; MIME type
  "src": uri,
  ? "song_title": text,
  ? "song_duration": text,     ; ISO 8601 duration or "MM:SS"
  ? "track_number": uint,
  ? "artists": [* artist],
  ? "genres": [* text],
  ? "copyright": text,
  ? "lyrics": text / uri
}
```

---

## API Design

### Core Structures

```go
package cip60

// MusicMetadataVersion identifies the CIP-60 version
type MusicMetadataVersion uint

const (
    MusicMetadataV1 MusicMetadataVersion = 1 // Deprecated
    MusicMetadataV2 MusicMetadataVersion = 2
    MusicMetadataV3 MusicMetadataVersion = 3
)

// CIP60Metadata represents music token metadata
type CIP60Metadata struct {
    // CIP-25 base fields
    Name        string   `json:"name" validate:"required"`
    Image       string   `json:"image" validate:"required"`
    MediaType   string   `json:"mediaType,omitempty"`
    Description string   `json:"description,omitempty"`

    // Music version identifier
    MusicMetadataVersion MusicMetadataVersion `json:"music_metadata_version" validate:"required,min=1,max=3"`

    // Music-specific fields
    Release             *ReleaseInfo `json:"release,omitempty"`
    Files               []MusicFile  `json:"files,omitempty"`
    Artists             []Artist     `json:"artists,omitempty"`
    Genres              []string     `json:"genres,omitempty"`
    Copyright           string       `json:"copyright,omitempty"`
    Lyrics              string       `json:"lyrics,omitempty"`
    ContributingArtists []Artist     `json:"contributing_artists,omitempty"`
    Series              string       `json:"series,omitempty"`
    Collection          string       `json:"collection,omitempty"`
    Set                 string       `json:"set,omitempty"`
    Mood                string       `json:"mood,omitempty"`
    Explicit            *bool        `json:"explicit,omitempty"`
    ISRC                string       `json:"isrc,omitempty"`
    ISWC                string       `json:"iswc,omitempty"`
    Links               Links        `json:"links,omitempty"`
}

// ReleaseInfo contains album/release metadata
type ReleaseInfo struct {
    ReleaseType     string `json:"release_type,omitempty"`
    ReleaseTitle    string `json:"release_title,omitempty"`
    ReleaseDate     string `json:"release_date,omitempty"`
    PublicationDate string `json:"publication_date,omitempty"`
    CatalogNumber   string `json:"catalog_number,omitempty"`
    Barcode         string `json:"barcode,omitempty"`
    Country         string `json:"country,omitempty"`
}

// Artist represents a music artist
type Artist struct {
    Name  string `json:"name" validate:"required"`
    ID    string `json:"id,omitempty"`
    Links Links  `json:"links,omitempty"`
}

// MusicFile represents an audio file with metadata
type MusicFile struct {
    Name         string   `json:"name" validate:"required"`
    MediaType    string   `json:"mediaType" validate:"required"`
    Src          string   `json:"src" validate:"required"`
    SongTitle    string   `json:"song_title,omitempty"`
    SongDuration string   `json:"song_duration,omitempty"`
    TrackNumber  *uint    `json:"track_number,omitempty"`
    Artists      []Artist `json:"artists,omitempty"`
    Genres       []string `json:"genres,omitempty"`
    Copyright    string   `json:"copyright,omitempty"`
    Lyrics       string   `json:"lyrics,omitempty"`
}

// Links is a map of platform names to URIs
type Links map[string]string
```

### Extraction Functions

```go
// ExtractCIP60Metadata extracts music metadata from 721 label
func ExtractCIP60Metadata(metadata map[uint64]any, policyId, assetName string) (*CIP60Metadata, error)

// IsCIP60Token checks if metadata contains music_metadata_version
func IsCIP60Token(metadata map[uint64]any, policyId, assetName string) bool

// GetMusicMetadataVersion returns the music metadata version if present
func GetMusicMetadataVersion(metadata map[uint64]any, policyId, assetName string) (MusicMetadataVersion, bool)
```

### Builder Pattern

```go
// NewCIP60Metadata creates a new music metadata builder
func NewCIP60Metadata() *CIP60MetadataBuilder

type CIP60MetadataBuilder struct {
    metadata *CIP60Metadata
}

// Base fields
func (b *CIP60MetadataBuilder) WithName(name string) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithImage(uri string) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithDescription(desc string) *CIP60MetadataBuilder

// Music version
func (b *CIP60MetadataBuilder) WithVersion(v MusicMetadataVersion) *CIP60MetadataBuilder

// Artists
func (b *CIP60MetadataBuilder) WithArtist(name string) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithArtistFull(artist Artist) *CIP60MetadataBuilder

// Release info
func (b *CIP60MetadataBuilder) WithRelease(release ReleaseInfo) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithReleaseType(releaseType string) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithReleaseDate(date string) *CIP60MetadataBuilder

// Music files
func (b *CIP60MetadataBuilder) WithMusicFile(file MusicFile) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithAudioFile(name, mediaType, src string) *CIP60MetadataBuilder

// Other fields
func (b *CIP60MetadataBuilder) WithGenres(genres ...string) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithCopyright(copyright string) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithLyrics(lyrics string) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithISRC(isrc string) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithExplicit(explicit bool) *CIP60MetadataBuilder
func (b *CIP60MetadataBuilder) WithLink(platform, uri string) *CIP60MetadataBuilder

func (b *CIP60MetadataBuilder) Build() (*CIP60Metadata, error)
```

### Validation

```go
// Validate checks music metadata validity
func (m *CIP60Metadata) Validate() error

// ValidateISRC validates ISRC format
func ValidateISRC(isrc string) error

// ValidateISWC validates ISWC format
func ValidateISWC(iswc string) error

// ValidateDuration validates duration format
func ValidateDuration(duration string) error
```

---

## Usage Examples

### Parse Music Token Metadata

```go
import "github.com/blinklabs-io/cardano-models/cip60"

// Check if it's a music token
if cip60.IsCIP60Token(tx.Metadata(), policyId, assetName) {
    meta, err := cip60.ExtractCIP60Metadata(tx.Metadata(), policyId, assetName)
    if err != nil {
        // Invalid music metadata
    }

    fmt.Println("Song:", meta.Name)
    fmt.Println("Version:", meta.MusicMetadataVersion)

    for _, artist := range meta.Artists {
        fmt.Println("Artist:", artist.Name)
    }

    for _, file := range meta.Files {
        fmt.Printf("Track %d: %s (%s)\n",
            *file.TrackNumber, file.SongTitle, file.SongDuration)
    }
}
```

### Build Music Metadata

```go
meta, err := cip60.NewCIP60Metadata().
    WithName("My Song").
    WithImage("ipfs://QmCover...").
    WithVersion(cip60.MusicMetadataV3).
    WithArtist("Artist Name").
    WithReleaseType("Single").
    WithReleaseDate("2026-01-17").
    WithAudioFile("My Song", "audio/flac", "ipfs://QmAudio...").
    WithGenres("Electronic", "Ambient").
    WithISRC("USRC12345678").
    WithCopyright("2026 Artist Name").
    WithExplicit(false).
    WithLink("spotify", "https://open.spotify.com/track/...").
    Build()
```

---

## Implementation Phases

### Phase 1: Core Types

- [ ] Define `MusicMetadataVersion` type
- [ ] Define `CIP60Metadata` struct
- [ ] Define `ReleaseInfo` struct
- [ ] Define `Artist` struct
- [ ] Define `MusicFile` struct

### Phase 2: Detection & Extraction

- [ ] Implement `IsCIP60Token()`
- [ ] Implement `GetMusicMetadataVersion()`
- [ ] Implement `ExtractCIP60Metadata()`
- [ ] Handle version differences (v1/v2/v3)

### Phase 3: Serialization

- [ ] Implement CBOR marshaling
- [ ] Implement JSON marshaling
- [ ] Handle flexible field types

### Phase 4: Validation

- [ ] Implement `Validate()`
- [ ] Implement ISRC validation
- [ ] Implement ISWC validation
- [ ] Implement duration format validation

### Phase 5: Builder

- [ ] Implement `CIP60MetadataBuilder`
- [ ] Add convenience methods
- [ ] Add validation on build

### Phase 6: Testing

- [ ] Unit tests with synthetic data
- [ ] Tests with real music NFTs from mainnet
- [ ] Cross-reference with NEWM, SickCity implementations

---

## Success Criteria

- [ ] Parse CIP-60 v1, v2, and v3 metadata
- [ ] Detect music tokens vs regular CIP-25 NFTs
- [ ] Validate ISRC/ISWC codes
- [ ] Generate valid CIP-60 metadata
- [ ] Comprehensive test coverage

---

## Dependencies

| Dependency | Purpose |
|------------|---------|
| gouroboros | Transaction metadata types |
| cbor/v2 | CBOR encoding |
| cip25 (existing) | Base NFT metadata |

---

## Related Plans

- [CIP25_PLAN.md](CIP25_PLAN.md) - Base NFT metadata (already implemented)
- [CIP68_PLAN.md](CIP68_PLAN.md) - Datum-based metadata alternative
