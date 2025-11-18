// Copyright 2025 Blink Labs Software
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

func TestNewCip25Metadata(t *testing.T) {
	// Test valid metadata
	policies := map[string]map[string]AssetMetadata{
		"policy1": {
			"asset1": {
				Name: "Test Asset",
				Image: UriField{
					Uris: []string{"https://example.com/image.png"},
				},
			},
		},
	}
	meta, err := NewCip25Metadata(1, policies)
	require.NoError(t, err)
	require.Equal(t, 1, meta.Num721.Version)
	require.Len(t, meta.Num721.Policies, 1)
}

func TestNewCip25Metadata_ValidationFailures(t *testing.T) {
	t.Run("missing name", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Image: UriField{
						Uris: []string{"https://example.com/image.png"},
					},
				},
			},
		}
		_, err := NewCip25Metadata(1, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Name")
	})

	t.Run("empty image", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Name:  "Test Asset",
					Image: UriField{Uris: []string{}},
				},
			},
		}
		_, err := NewCip25Metadata(1, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Image")
	})

	t.Run("invalid version", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Name: "Test Asset",
					Image: UriField{
						Uris: []string{"https://example.com/image.png"},
					},
				},
			},
		}
		_, err := NewCip25Metadata(3, policies) // Invalid version
		require.Error(t, err)
		require.Contains(t, err.Error(), "Version")
	})

	t.Run("empty policies", func(t *testing.T) {
		_, err := NewCip25Metadata(1, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Policies")
	})
}

func TestCip25Metadata_JSON_RoundTrip(t *testing.T) {
	// Create a CIP-25 metadata object
	policies := map[string]map[string]AssetMetadata{
		"policy1": {
			"asset1": {
				Name: "Test Asset",
				Image: UriField{
					Uris: []string{"https://example.com/image.png"},
				},
				MediaType:   "image/png",
				Description: DescField{Descriptions: []string{"A test asset"}},
				Files: []FileDetails{
					{
						Name:      "file1",
						MediaType: "image/png",
						Src: UriField{
							Uris: []string{"https://example.com/file1.png"},
						},
					},
				},
			},
		},
	}
	original, err := NewCip25Metadata(1, policies)
	require.NoError(t, err)

	// Marshal to JSON
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Cip25Metadata
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Marshal again
	data2, err := json.Marshal(unmarshaled)
	require.NoError(t, err)

	// Unmarshal data2
	var unmarshaled2 Cip25Metadata
	err = json.Unmarshal(data2, &unmarshaled2)
	require.NoError(t, err)

	// Check the structs are equal
	require.Equal(t, unmarshaled, unmarshaled2)
}

func TestCip25Metadata_JSON_RoundTrip_MultiplePolicies(t *testing.T) {
	// Create a CIP-25 metadata object with multiple policies and assets
	hexPolicies := map[HexBytes]map[HexBytes]AssetMetadata{
		HexBytes("policy1"): {
			HexBytes("asset1"): {
				Name: "Asset One",
				Image: UriField{
					Uris: []string{"https://example.com/asset1.png"},
				},
				MediaType:   "image/png",
				Description: DescField{Descriptions: []string{"First asset"}},
			},
			HexBytes("asset2"): {
				Name: "Asset Two",
				Image: UriField{
					Uris: []string{"https://example.com/asset2.png"},
				},
			},
		},
		HexBytes("policy2"): {
			HexBytes("asset3"): {
				Name: "Asset Three",
				Image: UriField{
					Uris: []string{"https://example.com/asset3.png"},
				},
				Description: DescField{
					Descriptions: []string{
						"Third asset",
						"With multiple descriptions",
					},
				},
				Files: []FileDetails{
					{
						Name:      "thumbnail",
						MediaType: "image/png",
						Src: UriField{
							Uris: []string{"https://example.com/thumb.png"},
						},
					},
				},
			},
		},
	}
	// Convert to string-keyed map for NewCip25Metadata
	policies := make(map[string]map[string]AssetMetadata)
	for policy, assets := range hexPolicies {
		policies[string(policy)] = make(map[string]AssetMetadata)
		for asset, meta := range assets {
			policies[string(policy)][string(asset)] = meta
		}
	}
	original, err := NewCip25Metadata(1, policies)
	require.NoError(t, err)

	// Marshal to JSON
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Cip25Metadata
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Check structural equality
	require.Equal(t, original.Num721.Version, unmarshaled.Num721.Version)
	require.Len(t, unmarshaled.Num721.Policies, 2)
	require.Len(t, unmarshaled.Num721.Policies[HexBytes("policy1")], 2)
	require.Len(t, unmarshaled.Num721.Policies[HexBytes("policy2")], 1)
	require.Equal(
		t,
		"Asset One",
		unmarshaled.Num721.Policies[HexBytes("policy1")][HexBytes("asset1")].Name,
	)
	require.Equal(
		t,
		[]string{"https://example.com/asset1.png"},
		unmarshaled.Num721.Policies[HexBytes("policy1")][HexBytes("asset1")].Image.Uris,
	)
	require.Equal(
		t,
		[]string{"Third asset", "With multiple descriptions"},
		unmarshaled.Num721.Policies[HexBytes("policy2")][HexBytes("asset3")].Description.Descriptions,
	)
}

