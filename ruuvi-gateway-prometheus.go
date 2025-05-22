// Copyright (c) 2025, Juho Juopperi, Joonas Kuorilehto
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

// Command ruuvi-gateway-prometheus is a Prometheus exporter that requests
// Data from Ruuvi Gateway HTTP API and exposes it in Prometheus format.
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	sloghttp "github.com/samber/slog-http"
)

const (
	// commandName is the name of this command used in help texts.
	commandName = "ruuvi-gateway-prometheus"
)

func main() {
	cmdline := parseSettings()

	// Logging
	logLevel := slog.LevelInfo
	if cmdline.debug {
		logLevel = slog.LevelDebug
	}
	logJsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(logJsonHandler).With(slog.String("service", commandName))
	slog.SetDefault(logger)
	config := sloghttp.Config{
		WithSpanID:  true,
		WithTraceID: true,
	}

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "0.0.0.0"
	}

	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		listenPort = "9090"
	}

	addr := fmt.Sprintf("%s:%s", listenAddress, listenPort)
	loggingHandler := sloghttp.NewWithConfig(logger, config)(sloghttp.Recovery(Handler))
	server := http.Server{
		Addr:    addr,
		Handler: loggingHandler,
	}

	slog.Info("Starting Ruuvi Gateway Prometheus Exporter", "listen_address", listenAddress, "listen_port", listenPort)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("HTTP server ListenAndServe", "error", err)
	}
}
