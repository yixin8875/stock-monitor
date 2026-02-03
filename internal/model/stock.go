package model

import "time"

// Stock 股票基本信息
type Stock struct {
	Code     string    `json:"code"`      // 股票代码 如 600519
	Name     string    `json:"name"`      // 股票名称
	Exchange string    `json:"exchange"`  // 交易所 sh/sz
	Price    float64   `json:"price"`     // 当前价格
	Open     float64   `json:"open"`      // 开盘价
	High     float64   `json:"high"`      // 最高价
	Low      float64   `json:"low"`       // 最低价
	Close    float64   `json:"close"`     // 收盘价
	PreClose float64   `json:"pre_close"` // 昨收价
	Volume   int64     `json:"volume"`    // 成交量
	Amount   float64   `json:"amount"`    // 成交额
	Time     time.Time `json:"time"`      // 行情时间
}

// FullCode 返回完整股票代码 (带交易所前缀)
func (s *Stock) FullCode() string {
	return s.Exchange + s.Code
}

// ChangePercent 涨跌幅
func (s *Stock) ChangePercent() float64 {
	if s.PreClose == 0 {
		return 0
	}
	return (s.Price - s.PreClose) / s.PreClose * 100
}
