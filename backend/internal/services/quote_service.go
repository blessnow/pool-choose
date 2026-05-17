package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yuchi/cycle-stock/internal/models"
)

// SinaQuoteClient 行情客户端
type SinaQuoteClient struct {
	BaseURL string
}

func NewSinaQuoteClient() *SinaQuoteClient {
	return &SinaQuoteClient{
		BaseURL: "https://web.sqt.gtimg.cn/q=",
	}
}

// GetQuotes 批量获取行情数据 (使用腾讯财经API)
func (c *SinaQuoteClient) GetQuotes(codes []string) (map[string]models.Quote, error) {
	if len(codes) == 0 {
		return make(map[string]models.Quote), nil
	}

	// 转换股票代码格式: 000707 -> sz000707, 600000 -> sh600000
	tencentCodes := make([]string, len(codes))
	for i, code := range codes {
		if len(code) == 6 {
			if code[0] == '6' || code[0] == '5' {
				tencentCodes[i] = "sh" + code
			} else {
				tencentCodes[i] = "sz" + code
			}
		} else {
			tencentCodes[i] = code
		}
	}

	url := c.BaseURL + strings.Join(tencentCodes, ",")
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求腾讯API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	quotes := make(map[string]models.Quote)
	lines := strings.Split(string(body), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// 解析格式: v_sz000707="51~双环科技~000707~5.69~5.74~..."
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		// 提取股票代码
		codePart := parts[0]
		codePart = strings.TrimPrefix(codePart, "v_")
		codePart = strings.TrimPrefix(codePart, "sh")
		codePart = strings.TrimPrefix(codePart, "sz")

		// 解析数据
		dataPart := strings.Trim(parts[1], "\";")
		if dataPart == "" {
			continue
		}

		fields := strings.Split(dataPart, "~")
		if len(fields) < 35 {
			continue
		}

		quote := models.Quote{
			Code: codePart,
		}

		// 腾讯财经数据格式: 序号~名称~代码~当前价~昨收~今开~成交量~...
		// fields[3] = 当前价, fields[4] = 昨收
		quote.Close = parseFloat(fields[3])
		quote.PrevClose = parseFloat(fields[4])
		quote.Open = parseFloat(fields[5])
		quote.Volume = parseInt(fields[6])
		quote.Amount = parseFloat(fields[7])
		quote.High = parseFloat(fields[33])
		quote.Low = parseFloat(fields[34])

		// 计算涨跌幅
		if quote.PrevClose > 0 && quote.Close > 0 {
			quote.Change = (quote.Close - quote.PrevClose) / quote.PrevClose
		}

		quote.QuoteTime = time.Now()
		quotes[codePart] = quote
	}

	return quotes, nil
}

