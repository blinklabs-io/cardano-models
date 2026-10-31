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
	"github.com/stretchr/testify/require"
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
					Ttl: models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
						models.CardanoDnsTtl(3600),
					),
				},
				{
					Lhs:  []byte("village.cardano"),
					Type: []byte("ns"),
					Rhs:  []byte("ns1.village.cardano"),
					Ttl: models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
						models.CardanoDnsTtl(28800),
					),
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
					Ttl: models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
						models.CardanoDnsTtl(3600),
					),
				},
				{
					Lhs:  []byte("enclave.cardano"),
					Type: []byte("ns"),
					Rhs:  []byte("ns1.enclave.cardano"),
					Ttl: models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
						models.CardanoDnsTtl(28800),
					),
				},
				{
					Lhs:  []byte("enclave.cardano"),
					Type: []byte("A"),
					Rhs:  []byte("172.28.0.2"),
					Ttl: models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
						models.CardanoDnsTtl(3600),
					),
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

// It ensures that encoding and then decoding a CardanoDnsDomain with CBOR preserves all data without any loss or mismatch.
func TestCardanoDnsRoundTripLossless(t *testing.T) {
	for _, testDef := range cardanoDnsTestDefs {
		enc, err := cbor.Encode(testDef.expectedObj)
		if err != nil {
			t.Fatalf("Failed to encode CardanoDnsDomain to CBOR: %v", err)
		}
		var dec models.CardanoDnsDomain
		if _, err := cbor.Decode(enc, &dec); err != nil {
			t.Fatalf(
				"Failed to decode CBOR back into CardanoDnsDomain: %v",
				err,
			)
		}
		if !reflect.DeepEqual(dec, testDef.expectedObj) {
			t.Fatalf("CBOR round-trip failed\n  got: %s\n  want: %s",
				dec.String(),
				testDef.expectedObj.String(),
			)
		}
	}
}

// It verifies that CardanoDnsMaybe encodes using constructor 0 for "Just(value)" and constructor 1 for "None".
func TestCardanoDnsMaybe_Encode_Constructors(t *testing.T) {
	just := models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
		models.CardanoDnsTtl(60),
	)
	none := models.NewCardanoDnsMaybe[models.CardanoDnsTtl](nil)

	// Encode Just(value) and verify whether it uses constructor 0
	justBytes, err := cbor.Encode(just)
	if err != nil {
		t.Fatalf("Failed to encode 'Just(value)': %v", err)
	}
	var jc cbor.ConstructorDecoder
	if _, err := cbor.Decode(justBytes, &jc); err != nil {
		t.Fatalf("Failed to decode CBOR for 'Just(value)': %v", err)
	}
	if jc.Tag() != 0 {
		t.Fatalf(
			"Invalid constructor for 'Just(value)': got %d want 0",
			jc.Tag(),
		)
	}

	// Encode "None" value and verify whether it uses constructor 1 or not
	noneBytes, err := cbor.Encode(none)
	if err != nil {
		t.Fatalf("Failed to encode 'None' value: %v", err)
	}
	var nc cbor.ConstructorDecoder
	if _, err := cbor.Decode(noneBytes, &nc); err != nil {
		t.Fatalf("Failed to decode CBOR for 'None' value: %v", err)
	}
	if nc.Tag() != 1 {
		t.Fatalf(
			"Invalid constructor for 'None' value: got %d want 1",
			nc.Tag(),
		)
	}
}

// It ensures that a single CardanoDnsDomainRecord encodes with constructor 1
// and decodes back correctly without data loss.
func TestCardanoDnsDomainRecord_Encode_RoundTrip(t *testing.T) {
	rec := models.CardanoDnsDomainRecord{
		Lhs:  []byte("enclave.cardano"),
		Type: []byte("A"),
		Rhs:  []byte("401.401.401.401"),
		Ttl: models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
			models.CardanoDnsTtl(28800),
		),
	}
	enc, err := cbor.Encode(rec)
	if err != nil {
		t.Fatalf("Failed to encode CardanoDnsDomainRecord to CBOR: %v", err)
	}

	var cons cbor.ConstructorDecoder
	if _, err := cbor.Decode(enc, &cons); err != nil {
		t.Fatalf("Failed to decode CBOR into Constructor: %v", err)
	}
	if cons.Tag() != 1 {
		t.Fatalf(
			"Invalid constructor tag for CardanoDnsDomainRecord: got %d want 1",
			cons.Tag(),
		)
	}
	var got models.CardanoDnsDomainRecord
	if _, err := cbor.Decode(enc, &got); err != nil {
		t.Fatalf(
			"Failed to decode constructor fields into CardanoDnsDomainRecord: %v",
			err,
		)
	}
	if !reflect.DeepEqual(got, rec) {
		t.Fatalf(
			"Round-trip mismatch for CardanoDnsDomainRecord:\n  got: %+v\n  want: %+v",
			got,
			rec,
		)
	}
}

