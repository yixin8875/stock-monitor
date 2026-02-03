package rule

import (
	"context"
	"sync"
	"time"

	"stock-monitor/internal/model"

	"github.com/google/uuid"
)

// Engine 规则引擎
type Engine struct {
	rules []Rule
	mu    sync.RWMutex
}

// NewEngine 创建规则引擎
func NewEngine() *Engine {
	return &Engine{
		rules: make([]Rule, 0),
	}
}

// AddRule 添加规则
func (e *Engine) AddRule(rule Rule) error {
	if err := rule.Validate(); err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = append(e.rules, rule)
	return nil
}

// Evaluate 评估所有规则
func (e *Engine) Evaluate(ctx context.Context, ruleCtx *RuleContext) ([]*model.Alert, error) {
	e.mu.RLock()
	rules := make([]Rule, len(e.rules))
	copy(rules, e.rules)
	e.mu.RUnlock()

	var alerts []*model.Alert
	for _, rule := range rules {
		result, err := rule.Evaluate(ctx, ruleCtx)
		if err != nil {
			continue
		}
		if result.Triggered {
			alert := &model.Alert{
				ID:        uuid.New().String(),
				StockCode: ruleCtx.Stock.Code,
				StockName: ruleCtx.Stock.Name,
				RuleName:  result.RuleName,
				Level:     result.Level,
				Message:   result.Message,
				Price:     ruleCtx.Stock.Price,
				Time:      time.Now(),
				Extra:     result.Extra,
			}
			alerts = append(alerts, alert)
		}
	}
	return alerts, nil
}
