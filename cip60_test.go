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
	"encoding/json"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool {
	return &b
}

func singleRelease() MusicAssetMetadata {
	return MusicAssetMetadata{
		Name:                 "Let's Turn it On",
		Image:                UriField{Uris: []string{"ipfs://image"}},
		MusicMetadataVersion: 3,
		Release: ReleaseDetails{
			ReleaseType:  "Single",
			ReleaseTitle: "Mr. Bad Guy",
			Distributor:  "newm.io",
		},
		Files: []MusicFileDetails{
			{
				Name:      "Let's Turn it On",
				MediaType: "audio/wav",
				Src:       UriField{Uris: []string{"ipfs://song"}},
				Song: SongDetails{
					SongTitle:    DescField{Descriptions: []string{"Let's Turn it On"}},
					SongDuration: "PT3M21S",
					TrackNumber:  1,
					CommonMusicDetails: CommonMusicDetails{
						Artists: []ArtistDetails{
							{Name: "Sick City", Isni: "xxxxxxxxxxxxx"},
						},
						Genres:    []string{"Rock", "Classic Rock"},
						Copyright: &CopyrightDetails{Master: "1985 Sony", Composition: "1985 Marvin Gaye"},
					},
				},
			},
		},
	}
}

func albumRelease() MusicAssetMetadata {
	return MusicAssetMetadata{
		Name:                 "Zipline From Outerspace",
		Image:                UriField{Uris: []string{"ipfs://image"}},
		MusicMetadataVersion: 3,
		Release: ReleaseDetails{
			ReleaseType:  "Album/EP",
			ReleaseTitle: "Zipline From Outerspace",
			CommonMusicDetails: CommonMusicDetails{
				Artists:   []ArtistDetails{{Name: "Refraktal", Isni: "0000000517483974"}},
				Genres:    []string{"Electronic", "Experimental", "Psychedelic"},
				Copyright: &CopyrightDetails{Master: "2024 Refraktal", Composition: "2024 Refraktal"},
			},
		},
		Files: []MusicFileDetails{
			{
				Name:      "Void",
				MediaType: "audio/wav",
				Src:       UriField{Uris: []string{"ipfs://void"}},
				Song: SongDetails{
					SongTitle:    DescField{Descriptions: []string{"Void"}},
					SongDuration: "PT4M21S",
					TrackNumber:  1,
					Isrc:         "US-SKG-22-12345",
					Explicit:     boolPtr(true),
				},
			},
			{
				Name:      "Bellywub",
				MediaType: "audio/wav",
				Src:       UriField{Uris: []string{"ipfs://bellywub"}},
				Song: SongDetails{
					SongTitle:    DescField{Descriptions: []string{"Bellywub"}},
					SongDuration: "PT5M31S",
					TrackNumber:  2,
				},
			},
		},
	}
}

func TestNewCip60Metadata_Single(t *testing.T) {
	policies := map[string]map[string]MusicAssetMetadata{
		"policy1": {"asset1": singleRelease()},
	}
	meta, err := NewCip60Metadata(policies)
	require.NoError(t, err)
	require.Len(t, meta.Num721.Policies, 1)
	asset := meta.Num721.Policies[HexBytes("policy1")][HexBytes("asset1")]
	require.Equal(t, "Single", asset.Release.ReleaseType)
}

func TestNewCip60Metadata_Album(t *testing.T) {
	policies := map[string]map[string]MusicAssetMetadata{
		"policy1": {"asset1": albumRelease()},
	}
	meta, err := NewCip60Metadata(policies)
	require.NoError(t, err)
	asset := meta.Num721.Policies[HexBytes("policy1")][HexBytes("asset1")]
	require.Equal(t, "Album/EP", asset.Release.ReleaseType)
	require.Len(t, asset.Files, 2)
}

func TestNewCip60Metadata_ValidationFailures(t *testing.T) {
	t.Run("missing music_metadata_version", func(t *testing.T) {
		asset := singleRelease()
		asset.MusicMetadataVersion = 0
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "MusicMetadataVersion")
	})

	t.Run("invalid release_type", func(t *testing.T) {
		asset := singleRelease()
		asset.Release.ReleaseType = "Box Set"
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "ReleaseType")
	})

	t.Run("no files", func(t *testing.T) {
		asset := singleRelease()
		asset.Files = nil
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Files")
	})

	t.Run("single release missing song artists", func(t *testing.T) {
		asset := singleRelease()
		asset.Files[0].Song.Artists = nil
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "song.artists")
	})

	t.Run("single release missing song genres", func(t *testing.T) {
		asset := singleRelease()
		asset.Files[0].Song.Genres = nil
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "song.genres")
	})

	t.Run("single release missing song copyright", func(t *testing.T) {
		asset := singleRelease()
		asset.Files[0].Song.Copyright = nil
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "song.copyright")
	})

	t.Run("album release missing release artists", func(t *testing.T) {
		asset := albumRelease()
		asset.Release.Artists = nil
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "release.artists")
	})

	t.Run("album release missing release genres", func(t *testing.T) {
		asset := albumRelease()
		asset.Release.Genres = nil
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "release.genres")
	})

	t.Run("album release missing release copyright", func(t *testing.T) {
		asset := albumRelease()
		asset.Release.Copyright = nil
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "release.copyright")
	})

	t.Run("too many genres", func(t *testing.T) {
		asset := singleRelease()
		asset.Files[0].Song.Genres = []string{"Rock", "Pop", "Jazz", "Blues"}
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Genres")
	})

	t.Run("invalid song_duration format", func(t *testing.T) {
		asset := singleRelease()
		asset.Files[0].Song.SongDuration = "three minutes"
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "SongDuration")
	})

	t.Run("missing song_title", func(t *testing.T) {
		asset := singleRelease()
		asset.Files[0].Song.SongTitle = DescField{}
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "song.song_title")
	})

	t.Run("single release with common details on release", func(t *testing.T) {
		asset := singleRelease()
		asset.Release.Artists = []ArtistDetails{{Name: "Someone"}}
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "release must not set")
	})

	t.Run("album release with common details on song", func(t *testing.T) {
		asset := albumRelease()
		asset.Files[0].Song.Genres = []string{"Electronic"}
		_, err := NewCip60Metadata(
			map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "song must not set")
	})
}

