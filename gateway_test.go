// Copyright (c) 2025, Juho Juopperi
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	_ "embed"
	"encoding/json"
	"testing"
)

//go:embed testdata/ruuvigw_history.json
var ruuvigwHistoryTestData []byte

func TestDecode(t *testing.T) {
	var d HistoryResponse
	err := json.Unmarshal(ruuvigwHistoryTestData, &d)

	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if d.Data.Timestamp != 1747751898 {
		t.Errorf("Expected timestamp 1747751898, got %d", d.Data.Timestamp)
	}

	if d.Data.Tags == nil {
		t.Error("Expected Tags to be non-nil")
	}

	if len(*d.Data.Tags) != 4 {
		t.Errorf("Expected 4 tags, got %d", len(*d.Data.Tags))
	}

	if tag, ok := (*d.Data.Tags)["D8:20:8F:5F:AB:4D"]; !ok {
		t.Error("Expected tag D8:20:8F:5F:AB:4D to be present")
	} else {
		if tag.Rssi != -79 {
			t.Errorf("Expected Rssi -79, got %d", tag.Rssi)
		}
		if tag.Temperature != 23.475 {
			t.Errorf("Expected Temperature 23.475, got %f", tag.Temperature)
		}
		if tag.Humidity != 32.1975 {
			t.Errorf("Expected Humidity 32.1975, got %f", tag.Humidity)
		}
	}
}
