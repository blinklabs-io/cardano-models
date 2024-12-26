// Copyright 2024 Blink Labs Software
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
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCip27Metadata_Success(t *testing.T) {
	// Single address, valid rate.
	meta, err := NewCip27Metadata("0.25", []string{"addr1xy..."})
	require.NoError(t, err)
	require.Equal(t, "0.25", meta.Num777.Rate)
	require.Len(t, meta.Num777.Addr.Addresses, 1)
}

func TestNewCip27Metadata_MultipleAddresses(t *testing.T) {
	addrs := []string{"addr1abc", "addr2def"}
	meta, err := NewCip27Metadata("0.25", addrs)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(addrs, meta.Num777.Addr.Addresses))
}

func TestRateBoundaries(t *testing.T) {
	// Valid boundary
	_, err := NewCip27Metadata("1.0", []string{"addr1..."})
	require.NoError(t, err)

	// Out-of-range: 1.1
	_, err = NewCip27Metadata("1.1", []string{"addr1..."})
	require.Error(t, err)

	// Negative
	_, err = NewCip27Metadata("-0.1", []string{"addr1..."})
	require.Error(t, err)

	// Not a float
	_, err = NewCip27Metadata("abc", []string{"addr1..."})
	require.Error(t, err)
}

func TestUnmarshal_LegacyPct(t *testing.T) {
	input := `{
      "777": {
        "pct": "0.125",
        "addr": "addr1legacy..."
      }
    }`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.NoError(t, err)
	require.Equal(t, "0.125", meta.Num777.Rate)
	require.Equal(t, "addr1legacy...", meta.Num777.Addr.Addresses[0])
}

func TestUnmarshal_Rate(t *testing.T) {
	input := `{
      "777": {
        "rate": "0.20",
        "addr": ["addr1qmodern","addr2other"]
      }
    }`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.NoError(t, err)
	require.Equal(t, "0.20", meta.Num777.Rate)
	require.Len(t, meta.Num777.Addr.Addresses, 2)
}

func TestUnmarshal_BothPctAndRate(t *testing.T) {
	// If both are present, 'rate' overrides 'pct'.
	input := `{
      "777": {
        "pct": "0.100",
        "rate": "0.200",
        "addr": "addr1override"
      }
    }`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.NoError(t, err)
	require.Equal(t, "0.200", meta.Num777.Rate)
	require.Equal(t, "addr1override", meta.Num777.Addr.Addresses[0])
}

func TestUnmarshal_MissingPctRate(t *testing.T) {
	// Neither 'pct' nor 'rate' is present -> error
	input := `{"777":{"addr":"addr1only"}}`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.Error(t, err)
}

func TestAddrField_SingleString(t *testing.T) {
	var af AddrField
	err := json.Unmarshal([]byte(`"addrSingle"`), &af)
	require.NoError(t, err)
	require.Equal(t, []string{"addrSingle"}, af.Addresses)
}

func TestAddrField_StringArray(t *testing.T) {
	var af AddrField
	err := json.Unmarshal([]byte(`["addr1","addr2"]`), &af)
	require.NoError(t, err)
	require.Equal(t, []string{"addr1", "addr2"}, af.Addresses)
}

func TestAddrField_InvalidType(t *testing.T) {
	var af AddrField
	err := json.Unmarshal([]byte("123"), &af)
	require.Error(t, err)
}

func TestAddrField_MarshalSingle(t *testing.T) {
	single := AddrField{Addresses: []string{"addr1single"}}
	b, err := json.Marshal(single)
	require.NoError(t, err)
	require.Equal(t, `"addr1single"`, string(b))
}

func TestAddrField_MarshalArray(t *testing.T) {
	multiple := AddrField{Addresses: []string{"addr1", "addr2"}}
	b, err := json.Marshal(multiple)
	require.NoError(t, err)
	require.Equal(t, `["addr1","addr2"]`, string(b))
}

func TestUnmarshal_SingleAddressExample(t *testing.T) {
	input := `{
      "777": {
          "rate": "0.2",
          "addr": "addr1v9nevxg9wunfck0gt7hpxuy0elnqygglme3u6l3nn5q5gnq5dc9un"
      }
  }`

	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.NoError(t, err, "Should unmarshal single address JSON without error")

	// Check that 'rate' is "0.2" and that we have exactly one address in our slice
	require.Equal(t, "0.2", meta.Num777.Rate)
	require.Len(t, meta.Num777.Addr.Addresses, 1)
	require.Equal(t,
		"addr1v9nevxg9wunfck0gt7hpxuy0elnqygglme3u6l3nn5q5gnq5dc9un",
		meta.Num777.Addr.Addresses[0],
	)
}

