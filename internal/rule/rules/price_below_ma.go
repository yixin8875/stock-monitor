package rules

import (
	"context"
	"fmt"

	"stock-monitor/internal/indicator"
	"stock-monitor/internal/model"
	"stock-monitor/internal/rule"
)

func init() {
	rule.GlobalRegistry.Register("price_below_ma", NewPriceBelowMARule, "跌破均线")
}

// PriceBelowMARule 价格跌破均线规则
type PriceBelowMARule struct {
	name      string
	period    int
	stockCode string
	klineType model.KLineType
	level     model.AlertLevel
}

// NewPriceBelowMARule 创建规则
func NewPriceBelowMARule(name string, level model.AlertLevel, params map[string]interface{}) (rule.Rule, error) {
	period := 60
	if p, ok := params["period"].(int); ok {
		period = p
	} else if p, ok := params["period"].(float64); ok {
		period = int(p)
	}

	stockCode, _ := params["stock_code"].(string)

	klineType := model.KLineDaily
	if kt, ok := params["kline_type"].(string); ok && kt != "" {
		klineType = model.KLineType(kt)
	}

	return &PriceBelowMARule{
		name:      name,
		period:    period,
		stockCode: stockCode,
		klineType: klineType,
		level:     level,
	}, nil
}

func (r *PriceBelowMARule) Name() string        { return r.name }
func (r *PriceBelowMARule) StockCode() string   { return r.stockCode }
func (r *PriceBelowMARule) KLineType() model.KLineType { return r.klineType }
func (r *PriceBelowMARule) Validate() error {
	if r.period <= 0 {
		return fmt.Errorf("period must be positive")
	}
	return nil
}

func (r *PriceBelowMARule) Description() string {
	return fmt.Sprintf("%s K线 MA%d 跌破", r.klineType, r.period)
}

func (r *PriceBelowMARule) Evaluate(ctx context.Context, ruleCtx *rule.RuleContext) (*rule.RuleResult, error) {
	if r.stockCode != "" && ruleCtx.Stock.Code != r.stockCode {
		return &rule.RuleResult{Triggered: false}, nil
	}

	if ruleCtx.KLines == nil || len(ruleCtx.KLines.Lines) < r.period {
		return &rule.RuleResult{Triggered: false}, nil
	}

	closes := make([]float64, len(ruleCtx.KLines.Lines))
	for i, kline := range ruleCtx.KLines.Lines {
		closes[i] = kline.Close
	}

	maValue := indicator.LastMA(closes, r.period)
	currentClose := ruleCtx.Stock.Close

	if currentClose < maValue {
		return &rule.RuleResult{
			Triggered: true,
			RuleName:  r.name,
			Level:     r.level,
			Message: fmt.Sprintf("%s 收盘价 %.2f 跌破 MA%d (%.2f)",
				ruleCtx.Stock.Name, currentClose, r.period, maValue),
			Extra: map[string]interface{}{
				"ma_value": maValue,
				"period":   r.period,
			},
		}, nil
	}

	return &rule.RuleResult{Triggered: false}, nil
}
