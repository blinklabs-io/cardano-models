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

// TunaV1State represents the datum format used by the $TUNA mining smart contract (v1)
type TunaV1State struct {
	BlockNumber      int64
	TargetHash       []byte
	LeadingZeros     int64
	DifficultyNumber int64
	EpochTime        int64
	RealTimeNow      int64
	Message          []byte
	Interlink        [][]byte
}
