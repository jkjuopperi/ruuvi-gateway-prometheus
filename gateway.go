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
	//"encoding/json"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type TagResponse struct {
	Rssi                      int     `json:"rssi"`
	Timestamp                 int     `json:"timestamp"`
	Data                      string  `json:"data"`
	DataFormat                int     `json:"dataFormat,omitempty"`
	Temperature               float64 `json:"temperature,omitempty"`
	Humidity                  float64 `json:"humidity,omitempty"`
	Pressure                  int     `json:"pressure,omitempty"`
	AccelX                    float64 `json:"accelX,omitempty"`
	AccelY                    float64 `json:"accelY,omitempty"`
	AccelZ                    float64 `json:"accelZ,omitempty"`
	MovementCounter           int     `json:"movementCounter,omitempty"`
	Voltage                   float64 `json:"voltage,omitempty"`
	TxPower                   int     `json:"txPower,omitempty"`
	MeasurementSequenceNumber int     `json:"measurementSequenceNumber,omitempty"`
	Id                        string  `json:"id,omitempty"`
}

type TagsResponse map[string]TagResponse

type HistoryDataResponse struct {
	Coordinates string        `json:"coordinates"`
	Timestamp   int           `json:"timestamp"`
	GwMac       string        `json:"gw_mac"`
	Tags        *TagsResponse `json:"tags"`
}

type HistoryResponse struct {
	Data HistoryDataResponse `json:"data"`
}

func ObserveRuuvi(o TagResponse) {
	addr := o.Id
	ruuviFrames.WithLabelValues(addr).Inc()
	signalRSSI.WithLabelValues(addr).Set(float64(o.Rssi))
	voltage.WithLabelValues(addr).Set(float64(o.Voltage))
	pressure.WithLabelValues(addr).Set(float64(o.Pressure) / 100)
	temperature.WithLabelValues(addr).Set(float64(o.Temperature))
	humidity.WithLabelValues(addr).Set(float64(o.Humidity) / 100)
	acceleration.WithLabelValues(addr, "X").Set(float64(o.AccelX))
	acceleration.WithLabelValues(addr, "Y").Set(float64(o.AccelY))
	acceleration.WithLabelValues(addr, "Z").Set(float64(o.AccelZ))
	format.WithLabelValues(addr).Set(float64(o.DataFormat))
	txPower.WithLabelValues(addr).Set(float64(o.TxPower))
	moveCount.WithLabelValues(addr).Set(float64(o.MovementCounter))
	seqno.WithLabelValues(addr).Set(float64(o.MeasurementSequenceNumber))
}

func UpdateMetrics() error {
	ResetMetrics()

	gatewayUrl := os.Getenv("GATEWAY_URL")
	req, err := http.NewRequest("GET", gatewayUrl, nil)
	if err != nil {
		return errors.New("failed to create http request for gateway")
	}

	gatewayApiKey := os.Getenv("GATEWAY_APIKEY")
	req.Header.Add("Authorization", "Bearer "+gatewayApiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("failed to send request to gateway")
	}

	defer resp.Body.Close()

	var d HistoryResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&d); err != nil {
		return errors.New("failed to decode json response from gateway")
	}

	if d.Data.Tags == nil {
		return errors.New("no tags found in response from gateway")
	}

	for _, tag := range *d.Data.Tags {
		ObserveRuuvi(tag)
	}

	return nil
}
