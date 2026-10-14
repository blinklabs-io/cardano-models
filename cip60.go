// Copyright 2026 Blink Labs Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/fxamacker/cbor/v2"
	"github.com/go-playground/validator/v10"
)

// cip60SongDurationPattern matches the ISO-8601 duration subset required by
// CIP-60 for song_duration: "^P(?!$)(T(?=\d)(\d+H)?(\d+M)?(\d+S)?)$". Go's
// RE2 engine doesn't support the lookaheads in that pattern, so this is
// expressed as "PT" followed by at least one of hours/minutes/seconds,
// enforced by also rejecting the literal "PT" (no components) below.
var cip60SongDurationPattern = regexp.MustCompile(`^PT(\d+H)?(\d+M)?(\d+S)?$`)

// validateCip60SongDuration implements the "cip60duration" validator tag.
func validateCip60SongDuration(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return s != "PT" && cip60SongDurationPattern.MatchString(s)
}

func init() {
	if err := validate.RegisterValidation("cip60duration", validateCip60SongDuration); err != nil {
		panic(err)
	}
}

// Cip60Metadata is the top-level container for CIP-60 music token metadata
// under the "721" tag. CIP-60 extends CIP-25 (and optionally CIP-68) to add
// a common vocabulary for music-specific fields.
type Cip60Metadata struct {
	Num721 MusicNum721 `cbor:"721,keyasint" json:"721" validate:"required"`
}

// MusicNum721 represents the "721" label_metadata structure for CIP-60,
// mapping policy IDs to asset names to music asset metadata.
//
//nolint:recvcheck
type MusicNum721 struct {
	Policies map[HexBytes]map[HexBytes]MusicAssetMetadata `json:"-" cbor:"" validate:"required,min=1,dive,dive"`
}

// UnmarshalJSON handles unmarshaling the MusicNum721 structure.
func (n *MusicNum721) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	n.Policies = make(map[HexBytes]map[HexBytes]MusicAssetMetadata)

	for key, value := range raw {
		// key is policy_id, value is map[asset_name]MusicAssetMetadata
		policyMap, ok := value.(map[string]any)
		if !ok {
			// Skip non-object entries instead of failing
			continue
		}
		assets := make(map[HexBytes]MusicAssetMetadata)
		for assetKey, assetValue := range policyMap {
			var asset MusicAssetMetadata
			assetBytes, err := json.Marshal(assetValue)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(assetBytes, &asset); err != nil {
				return err
			}
			assets[HexBytes(assetKey)] = asset
		}
		n.Policies[HexBytes(key)] = assets
	}

	return nil
}

// MarshalJSON handles marshaling the MusicNum721 structure.
func (n MusicNum721) MarshalJSON() ([]byte, error) {
	out := make(map[string]any)
	for policy, assets := range n.Policies {
		out[string(policy)] = assets
	}
	return json.Marshal(out)
}

// decodeCborMapKey converts a CBOR map key (string, []byte, or
// cbor.ByteString) into HexBytes.
func decodeCborMapKey(key any) (HexBytes, error) {
	switch v := key.(type) {
	case string:
		return HexBytes(v), nil
	case []byte:
		return HexBytes(hex.EncodeToString(v)), nil
	case cbor.ByteString:
		return HexBytes(hex.EncodeToString(v.Bytes())), nil
	default:
		return "", fmt.Errorf("unsupported map key type: %T", key)
	}
}

// UnmarshalCBOR handles unmarshaling the MusicNum721 structure from CBOR.
func (n *MusicNum721) UnmarshalCBOR(data []byte) error {
	var raw map[any]any
	if err := cbor.Unmarshal(data, &raw); err != nil {
		return err
	}

	n.Policies = make(map[HexBytes]map[HexBytes]MusicAssetMetadata)

	for key, value := range raw {
		k, err := decodeCborMapKey(key)
		if err != nil {
			return err
		}

		policyMap, ok := value.(map[any]any)
		if !ok {
			return errors.New("invalid policy structure")
		}
		assets := make(map[HexBytes]MusicAssetMetadata)
		for assetKey, assetValue := range policyMap {
			ak, err := decodeCborMapKey(assetKey)
			if err != nil {
				return err
			}

			var asset MusicAssetMetadata
			assetBytes, err := cbor.Marshal(assetValue)
			if err != nil {
				return err
			}
			if err := cbor.Unmarshal(assetBytes, &asset); err != nil {
				return err
			}
			assets[ak] = asset
		}
		n.Policies[k] = assets
	}

	return nil
}

