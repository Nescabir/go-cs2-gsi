package cs2gsi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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
	if r.Method != http.MethodPost {
		gsi.logger.Warn("invalid HTTP method", "method", r.Method, "remote_addr", r.RemoteAddr)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "application/json") {
		gsi.logger.Warn("invalid content type", "content_type", contentType, "remote_addr", r.RemoteAddr)
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return nil
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		gsi.logger.Error("failed to read request body", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return fmt.Errorf("failed to read request body: %w", err)
	}

	if err := gsi.Digest(body); err != nil {
		if IsInvalidToken(err) {
			gsi.logger.Warn("invalid auth token", "remote_addr", r.RemoteAddr)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return nil
		}
		gsi.logger.Error("failed to process game state", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "Invalid game state data", http.StatusBadRequest)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return nil
}
