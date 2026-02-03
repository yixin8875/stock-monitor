package datasource

import (
	"context"

	"stock-monitor/internal/model"
)

// DataSource 数据源接口
type DataSource interface {
	// Name 数据源名称
	Name() string

	// GetRealTimeQuote 获取实时行情
	GetRealTimeQuote(ctx context.Context, codes []string) ([]*model.Stock, error)

	// GetKLine 获取K线数据
	GetKLine(ctx context.Context, code string, ktype model.KLineType, count int) (*model.KLineData, error)
}
