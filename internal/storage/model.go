package storage

// Data 持久化数据结构
type Data struct {
	Stocks    []StockItem    `json:"stocks"`
	Rules     []RuleItem     `json:"rules"`
	Notifiers NotifierConfig `json:"notifiers"`
}

// StockItem 股票配置
type StockItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// RuleItem 规则配置
type RuleItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Enabled   bool   `json:"enabled"`
	Level     string `json:"level"`
	StockCode string `json:"stock_code"`
	KLineType string `json:"kline_type"`
	Period    int    `json:"period"`
}

// NotifierConfig 通知配置
type NotifierConfig struct {
	ServerChan ServerChanConfig `json:"serverchan"`
	Feishu     FeishuConfig     `json:"feishu"`
	DingTalk   DingTalkConfig   `json:"dingtalk"`
}

type ServerChanConfig struct {
	Enabled bool   `json:"enabled"`
	SendKey string `json:"send_key"`
}

type FeishuConfig struct {
	Enabled bool   `json:"enabled"`
	Webhook string `json:"webhook"`
}

type DingTalkConfig struct {
	Enabled bool   `json:"enabled"`
	Webhook string `json:"webhook"`
}
