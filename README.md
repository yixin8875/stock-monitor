# Stock Monitor - 股票监控系统

一个基于 Go 语言的可扩展股票监控系统，支持自定义规则和多渠道通知。

## 功能特性

- **实时行情监控**: 支持 A 股市场数据（新浪财经数据源）
- **可扩展规则引擎**: 支持自定义监控规则，如均线突破、涨跌幅等
- **多渠道通知**: 支持 Server酱、飞书、钉钉等多种通知方式
- **Web 管理后台**: 可视化管理股票、规则和通知配置
- **多周期 K 线**: 支持 5分钟/15分钟/30分钟/60分钟/日K 等多种周期

## 快速开始

### 环境要求

- Go 1.21+（源码编译）
- Docker & Docker Compose（容器运行）

### 方式一：源码编译

```bash
# 克隆项目
git clone <repository-url>
cd stock-monitor

# 安装依赖
go mod tidy

# 编译
go build -o monitor ./cmd/monitor

# 运行（启动 Web 后台 + 监控服务）
./monitor

# 单次检查（不启动 Web）
./monitor -once
```

### 方式二：Docker 运行

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## Web 管理后台

启动程序后访问 http://localhost:8080 进入管理后台。

### 功能说明

1. **股票管理**: 添加/删除监控的股票
2. **规则管理**: 创建/启用/禁用监控规则，支持选择 K 线周期
3. **通知配置**: 配置飞书、Server酱、钉钉等通知渠道

### K 线周期选项

| 周期 | 说明 |
|------|------|
| 5min | 5分钟K线 |
| 15min | 15分钟K线 |
| 30min | 30分钟K线 |
| 60min | 60分钟K线 |
| daily | 日K线 |

### 示例：创建 15 分钟 MA60 规则

1. 先添加股票（如 600519 贵州茅台）
2. 在规则管理中：
   - 规则名称：茅台15分钟MA60
   - 选择股票：贵州茅台
   - K线周期：15分钟
   - MA周期：60
   - 点击「添加规则」

## 命令行参数

```bash
./monitor [options]

选项：
  -data string    数据文件路径 (默认 "data/config.json")
  -addr string    Web服务地址 (默认 ":8080")
  -once           只运行一次后退出（不启动Web服务）
```

**示例：**

```bash
# 启动 Web 后台 + 监控服务
./monitor

# 指定端口
./monitor -addr :9090

# 单次检查
./monitor -once
```

## 数据存储

配置数据存储在 `data/config.json`，包含：
- 股票列表
- 规则配置
- 通知渠道配置

通过 Web 后台修改的配置会自动持久化到该文件。

## 通知渠道配置

通过 Web 后台配置通知渠道，支持以下方式：

### 飞书机器人

1. 打开飞书，进入目标群聊
2. 点击右上角「...」→「设置」→「群机器人」→「添加机器人」
3. 选择「自定义机器人」，复制 Webhook 地址
4. 在 Web 后台填入 Webhook 地址并启用

### Server酱（微信通知）

1. 访问 [Server酱官网](https://sct.ftqq.com/) 并登录
2. 获取 SendKey
3. 在 Web 后台填入 SendKey 并启用

### 钉钉机器人

1. 打开钉钉群 → 群设置 → 智能群助手 → 添加机器人
2. 选择「自定义」机器人，获取 Webhook
3. 在 Web 后台填入 Webhook 地址并启用

## 项目结构

```
stock-monitor/
├── cmd/monitor/main.go        # 程序入口
├── data/config.json           # 持久化数据（自动生成）
├── internal/
│   ├── api/                   # Web API 服务
│   ├── datasource/            # 数据源（新浪）
│   ├── model/                 # 数据模型
│   ├── rule/                  # 规则引擎
│   │   └── rules/             # 具体规则实现
│   ├── indicator/             # 技术指标
│   ├── notifier/              # 通知模块
│   ├── storage/               # 数据持久化
│   └── monitor/               # 监控主逻辑
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 扩展开发

### 添加新规则

1. 在 `internal/rule/rules/` 创建新文件，如 `price_change.go`
2. 实现 `Rule` 接口
3. 在 `init()` 中注册规则

```go
package rules

import (
    "stock-monitor/internal/rule"
)

func init() {
    rule.GlobalRegistry.Register("price_change", NewPriceChangeRule)
}

type PriceChangeRule struct {
    // 字段定义
}

func NewPriceChangeRule(name string, level model.AlertLevel, params map[string]interface{}) (rule.Rule, error) {
    // 创建规则实例
}

func (r *PriceChangeRule) Name() string { return r.name }
func (r *PriceChangeRule) Description() string { return "涨跌幅监控" }
func (r *PriceChangeRule) Validate() error { return nil }
func (r *PriceChangeRule) Evaluate(ctx context.Context, ruleCtx *rule.RuleContext) (*rule.RuleResult, error) {
    // 规则逻辑
}
```

### 添加新通知渠道

1. 在 `internal/notifier/` 创建新文件
2. 实现 `Notifier` 接口
3. 在 `monitor.go` 的 `Setup()` 中加载

```go
package notifier

type MyNotifier struct {
    // 配置字段
}

func (n *MyNotifier) Name() string { return "my_notifier" }
func (n *MyNotifier) Send(ctx context.Context, alert *model.Alert) error {
    // 发送逻辑
}
```

## 常见问题

**Q: 为什么没有收到通知？**
- 检查 Web 后台中对应通知渠道是否已启用
- 检查 webhook 地址是否正确
- 查看程序日志确认规则是否触发

**Q: 如何调整监控频率？**
- 默认每 5 分钟检查一次
- 修改 `cmd/monitor/main.go` 中的 `time.NewTicker(5 * time.Minute)`

**Q: 支持哪些股票市场？**
- 目前支持 A 股（上证、深证）
- 股票代码：6 开头为上证，0/3 开头为深证

## License

MIT
