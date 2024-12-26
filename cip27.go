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
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// Cip27Metadata is the top-level container for royalties data under the "777" tag.
type Cip27Metadata struct {
	Num777 Cip777 `cbor:"777,keyasint" json:"777" validate:"required"`
}

// Cip777 represents the actual royalty info. It handles both modern "rate" and legacy "pct."
type Cip777 struct {
	// Internally, Rate is our main numeric string (e.g., "0.20").
	Rate string `json:"-"`

	// We store the raw strings for either "pct" or "rate."
	// We only expose "rate" in the final JSON.
	pctRaw  *string
	rateRaw *string

	// 'addr' can be either a string or array of strings, so we wrap it in AddrField.
	Addr AddrField `cbor:"addr" json:"addr" validate:"required"`
}

func (c *Cip27Metadata) UnmarshalJSON(data []byte) error {
	// Unmarshal into a map so we can check for "777" explicitly.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Verify the "777" key exists at the top level.
	val, ok := raw["777"]
	if !ok {
		return errors.New(`missing "777" key in CIP-27 metadata`)
	}

	// Unmarshal the contents of "777" into c.Num777.
	if err := json.Unmarshal(val, &c.Num777); err != nil {
		return err
	}

	// Run validation (so any "required" or numeric checks fail immediately).
	if err := c.Validate(); err != nil {
		return err
	}
	return nil
}

// UnmarshalJSON checks which field ("rate" or "pct") is present, giving precedence to "rate."
func (c *Cip777) UnmarshalJSON(data []byte) error {
	// Temporary structure for decoding both fields plus 'addr.'
	var raw struct {
		Pct  *string   `json:"pct"`
		Rate *string   `json:"rate"`
		Addr AddrField `json:"addr"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch {
	case raw.Rate != nil:
		c.Rate = *raw.Rate
	case raw.Pct != nil:
		c.Rate = *raw.Pct
	default:
		return errors.New("missing both 'rate' and 'pct' fields")
	}

	c.pctRaw = raw.Pct
	c.rateRaw = raw.Rate
	c.Addr = raw.Addr
	return nil
}

// MarshalJSON outputs "rate" as our canonical field.
func (c Cip777) MarshalJSON() ([]byte, error) {
	// We only expose "rate" in the final JSON.
	var out struct {
		Rate string    `json:"rate"`
		Addr AddrField `json:"addr"`
	}
	out.Rate = c.Rate
	out.Addr = c.Addr
	return json.Marshal(out)
}

// AddrField supports either a single string or an array of strings in JSON.
type AddrField struct {
	Addresses []string
}

// UnmarshalJSON attempts to parse 'addr' as a single string; if that fails, it tries an array of strings.
func (af *AddrField) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		af.Addresses = []string{single}
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		af.Addresses = arr
		return nil
	}

	return errors.New("addr must be a string or an array of strings")
}

// MarshalJSON returns 'addr' as a single string if only one address is present, otherwise an array.
func (af AddrField) MarshalJSON() ([]byte, error) {
	if len(af.Addresses) == 1 {
		return json.Marshal(af.Addresses[0])
	}
	return json.Marshal(af.Addresses)
}

// NewCip27Metadata creates a new CIP-027 metadata object with the given rate and addresses.
func NewCip27Metadata(rate string, addresses []string) (*Cip27Metadata, error) {
	meta := &Cip27Metadata{
		Num777: Cip777{
			Rate: rate,
			Addr: AddrField{Addresses: addresses},
		},
	}
	if err := meta.Validate(); err != nil {
		return nil, err
	}
	return meta, nil
}

// Validate checks that Rate is within [0..1] and there's at least one address.
func (c *Cip27Metadata) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}
	val, err := strconv.ParseFloat(c.Num777.Rate, 64)
	if err != nil {
		return errors.New("rate must be a valid floating point number")
	}
	if val < 0 || val > 1 {
		return errors.New("rate must be between 0.0 and 1.0")
	}
	if len(c.Num777.Addr.Addresses) == 0 {
		return errors.New("at least one address is required")
	}
	return nil
}
