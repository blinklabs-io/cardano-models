// Copyright 2026 Blink Labs Software
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
	"testing"

	models "github.com/blinklabs-io/cardano-models"
	"github.com/blinklabs-io/gouroboros/cbor"
	"github.com/stretchr/testify/require"
)

func TestTunaV1State_DecodeConstructorValidation(t *testing.T) {
	valid := models.TunaV1State{
		BlockNumber: 1, CurrentHash: []byte{2}, LeadingZeros: 3,
		DifficultyNumber: 4, EpochTime: 5, RealTimeNow: 6,
		Extra: "extra", Interlink: [][]byte{{7}},
	}
	data, err := cbor.Encode(valid)
	require.NoError(t, err)
	var decoded models.TunaV1State
	_, err = cbor.Decode(data, &decoded)
	require.NoError(t, err)
	require.Equal(t, valid, decoded)

	for _, tag := range []uint{1, 2} {
		bad, err := cbor.Encode(cbor.NewConstructorEncoder(tag, cbor.IndefLengthList{
			valid.BlockNumber, valid.CurrentHash, valid.LeadingZeros,
			valid.DifficultyNumber, valid.EpochTime, valid.RealTimeNow,
			valid.Extra, cbor.IndefLengthList{[]byte{7}},
		}))
		require.NoError(t, err)
		want := models.TunaV1State{
			BlockNumber: 99, CurrentHash: []byte{98}, LeadingZeros: 97,
			DifficultyNumber: 96, EpochTime: 95, RealTimeNow: 94,
			Extra: "sentinel", Interlink: [][]byte{{93}},
		}
		got := want
		_, err = cbor.Decode(bad, &got)
		require.Error(t, err)
		require.Equal(t, want, got)
	}
	bad, err := cbor.Encode(cbor.NewConstructorEncoder(0, cbor.IndefLengthList{}))
	require.NoError(t, err)
	want := models.TunaV1State{BlockNumber: 99, CurrentHash: []byte{98}}
	got := want
	_, err = cbor.Decode(bad, &got)
	require.Error(t, err)
	require.Equal(t, want, got)
}

func TestTunaV2State_DecodeConstructorValidation(t *testing.T) {
	valid := models.TunaV2State{
		BlockNumber: 1, CurrentHash: []byte{2}, LeadingZeros: 3,
		DifficultyNumber: 4, EpochTime: 5, CurrentPosixTime: 6,
		MerkleRoot: []byte{7},
	}
	data, err := cbor.Encode(valid)
	require.NoError(t, err)
	var decoded models.TunaV2State
	_, err = cbor.Decode(data, &decoded)
	require.NoError(t, err)
	require.Equal(t, valid, decoded)

	for _, tag := range []uint{1, 2} {
		bad, err := cbor.Encode(cbor.NewConstructorEncoder(tag, cbor.IndefLengthList{
			valid.BlockNumber, valid.CurrentHash, valid.LeadingZeros,
			valid.DifficultyNumber, valid.EpochTime, valid.CurrentPosixTime,
			valid.MerkleRoot,
		}))
		require.NoError(t, err)
		want := models.TunaV2State{
			BlockNumber: 99, CurrentHash: []byte{98}, LeadingZeros: 97,
			DifficultyNumber: 96, EpochTime: 95, CurrentPosixTime: 94,
			MerkleRoot: []byte{93},
		}
		got := want
		_, err = cbor.Decode(bad, &got)
		require.Error(t, err)
		require.Equal(t, want, got)
	}
	bad, err := cbor.Encode(cbor.NewConstructorEncoder(0, cbor.IndefLengthList{}))
	require.NoError(t, err)
	want := models.TunaV2State{BlockNumber: 99, CurrentHash: []byte{98}}
	got := want
	_, err = cbor.Decode(bad, &got)
	require.Error(t, err)
	require.Equal(t, want, got)
}