func TestCip60Metadata_Description_OmittedWhenUnset(t *testing.T) {
	policies := map[string]map[string]MusicAssetMetadata{
		"policy1": {"asset1": singleRelease()},
	}
	original, err := NewCip60Metadata(policies)
	require.NoError(t, err)

	data, err := json.Marshal(original)
	require.NoError(t, err)
	require.NotContains(t, string(data), "description")

	asset := singleRelease()
	asset.Description = &DescField{Descriptions: []string{"A great song"}}
	withDesc, err := NewCip60Metadata(
		map[string]map[string]MusicAssetMetadata{"policy1": {"asset1": asset}},
	)
	require.NoError(t, err)

	data, err = json.Marshal(withDesc)
	require.NoError(t, err)
	require.Contains(t, string(data), `"description":"A great song"`)
}

func TestCip60Metadata_JSON_RoundTrip(t *testing.T) {
	policies := map[string]map[string]MusicAssetMetadata{
		"policy1": {"asset1": albumRelease()},
	}
	original, err := NewCip60Metadata(policies)
	require.NoError(t, err)

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var unmarshaled Cip60Metadata
	require.NoError(t, json.Unmarshal(data, &unmarshaled))

	data2, err := json.Marshal(unmarshaled)
	require.NoError(t, err)

	var unmarshaled2 Cip60Metadata
	require.NoError(t, json.Unmarshal(data2, &unmarshaled2))

	require.Equal(t, unmarshaled, unmarshaled2)

	asset := unmarshaled.Num721.Policies[HexBytes("policy1")][HexBytes("asset1")]
	require.Equal(t, "Album/EP", asset.Release.ReleaseType)
	require.Equal(t, 3, asset.MusicMetadataVersion)
	require.Len(t, asset.Files, 2)
	require.True(t, *asset.Files[0].Song.Explicit)
	require.Equal(t, "US-SKG-22-12345", asset.Files[0].Song.Isrc)
}

func TestCip60Metadata_CBOR_RoundTrip(t *testing.T) {
	policies := map[string]map[string]MusicAssetMetadata{
		"policy1": {"asset1": singleRelease()},
	}
	original, err := NewCip60Metadata(policies)
	require.NoError(t, err)

	data, err := cbor.Marshal(original)
	require.NoError(t, err)

	var unmarshaled Cip60Metadata
	require.NoError(t, cbor.Unmarshal(data, &unmarshaled))

	data2, err := cbor.Marshal(unmarshaled)
	require.NoError(t, err)

	var unmarshaled2 Cip60Metadata
	require.NoError(t, cbor.Unmarshal(data2, &unmarshaled2))

	require.Equal(t, unmarshaled, unmarshaled2)

	asset := unmarshaled.Num721.Policies[HexBytes("policy1")][HexBytes("asset1")]
	require.Equal(t, "Single", asset.Release.ReleaseType)
	require.Equal(t, "Mr. Bad Guy", asset.Release.ReleaseTitle)
	require.Equal(t, []string{"Rock", "Classic Rock"}, asset.Files[0].Song.Genres)
}

// TestCip60Metadata_CBOR_TextKeysByDefault asserts that, under the default
// (absent) CIP-25 "721" version, policy_id/asset_name map keys are encoded
// as CBOR text strings on the wire even when they look like hex: a real
// policy ID (56 hex chars) and a hex-looking UTF-8 asset name ("AB", valid
// hex digits) must not be silently converted to byte strings, which would
// change the asset identity/type seen by indexers.
func TestCip60Metadata_CBOR_TextKeysByDefault(t *testing.T) {
	const policyID = "c00d776a22ca5db986039420b2a9b3f880d593136a9e2262fabeeb58"
	const assetName = "AB"

	policies := map[string]map[string]MusicAssetMetadata{
		policyID: {assetName: singleRelease()},
	}
	meta, err := NewCip60Metadata(policies)
	require.NoError(t, err)

	data, err := cbor.Marshal(meta)
	require.NoError(t, err)

	var raw map[any]any
	require.NoError(t, cbor.Unmarshal(data, &raw))
	inner, ok := raw[uint64(721)].(map[any]any)
	require.True(t, ok, "expected a map under the 721 label")

	for policyKey, policyValue := range inner {
		_, isString := policyKey.(string)
		require.Truef(
			t,
			isString,
			"policy key %v (%T) should be encoded as CBOR text, not bytes",
			policyKey,
			policyKey,
		)
		require.Equal(t, policyID, policyKey)

		assetMap, ok := policyValue.(map[any]any)
		require.True(t, ok, "expected a map of asset names")
		for assetKey := range assetMap {
			_, isString := assetKey.(string)
			require.Truef(
				t,
				isString,
				"asset key %v (%T) should be encoded as CBOR text, not bytes",
				assetKey,
				assetKey,
			)
			require.Equal(t, assetName, assetKey)
		}
	}

	var unmarshaled Cip60Metadata
	require.NoError(t, cbor.Unmarshal(data, &unmarshaled))
	_, ok = unmarshaled.Num721.Policies[HexBytes(policyID)][HexBytes(assetName)]
	require.True(t, ok)
}
