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
	"fmt"
	"strings"

	"github.com/blinklabs-io/gouroboros/cbor"
)

type CardanoDnsTtl uint

type CardanoDnsDomain struct {
	// This allows the type to be reused for the inner content during decoding
	cbor.StructAsArray
	Origin         []byte
	Records        []CardanoDnsDomainRecord
	AdditionalData CardanoDnsMaybe[any]
}

func (c *CardanoDnsDomain) String() string {
	var sb strings.Builder
	sb.WriteString("CardanoDnsDomain { Origin = ")
	sb.WriteString(string(c.Origin))
	sb.WriteString(", Records = [ ")
	for idx, record := range c.Records {
		sb.WriteString(record.String())
		if idx == len(c.Records)-1 {
			sb.WriteString(" ")
		} else {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("] }")
	return sb.String()
}

func (c *CardanoDnsDomain) UnmarshalCBOR(cborData []byte) error {
	var tmpData cbor.Constructor
	if _, err := cbor.Decode(cborData, &tmpData); err != nil {
		return err
	}
	if tmpData.Constructor() != 1 {
		return fmt.Errorf(
			"unexpected constructor index: %d",
			tmpData.Constructor(),
		)
	}
	type tCardanoDnsDomain CardanoDnsDomain
	var tmpCardanoDnsDomain tCardanoDnsDomain
	if _, err := cbor.Decode(tmpData.FieldsCbor(), &tmpCardanoDnsDomain); err != nil {
		return err
	}
	*c = CardanoDnsDomain(tmpCardanoDnsDomain)
	return nil
}

func (c *CardanoDnsDomain) MarshalCBOR() ([]byte, error) {
	tmpData := cbor.NewConstructor(
		1,
		[]any{
			c.Origin,
			c.Records,
			c.AdditionalData,
		},
	)
	return cbor.Encode(tmpData)
}

type CardanoDnsDomainRecord struct {
	// This allows the type to be reused for the inner content during decoding
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
		return fmt.Errorf(
			"unexpected constructor index: %d",
			tmpConstr.Constructor(),
		)
	}
	type tCardanoDnsDomainRecord CardanoDnsDomainRecord
	var tmpCardanoDnsDomainRecord tCardanoDnsDomainRecord
	if _, err := cbor.Decode(tmpConstr.FieldsCbor(), &tmpCardanoDnsDomainRecord); err != nil {
		return err
	}
	*c = CardanoDnsDomainRecord(tmpCardanoDnsDomainRecord)
	return nil
}

func (r *CardanoDnsDomainRecord) MarshalCBOR() ([]byte, error) {
	tmpData := cbor.NewConstructor(
		1,
		[]any{
			r.Lhs,
			r.Ttl,
			r.Type,
			r.Rhs,
		},
	)
	return cbor.Encode(tmpData)
}

func (c *CardanoDnsDomainRecord) String() string {
	return fmt.Sprintf(
		"CardanoDnsDomainRecord { Lhs = %s, Ttl = %d, Type = %s, Rhs = %s }",
		c.Lhs,
		c.Ttl.Value,
		c.Type,
		c.Rhs,
	)
}

type CardanoDnsMaybe[T any] struct {
	// This allows the type to be reused when decoding
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
		type tCardanoDnsMaybe CardanoDnsMaybe[T]
		var tmpCardanoDnsMaybe tCardanoDnsMaybe
		if _, err := cbor.Decode(tmpConstr.FieldsCbor(), &tmpCardanoDnsMaybe); err != nil {
			return err
		}
		*c = CardanoDnsMaybe[T](tmpCardanoDnsMaybe)
		c.hasValue = true
	}
	return nil
}

func (m *CardanoDnsMaybe[T]) MarshalCBOR() ([]byte, error) {
	if !m.HasValue() {
		// None: constructor(1) with empty field array
		return cbor.Encode(cbor.NewConstructor(1, []any{}))
	}
	// Some(Value): constructor(0) with single-field array
	return cbor.Encode(cbor.NewConstructor(0, []any{m.Value}))
}
