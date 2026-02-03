package api

import (
	"encoding/json"
	"net/http"

	"stock-monitor/internal/rule"
	"stock-monitor/internal/storage"

	"github.com/google/uuid"
)

// Server API服务器
type Server struct {
	store *storage.Store
	mux   *http.ServeMux
}

// NewServer 创建API服务器
func NewServer(store *storage.Store) *Server {
	s := &Server{
		store: store,
		mux:   http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/stocks", s.handleStocks)
	s.mux.HandleFunc("/api/rules", s.handleRules)
	s.mux.HandleFunc("/api/rule-types", s.handleRuleTypes)
	s.mux.HandleFunc("/api/notifiers", s.handleNotifiers)
	s.mux.HandleFunc("/", s.handleIndex)
}

func (s *Server) json(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) handleStocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.json(w, s.store.GetStocks())
	case http.MethodPost:
		var stock storage.StockItem
		json.NewDecoder(r.Body).Decode(&stock)
		s.store.AddStock(stock)
		s.json(w, stock)
	case http.MethodDelete:
		code := r.URL.Query().Get("code")
		s.store.DeleteStock(code)
		s.json(w, map[string]bool{"ok": true})
	}
}

func (s *Server) handleRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.json(w, s.store.GetRules())
	case http.MethodPost:
		var rule storage.RuleItem
		json.NewDecoder(r.Body).Decode(&rule)
		rule.ID = uuid.New().String()
		s.store.AddRule(rule)
		s.json(w, rule)
	case http.MethodPut:
		var rule storage.RuleItem
		json.NewDecoder(r.Body).Decode(&rule)
		s.store.UpdateRule(rule)
		s.json(w, rule)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		s.store.DeleteRule(id)
		s.json(w, map[string]bool{"ok": true})
	}
}

func (s *Server) handleNotifiers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.json(w, s.store.GetNotifiers())
	case http.MethodPut:
		var cfg storage.NotifierConfig
		json.NewDecoder(r.Body).Decode(&cfg)
		s.store.UpdateNotifiers(cfg)
		s.json(w, cfg)
	}
}

func (s *Server) handleRuleTypes(w http.ResponseWriter, r *http.Request) {
	s.json(w, rule.GlobalRegistry.Types())
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}
