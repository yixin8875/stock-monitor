package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"stock-monitor/internal/model"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// klineAPIResponse K线API响应结构
type klineAPIResponse struct {
	Data struct {
		Symbol string     `json:"symbol"`
		Kline  [][]string `json:"kline"`
	} `json:"data"`
}

// SinaDataSource 新浪数据源
type SinaDataSource struct {
	client *http.Client
}

// NewSinaDataSource 创建新浪数据源
func NewSinaDataSource() *SinaDataSource {
	return &SinaDataSource{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SinaDataSource) Name() string {
	return "sina"
}

// formatCode 格式化股票代码
func (s *SinaDataSource) formatCode(code string) string {
	if strings.HasPrefix(code, "6") {
		return "sh" + code
	}
	return "sz" + code
}

// GetRealTimeQuote 获取实时行情
func (s *SinaDataSource) GetRealTimeQuote(ctx context.Context, codes []string) ([]*model.Stock, error) {
	symbols := make([]string, len(codes))
	for i, code := range codes {
		symbols[i] = s.formatCode(code)
	}

	url := fmt.Sprintf("https://hq.sinajs.cn/list=%s", strings.Join(symbols, ","))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://finance.sina.com.cn")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// GBK 转 UTF-8
	reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return s.parseQuoteResponse(string(body), codes)
}

// parseQuoteResponse 解析行情响应
func (s *SinaDataSource) parseQuoteResponse(body string, codes []string) ([]*model.Stock, error) {
	re := regexp.MustCompile(`var hq_str_(\w+)="([^"]*)"`)
	matches := re.FindAllStringSubmatch(body, -1)

	var stocks []*model.Stock
	for _, match := range matches {
		if len(match) < 3 || match[2] == "" {
			continue
		}

		symbol := match[1]
		data := strings.Split(match[2], ",")
		if len(data) < 32 {
			continue
		}

		stock := &model.Stock{
			Code:     symbol[2:],
			Exchange: symbol[:2],
			Name:     data[0],
		}

		stock.Open, _ = strconv.ParseFloat(data[1], 64)
		stock.PreClose, _ = strconv.ParseFloat(data[2], 64)
		stock.Price, _ = strconv.ParseFloat(data[3], 64)
		stock.High, _ = strconv.ParseFloat(data[4], 64)
		stock.Low, _ = strconv.ParseFloat(data[5], 64)
		stock.Close = stock.Price
		vol, _ := strconv.ParseFloat(data[8], 64)
		stock.Volume = int64(vol)
		stock.Amount, _ = strconv.ParseFloat(data[9], 64)

		timeStr := data[30] + " " + data[31]
		stock.Time, _ = time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)

		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetKLine 获取K线数据
func (s *SinaDataSource) GetKLine(ctx context.Context, code string, ktype model.KLineType, count int) (*model.KLineData, error) {
	symbol := s.formatCode(code)
	scale := s.klineTypeToScale(ktype)

	url := fmt.Sprintf("https://quotes.sina.cn/cn/api/jsonp_v2.php/var%%20_%s_%s=/CN_MarketDataService.getKLineData?symbol=%s&scale=%s&datalen=%d",
		symbol, scale, symbol, scale, count)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://finance.sina.com.cn")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return s.parseKLineResponse(string(body), code, ktype)
}

func (s *SinaDataSource) klineTypeToScale(ktype model.KLineType) string {
	switch ktype {
	case model.KLine5Min:
		return "5"
	case model.KLine15Min:
		return "15"
	case model.KLine30Min:
		return "30"
	case model.KLine60Min:
		return "60"
	case model.KLineDaily:
		return "240"
	case model.KLineWeekly:
		return "1200"
	case model.KLineMonthly:
		return "7200"
	default:
		return "240"
	}
}

func (s *SinaDataSource) parseKLineResponse(body string, code string, ktype model.KLineType) (*model.KLineData, error) {
	start := strings.Index(body, "[")
	end := strings.LastIndex(body, "]")
	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("invalid kline response")
	}

	jsonStr := body[start : end+1]
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &items); err != nil {
		return nil, err
	}

	klineData := &model.KLineData{
		Code:  code,
		Type:  ktype,
		Lines: make([]model.KLine, 0, len(items)),
	}

	for _, item := range items {
		kline := model.KLine{}
		if day, ok := item["day"].(string); ok {
			kline.Time, _ = time.ParseInLocation("2006-01-02", day, time.Local)
		}
		kline.Open = parseFloat(item["open"])
		kline.High = parseFloat(item["high"])
		kline.Low = parseFloat(item["low"])
		kline.Close = parseFloat(item["close"])
		kline.Volume = int64(parseFloat(item["volume"]))
		klineData.Lines = append(klineData.Lines, kline)
	}

	return klineData, nil
}

func parseFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}
	return 0
}
