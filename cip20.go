// Copyright 2023 Blink Labs Software
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
	"github.com/fxamacker/cbor/v2"
	"github.com/go-playground/validator/v10"
)

type Cip20Metadata struct {
	Num674 Num674 `cbor:"674,keyasint" json:"674" validate:"required"`
}

type Num674 struct {
	Msg []string `cbor:"msg" json:"msg" validate:"required,gt=0,dive,max=64"`
}

func (c *Cip20Metadata) UnmarshalCBOR(data []byte) error {
	var raw map[any]cbor.RawMessage
	if err := cbor.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		if !isCip20Label(key) {
			continue
		}
		var num674 Num674
		if err := cbor.Unmarshal(value, &num674); err != nil {
			return err
		}
		c.Num674 = num674
		return nil
	}
	return nil
}

func isCip20Label(key any) bool {
	switch v := key.(type) {
	case string:
		return v == "674"
	case uint64:
		return v == 674
	case int64:
		return v == 674
	case int:
		return v == 674
	}
	return false
}

func NewCip20Metadata(messages []string) (*Cip20Metadata, error) {
	validate := validator.New()

	metadata := &Cip20Metadata{Num674: Num674{Msg: messages}}

	if err := validate.Struct(metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func (c *Cip20Metadata) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
