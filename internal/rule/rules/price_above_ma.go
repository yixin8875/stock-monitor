package rules

import (
	"context"
	"fmt"

	"stock-monitor/internal/indicator"
	"stock-monitor/internal/model"
	"stock-monitor/internal/rule"
)

func init() {
	rule.GlobalRegistry.Register("price_above_ma", NewPriceAboveMARule, "突破均线")
}

// PriceAboveMARule 价格突破均线规则
type PriceAboveMARule struct {
	name      string
	period    int
	stockCode string
	klineType model.KLineType
	level     model.AlertLevel
}

// NewPriceAboveMARule 创建规则
func NewPriceAboveMARule(name string, level model.AlertLevel, params map[string]interface{}) (rule.Rule, error) {
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

	return &PriceAboveMARule{
		name:      name,
		period:    period,
		stockCode: stockCode,
		klineType: klineType,
		level:     level,
	}, nil
}

func (r *PriceAboveMARule) Name() string {
	return r.name
}

func (r *PriceAboveMARule) Description() string {
	return fmt.Sprintf("%s K线 MA%d 突破", r.klineType, r.period)
}

func (r *PriceAboveMARule) KLineType() model.KLineType {
	return r.klineType
}

func (r *PriceAboveMARule) StockCode() string {
	return r.stockCode
}

func (r *PriceAboveMARule) Validate() error {
	if r.period <= 0 {
		return fmt.Errorf("period must be positive")
	}
	return nil
}

func (r *PriceAboveMARule) Evaluate(ctx context.Context, ruleCtx *rule.RuleContext) (*rule.RuleResult, error) {
	// 检查是否匹配股票代码
	if r.stockCode != "" && ruleCtx.Stock.Code != r.stockCode {
		return &rule.RuleResult{Triggered: false}, nil
	}

	// 检查K线数据
	if ruleCtx.KLines == nil || len(ruleCtx.KLines.Lines) < r.period {
		return &rule.RuleResult{Triggered: false}, nil
	}

	// 提取收盘价
	closes := make([]float64, len(ruleCtx.KLines.Lines))
	for i, kline := range ruleCtx.KLines.Lines {
		closes[i] = kline.Close
	}

	// 计算MA
	maValue := indicator.LastMA(closes, r.period)
	currentClose := ruleCtx.Stock.Close

	// 判断是否突破
	if currentClose > maValue {
		return &rule.RuleResult{
			Triggered: true,
			RuleName:  r.name,
			Level:     r.level,
			Message: fmt.Sprintf("%s 收盘价 %.2f 突破 MA%d (%.2f)",
				ruleCtx.Stock.Name, currentClose, r.period, maValue),
			Extra: map[string]interface{}{
				"ma_value": maValue,
				"period":   r.period,
			},
		}, nil
	}

	return &rule.RuleResult{Triggered: false}, nil
}
