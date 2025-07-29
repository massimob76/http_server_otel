package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var identity int = getIdentity()

func server() {
	_, err := setupOtel()
	if err != nil {
		log.Error("Couldn't setup OTEL", "error", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("GET /countdown/{counter}", otelhttp.NewHandler(countDownHandler(), "COUNTDOWN ENDPOINT"))

	port := serverPort(identity)
	log.Info("Starting server with identity", "identity", identity, "port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Error("Error starting server", "port", port, "error", err.Error())
		return
	}

}

func serverPort(identity int) int {
	return 3000 + identity
}

func countDownHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentCounter, err := strconv.Atoi(r.PathValue("counter"))
		if err != nil {
			log.Error("Cannot convert counter to integer", "counter", r.PathValue("counter"), "error", err.Error())
		}

		log.Info("Server received countdown call", "identity", identity, "countdown", currentCounter)

		if currentCounter == 0 {
			log.Info("Reached the end of the countdown")
			return
		}

		nextServerPort := serverPort(identity + 1)
		nextCountdown := currentCounter - 1
		nextURL := fmt.Sprintf("http://localhost:%d/countdown/%d", nextServerPort, nextCountdown)
		log.Info("Trying to contact server", "URL", nextURL)
		// ctx := context.Background()
		ctx := r.Context()

		req, err := http.NewRequestWithContext(ctx, "GET", nextURL, nil)
		if err != nil {
			log.Error("Error creating request", "URL", nextURL, "error", err.Error())
			return
		}

		// client := http.DefaultClient
		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
		resp, err := client.Do(req)
		if err != nil {
			log.Error("Error contacting next server", "error", err.Error())
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("Error reading response body", "error", err.Error())
			return
		}
		log.Info("Returned payload", "body", body)
	})
}