func TestCip25Metadata_CBOR_RoundTrip(t *testing.T) {
	// Create a CIP-25 metadata object
	policies := map[string]map[string]AssetMetadata{
		"policy1": {
			"asset1": {
				Name: "Test Asset",
				Image: UriField{
					Uris: []string{"https://example.com/image.png"},
				},
				MediaType:   "image/png",
				Description: DescField{Descriptions: []string{"A test asset"}},
				Files: []FileDetails{
					{
						Name:      "file1",
						MediaType: "image/png",
						Src: UriField{
							Uris: []string{"https://example.com/file1.png"},
						},
					},
				},
			},
		},
	}
	original, err := NewCip25Metadata(1, policies)
	require.NoError(t, err)

	// Marshal to CBOR
	data, err := cbor.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Cip25Metadata
	err = cbor.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Marshal again
	data2, err := cbor.Marshal(unmarshaled)
	require.NoError(t, err)

	// Unmarshal data2
	var unmarshaled2 Cip25Metadata
	err = cbor.Unmarshal(data2, &unmarshaled2)
	require.NoError(t, err)

	// Check the unmarshaled structs are equal
	require.Equal(t, unmarshaled, unmarshaled2)
}

func TestCip25Metadata_Version2_JSON_RoundTrip(t *testing.T) {
	// Create a CIP-25 metadata object for version 2
	policies := map[string]map[string]AssetMetadata{
		"706f6c69637931": { // hex of "policy1"
			"617373657431": { // hex of "asset1"
				Name: "Test Asset",
				Image: UriField{
					Uris: []string{"https://example.com/image.png"},
				},
			},
		},
	}
	original, err := NewCip25Metadata(2, policies)
	require.NoError(t, err)

	// Marshal to JSON
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Cip25Metadata
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Marshal again
	data2, err := json.Marshal(unmarshaled)
	require.NoError(t, err)

	// Unmarshal data2
	var unmarshaled2 Cip25Metadata
	err = json.Unmarshal(data2, &unmarshaled2)
	require.NoError(t, err)

	// Check the structs are equal
	require.Equal(t, unmarshaled, unmarshaled2)
}

// TODO: Add a fixture-based test that decodes a CBOR blob with raw byte keys (canonical v2 on-chain form)
// to ensure interoperability with external CIP-25 v2 implementations.
func TestCip25Metadata_Version2_CBOR_RoundTrip(t *testing.T) {
	// Create a CIP-25 metadata object for version 2
	policies := map[string]map[string]AssetMetadata{
		"706f6c69637931": { // hex of "policy1"
			"617373657431": { // hex of "asset1"
				Name: "Test Asset",
				Image: UriField{
					Uris: []string{"https://example.com/image.png"},
				},
			},
		},
	}
	original, err := NewCip25Metadata(2, policies)
	require.NoError(t, err)

	// Marshal to CBOR
	data, err := cbor.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled Cip25Metadata
	err = cbor.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Marshal again
	data2, err := cbor.Marshal(unmarshaled)
	require.NoError(t, err)

	// Unmarshal data2
	var unmarshaled2 Cip25Metadata
	err = cbor.Unmarshal(data2, &unmarshaled2)
	require.NoError(t, err)

	// Check the structs are equal
	require.Equal(t, unmarshaled, unmarshaled2)
}

func TestUriField_JSON(t *testing.T) {
	// Single URI
	field := UriField{Uris: []string{"https://example.com"}}
	data, err := json.Marshal(field)
	require.NoError(t, err)
	require.Equal(t, `"https://example.com"`, string(data))

	var unmarshaled UriField
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	require.Equal(t, field, unmarshaled)

	// Multiple URIs
	field = UriField{
		Uris: []string{"https://example.com", "https://example2.com"},
	}
	data, err = json.Marshal(field)
	require.NoError(t, err)
	require.Equal(
		t,
		`["https://example.com","https://example2.com"]`,
		string(data),
	)

	unmarshaled = UriField{}
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	require.Equal(t, field, unmarshaled)
}

