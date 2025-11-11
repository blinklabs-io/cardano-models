package models

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
)

func TestValidCip20Metadata(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                     string
		cborHex                  string
		expectedObj              Cip20Metadata
		jsonData                 string
		expectCBORUnmarshalError bool
		expectJSONUnmarshalError bool
		expectValidationError    bool
		expectCBORDeepEqualError bool
		expectJSONDeepEqualError bool
	}{
		{
			name: "Valid - Single Message",
			expectedObj: Cip20Metadata{
				Num674: Num674{
					Msg: []string{
						"This is a comment for the transaction xyz, thank you very much!",
					},
				},
			},
			cborHex: "A163363734A1636D736781783F54686973206973206120636F6D6D656E7420666F7220746865207472616E73616374696F6E2078797A2C207468616E6B20796F752076657279206D75636821",
			jsonData: `{
		    "674":
		    {
		      "msg":
		      [
		        "This is a comment for the transaction xyz, thank you very much!"
		      ]
		    }
		  }`,
			expectCBORUnmarshalError: false,
			expectJSONUnmarshalError: false,
			expectValidationError:    false,
			expectCBORDeepEqualError: false,
			expectJSONDeepEqualError: false,
		},
		{
			name: "Valid - Multiple Messages",
			expectedObj: Cip20Metadata{
				Num674: Num674{
					Msg: []string{
						"Invoice-No: 1234567890",
						"Customer-No: 555-1234",
						"P.S.: i will shop again at your store :-)",
					},
				},
			},
			jsonData: `{
		    "674":
		           {
		             "msg":
		                    [
		                      "Invoice-No: 1234567890",
		                      "Customer-No: 555-1234",
		                      "P.S.: i will shop again at your store :-)"
		                    ]
		           }
		  }`,
			cborHex:                  "A163363734A1636D73678376496E766F6963652D4E6F3A203132333435363738393075437573746F6D65722D4E6F3A203535352D313233347829502E532E3A20692077696C6C2073686F7020616761696E20617420796F75722073746F7265203A2D29",
			expectCBORUnmarshalError: false,
			expectJSONUnmarshalError: false,
			expectValidationError:    false,
			expectCBORDeepEqualError: false,
			expectJSONDeepEqualError: false,
		},
		{
			name: "Invalid - Message Exceeds 64 Bytes",
			expectedObj: Cip20Metadata{
				Num674: Num674{
					Msg: []string{
						strings.Repeat("a", 65),
					}, // 65 'a's to exceed the limit
				},
			},
			jsonData: `{
        "674": {
          "msg": ["aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"]
		    }
        }`,
			cborHex:                  "A163363734A1636D73678178416161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161",
			expectCBORUnmarshalError: false,
			expectJSONUnmarshalError: false,
			expectValidationError:    true,
			expectCBORDeepEqualError: false,
			expectJSONDeepEqualError: false,
		},
		{
			name: "Invalid - Message Deep error",
			expectedObj: Cip20Metadata{
				Num674: Num674{
					Msg: []string{strings.Repeat("b", 65)},
				},
			},
			jsonData: `{
        "674": {
          "msg": ["aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"]
		    }
        }`,
			cborHex:                  "A163363734A1636D73678178416161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161",
			expectCBORUnmarshalError: false,
			expectJSONUnmarshalError: false,
			expectValidationError:    true,
			expectCBORDeepEqualError: true,
			expectJSONDeepEqualError: true,
		},
		{
			name: "Invalid - Message Deep error 675",
			expectedObj: Cip20Metadata{
				Num674: Num674{
					Msg: []string{strings.Repeat("b", 65)},
				},
			},
			jsonData: `{
        "675": {
          "msg": ["aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"]
		    }
        }`,
			cborHex:                  "A163363734A1636D73678178416161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161",
			expectCBORUnmarshalError: false,
			expectJSONUnmarshalError: false,
			expectValidationError:    true,
			expectCBORDeepEqualError: true,
			expectJSONDeepEqualError: true,
		},
	}

	for _, tc := range testCases {
		// capture range variable for goroutines
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Decode the CBOR hex string
			cborData, err := hex.DecodeString(tc.cborHex)
			if err != nil {
				t.Errorf(
					"failed to decode CBOR hex string for %s, error: %v",
					tc.name,
					err,
				)
			}

			// CBOR Unmarshal
			var decodedMetadata Cip20Metadata
			err = cbor.Unmarshal(cborData, &decodedMetadata)
			if tc.expectCBORUnmarshalError {
				if err == nil {
					t.Errorf(
						"expected CBOR unmarshal error but got none for test: %s",
						tc.name,
					)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect CBOR unmarshal error but got one for test: %s, error: %v", tc.name, err)
				}
			}

			// CBOR Validate
			err = decodedMetadata.Validate()
			if tc.expectValidationError {
				if err == nil {
					t.Errorf(
						"expected validation error but got none for test: %s",
						tc.name,
					)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect validation error but got one for test: %s, error: %v", tc.name, err)
				}
			}

			// Deep Equality Check for CBOR
			deepEqual := reflect.DeepEqual(decodedMetadata, tc.expectedObj)
			if tc.expectCBORDeepEqualError && deepEqual {
				t.Errorf(
					"Cbor expected deep equal error but objects are identical for test: %s",
					tc.name,
				)
			}
			if !tc.expectCBORDeepEqualError && !deepEqual {
				t.Errorf(
					"Cbor expected objects to be identical but they are not for test: %s, expected: %v, got: %v",
					tc.name,
					tc.expectedObj,
					decodedMetadata,
				)
			}

			// Reset the object for JSON testing
			decodedMetadata = Cip20Metadata{}
			// Decode the JSON string
			err = json.Unmarshal([]byte(tc.jsonData), &decodedMetadata)
			if err != nil {
				t.Errorf(
					"unexpected result unmarshaling JSON to Cip20Metadata for test %s, error: %v",
					tc.name,
					err,
				)
			}
			if tc.expectJSONUnmarshalError {
				if err == nil {
					t.Errorf(
						"expected JSON unmarshal error but got none for test: %s",
						tc.name,
					)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect JSON unmarshal error but got one for test: %s, error: %v", tc.name, err)
				}
			}

			// JSON Validate
			err = decodedMetadata.Validate()
			if tc.expectValidationError {
				if err == nil {
					t.Errorf(
						"expected validation error but got none for test: %s",
						tc.name,
					)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect validation error but got one for test: %s, error: %v", tc.name, err)
				}
			}

			// Deep Equality Check for JSON
			deepEqual = reflect.DeepEqual(decodedMetadata, tc.expectedObj)
			if tc.expectJSONDeepEqualError && deepEqual {
				t.Errorf(
					"Json expected deep equal error but objects are identical for test: %s",
					tc.name,
				)
			}
			if !tc.expectJSONDeepEqualError && !deepEqual {
				t.Errorf(
					"Json expected objects to be identical but they are not for test: %s, expected: %v, got: %v",
					tc.name,
					tc.expectedObj,
					decodedMetadata,
				)
			}
		})
	}
}

func TestNewCip20Metadata(t *testing.T) {
	t.Parallel()
	// Test case: Non-empty messages
	messages := []string{"message1", "message2"}
	expectedJSON := `{"674":{"msg":["message1","message2"]}}`

	metadata, err := NewCip20Metadata(messages)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	actualJSON, err := json.Marshal(metadata)
	if err != nil {
		t.Errorf("Failed to marshal metadata to JSON: %v", err)
	}

	if string(actualJSON) != expectedJSON {
		t.Errorf(
			"Expected JSON: %s, but got: %s",
			expectedJSON,
			string(actualJSON),
		)
	}

	// Test case: Invalid metadata
	messages = []string{}
	expectedErr := "Key: 'Cip20Metadata.Num674.Msg' Error:Field validation for 'Msg' failed on the 'gt' tag"

	_, err = NewCip20Metadata(messages)
	if err == nil {
		t.Errorf("Expected validation error, but got no error")
	} else {
		actualErr := err.Error()
		if actualErr != expectedErr {
			t.Errorf("Expected error: %s, but got: %s", expectedErr, actualErr)
		}
	}
}