// MarshalCBOR handles marshaling the MusicNum721 structure to CBOR.
func (n MusicNum721) MarshalCBOR() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := cbor.NewEncoder(buf)
	if err := enc.StartIndefiniteMap(); err != nil {
		return nil, err
	}

	keys := make([]HexBytes, 0, len(n.Policies))
	for policy := range n.Policies {
		keys = append(keys, policy)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})

	for _, k := range keys {
		if err := enc.Encode(k); err != nil {
			return nil, err
		}

		inner := n.Policies[k]
		innerKeys := make([]HexBytes, 0, len(inner))
		for asset := range inner {
			innerKeys = append(innerKeys, asset)
		}
		sort.Slice(innerKeys, func(i, j int) bool {
			return string(innerKeys[i]) < string(innerKeys[j])
		})

		if err := enc.StartIndefiniteMap(); err != nil {
			return nil, err
		}
		for _, ik := range innerKeys {
			if err := enc.Encode(ik); err != nil {
				return nil, err
			}
			if err := enc.Encode(inner[ik]); err != nil {
				return nil, err
			}
		}
		if err := enc.EndIndefinite(); err != nil {
			return nil, err
		}
	}

	if err := enc.EndIndefinite(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MusicAssetMetadata represents the CIP-60 "metadata_details" for a single
// music asset.
type MusicAssetMetadata struct {
	Name                 string             `json:"name"                   cbor:"name"                   validate:"required"`
	Image                UriField           `json:"image"                  cbor:"image"                  validate:"required"`
	MusicMetadataVersion int                `json:"music_metadata_version" cbor:"music_metadata_version" validate:"required,eq=3"`
	Release              ReleaseDetails     `json:"release"                cbor:"release"                validate:"required"`
	Files                []MusicFileDetails `json:"files"                  cbor:"files"                  validate:"required,min=1,dive"`
	Version              int                `json:"version,omitempty"      cbor:"version,omitempty"      validate:"omitempty,oneof=1 2"`
	MediaType            string             `json:"mediaType,omitempty"    cbor:"mediaType,omitempty"`
	Description          *DescField         `json:"description,omitempty"  cbor:"description,omitempty"`
}

// ReleaseDetails represents release-level metadata ("release_details").
// Artists, Genres, Copyright, and ContributingArtists (CommonMusicDetails)
// are expected here when ReleaseType is "Album/EP", and in each song
// otherwise.
type ReleaseDetails struct {
	ReleaseType     string `json:"release_type"              cbor:"release_type"              validate:"required,oneof=Single Multiple Album/EP"`
	ReleaseTitle    string `json:"release_title"              cbor:"release_title"              validate:"required"`
	Distributor     string `json:"distributor,omitempty"      cbor:"distributor,omitempty"`
	VisualArtist    string `json:"visual_artist,omitempty"    cbor:"visual_artist,omitempty"`
	ReleaseDate     string `json:"release_date,omitempty"     cbor:"release_date,omitempty"`
	PublicationDate string `json:"publication_date,omitempty" cbor:"publication_date,omitempty"`
	CatalogNumber   string `json:"catalog_number,omitempty"   cbor:"catalog_number,omitempty"`
	Series          string `json:"series,omitempty"           cbor:"series,omitempty"`
	Collection      string `json:"collection,omitempty"       cbor:"collection,omitempty"`

	CommonMusicDetails
}

// SongDetails represents the metadata attached to an individual song
// ("song_details"). Artists, Genres, Copyright, and ContributingArtists
// (CommonMusicDetails) are expected here unless the release's ReleaseType is
// "Album/EP", in which case they belong on the release instead.
type SongDetails struct {
	SongTitle    DescField `json:"song_title"    cbor:"song_title"    validate:"required"`
	SongDuration string    `json:"song_duration" cbor:"song_duration" validate:"required,cip60duration"`
	TrackNumber  uint      `json:"track_number"  cbor:"track_number"  validate:"required"`

	FeaturedArtists   []ArtistDetails `json:"featured_artists,omitempty"   cbor:"featured_artists,omitempty"   validate:"omitempty,dive"`
	Authors           []AuthorDetails `json:"authors,omitempty"            cbor:"authors,omitempty"            validate:"omitempty,dive"`
	Mood              string          `json:"mood,omitempty"               cbor:"mood,omitempty"`
	Set               string          `json:"set,omitempty"                cbor:"set,omitempty"`
	Lyrics            string          `json:"lyrics,omitempty"             cbor:"lyrics,omitempty"`
	SpecialThanks     []string        `json:"special_thanks,omitempty"     cbor:"special_thanks,omitempty"`
	Bitrate           string          `json:"bitrate,omitempty"            cbor:"bitrate,omitempty"`
	Bpm               string          `json:"bpm,omitempty"                cbor:"bpm,omitempty"`
	MixEngineer       string          `json:"mix_engineer,omitempty"       cbor:"mix_engineer,omitempty"`
	MasteringEngineer string          `json:"mastering_engineer,omitempty" cbor:"mastering_engineer,omitempty"`
	Producer          string          `json:"producer,omitempty"           cbor:"producer,omitempty"`
	CoProducer        string          `json:"co_producer,omitempty"        cbor:"co_producer,omitempty"`
	RecordingEngineer string          `json:"recording_engineer,omitempty" cbor:"recording_engineer,omitempty"`
	Explicit          *bool           `json:"explicit,omitempty"           cbor:"explicit,omitempty"`
	Isrc              string          `json:"isrc,omitempty"               cbor:"isrc,omitempty"`
	Iswc              string          `json:"iswc,omitempty"               cbor:"iswc,omitempty"`
	MetadataLanguage  string          `json:"metadata_language,omitempty"  cbor:"metadata_language,omitempty"`
	CountryOfOrigin   string          `json:"country_of_origin,omitempty"  cbor:"country_of_origin,omitempty"`
	Language          string          `json:"language,omitempty"           cbor:"language,omitempty"`
	DerivedFrom       string          `json:"derived_from,omitempty"       cbor:"derived_from,omitempty"`
	AiGenerated       *bool           `json:"ai_generated,omitempty"       cbor:"ai_generated,omitempty"`

	CommonMusicDetails
}

// CommonMusicDetails holds the fields shared between release-level and
// song-level metadata ("common_music_details" in the CIP-60 CDDL). They
// belong on the release when ReleaseType is "Album/EP", and on each song
// otherwise.
type CommonMusicDetails struct {
	Artists             []ArtistDetails             `json:"artists,omitempty"              cbor:"artists,omitempty"              validate:"omitempty,dive"`
	Genres              []string                    `json:"genres,omitempty"               cbor:"genres,omitempty"               validate:"omitempty,min=1,max=3"`
	Copyright           *CopyrightDetails           `json:"copyright,omitempty"            cbor:"copyright,omitempty"`
	ContributingArtists []ContributingArtistDetails `json:"contributing_artists,omitempty" cbor:"contributing_artists,omitempty" validate:"omitempty,dive"`
}

// hasCommonMusicDetails reports whether any of the shared release/song
// fields (artists, genres, copyright, contributing_artists) are set.
func (c CommonMusicDetails) hasCommonMusicDetails() bool {
	return len(c.Artists) > 0 || len(c.Genres) > 0 || c.Copyright != nil ||
		len(c.ContributingArtists) > 0
}

// MusicFileDetails represents a file entry extended with song metadata
// ("files_details").
type MusicFileDetails struct {
	Name      string      `json:"name"      cbor:"name"      validate:"required"`
	MediaType string      `json:"mediaType" cbor:"mediaType" validate:"required"`
	Src       UriField    `json:"src"       cbor:"src"       validate:"required"`
	Song      SongDetails `json:"song"      cbor:"song"      validate:"required"`
}

// ArtistDetails represents an artist entry ("artist_details").
type ArtistDetails struct {
	Name  string            `json:"name"            cbor:"name"            validate:"required"`
	Isni  string            `json:"isni,omitempty"  cbor:"isni,omitempty"`
	Links map[string]string `json:"links,omitempty" cbor:"links,omitempty"`
}

// AuthorDetails represents an author/publisher entry ("author_details").
type AuthorDetails struct {
	Name  string `json:"name"            cbor:"name"            validate:"required"`
	Ipi   string `json:"ipi,omitempty"   cbor:"ipi,omitempty"`
	Share string `json:"share,omitempty" cbor:"share,omitempty"`
}

// ContributingArtistDetails represents a contributing artist entry
// ("contributing_artist_details").
type ContributingArtistDetails struct {
	Name  string            `json:"name"            cbor:"name"            validate:"required"`
	Ipn   string            `json:"ipn,omitempty"   cbor:"ipn,omitempty"`
	Ipi   string            `json:"ipi,omitempty"   cbor:"ipi,omitempty"`
	Role  []string          `json:"role,omitempty"  cbor:"role,omitempty"`
	Links map[string]string `json:"links,omitempty" cbor:"links,omitempty"`
}

// CopyrightDetails represents the master recording and composition
// copyright holders.
type CopyrightDetails struct {
	Master      string `json:"master"      cbor:"master"      validate:"required"`
	Composition string `json:"composition" cbor:"composition" validate:"required"`
}

// NewCip60Metadata creates a new CIP-60 metadata object from the given
// policy ID -> asset name -> music asset metadata map.
func NewCip60Metadata(
	policies map[string]map[string]MusicAssetMetadata,
) (*Cip60Metadata, error) {
	convertedPolicies := make(map[HexBytes]map[HexBytes]MusicAssetMetadata)
	for k, v := range policies {
		convertedInner := make(map[HexBytes]MusicAssetMetadata)
		for ik, iv := range v {
			convertedInner[HexBytes(ik)] = iv
		}
		convertedPolicies[HexBytes(k)] = convertedInner
	}

	metadata := &Cip60Metadata{
		Num721: MusicNum721{Policies: convertedPolicies},
	}

	if err := metadata.Validate(); err != nil {
		return nil, err
	}

	return metadata, nil
}

// Validate checks the CIP-60 metadata for validity, including the struct
// tags and the conditional placement of artists, genres, and copyright,
// which depends on each release's ReleaseType (see CIP-60 "Required
// Fields").
func (c *Cip60Metadata) Validate() error {
	if err := validate.Struct(c); err != nil {
		return err
	}
	for policyID, assets := range c.Num721.Policies {
		for assetName, asset := range assets {
			if err := asset.validatePlacement(); err != nil {
				return fmt.Errorf(
					"policy %s asset %s: %w",
					policyID,
					assetName,
					err,
				)
			}
		}
	}
	return nil
}

// validatePlacement enforces that artists, genres, and copyright appear on
// the release when ReleaseType is "Album/EP", and on each song otherwise,
// rejecting them from the disallowed level. It also requires a non-empty
// song_title on every song, which the "required" struct tag cannot enforce
// reliably on a DescField value.
func (a MusicAssetMetadata) validatePlacement() error {
	isAlbum := a.Release.ReleaseType == "Album/EP"

	if isAlbum {
		if len(a.Release.Artists) == 0 {
			return errors.New(
				"release.artists is required when release_type is Album/EP",
			)
		}
		if len(a.Release.Genres) == 0 {
			return errors.New(
				"release.genres is required when release_type is Album/EP",
			)
		}
		if a.Release.Copyright == nil {
			return errors.New(
				"release.copyright is required when release_type is Album/EP",
			)
		}
	} else if a.Release.hasCommonMusicDetails() {
		return fmt.Errorf(
			"release must not set artists, genres, copyright, or contributing_artists when release_type is %s (use song instead)",
			a.Release.ReleaseType,
		)
	}

	for i, file := range a.Files {
		if len(file.Song.SongTitle.Descriptions) == 0 {
			return fmt.Errorf("files[%d].song.song_title is required", i)
		}

		if isAlbum {
			if file.Song.hasCommonMusicDetails() {
				return fmt.Errorf(
					"files[%d].song must not set artists, genres, copyright, or contributing_artists when release_type is Album/EP (use release instead)",
					i,
				)
			}
			continue
		}

		if len(file.Song.Artists) == 0 {
			return fmt.Errorf(
				"files[%d].song.artists is required when release_type is %s",
				i,
				a.Release.ReleaseType,
			)
		}
		if len(file.Song.Genres) == 0 {
			return fmt.Errorf(
				"files[%d].song.genres is required when release_type is %s",
				i,
				a.Release.ReleaseType,
			)
		}
		if file.Song.Copyright == nil {
			return fmt.Errorf(
				"files[%d].song.copyright is required when release_type is %s",
				i,
				a.Release.ReleaseType,
			)
		}
	}
	return nil
}
