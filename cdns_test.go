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
		cborHex: "d87a9f4776696c6c6167659fd87a9f4f76696c6c6167652e63617264616e6fd8799f190e10ff41414a3137322e32382e302e32ffd87a9f4f76696c6c6167652e63617264616e6fd8799f197080ff426e73536e73312e76696c6c6167652e63617264616e6fffffd87a80ff",
		expectedObj: models.CardanoDnsDomain{
			Origin: []byte("village"),
			Records: []models.CardanoDnsDomainRecord{
				{
					Lhs:  []byte("village.cardano"),
					Type: []byte("A"),
					Rhs:  []byte("172.28.0.2"),
					Ttl:  models.NewCardanoDnsMaybe[models.CardanoDnsTtl](models.CardanoDnsTtl(3600)),
				},
				{
					Lhs:  []byte("village.cardano"),
					Type: []byte("ns"),
					Rhs:  []byte("ns1.village.cardano"),
					Ttl:  models.NewCardanoDnsMaybe[models.CardanoDnsTtl](models.CardanoDnsTtl(28800)),
				},
			},
		},
	},
	{
		cborHex: "d87a9f47656e636c6176659fd87a9f4f656e636c6176652e63617264616e6fd8799f190e10ff41414f3430312e3430312e3430312e343031ffd87a9f4f656e636c6176652e63617264616e6fd8799f197080ff426e73536e73312e656e636c6176652e63617264616e6fffd87a9f4f656e636c6176652e63617264616e6fd8799f190e10ff41414a3137322e32382e302e32ffd87a9f4f656e636c6176652e63617264616e6fd87a80426e73536e73322e656e636c6176652e63617264616e6fffffd87a80ff",
		expectedObj: models.CardanoDnsDomain{
			Origin: []byte("enclave"),
			Records: []models.CardanoDnsDomainRecord{
				{
					Lhs:  []byte("enclave.cardano"),
					Type: []byte("A"),
					Rhs:  []byte("401.401.401.401"),
					Ttl:  models.NewCardanoDnsMaybe[models.CardanoDnsTtl](models.CardanoDnsTtl(3600)),
				},
				{
					Lhs:  []byte("enclave.cardano"),
					Type: []byte("ns"),
					Rhs:  []byte("ns1.enclave.cardano"),
					Ttl:  models.NewCardanoDnsMaybe[models.CardanoDnsTtl](models.CardanoDnsTtl(28800)),
				},
				{
					Lhs:  []byte("enclave.cardano"),
					Type: []byte("A"),
					Rhs:  []byte("172.28.0.2"),
					Ttl:  models.NewCardanoDnsMaybe[models.CardanoDnsTtl](models.CardanoDnsTtl(3600)),
				},
				{
					Lhs:  []byte("enclave.cardano"),
					Type: []byte("ns"),
					Rhs:  []byte("ns2.enclave.cardano"),
					Ttl:  models.NewCardanoDnsMaybe[models.CardanoDnsTtl](nil),
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
