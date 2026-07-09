package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"compatgate/internal/domain"
	"compatgate/internal/policy"
	"compatgate/internal/compat"
)

type Server struct {
	evaluator *compat.Evaluator
}

func NewServer(catalog policy.Catalog) *Server {
	return &Server{evaluator: compat.NewEvaluator(catalog)}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/v1/compatibility/assess", s.handleAssess)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleAssess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req domain.AssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json payload"})
		return
	}
	resp, err := s.evaluator.Assess(req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrBadRequest) {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
