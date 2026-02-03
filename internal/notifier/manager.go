package notifier

import (
	"context"
	"sync"

	"stock-monitor/internal/model"
)

// Manager 通知管理器
type Manager struct {
	notifiers []Notifier
	mu        sync.RWMutex
}

// NewManager 创建通知管理器
func NewManager() *Manager {
	return &Manager{
		notifiers: make([]Notifier, 0),
	}
}

// Add 添加通知渠道
func (m *Manager) Add(n Notifier) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifiers = append(m.notifiers, n)
}

// Notify 发送通知到所有渠道
func (m *Manager) Notify(ctx context.Context, alert *model.Alert) error {
	m.mu.RLock()
	notifiers := make([]Notifier, len(m.notifiers))
	copy(notifiers, m.notifiers)
	m.mu.RUnlock()

	var wg sync.WaitGroup
	for _, n := range notifiers {
		wg.Add(1)
		go func(notifier Notifier) {
			defer wg.Done()
			_ = notifier.Send(ctx, alert)
		}(n)
	}
	wg.Wait()
	return nil
}