func TestUriField_JSON_Invalid(t *testing.T) {
	// Invalid type (object)
	var field UriField
	err := json.Unmarshal([]byte(`{"invalid": "object"}`), &field)
	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"URI field must be a string or an array of strings",
	)
}

func TestDescField_JSON(t *testing.T) {
	// Single description
	field := DescField{Descriptions: []string{"A description"}}
	data, err := json.Marshal(field)
	require.NoError(t, err)
	require.Equal(t, `"A description"`, string(data))

	var unmarshaled DescField
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	require.Equal(t, field, unmarshaled)

	// Multiple descriptions
	field = DescField{Descriptions: []string{"Desc1", "Desc2"}}
	data, err = json.Marshal(field)
	require.NoError(t, err)
	require.Equal(t, `["Desc1","Desc2"]`, string(data))

	unmarshaled = DescField{}
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	require.Equal(t, field, unmarshaled)
}

func TestDescField_JSON_Invalid(t *testing.T) {
	// Invalid type (number)
	var field DescField
	err := json.Unmarshal([]byte(`123`), &field)
	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"description must be a string or an array of strings",
	)
}

func TestUriField_CBOR(t *testing.T) {
	// Test single URI
	uf := UriField{Uris: []string{"https://example.com"}}
	data, err := cbor.Marshal(uf)
	require.NoError(t, err)

	var uf2 UriField
	err = cbor.Unmarshal(data, &uf2)
	require.NoError(t, err)
	require.Equal(t, uf, uf2)

	// Test multiple URIs
	uf3 := UriField{
		Uris: []string{"https://example.com/1", "https://example.com/2"},
	}
	data3, err := cbor.Marshal(uf3)
	require.NoError(t, err)

	var uf4 UriField
	err = cbor.Unmarshal(data3, &uf4)
	require.NoError(t, err)
	require.Equal(t, uf3, uf4)
}

func TestDescField_CBOR(t *testing.T) {
	// Test single description
	df := DescField{Descriptions: []string{"A test"}}
	data, err := cbor.Marshal(df)
	require.NoError(t, err)

	var df2 DescField
	err = cbor.Unmarshal(data, &df2)
	require.NoError(t, err)
	require.Equal(t, df, df2)

	// Test multiple descriptions
	df3 := DescField{Descriptions: []string{"Desc 1", "Desc 2"}}
	data3, err := cbor.Marshal(df3)
	require.NoError(t, err)

	var df4 DescField
	err = cbor.Unmarshal(data3, &df4)
	require.NoError(t, err)
	require.Equal(t, df3, df4)
}

func TestCip25Metadata_VersionDefault(t *testing.T) {
	// Test that version defaults to 1 when not specified
	jsonData := `{
		"721": {
			"policy1": {
				"asset1": {
					"name": "Test Asset",
					"image": "https://example.com/image.png"
				}
			}
		}
	}`

	var meta Cip25Metadata
	err := json.Unmarshal([]byte(jsonData), &meta)
	require.NoError(t, err)
	require.Equal(t, 1, meta.Num721.Version)
}

func TestCip25Metadata_RequiredFields(t *testing.T) {
	t.Run("missing name", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Image: UriField{Uris: []string{"https://example.com/image.png"}},
				},
			},
		}
		_, err := NewCip25Metadata(1, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Name")
	})

	t.Run("missing image", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Name: "Test Asset",
				},
			},
		}
		_, err := NewCip25Metadata(1, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Image")
	})

	t.Run("empty image URIs", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Name:  "Test Asset",
					Image: UriField{Uris: []string{}},
				},
			},
		}
		_, err := NewCip25Metadata(1, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Image")
	})
}

func TestCip25Metadata_FilesValidation(t *testing.T) {
	t.Run("missing mediaType in file", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Name: "Test Asset",
					Image: UriField{Uris: []string{"https://example.com/image.png"}},
					Files: []FileDetails{
						{
							Name: "file1",
							Src:  UriField{Uris: []string{"https://example.com/file1.png"}},
							// Missing MediaType
						},
					},
				},
			},
		}
		_, err := NewCip25Metadata(1, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "MediaType")
	})

	t.Run("missing name in file", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Name: "Test Asset",
					Image: UriField{Uris: []string{"https://example.com/image.png"}},
					Files: []FileDetails{
						{
							MediaType: "image/png",
							Src:       UriField{Uris: []string{"https://example.com/file1.png"}},
							// Missing Name
						},
					},
				},
			},
		}
		_, err := NewCip25Metadata(1, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Name")
	})

	t.Run("missing src in file", func(t *testing.T) {
		policies := map[string]map[string]AssetMetadata{
			"policy1": {
				"asset1": {
					Name: "Test Asset",
					Image: UriField{Uris: []string{"https://example.com/image.png"}},
					Files: []FileDetails{
						{
							Name:      "file1",
							MediaType: "image/png",
							// Missing Src
						},
					},
				},
			},
		}
		_, err := NewCip25Metadata(1, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Src")
	})
}

