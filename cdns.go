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
	"fmt"

	"github.com/blinklabs-io/gouroboros/cbor"
)

type CardanoDnsTtl uint

type CardanoDnsDomain struct {
	Origin         []byte
	Records        []CardanoDnsDomainRecord
	AdditionalData CardanoDnsMaybe[any]
}

func (c CardanoDnsDomain) String() string {
	ret := fmt.Sprintf(
		"CardanoDnsDomain { Origin = %s, Records = [ ",
		c.Origin,
	)
	for idx, record := range c.Records {
		ret += record.String()
		if idx == len(c.Records)-1 {
			ret += " "
		} else {
			ret += ", "
		}
	}
	ret += "] }"
	return ret
}

func (c *CardanoDnsDomain) UnmarshalCBOR(cborData []byte) error {
	var tmpData cbor.Constructor
	if _, err := cbor.Decode(cborData, &tmpData); err != nil {
		return err
	}
	if tmpData.Constructor() != 1 {
		return fmt.Errorf("unexpected constructor index: %d", tmpData.Constructor())
	}
	tmpFields := tmpData.Fields()
	c.Origin = tmpFields[0].(cbor.ByteString).Bytes()
	for _, record := range tmpFields[1].([]any) {
		recordConstr := record.(cbor.Constructor)
		var tmpRecord CardanoDnsDomainRecord
		if _, err := cbor.Decode(recordConstr.Cbor(), &tmpRecord); err != nil {
			return err
		}
		c.Records = append(c.Records, tmpRecord)
	}
	return nil
}

type CardanoDnsDomainRecord struct {
	// This allows the type to be used with cbor.DecodeGeneric
	cbor.StructAsArray
	Lhs  []byte
	Ttl  CardanoDnsMaybe[CardanoDnsTtl]
	Type []byte
	Rhs  []byte
}

func (c *CardanoDnsDomainRecord) UnmarshalCBOR(data []byte) error {
	var tmpConstr cbor.Constructor
	if _, err := cbor.Decode(data, &tmpConstr); err != nil {
		return err
	}
	if tmpConstr.Constructor() != 1 {
		return fmt.Errorf("unexpected constructor index: %d", tmpConstr.Constructor())
	}
	return cbor.DecodeGeneric(tmpConstr.FieldsCbor(), c)
}

func (c CardanoDnsDomainRecord) String() string {
	return fmt.Sprintf(
		"CardanoDnsDomainRecord { Lhs = %s, Ttl = %d, Type = %s, Rhs = %s }",
		c.Lhs,
		c.Ttl.Value,
		c.Type,
		c.Rhs,
	)
}

type CardanoDnsMaybe[T any] struct {
	// This allows the type to be used with cbor.DecodeGeneric
	cbor.StructAsArray
	Value    T
	hasValue bool
}

func NewCardanoDnsMaybe[T any](v any) CardanoDnsMaybe[T] {
	if v == nil {
		return CardanoDnsMaybe[T]{}
	}
	return CardanoDnsMaybe[T]{
		Value:    v.(T),
		hasValue: true,
	}
}

func (c CardanoDnsMaybe[T]) HasValue() bool {
	return c.hasValue
}

func (c *CardanoDnsMaybe[T]) UnmarshalCBOR(data []byte) error {
	var tmpConstr cbor.Constructor
	if _, err := cbor.Decode(data, &tmpConstr); err != nil {
		return err
	}
	if tmpConstr.Constructor() == 0 {
		if err := cbor.DecodeGeneric(tmpConstr.FieldsCbor(), c); err != nil {
			return err
		}
		c.hasValue = true
	}
	return nil
}
