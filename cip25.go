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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"sort"
	"strconv"

	"github.com/fxamacker/cbor/v2"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

//nolint:recvcheck
type HexBytes string

func (h HexBytes) MarshalCBOR() ([]byte, error) {
	// Check if the string is valid hex
	if len(h)%2 == 0 && len(h) > 0 {
		if _, err := hex.DecodeString(string(h)); err == nil {
			// Valid hex, encode as byte string
			b, err := hex.DecodeString(string(h))
			if err != nil {
				return nil, err
			}
			return cbor.Marshal(b)
		}
	}
	// Not valid hex, encode as string
	return cbor.Marshal(string(h))
}

func (h *HexBytes) UnmarshalCBOR(data []byte) error { //nolint:recvcheck
	var b []byte
	if err := cbor.Unmarshal(data, &b); err != nil {
		// Try as string
		var s string
		if err := cbor.Unmarshal(data, &s); err != nil {
			return err
		}
		*h = HexBytes(s)
		return nil
	}
	*h = HexBytes(hex.EncodeToString(b))
	return nil
}

func (h HexBytes) String() string {
	return string(h)
}

func (h HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(h))
}

// Cip25Metadata is the top-level container for CIP-25 NFT metadata under the "721" tag.
type Cip25Metadata struct {
	Num721 Num721 `cbor:"721,keyasint" json:"721" validate:"required"`
}

// Num721 represents the metadata structure for CIP-25, supporting versions 1 and 2.
//
//nolint:recvcheck
type Num721 struct {
	Version  int                                     `json:"version,omitempty" cbor:"version,omitempty" validate:"min=1,max=2"`
	Policies map[HexBytes]map[HexBytes]AssetMetadata `json:"-"                 cbor:""                  validate:"required,min=1,dive,dive"`
}

// AssetMetadata represents the metadata for a single asset.
type AssetMetadata struct {
	Name        string         `json:"name"                cbor:"name"                  validate:"required"`
	Image       UriField       `json:"image"               cbor:"image"                 validate:"required"`
	MediaType   string         `json:"mediaType,omitempty" cbor:"mediaType,omitempty"`
	Description DescField      `json:"description"         cbor:"description,omitempty"`
	Files       []FileDetails  `json:"files,omitempty"     cbor:"files,omitempty"       validate:"dive"`
	Extra       map[string]any `json:"-"                   cbor:"-"` // Preserve unknown properties
}

// UnmarshalJSON custom unmarshals AssetMetadata, preserving unknown properties.
func (a *AssetMetadata) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	a.Extra = make(map[string]any)

	for k, v := range raw {
		switch k {
		case "name":
			if s, ok := v.(string); ok {
				a.Name = s
			}
		case "image":
			imgData, err := json.Marshal(v)
			if err != nil {
				return err
			}
			if err := a.Image.UnmarshalJSON(imgData); err != nil {
				return err
			}
		case "mediaType":
			if s, ok := v.(string); ok {
				a.MediaType = s
			}
		case "description":
			descData, err := json.Marshal(v)
			if err != nil {
				return err
			}
			if err := a.Description.UnmarshalJSON(descData); err != nil {
				return err
			}
		case "files":
			if files, ok := v.([]any); ok {
				for i, f := range files {
					fData, err := json.Marshal(f)
					if err != nil {
						return fmt.Errorf("invalid files[%d]: %w", i, err)
					}
					var file FileDetails
					if err := json.Unmarshal(fData, &file); err != nil {
						return fmt.Errorf("invalid files[%d]: %w", i, err)
					}
					a.Files = append(a.Files, file)
				}
			}
		default:
			a.Extra[k] = v
		}
	}

	return nil
}

// MarshalJSON custom marshals AssetMetadata, including unknown properties.
func (a AssetMetadata) MarshalJSON() ([]byte, error) {
	out := make(map[string]any)
	out["name"] = a.Name
	out["image"] = a.Image
	if a.MediaType != "" {
		out["mediaType"] = a.MediaType
	}
	if len(a.Description.Descriptions) > 0 {
		out["description"] = a.Description
	}
	if len(a.Files) > 0 {
		out["files"] = a.Files
	}
	maps.Copy(out, a.Extra)
	return json.Marshal(out)
}

// UriField supports either a single URI string or an array of URI strings.
//
//nolint:recvcheck
type UriField struct {
	Uris []string `validate:"required,min=1"`
}

// UnmarshalJSON attempts to parse 'image' or 'src' as a single string; if that fails, it tries an array of strings.
func (u *UriField) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		u.Uris = []string{single}
		return nil
	}

	arr := []string{}
	if err := json.Unmarshal(data, &arr); err == nil {
		u.Uris = arr
		return nil
	}

	return errors.New("URI field must be a string or an array of strings")
}

// MarshalJSON returns the URI field as a single string if only one URI is present, otherwise an array.
func (u UriField) MarshalJSON() ([]byte, error) {
	if len(u.Uris) == 1 {
		return json.Marshal(u.Uris[0])
	}
	return json.Marshal(u.Uris)
}

