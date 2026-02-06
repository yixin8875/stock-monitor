package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"stock-monitor/internal/rule"
	"stock-monitor/internal/storage"

	"github.com/google/uuid"
)

var stockCodeRegexp = regexp.MustCompile(`^\d{6}$`)

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

func (s *Server) errJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (s *Server) handleStocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.json(w, s.store.GetStocks())
	case http.MethodPost:
		var stock storage.StockItem
		if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
			s.errJSON(w, http.StatusBadRequest, "请求格式错误")
			return
		}
		if !stockCodeRegexp.MatchString(stock.Code) {
			s.errJSON(w, http.StatusBadRequest, "股票代码必须为6位数字")
			return
		}
		if stock.Name == "" {
			s.errJSON(w, http.StatusBadRequest, "股票名称不能为空")
			return
		}
		for _, existing := range s.store.GetStocks() {
			if existing.Code == stock.Code {
				s.errJSON(w, http.StatusConflict, "股票已存在")
				return
			}
		}
		if err := s.store.AddStock(stock); err != nil {
			s.errJSON(w, http.StatusInternalServerError, "保存失败")
			return
		}
		s.json(w, stock)
	case http.MethodDelete:
		code := r.URL.Query().Get("code")
		if code == "" {
			s.errJSON(w, http.StatusBadRequest, "缺少code参数")
			return
		}
		s.store.DeleteStock(code)
		s.json(w, map[string]bool{"ok": true})
	}
}

func (s *Server) handleRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.json(w, s.store.GetRules())
	case http.MethodPost:
		var ri storage.RuleItem
		if err := json.NewDecoder(r.Body).Decode(&ri); err != nil {
			s.errJSON(w, http.StatusBadRequest, "请求格式错误")
			return
		}
		if err := s.validateRule(ri); err != nil {
			s.errJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		ri.ID = uuid.New().String()
		if err := s.store.AddRule(ri); err != nil {
			s.errJSON(w, http.StatusInternalServerError, "保存失败")
			return
		}
		s.json(w, ri)
	case http.MethodPut:
		var ri storage.RuleItem
		if err := json.NewDecoder(r.Body).Decode(&ri); err != nil {
			s.errJSON(w, http.StatusBadRequest, "请求格式错误")
			return
		}
		if ri.ID == "" {
			s.errJSON(w, http.StatusBadRequest, "缺少规则ID")
			return
		}
		if err := s.validateRule(ri); err != nil {
			s.errJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := s.store.UpdateRule(ri); err != nil {
			s.errJSON(w, http.StatusInternalServerError, "保存失败")
			return
		}
		s.json(w, ri)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			s.errJSON(w, http.StatusBadRequest, "缺少id参数")
			return
		}
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
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			s.errJSON(w, http.StatusBadRequest, "请求格式错误")
			return
		}
		if err := s.store.UpdateNotifiers(cfg); err != nil {
			s.errJSON(w, http.StatusInternalServerError, "保存失败")
			return
		}
		s.json(w, cfg)
	}
}

func (s *Server) handleRuleTypes(w http.ResponseWriter, r *http.Request) {
	s.json(w, rule.GlobalRegistry.Types())
}

var validLevels = map[string]bool{"info": true, "warning": true, "critical": true}

var validKLineTypes = map[string]bool{
	"5min": true, "15min": true, "30min": true, "60min": true,
	"daily": true, "weekly": true, "monthly": true,
}

func (s *Server) validateRule(ri storage.RuleItem) error {
	if ri.Name == "" {
		return fmt.Errorf("规则名称不能为空")
	}
	if ri.Type == "" {
		return fmt.Errorf("规则类型不能为空")
	}
	if !validLevels[ri.Level] {
		return fmt.Errorf("无效的告警级别: %s", ri.Level)
	}
	if ri.StockCode != "" && !stockCodeRegexp.MatchString(ri.StockCode) {
		return fmt.Errorf("股票代码必须为6位数字")
	}
	if ri.KLineType != "" && !validKLineTypes[ri.KLineType] {
		return fmt.Errorf("无效的K线类型: %s", ri.KLineType)
	}
	if ri.Period < 0 {
		return fmt.Errorf("周期不能为负数")
	}
	return nil
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}