func TestCip25Metadata_VersionValidation(t *testing.T) {
	policies := map[string]map[string]AssetMetadata{
		"policy1": {
			"asset1": {
				Name: "Test Asset",
				Image: UriField{Uris: []string{"https://example.com/image.png"}},
			},
		},
	}

	t.Run("valid version 1", func(t *testing.T) {
		_, err := NewCip25Metadata(1, policies)
		require.NoError(t, err)
	})

	t.Run("valid version 2", func(t *testing.T) {
		_, err := NewCip25Metadata(2, policies)
		require.NoError(t, err)
	})

	t.Run("invalid version 0", func(t *testing.T) {
		_, err := NewCip25Metadata(0, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Version")
	})

	t.Run("invalid version 3", func(t *testing.T) {
		_, err := NewCip25Metadata(3, policies)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Version")
	})
}

func TestCip25Metadata_MinimalValid(t *testing.T) {
	// Test the minimal valid metadata as per spec
	policies := map[string]map[string]AssetMetadata{
		"policy1": {
			"asset1": {
				Name:  "Test Asset",
				Image: UriField{Uris: []string{"https://example.com/image.png"}},
			},
		},
	}

	meta, err := NewCip25Metadata(1, policies)
	require.NoError(t, err)
	require.Equal(t, 1, meta.Num721.Version)
	require.Len(t, meta.Num721.Policies, 1)
	require.Len(t, meta.Num721.Policies[HexBytes("policy1")], 1)
	require.Equal(t, "Test Asset", meta.Num721.Policies[HexBytes("policy1")][HexBytes("asset1")].Name)
	require.Equal(t, []string{"https://example.com/image.png"}, meta.Num721.Policies[HexBytes("policy1")][HexBytes("asset1")].Image.Uris)
}

func TestCip25Metadata_Version1Vs2Keys(t *testing.T) {
	policies := map[string]map[string]AssetMetadata{
		"policy1": {
			"asset1": {
				Name:  "Test Asset",
				Image: UriField{Uris: []string{"https://example.com/image.png"}},
			},
		},
	}

	t.Run("version 1 JSON", func(t *testing.T) {
		meta, err := NewCip25Metadata(1, policies)
		require.NoError(t, err)

		data, err := json.Marshal(meta)
		require.NoError(t, err)

		// Should contain string keys
		require.Contains(t, string(data), `"policy1"`)
		require.Contains(t, string(data), `"asset1"`)
	})

	t.Run("version 2 JSON", func(t *testing.T) {
		meta, err := NewCip25Metadata(2, policies)
		require.NoError(t, err)

		data, err := json.Marshal(meta)
		require.NoError(t, err)

		// Should still contain string keys in JSON (hex encoded)
		require.Contains(t, string(data), `"policy1"`)
		require.Contains(t, string(data), `"asset1"`)
	})

	t.Run("version 1 CBOR keys are strings", func(t *testing.T) {
		meta, err := NewCip25Metadata(1, policies)
		require.NoError(t, err)

		data, err := cbor.Marshal(meta.Num721)
		require.NoError(t, err)

		// Unmarshal to check key types
		var raw map[any]any
		err = cbor.Unmarshal(data, &raw)
		require.NoError(t, err)

		// Keys should be strings for v1
		for k := range raw {
			if keyStr, ok := k.(string); ok {
				if keyStr == "version" {
					continue
				}
				require.Equal(t, "policy1", keyStr)
			}
		}
	})

	t.Run("version 2 CBOR keys are bytes", func(t *testing.T) {
		meta, err := NewCip25Metadata(2, policies)
		require.NoError(t, err)

		data, err := cbor.Marshal(meta.Num721)
		require.NoError(t, err)

		// Unmarshal to check key types
		var raw map[any]any
		err = cbor.Unmarshal(data, &raw)
		require.NoError(t, err)

		// Keys should be byte strings for v2
		foundByteKey := false
		for k := range raw {
			if _, ok := k.(cbor.ByteString); ok {
				foundByteKey = true
				break
			}
		}
		require.True(t, foundByteKey, "Expected byte keys for version 2")
	})
}