func TestCardanoDnsMaybe_DecodeValidation(t *testing.T) {
	encode := func(t *testing.T, tag uint, fields any) []byte {
		t.Helper()
		data, err := cbor.Encode(cbor.NewConstructorEncoder(tag, fields))
		require.NoError(t, err)
		return data
	}

	someBytes := encode(t, 0, []any{models.CardanoDnsTtl(60)})
	noneBytes := encode(t, 1, []any{})

	t.Run("valid Some", func(t *testing.T) {
		var got models.CardanoDnsMaybe[models.CardanoDnsTtl]
		_, err := cbor.Decode(someBytes, &got)
		require.NoError(t, err)
		require.True(t, got.HasValue())
		require.Equal(t, models.CardanoDnsTtl(60), got.Value)
	})

	t.Run("None clears reused Some", func(t *testing.T) {
		got := models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
			models.CardanoDnsTtl(60),
		)
		_, err := cbor.Decode(noneBytes, &got)
		require.NoError(t, err)
		require.False(t, got.HasValue())
		require.Zero(t, got.Value)
	})

	for name, data := range map[string][]byte{
		"unknown constructor":        encode(t, 2, []any{}),
		"Some without a field":       encode(t, 0, []any{}),
		"Some with two fields":       encode(t, 0, []any{1, 2}),
		"Some with wrong value type": encode(t, 0, []any{"wrong"}),
		"None with a field":          encode(t, 1, []any{1}),
		"Some with non-array fields": encode(t, 0, nil),
		"None with non-array fields": encode(t, 1, nil),
	} {
		t.Run(name, func(t *testing.T) {
			got := models.NewCardanoDnsMaybe[models.CardanoDnsTtl](
				models.CardanoDnsTtl(99),
			)
			_, err := cbor.Decode(data, &got)
			require.Error(t, err)
			require.True(t, got.HasValue())
			require.Equal(t, models.CardanoDnsTtl(99), got.Value)
		})
	}
}

func TestCardanoDnsDomain_DecodeValidation(t *testing.T) {
	encode := func(t *testing.T, value any) []byte {
		t.Helper()
		data, err := cbor.Encode(value)
		require.NoError(t, err)
		return data
	}
	none := cbor.NewConstructorEncoder(1, []any{})
	validRecord := cbor.NewConstructorEncoder(
		1,
		[]any{[]byte("example.cardano"), none, []byte("A"), []byte("192.0.2.1")},
	)
	domain := func(records any) cbor.ConstructorEncoder {
		return cbor.NewConstructorEncoder(
			1,
			[]any{[]byte("example"), records, none},
		)
	}

	for name, testDef := range map[string]struct {
		data       []byte
		errorParts []string
	}{
		"invalid CBOR": {
			data:       []byte{0xff},
			errorParts: []string{"unexpected \"break\" code"},
		},
		"wrong domain constructor": {
			data: encode(t, cbor.NewConstructorEncoder(2, []any{})),
			errorParts: []string{
				"CardanoDnsDomain",
				"unexpected constructor index: 2",
			},
		},
		"wrong origin type": {
			data: encode(t, cbor.NewConstructorEncoder(
				1,
				[]any{"example", []any{validRecord}, none},
			)),
			errorParts: []string{"decode CardanoDnsDomain fields", "Origin"},
		},
		"wrong records type": {
			data: encode(t, domain(uint64(0))),
			errorParts: []string{
				"decode CardanoDnsDomain fields",
				"Records",
			},
		},
		"record is not a constructor": {
			data: encode(t, domain([]any{uint64(0)})),
			errorParts: []string{
				"decode CardanoDnsDomain record 0",
				"decode CardanoDnsDomainRecord constructor",
			},
		},
		"wrong record constructor": {
			data: encode(t, domain([]any{
				cbor.NewConstructorEncoder(2, []any{}),
			})),
			errorParts: []string{
				"decode CardanoDnsDomain record 0",
				"CardanoDnsDomainRecord",
				"unexpected constructor index: 2",
			},
		},
		"wrong record field type": {
			data: encode(t, domain([]any{
				cbor.NewConstructorEncoder(
					1,
					[]any{"example.cardano", none, []byte("A"), []byte("192.0.2.1")},
				),
			})),
			errorParts: []string{
				"decode CardanoDnsDomain record 0",
				"decode CardanoDnsDomainRecord fields",
				"Lhs",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			original := models.CardanoDnsDomain{
				Origin:  []byte("unchanged"),
				Records: []models.CardanoDnsDomainRecord{{Lhs: []byte("unchanged")}},
			}
			got := original
			_, err := cbor.Decode(testDef.data, &got)
			require.Error(t, err)
			for _, errorPart := range testDef.errorParts {
				require.ErrorContains(t, err, errorPart)
			}
			require.Equal(t, original, got)
		})
	}
}

func TestCardanoDnsDomainRecord_DecodeValidation(t *testing.T) {
	encode := func(t *testing.T, value any) []byte {
		t.Helper()
		data, err := cbor.Encode(value)
		require.NoError(t, err)
		return data
	}
	none := cbor.NewConstructorEncoder(1, []any{})

	for name, testDef := range map[string]struct {
		data       []byte
		errorParts []string
	}{
		"wrong constructor": {
			data: encode(t, cbor.NewConstructorEncoder(0, []any{})),
			errorParts: []string{
				"CardanoDnsDomainRecord",
				"unexpected constructor index: 0",
			},
		},
		"missing field": {
			data: encode(t, cbor.NewConstructorEncoder(
				1,
				[]any{[]byte("example.cardano"), none, []byte("A")},
			)),
			errorParts: []string{"decode CardanoDnsDomainRecord fields"},
		},
		"wrong TTL value type": {
			data: encode(t, cbor.NewConstructorEncoder(
				1,
				[]any{
					[]byte("example.cardano"),
					cbor.NewConstructorEncoder(0, []any{"sixty"}),
					[]byte("A"),
					[]byte("192.0.2.1"),
				},
			)),
			errorParts: []string{
				"decode CardanoDnsDomainRecord fields",
				"Ttl",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			original := models.CardanoDnsDomainRecord{Lhs: []byte("unchanged")}
			got := original
			_, err := cbor.Decode(testDef.data, &got)
			require.Error(t, err)
			for _, errorPart := range testDef.errorParts {
				require.ErrorContains(t, err, errorPart)
			}
			require.Equal(t, original, got)
		})
	}
}