func TestUnmarshal_ArrayAddressExample(t *testing.T) {
	input := `{
      "777": {
          "rate": "0.2",
          "addr": [
              "addr1q8g3dv6ptkgsafh7k5muggrvfde2szzmc2mqkcxpxn7c63l9znc9e3xa82h",
              "pf39scc37tcu9ggy0l89gy2f9r2lf7husfvu8wh"
          ]
      }
  }`

	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.NoError(t, err, "Should unmarshal multiple address JSON without error")

	// Check that 'rate' is "0.2" and that we have two addresses in our slice
	require.Equal(t, "0.2", meta.Num777.Rate)
	require.Len(t, meta.Num777.Addr.Addresses, 2)
	require.Equal(t,
		"addr1q8g3dv6ptkgsafh7k5muggrvfde2szzmc2mqkcxpxn7c63l9znc9e3xa82h",
		meta.Num777.Addr.Addresses[0],
	)
	require.Equal(t,
		"pf39scc37tcu9ggy0l89gy2f9r2lf7husfvu8wh",
		meta.Num777.Addr.Addresses[1],
	)
}

func TestUnmarshal_Cip27Metadata_InvalidRootType(t *testing.T) {
	// For example: top-level is a string, not an object.
	input := `"this is not an object"`

	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.Error(t, err, "Should fail unmarshaling top-level non-object JSON")
}

func TestUnmarshal_Cip27Metadata_No777Key(t *testing.T) {
	input := `{"someOtherKey":{"rate":"0.2"}}`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.Error(t, err, "Should fail if '777' key is missing")
}

func TestUnmarshal_EmptyRateString(t *testing.T) {
	input := `{
    "777": {
      "rate": "",
      "addr": "addr1..."
    }
  }`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.Error(t, err, "Empty rate string should not parse as a valid float")
}

func TestUnmarshal_NoAddrKey(t *testing.T) {
	input := `{
    "777": {
      "rate": "0.2"
    }
  }`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.Error(t, err, "Should error because 'addr' key is required")
}

func TestUnmarshal_AddrAsObject(t *testing.T) {
	input := `{
    "777": {
      "rate": "0.2",
      "addr": { "some": "thing" }
    }
  }`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.Error(t, err, "Should fail because 'addr' must be string or array of strings")
}

func TestUnmarshal_EmptyAddrArray(t *testing.T) {
	input := `{
    "777": {
      "rate": "0.25",
      "addr": []
    }
  }`

	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.Error(t, err, "Should fail because we have 0 addresses in 'addr'")
	require.Contains(t, err.Error(), "at least one address is required")
}

func TestUnmarshal_BothPctAndRate_RateEmpty(t *testing.T) {
	input := `{
    "777": {
      "pct": "0.1",
      "rate": "",
      "addr": "addr1..."
    }
  }`
	var meta Cip27Metadata
	err := json.Unmarshal([]byte(input), &meta)
	require.Error(t, err, "Empty rate should fail even if pct is valid, because rate takes precedence")
}

func TestCip27Metadata_MarshalJSON(t *testing.T) {
	// Create a CIP-27 metadata object.
	meta, err := NewCip27Metadata("0.25", []string{"addr1xy...", "addr2zzz"})
	require.NoError(t, err, "Should create CIP-27 metadata without error")

	data, err := json.Marshal(meta)
	require.NoError(t, err, "Should marshal CIP-27 metadata to JSON without error")

	// Check that the resulting JSON contains expected fields/values.
	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result), "Should unmarshal JSON back into map")

	// Expect top-level key "777"
	topLevel, ok := result["777"].(map[string]interface{})
	require.True(t, ok, "Should have a '777' key in marshaled JSON")

	// Ensure "rate" is "0.25"
	require.Equal(t, "0.25", topLevel["rate"], "Rate should match the input string")

	// Since addresses can be single string or an array,
	// we expect them as an array (2 addresses).
	addrs, ok := topLevel["addr"].([]interface{})
	require.True(t, ok, "Should have an array of addresses")
	require.Len(t, addrs, 2, "Should have exactly 2 addresses")
	require.Equal(t, "addr1xy...", addrs[0])
	require.Equal(t, "addr2zzz", addrs[1])
}
