package storage

import (
	"encoding/json"
	"os"
	"sync"
)

// Store JSON文件存储
type Store struct {
	path string
	data *Data
	mu   sync.RWMutex
}

// NewStore 创建存储
func NewStore(path string) *Store {
	return &Store{
		path: path,
		data: &Data{
			Stocks: []StockItem{},
			Rules:  []RuleItem{},
		},
	}
}

// Load 加载数据
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return s.saveUnsafe()
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s.data)
}

// Save 保存数据
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveUnsafe()
}

func (s *Store) saveUnsafe() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// GetStocks 获取股票列表
func (s *Store) GetStocks() []StockItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]StockItem, len(s.data.Stocks))
	copy(result, s.data.Stocks)
	return result
}

// AddStock 添加股票
func (s *Store) AddStock(stock StockItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Stocks = append(s.data.Stocks, stock)
	return s.saveUnsafe()
}

// DeleteStock 删除股票
func (s *Store) DeleteStock(code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, st := range s.data.Stocks {
		if st.Code == code {
			s.data.Stocks = append(s.data.Stocks[:i], s.data.Stocks[i+1:]...)
			return s.saveUnsafe()
		}
	}
	return nil
}

// GetRules 获取规则列表
func (s *Store) GetRules() []RuleItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]RuleItem, len(s.data.Rules))
	copy(result, s.data.Rules)
	return result
}

// AddRule 添加规则
func (s *Store) AddRule(rule RuleItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Rules = append(s.data.Rules, rule)
	return s.saveUnsafe()
}

// UpdateRule 更新规则
func (s *Store) UpdateRule(rule RuleItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.data.Rules {
		if r.ID == rule.ID {
			s.data.Rules[i] = rule
			return s.saveUnsafe()
		}
	}
	return nil
}

// DeleteRule 删除规则
func (s *Store) DeleteRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.data.Rules {
		if r.ID == id {
			s.data.Rules = append(s.data.Rules[:i], s.data.Rules[i+1:]...)
			return s.saveUnsafe()
		}
	}
	return nil
}

// GetNotifiers 获取通知配置
func (s *Store) GetNotifiers() NotifierConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.Notifiers
}

// UpdateNotifiers 更新通知配置
func (s *Store) UpdateNotifiers(cfg NotifierConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Notifiers = cfg
	return s.saveUnsafe()
}
