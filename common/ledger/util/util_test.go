/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

func TestBasicEncodingDecoding(t *testing.T) {
	for i := range 10000 {
		value := EncodeOrderPreservingVarUint64(uint64(i))
		nextValue := EncodeOrderPreservingVarUint64(uint64(i + 1))
		if !(bytes.Compare(value, nextValue) < 0) {
			t.Fatalf("A smaller integer should result into smaller bytes. Encoded bytes for [%d] is [%x] and for [%d] is [%x]",
				i, i+1, value, nextValue)
		}
		decodedValue, _, err := DecodeOrderPreservingVarUint64(value)
		require.NoError(t, err, "Error via calling DecodeOrderPreservingVarUint64")
		if decodedValue != uint64(i) {
			t.Fatalf("Value not same after decoding. Original value = [%d], decode value = [%d]", i, decodedValue)
		}
	}
}

func TestDecodingAppendedValues(t *testing.T) {
	appendedValues := []byte{}
	for i := range 1000 {
		appendedValues = append(appendedValues, EncodeOrderPreservingVarUint64(uint64(i))...)
	}

	l := 0
	value := uint64(0)
	var err error
	for i := range 1000 {
		appendedValues = appendedValues[l:]
		value, l, err = DecodeOrderPreservingVarUint64(appendedValues)
		require.NoError(t, err, "Error via calling DecodeOrderPreservingVarUint64")
		if value != uint64(i) {
			t.Fatalf("expected value = [%d], decode value = [%d]", i, value)
		}
	}
}

func TestDecodingBadInputBytes(t *testing.T) {
	// error case when num consumed bytes > 1
	sizeBytes := protowire.AppendVarint(nil, uint64(1000))
	_, _, err := DecodeOrderPreservingVarUint64(sizeBytes)
	require.Equal(t, fmt.Sprintf("number of consumed bytes from DecodeVarint is invalid, expected 1, but got %d", len(sizeBytes)), err.Error())

	// error case when decoding invalid bytes - trim off last byte
	invalidSizeBytes := sizeBytes[0 : len(sizeBytes)-1]
	_, _, err = DecodeOrderPreservingVarUint64(invalidSizeBytes)
	require.Equal(t, "number of consumed bytes from DecodeVarint is invalid, expected 1, but got 0", err.Error())

	// error case when size is more than available bytes
	inputBytes := protowire.AppendVarint(nil, uint64(8))
	_, _, err = DecodeOrderPreservingVarUint64(inputBytes)
	require.Equal(t, "decoded size (8) from DecodeVarint is more than available bytes (0)", err.Error())

	// error case when size is greater than 8
	bigSizeBytes := protowire.AppendVarint(nil, uint64(12))
	_, _, err = DecodeOrderPreservingVarUint64(bigSizeBytes)
	require.Equal(t, "decoded size from DecodeVarint is invalid, expected <=8, but got 12", err.Error())
}