// MarshalCBOR returns the URI field as a single string if only one URI is present, otherwise an array.
func (u UriField) MarshalCBOR() ([]byte, error) {
	if len(u.Uris) == 1 {
		return cbor.Marshal(u.Uris[0])
	}
	return cbor.Marshal(u.Uris)
}

// UnmarshalCBOR attempts to parse as a single string; if that fails, it tries an array of strings.
func (u *UriField) UnmarshalCBOR(data []byte) error {
	// Check for null
	var nilVal any
	if err := cbor.Unmarshal(data, &nilVal); err == nil && nilVal == nil {
		u.Uris = nil
		return nil
	}

	var single string
	if err := cbor.Unmarshal(data, &single); err == nil {
		u.Uris = []string{single}
		return nil
	}

	arr := []string{}
	if err := cbor.Unmarshal(data, &arr); err == nil {
		u.Uris = arr
		return nil
	}

	return errors.New("URI field must be a string or an array of strings")
}

// DescField supports either a single description string or an array of description strings.
//
//nolint:recvcheck
type DescField struct {
	Descriptions []string
}

// UnmarshalJSON attempts to parse 'description' as a single string; if that fails, it tries an array of strings.
func (d *DescField) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		d.Descriptions = []string{single}
		return nil
	}

	arr := []string{}
	if err := json.Unmarshal(data, &arr); err == nil {
		d.Descriptions = arr
		return nil
	}

	return errors.New("description must be a string or an array of strings")
}

// MarshalJSON returns the description as a single string if only one description is present, otherwise an array.
// If no descriptions, returns empty string.
func (d DescField) MarshalJSON() ([]byte, error) {
	if len(d.Descriptions) == 0 {
		return json.Marshal("")
	}
	if len(d.Descriptions) == 1 {
		return json.Marshal(d.Descriptions[0])
	}
	return json.Marshal(d.Descriptions)
}

// MarshalCBOR returns the description as a single string if only one description is present, otherwise an array.
// If no descriptions, returns null.
func (d DescField) MarshalCBOR() ([]byte, error) {
	if len(d.Descriptions) == 0 {
		return cbor.Marshal(nil)
	}
	if len(d.Descriptions) == 1 {
		return cbor.Marshal(d.Descriptions[0])
	}
	return cbor.Marshal(d.Descriptions)
}

// UnmarshalCBOR attempts to parse as a single string; if that fails, it tries an array of strings.
func (d *DescField) UnmarshalCBOR(data []byte) error {
	// Check for null
	var nilVal any
	if err := cbor.Unmarshal(data, &nilVal); err == nil && nilVal == nil {
		d.Descriptions = nil
		return nil
	}

	var single string
	if err := cbor.Unmarshal(data, &single); err == nil {
		d.Descriptions = []string{single}
		return nil
	}

	arr := []string{}
	if err := cbor.Unmarshal(data, &arr); err == nil {
		d.Descriptions = arr
		return nil
	}

	return errors.New("description must be a string or an array of strings")
}

// FileDetails represents the details of a file associated with the asset.
type FileDetails struct {
	Name      string         `json:"name"      cbor:"name"      validate:"required"`
	MediaType string         `json:"mediaType" cbor:"mediaType" validate:"required"`
	Src       UriField       `json:"src"       cbor:"src"       validate:"required"`
	Extra     map[string]any `json:"-"         cbor:"-"` // Preserve unknown properties
}

// UnmarshalJSON custom unmarshals FileDetails, preserving unknown properties.
func (f *FileDetails) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	f.Extra = make(map[string]any)

	for k, v := range raw {
		switch k {
		case "name":
			if s, ok := v.(string); ok {
				f.Name = s
			}
		case "mediaType":
			if s, ok := v.(string); ok {
				f.MediaType = s
			}
		case "src":
			if srcData, err := json.Marshal(v); err == nil {
				if err := f.Src.UnmarshalJSON(srcData); err != nil {
					return err
				}
			} else {
				return err
			}
		default:
			f.Extra[k] = v
		}
	}

	return nil
}

