package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"stock-monitor/internal/model"
)

// ServerChan Server酱通知
type ServerChan struct {
	sendKey string
	client  *http.Client
}

// NewServerChan 创建Server酱通知
func NewServerChan(sendKey string) *ServerChan {
	return &ServerChan{
		sendKey: sendKey,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *ServerChan) Name() string {
	return "serverchan"
}

func (s *ServerChan) Send(ctx context.Context, alert *model.Alert) error {
	url := fmt.Sprintf("https://sctapi.ftqq.com/%s.send", s.sendKey)

	title := fmt.Sprintf("[%s] %s", alert.Level, alert.RuleName)
	desp := fmt.Sprintf("**股票**: %s (%s)\n\n**价格**: %.2f\n\n**消息**: %s\n\n**时间**: %s",
		alert.StockName, alert.StockCode, alert.Price, alert.Message,
		alert.Time.Format("2006-01-02 15:04:05"))

	data := map[string]string{"title": title, "desp": desp}
	body, _ := json.Marshal(data)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
