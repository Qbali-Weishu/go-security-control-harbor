package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"compatgate/internal/domain"
	"compatgate/internal/policy"
	"compatgate/internal/compat"
)

// Server 是 HTTP API 服务器
type Server struct {
	evaluator *compat.Evaluator
}

// NewServer 创建新的服务器实例
func NewServer(catalog policy.Catalog) *Server {
	return &Server{evaluator: compat.NewEvaluator(catalog)}
}

// Routes 返回配置好的 HTTP 路由处理器
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/v1/compatibility/assess", s.handleAssess)
	return mux
}

// handleHealth 处理健康检查请求
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleAssess 处理兼容性评估请求
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

// writeJSON 将数据序列化为 JSON 并写入响应
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
