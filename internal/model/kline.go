package model

import "time"

// KLineType K线类型
type KLineType string

const (
	KLine5Min    KLineType = "5min"    // 5分钟K
	KLine15Min   KLineType = "15min"   // 15分钟K
	KLine30Min   KLineType = "30min"   // 30分钟K
	KLine60Min   KLineType = "60min"   // 60分钟K
	KLineDaily   KLineType = "daily"   // 日K
	KLineWeekly  KLineType = "weekly"  // 周K
	KLineMonthly KLineType = "monthly" // 月K
)

// KLine K线数据
type KLine struct {
	Time   time.Time `json:"time"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int64     `json:"volume"`
}

// KLineData K线数据集合
type KLineData struct {
	Code  string    `json:"code"`
	Type  KLineType `json:"type"`
	Lines []KLine   `json:"lines"`
}
