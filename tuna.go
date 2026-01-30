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
	"github.com/blinklabs-io/gouroboros/cbor"
)

// TunaV1State represents the datum format used by the $TUNA mining smart contract (v1)
type TunaV1State struct {
	// This allows the type to be reused during decoding
	cbor.StructAsArray
	BlockNumber      int64
	CurrentHash      []byte
	LeadingZeros     int64
	DifficultyNumber int64
	EpochTime        int64
	RealTimeNow      int64
	Extra            any
	Interlink        [][]byte
}

func (t *TunaV1State) MarshalCBOR() ([]byte, error) {
	tmpInterlink := make([]any, 0, len(t.Interlink))
	for _, item := range t.Interlink {
		tmpInterlink = append(tmpInterlink, item)
	}
	tmp := cbor.NewConstructor(
		0,
		cbor.IndefLengthList{
			t.BlockNumber,
			t.CurrentHash,
			t.LeadingZeros,
			t.DifficultyNumber,
			t.EpochTime,
			t.RealTimeNow,
			t.Extra,
			cbor.IndefLengthList(tmpInterlink),
		},
	)
	return cbor.Encode(&tmp)
}

func (t *TunaV1State) UnmarshalCBOR(cborData []byte) error {
	var tmpConstr cbor.Constructor
	if _, err := cbor.Decode(cborData, &tmpConstr); err != nil {
		return err
	}
	type tTunaV1State TunaV1State
	var tmpTunaV1State tTunaV1State
	if _, err := cbor.Decode(tmpConstr.FieldsCbor(), &tmpTunaV1State); err != nil {
		return err
	}
	*t = TunaV1State(tmpTunaV1State)
	return nil
}

// TunaV2State represents the datum format used by the $TUNA mining smart contract (v2)
type TunaV2State struct {
	// This allows the type to be reused during decoding
	cbor.StructAsArray
	BlockNumber      int64
	CurrentHash      []byte
	LeadingZeros     int64
	DifficultyNumber int64
	EpochTime        int64
	CurrentPosixTime int64
	MerkleRoot       []byte
}

func (t *TunaV2State) MarshalCBOR() ([]byte, error) {
	tmp := cbor.NewConstructor(
		0,
		cbor.IndefLengthList{
			t.BlockNumber,
			t.CurrentHash,
			t.LeadingZeros,
			t.DifficultyNumber,
			t.EpochTime,
			t.CurrentPosixTime,
			t.MerkleRoot,
		},
	)
	return cbor.Encode(&tmp)
}

func (t *TunaV2State) UnmarshalCBOR(cborData []byte) error {
	var tmpConstr cbor.Constructor
	if _, err := cbor.Decode(cborData, &tmpConstr); err != nil {
		return err
	}
	type tTunaV2State TunaV2State
	var tmpTunaV2State tTunaV2State
	if _, err := cbor.Decode(tmpConstr.FieldsCbor(), &tmpTunaV2State); err != nil {
		return err
	}
	*t = TunaV2State(tmpTunaV2State)
	return nil
}
