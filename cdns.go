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

type CardanoDnsDomain struct {
	Origin         []byte
	Records        []CardanoDnsDomainRecord
	AdditionalData any
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
	switch len(tmpFields) {
	case 2, 3:
		c.Origin = tmpFields[0].(cbor.ByteString).Bytes()
		for _, record := range tmpFields[1].([]any) {
			recordConstr := record.(cbor.Constructor)
			recordFields := recordConstr.Fields()
			if recordConstr.Constructor() != 1 {
				return fmt.Errorf("unexpected constructor index: %d", recordConstr.Constructor())
			}
			var tmpRecord CardanoDnsDomainRecord
			switch len(recordFields) {
			case 3:
				tmpRecord.Lhs = recordFields[0].(cbor.ByteString).Bytes()
				tmpRecord.Type = recordFields[1].(cbor.ByteString).Bytes()
				tmpRecord.Rhs = recordFields[2].(cbor.ByteString).Bytes()
			case 4:
				tmpRecord.Lhs = recordFields[0].(cbor.ByteString).Bytes()
				tmpRecord.Ttl = uint(recordFields[1].(uint64))
				tmpRecord.Type = recordFields[2].(cbor.ByteString).Bytes()
				tmpRecord.Rhs = recordFields[3].(cbor.ByteString).Bytes()
			default:
				return fmt.Errorf("unexpected constructor field length: %d", len(recordFields))
			}
			c.Records = append(c.Records, tmpRecord)
		}
		if len(tmpData.Fields()) == 3 {
			c.AdditionalData = tmpData.Fields()[2]
		}
	default:
		return fmt.Errorf("unexpected constructor field length: %d", len(tmpData.Fields()))
	}
	return nil
}

type CardanoDnsDomainRecord struct {
	Lhs  []byte
	Ttl  uint
	Type []byte
	Rhs  []byte
}

func (c CardanoDnsDomainRecord) String() string {
	return fmt.Sprintf(
		"CardanoDnsDomainRecord { Lhs = %s, Ttl = %d, Type = %s, Rhs = %s }",
		c.Lhs,
		c.Ttl,
		c.Type,
		c.Rhs,
	)
}
