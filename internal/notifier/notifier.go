package notifier

import (
	"context"

	"stock-monitor/internal/model"
)

// Notifier 通知接口
type Notifier interface {
	Name() string
	Send(ctx context.Context, alert *model.Alert) error
}