// GetKLineData 获取K线数据
func (c *SinaQuoteClient) GetKLineData(code string, period string) ([]map[string]interface{}, error) {
	// 转换股票代码格式
	var sinaCode string
	if len(code) == 6 {
		if code[0] == '6' || code[0] == '5' {
			sinaCode = "sh" + code
		} else {
			sinaCode = "sz" + code
		}
	} else {
		sinaCode = code
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// 对于日线和周线，直接从新浪获取
	// 对于月线，从日线数据聚合生成
	var scale string
	var dataLen int
	switch period {
	case "daily":
		scale = "240"
		dataLen = 300
	case "weekly":
		scale = "1200"
		dataLen = 300
	case "monthly":
		// 月线需要更多日线数据来聚合
		scale = "240"
		dataLen = 600
	default:
		scale = "240"
		dataLen = 300
	}

	url := fmt.Sprintf("https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=%s&scale=%s&datalen=%d",
		sinaCode, scale, dataLen)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求K线数据失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	result := make([]map[string]interface{}, 0)

	var klines []struct {
		Day     string `json:"day"`
		Open    string `json:"open"`
		Close   string `json:"close"`
		High    string `json:"high"`
		Low     string `json:"low"`
		Volume  string `json:"volume"`
	}
	if err := json.Unmarshal(body, &klines); err != nil {
		return result, nil
	}

	// 如果是月线，从日线数据聚合
	if period == "monthly" {
		return aggregateToMonthly(klines), nil
	}

	for _, k := range klines {
		item := map[string]interface{}{
			"date":   k.Day,
			"open":   parseFloat(k.Open),
			"close":  parseFloat(k.Close),
			"high":   parseFloat(k.High),
			"low":    parseFloat(k.Low),
			"volume": parseFloat(k.Volume),
		}
		result = append(result, item)
	}

	return result, nil
}

// aggregateToMonthly 从日线数据聚合生成月线数据
func aggregateToMonthly(dailyKlines []struct {
	Day     string `json:"day"`
	Open    string `json:"open"`
	Close   string `json:"close"`
	High    string `json:"high"`
	Low     string `json:"low"`
	Volume  string `json:"volume"`
}) []map[string]interface{} {
	if len(dailyKlines) == 0 {
		return []map[string]interface{}{}
	}

	// 按月份分组
	monthlyData := make(map[string][]struct {
		Day     string `json:"day"`
		Open    string `json:"open"`
		Close   string `json:"close"`
		High    string `json:"high"`
		Low     string `json:"low"`
		Volume  string `json:"volume"`
	})

	for _, k := range dailyKlines {
		// 日期格式: 2024-11-11，取年月作为key
		if len(k.Day) >= 7 {
			monthKey := k.Day[:7] // "2024-11"
			monthlyData[monthKey] = append(monthlyData[monthKey], k)
		}
	}

	// 聚合每个月的数据
	result := make([]map[string]interface{}, 0)
	for _, days := range monthlyData {
		if len(days) == 0 {
			continue
		}

		// 月开盘价 = 该月第一个交易日的开盘价
		// 月收盘价 = 该月最后一个交易日的收盘价
		// 月最高价 = 该月所有交易日最高价的最大值
		// 月最低价 = 该月所有交易日最低价的最小值
		// 月成交量 = 该月所有交易日成交量的总和
		monthOpen := parseFloat(days[0].Open)
		monthClose := parseFloat(days[len(days)-1].Close)
		monthHigh := 0.0
		monthLow := 999999.0
		monthVolume := 0.0

		for _, d := range days {
			high := parseFloat(d.High)
			low := parseFloat(d.Low)
			volume := parseFloat(d.Volume)
			if high > monthHigh {
				monthHigh = high
			}
			if low < monthLow {
				monthLow = low
			}
			monthVolume += volume
		}

		// 使用该月最后一个交易日作为日期
		lastDay := days[len(days)-1].Day

		result = append(result, map[string]interface{}{
			"date":   lastDay,
			"open":   monthOpen,
			"close":  monthClose,
			"high":   monthHigh,
			"low":    monthLow,
			"volume": monthVolume,
		})
	}

	// 按日期排序（从旧到新）
	sort.Slice(result, func(i, j int) bool {
		return result[i]["date"].(string) < result[j]["date"].(string)
	})

	return result
}

func parseFloat(s string) float64 {
	val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return val
}

func parseInt(s string) int64 {
	val, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// EastMoneyClient 东方财富客户端
type EastMoneyClient struct {
	BaseURL string
}

func NewEastMoneyClient() *EastMoneyClient {
	return &EastMoneyClient{
		BaseURL: "https://datacenter.eastmoney.com",
	}
}

// GetCompanyInfo 获取公司基本信息
func (c *EastMoneyClient) GetCompanyInfo(code string) (*models.CompanySummary, error) {
	// 东方财富API获取公司财务数据
	// 实际项目中需要完整实现
	return nil, nil
}

// GetStockList 获取股票列表
func (c *EastMoneyClient) GetStockList() ([]models.Stock, error) {
	// 从东方财富获取股票列表
	return nil, nil
}

// GetPEPBHistory 获取PE/PB历史数据
func (c *SinaQuoteClient) GetPEPBHistory(code string) ([]map[string]interface{}, error) {
	// 转换股票代码格式
	var secucode string
	if len(code) == 6 {
		if code[0] == '6' || code[0] == '5' {
			secucode = code + ".SH"
		} else {
			secucode = code + ".SZ"
		}
	} else {
		secucode = code
	}

	// 使用东方财富API获取财务数据
	url := fmt.Sprintf("https://emweb.eastmoney.com/PC_HSF10/NewFinanceAnalysis/ZYZBAjaxNew?type=0&code=%s", strings.ToLower(secucode))

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://emweb.eastmoney.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求财务数据失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON
	var result struct {
		Data []struct {
			ReportDate string  `json:"REPORT_DATE"`
			ROEJQ      float64 `json:"ROEJQ"`
			BPS        float64 `json:"BPS"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 转换为输出格式
	data := make([]map[string]interface{}, 0)
	for _, item := range result.Data {
		if len(item.ReportDate) >= 10 {
			data = append(data, map[string]interface{}{
				"date": item.ReportDate[:10],
				"roe":  item.ROEJQ,
				"bps":  item.BPS,
			})
		}
	}

	return data, nil
}

// GetValuationHistory 获取估值历史数据（每日PE/PB）
func (c *SinaQuoteClient) GetValuationHistory(code string, days int) ([]map[string]interface{}, error) {
	// 使用新浪财经API获取K线数据
	var scale string
	if days <= 100 {
		scale = "240" // 日线
	} else if days <= 300 {
		scale = "240"
	} else {
		scale = "240"
	}

	// 转换股票代码格式 - 新浪格式
	var sinaCode string
	if len(code) == 6 {
		if code[0] == '6' || code[0] == '5' {
			sinaCode = "sh" + code
		} else {
			sinaCode = "sz" + code
		}
	} else {
		sinaCode = code
	}

	url := fmt.Sprintf("https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=%s&scale=%s&datalen=%d",
		sinaCode, scale, days)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求K线数据失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON数组
	var klines []struct {
		Day    string `json:"day"`
		Open   string `json:"open"`
		Close  string `json:"close"`
		High   string `json:"high"`
		Low    string `json:"low"`
		Volume string `json:"volume"`
	}
	if err := json.Unmarshal(body, &klines); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 转换为输出格式
	data := make([]map[string]interface{}, 0)
	for _, k := range klines {
		data = append(data, map[string]interface{}{
			"date":   k.Day,
			"open":   parseFloat(k.Open),
			"close":  parseFloat(k.Close),
			"high":   parseFloat(k.High),
			"low":    parseFloat(k.Low),
			"volume": parseFloat(k.Volume),
		})
	}

	return data, nil
}

// GetCurrentValuation 获取当前估值数据
func (c *SinaQuoteClient) GetCurrentValuation(code string) (map[string]interface{}, error) {
	// 转换股票代码格式
	var secid string
	if len(code) == 6 {
		if code[0] == '6' || code[0] == '5' {
			secid = "1." + code
		} else {
			secid = "0." + code
		}
	} else {
		secid = "0." + code
	}

	url := fmt.Sprintf("https://push2.eastmoney.com/api/qt/stock/get?secid=%s&fields=f57,f58,f162,f167,f92,f173,f187,f43,f44,f45,f46", secid)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求估值数据失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON
	var result struct {
		Data map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	data := result.Data
	if data == nil {
		return nil, fmt.Errorf("未获取到估值数据")
	}

	// f162: PE(TTM), f167: PE(静), f92: PB, f173: ROE, f187: 市销率
	// f43: 最新价(分), f44: 最高价(分), f45: 最低价(分), f46: 开盘价(分)
	return map[string]interface{}{
		"code":    data["f57"],
		"name":    data["f58"],
		"pe_ttm":  parseFloat(fmt.Sprintf("%v", data["f162"])) / 100,
		"pe":      parseFloat(fmt.Sprintf("%v", data["f167"])) / 100,
		"pb":      parseFloat(fmt.Sprintf("%v", data["f92"])),
		"roe":     parseFloat(fmt.Sprintf("%v", data["f173"])),
		"ps":      parseFloat(fmt.Sprintf("%v", data["f187"])),
		"close":   parseFloat(fmt.Sprintf("%v", data["f43"])) / 100,
		"high":    parseFloat(fmt.Sprintf("%v", data["f44"])) / 100,
		"low":     parseFloat(fmt.Sprintf("%v", data["f45"])) / 100,
		"open":    parseFloat(fmt.Sprintf("%v", data["f46"])) / 100,
	}, nil
}

// ParseCSVFromPage 从页面CSV解析数据
func ParseCSVFromPage(csvData string) ([]models.Stock, error) {
	reader := csv.NewReader(strings.NewReader(csvData))
	var stocks []models.Stock

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if len(record) < 8 {
			continue
		}

		stock := models.Stock{
			Code:        record[0],
			Name:        record[1],
			Sector:      record[2],
			IsRecommend: record[3] == "推荐" || record[3] == "true",
			EntryPrice:  parseFloat(record[4]),
			HeavyPrice:  parseFloat(record[5]),
			TargetPrice: parseFloat(record[6]),
			CoreLogic:   record[7],
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}