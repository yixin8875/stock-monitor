package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	DataSource DataSourceConfig  `yaml:"datasource"`
	Stocks     []StockConfig     `yaml:"stocks"`
	Rules      []RuleConfig      `yaml:"rules"`
	Notifiers  NotifiersConfig   `yaml:"notifiers"`
	Schedule   ScheduleConfig    `yaml:"schedule"`
}

// DataSourceConfig 数据源配置
type DataSourceConfig struct {
	Name string `yaml:"name"`
}

// StockConfig 股票配置
type StockConfig struct {
	Code string `yaml:"code"`
	Name string `yaml:"name"`
}

// RuleConfig 规则配置
type RuleConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Enabled bool                   `yaml:"enabled"`
	Level   string                 `yaml:"level"`
	Params  map[string]interface{} `yaml:"params"`
}

// NotifiersConfig 通知配置
type NotifiersConfig struct {
	ServerChan ServerChanConfig `yaml:"serverchan"`
	DingTalk   DingTalkConfig   `yaml:"dingtalk"`
	Feishu     FeishuConfig     `yaml:"feishu"`
}

// ServerChanConfig Server酱配置
type ServerChanConfig struct {
	Enabled bool   `yaml:"enabled"`
	SendKey string `yaml:"send_key"`
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	Enabled bool   `yaml:"enabled"`
	Webhook string `yaml:"webhook"`
}

// FeishuConfig 飞书配置
type FeishuConfig struct {
	Enabled bool   `yaml:"enabled"`
	Webhook string `yaml:"webhook"`
}

// ScheduleConfig 调度配置
type ScheduleConfig struct {
	Cron string `yaml:"cron"`
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
