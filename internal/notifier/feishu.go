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

// Feishu 飞书通知
type Feishu struct {
	webhook string
	client  *http.Client
}

// NewFeishu 创建飞书通知
func NewFeishu(webhook string) *Feishu {
	return &Feishu{
		webhook: webhook,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (f *Feishu) Name() string {
	return "feishu"
}

func (f *Feishu) Send(ctx context.Context, alert *model.Alert) error {
	msg := f.buildMessage(alert)
	body, _ := json.Marshal(msg)

	req, err := http.NewRequestWithContext(ctx, "POST", f.webhook, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("飞书通知失败, 状态码: %d", resp.StatusCode)
	}

	return nil
}

func (f *Feishu) buildMessage(alert *model.Alert) map[string]interface{} {
	return map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"header": map[string]interface{}{
				"title": map[string]string{
					"tag":     "plain_text",
					"content": "📈 股票监控告警",
				},
				"template": f.levelToColor(alert.Level),
			},
			"elements": []interface{}{
				map[string]interface{}{
					"tag": "div",
					"fields": []map[string]interface{}{
						{"is_short": true, "text": map[string]string{"tag": "lark_md", "content": "**股票**\n" + alert.StockName}},
						{"is_short": true, "text": map[string]string{"tag": "lark_md", "content": "**代码**\n" + alert.StockCode}},
						{"is_short": true, "text": map[string]string{"tag": "lark_md", "content": "**价格**\n" + formatPrice(alert.Price)}},
						{"is_short": true, "text": map[string]string{"tag": "lark_md", "content": "**级别**\n" + string(alert.Level)}},
					},
				},
				map[string]interface{}{
					"tag":     "div",
					"text":    map[string]string{"tag": "lark_md", "content": "**消息**\n" + alert.Message},
				},
				map[string]interface{}{
					"tag":     "note",
					"elements": []map[string]string{{"tag": "plain_text", "content": alert.Time.Format("2006-01-02 15:04:05")}},
				},
			},
		},
	}
}

func (f *Feishu) levelToColor(level model.AlertLevel) string {
	switch level {
	case model.AlertLevelCritical:
		return "red"
	case model.AlertLevelWarning:
		return "orange"
	default:
		return "blue"
	}
}

func formatPrice(price float64) string {
	return fmt.Sprintf("%.2f", price)
}

