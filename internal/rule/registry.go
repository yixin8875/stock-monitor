package rule

import (
	"fmt"
	"sync"

	"stock-monitor/internal/model"
)

// RuleFactory 规则工厂函数
type RuleFactory func(name string, level model.AlertLevel, params map[string]interface{}) (Rule, error)

// RuleTypeInfo 规则类型信息
type RuleTypeInfo struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Registry 规则注册表
type Registry struct {
	mu        sync.RWMutex
	factories map[string]RuleFactory
	typeNames map[string]string
}

// GlobalRegistry 全局规则注册表
var GlobalRegistry = NewRegistry()

// NewRegistry 创建规则注册表
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]RuleFactory),
		typeNames: make(map[string]string),
	}
}

// Register 注册规则类型
func (r *Registry) Register(ruleType string, factory RuleFactory, displayName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[ruleType] = factory
	r.typeNames[ruleType] = displayName
}

// Types 获取所有规则类型
func (r *Registry) Types() []RuleTypeInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	types := make([]RuleTypeInfo, 0, len(r.typeNames))
	for t, name := range r.typeNames {
		types = append(types, RuleTypeInfo{Type: t, Name: name})
	}
	return types
}

// Create 创建规则实例
func (r *Registry) Create(ruleType, name string, level model.AlertLevel, params map[string]interface{}) (Rule, error) {
	r.mu.RLock()
	factory, ok := r.factories[ruleType]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown rule type: %s", ruleType)
	}
	return factory(name, level, params)
}
