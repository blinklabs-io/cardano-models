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

package models_test

import (
	"encoding/hex"
	"reflect"
	"testing"

	models "github.com/blinklabs-io/cardano-models"

	"github.com/blinklabs-io/gouroboros/cbor"
)

var cardanoDnsTestDefs = []struct {
	cborHex     string
	expectedObj models.CardanoDnsDomain
}{
	{
		cborHex: "d87a9f4b666f6f2e63617264616e6f9fd87a9f4b666f6f2e63617264616e6f426e734f6e73312e666f6f2e63617264616e6fffd87a9f4b666f6f2e63617264616e6f426e734f6e73322e666f6f2e63617264616e6fffd87a9f4f6e73312e666f6f2e63617264616e6f41614a3137322e32382e302e32ffd87a9f4f6e73322e666f6f2e63617264616e6f187b416147312e322e332e34ffffff",
		expectedObj: models.CardanoDnsDomain{
			Origin: []byte("foo.cardano"),
			Records: []models.CardanoDnsDomainRecord{
				{
					Lhs:  []byte("foo.cardano"),
					Type: []byte("ns"),
					Rhs:  []byte("ns1.foo.cardano"),
				},
				{
					Lhs:  []byte("foo.cardano"),
					Type: []byte("ns"),
					Rhs:  []byte("ns2.foo.cardano"),
				},
				{
					Lhs:  []byte("ns1.foo.cardano"),
					Type: []byte("a"),
					Rhs:  []byte("172.28.0.2"),
				},
				{
					Lhs:  []byte("ns2.foo.cardano"),
					Ttl:  123,
					Type: []byte("a"),
					Rhs:  []byte("1.2.3.4"),
				},
			},
		},
	},
}

func TestCardanoDnsDecode(t *testing.T) {
	for _, testDef := range cardanoDnsTestDefs {
		testDatumBytes, err := hex.DecodeString(testDef.cborHex)
		if err != nil {
			t.Fatalf("unexpected error decoding test datum hex: %s", err)
		}
		// Decode CBOR into object
		var testObj models.CardanoDnsDomain
		if _, err := cbor.Decode(testDatumBytes, &testObj); err != nil {
			t.Fatalf("unexpected error decoding test datum CBOR: %s", err)
		}
		if !reflect.DeepEqual(testObj, testDef.expectedObj) {
			t.Fatalf(
				"CBOR did not decode to expected object\n  got: %s\n  wanted: %s",
				testObj.String(),
				testDef.expectedObj.String(),
			)
		}
	}
}
