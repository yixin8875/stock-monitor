package rule

import (
	"context"

	"stock-monitor/internal/model"
)

// RuleContext 规则执行上下文
type RuleContext struct {
	Stock  *model.Stock
	KLines *model.KLineData
}

// RuleResult 规则执行结果
type RuleResult struct {
	Triggered bool
	RuleName  string
	Message   string
	Level     model.AlertLevel
	Extra     map[string]interface{}
}

// Rule 规则接口
type Rule interface {
	Name() string
	Description() string
	Evaluate(ctx context.Context, ruleCtx *RuleContext) (*RuleResult, error)
	Validate() error
}

// KLineRule 需要K线数据的规则接口
type KLineRule interface {
	Rule
	KLineType() model.KLineType
	StockCode() string
}