// UnmarshalJSON handles unmarshaling the Num721 structure, supporting versions 1 and 2.
func (n *Num721) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Check for version (handle both number and string)
	if v, ok := raw["version"]; ok {
		switch version := v.(type) {
		case float64:
			n.Version = int(version)
		case string:
			// Try to parse as float first, then int
			if f, err := strconv.ParseFloat(version, 64); err == nil {
				n.Version = int(f)
			}
		}
	} else {
		// Default to version 1 if not specified
		n.Version = 1
	}

	n.Policies = make(map[HexBytes]map[HexBytes]AssetMetadata)

	for key, value := range raw {
		if key == "version" {
			continue
		}
		// key is policy_id, value is map[asset_name]AssetMetadata
		policyMap, ok := value.(map[string]any)
		if !ok {
			// Skip non-object entries instead of failing
			continue
		}
		assets := make(map[HexBytes]AssetMetadata)
		for assetKey, assetValue := range policyMap {
			var asset AssetMetadata
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

// MarshalJSON handles marshaling the Num721 structure.
func (n Num721) MarshalJSON() ([]byte, error) {
	out := make(map[string]any)
	if n.Version != 1 {
		out["version"] = n.Version
	}
	for policy, assets := range n.Policies {
		out[string(policy)] = assets
	}
	return json.Marshal(out)
}

// UnmarshalCBOR handles unmarshaling the Num721 structure from CBOR.
func (n *Num721) UnmarshalCBOR(data []byte) error {
	var raw map[any]any
	if err := cbor.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Check for version
	for key, value := range raw {
		if keyStr, ok := key.(string); ok && keyStr == "version" {
			if version, ok := value.(uint64); ok {
				if version > 100 {
					return errors.New("version number too large")
				}
				n.Version = int(version)
			}
			break
		}
	}
	if n.Version == 0 {
		// Default to version 1
		n.Version = 1
	}

	n.Policies = make(map[HexBytes]map[HexBytes]AssetMetadata)

	for key, value := range raw {
		var isVersion bool
		if keyStr, ok := key.(string); ok && keyStr == "version" {
			isVersion = true
		}
		if isVersion {
			continue
		}
		// Convert key to HexBytes
		var k HexBytes
		if n.Version == 2 {
			if s, ok := key.(string); ok {
				k = HexBytes(s)
			} else {
				var b []byte
				if bb, ok := key.([]byte); ok {
					b = bb
				} else if bs, ok := key.(cbor.ByteString); ok {
					b = []byte(bs)
				} else {
					return errors.New("expected key for v2")
				}
				k = HexBytes(hex.EncodeToString(b))
			}
		} else {
			if s, ok := key.(string); ok {
				k = HexBytes(s)
			} else {
				return errors.New("expected string key for v1")
			}
		}

		// value is map[asset_name]AssetMetadata
		policyMap, ok := value.(map[any]any)
		if !ok {
			return errors.New("invalid policy structure")
		}
		assets := make(map[HexBytes]AssetMetadata)
		for assetKey, assetValue := range policyMap {
			// Convert assetKey to HexBytes
			var ak HexBytes
			if n.Version == 2 {
				if s, ok := assetKey.(string); ok {
					ak = HexBytes(s)
				} else {
					var b []byte
					if bb, ok := assetKey.([]byte); ok {
						b = bb
					} else if bs, ok := assetKey.(cbor.ByteString); ok {
						b = []byte(bs)
					} else {
						return errors.New("expected asset key for v2")
					}
					ak = HexBytes(hex.EncodeToString(b))
				}
			} else {
				if s, ok := assetKey.(string); ok {
					ak = HexBytes(s)
				} else {
					return errors.New("expected string asset key for v1")
				}
			}

			var asset AssetMetadata
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

// MarshalCBOR handles marshaling the Num721 structure to CBOR.
func (n Num721) MarshalCBOR() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := cbor.NewEncoder(buf)
	if err := enc.StartIndefiniteMap(); err != nil {
		return nil, err
	}
	if n.Version != 1 {
		if err := enc.Encode("version"); err != nil {
			return nil, err
		}
		if err := enc.Encode(n.Version); err != nil {
			return nil, err
		}
	}
	keys := make([]HexBytes, 0, len(n.Policies))
	for policy := range n.Policies {
		keys = append(keys, policy)
	}
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i]) < string(keys[j])
	})
	for _, k := range keys {
		if n.Version == 2 {
			if b, err := hex.DecodeString(string(k)); err == nil {
				if err := enc.Encode(b); err != nil {
					return nil, err
				}
			} else {
				if err := enc.Encode([]byte(string(k))); err != nil {
					return nil, err
				}
			}
		} else {
			if err := enc.Encode(string(k)); err != nil {
				return nil, err
			}
		}
		// Encode the inner map
		inner := n.Policies[k]
		var innerKeys []HexBytes
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
			if n.Version == 2 {
				if b, err := hex.DecodeString(string(ik)); err == nil {
					if err := enc.Encode(b); err != nil {
						return nil, err
					}
				} else {
					if err := enc.Encode([]byte(string(ik))); err != nil {
						return nil, err
					}
				}
			} else {
				if err := enc.Encode(string(ik)); err != nil {
					return nil, err
				}
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

// NewCip25Metadata creates a new CIP-25 metadata object.
func NewCip25Metadata(
	version int,
	policies map[string]map[string]AssetMetadata,
) (*Cip25Metadata, error) {
	convertedPolicies := make(map[HexBytes]map[HexBytes]AssetMetadata)
	for k, v := range policies {
		convertedInner := make(map[HexBytes]AssetMetadata)
		for ik, iv := range v {
			convertedInner[HexBytes(ik)] = iv
		}
		convertedPolicies[HexBytes(k)] = convertedInner
	}

	metadata := &Cip25Metadata{
		Num721: Num721{Version: version, Policies: convertedPolicies},
	}

	if err := validate.Struct(metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

// Validate checks the CIP-25 metadata for validity.
func (c *Cip25Metadata) Validate() error {
	return validate.Struct(c)
}
