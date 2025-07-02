package cs2gsi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	rawModels "github.com/nescabir/go-cs2-gsi/raw"
)

// Listen starts the HTTP server and handles incoming game state requests
func (gsi *CS2GSI) Listen() error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		if err := gsi.handleGameStateRequest(w, r); err != nil {
			gsi.logger.Error("failed to handle game state request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})

	gsi.logger.Info("starting CS2 GSI server", "address", gsi.config.ServerAddr)
	return http.ListenAndServe(gsi.config.ServerAddr, mux)
}

// handleGameStateRequest processes incoming game state requests from CS2
func (gsi *CS2GSI) handleGameStateRequest(w http.ResponseWriter, r *http.Request) error {
	// Validate HTTP method
	if r.Method != http.MethodPost {
		gsi.logger.Warn("invalid HTTP method", "method", r.Method, "remote_addr", r.RemoteAddr)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	// Validate content type
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "application/json") {
		gsi.logger.Warn("invalid content type", "content_type", contentType, "remote_addr", r.RemoteAddr)
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return nil
	}

	// Limit request body size (1MB should be sufficient for GSI data)
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	// Parse request body
	stateRaw := &rawModels.State{}
	if err := json.NewDecoder(r.Body).Decode(stateRaw); err != nil {
		gsi.logger.Error("failed to decode JSON", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Validate request data
	if err := gsi.digest(stateRaw); err != nil {
		gsi.logger.Error("failed to process game state", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid game state data", http.StatusBadRequest)
		return nil
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return nil
}
