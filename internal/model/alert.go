package model

import "time"

// AlertLevel 告警级别
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelCritical AlertLevel = "critical"
)

// Alert 告警信息
type Alert struct {
	ID        string                 `json:"id"`
	StockCode string                 `json:"stock_code"`
	StockName string                 `json:"stock_name"`
	RuleName  string                 `json:"rule_name"`
	Level     AlertLevel             `json:"level"`
	Message   string                 `json:"message"`
	Price     float64                `json:"price"`
	Time      time.Time              `json:"time"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}
