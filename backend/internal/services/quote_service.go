package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yuchi/cycle-stock/internal/models"
)

// K线内存缓存 —— 东财对突发请求有限流，缓存可以大幅降低 500 EOF 率，
// 也匹配真实用法（用户点单只股票，反复切换 daily/weekly/monthly）。
type klineCacheEntry struct {
	data     []map[string]interface{}
	storedAt time.Time
}

var (
	klineCache   = map[string]klineCacheEntry{}
	klineCacheMu sync.RWMutex
)

// 不同周期的 TTL：日线 2h（盘中可手动重启拿最新），周/月线 12h。
func klineCacheTTL(period string) time.Duration {
	switch period {
	case "weekly", "monthly":
		return 12 * time.Hour
	default:
		return 2 * time.Hour
	}
}

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

// GetKLineData 获取K线数据（腾讯 fqkline，前复权）
// 腾讯日线上限 ~640 条 (~2.5 年)，周线 ~12 年，月线 ~20+ 年，全量股票稳定命中。
func (c *SinaQuoteClient) GetKLineData(code string, period string) ([]map[string]interface{}, error) {
	cacheKey := code + ":" + period
	ttl := klineCacheTTL(period)
	klineCacheMu.RLock()
	if e, ok := klineCache[cacheKey]; ok && time.Since(e.storedAt) < ttl {
		klineCacheMu.RUnlock()
		return e.data, nil
	}
	klineCacheMu.RUnlock()

	result, err := fetchTencentKLine(code, period)
	if err != nil {
		// 失败时回退到（即使过期的）旧缓存，避免页面全空
		klineCacheMu.RLock()
		if e, ok := klineCache[cacheKey]; ok && len(e.data) > 0 {
			klineCacheMu.RUnlock()
			return e.data, nil
		}
		klineCacheMu.RUnlock()
		return nil, err
	}

	if len(result) > 0 {
		klineCacheMu.Lock()
		klineCache[cacheKey] = klineCacheEntry{data: result, storedAt: time.Now()}
		klineCacheMu.Unlock()
	}
	return result, nil
}

// fetchTencentKLine 从腾讯 web.ifzq.gtimg.cn 拉 K 线（前复权）
// 日线上限约 640 条 (~2.5 年)，周/月线足够长。EM 拿不到的股票用这个兜底。
func fetchTencentKLine(code, period string) ([]map[string]interface{}, error) {
	var symbol string
	if len(code) == 6 {
		if code[0] == '6' || code[0] == '5' {
			symbol = "sh" + code
		} else {
			symbol = "sz" + code
		}
	} else {
		symbol = code
	}

	var qqPeriod string
	var n int
	switch period {
	case "weekly":
		qqPeriod, n = "week", 640
	case "monthly":
		qqPeriod, n = "month", 640
	default:
		qqPeriod, n = "day", 640
	}

	url := fmt.Sprintf("https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?param=%s,%s,,,%d,qfq", symbol, qqPeriod, n)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://gu.qq.com/")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 响应结构: {"code":0,"data":{"<symbol>":{"qfqday":[["date","open","close","high","low","volume"], ...]}}}
	// 异常时 data.<symbol> 会变成 list 而不是 dict，用 RawMessage 容错
	var top struct {
		Code int                                   `json:"code"`
		Data map[string]map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &top); err != nil {
		return nil, fmt.Errorf("解析腾讯K线失败: %v", err)
	}
	sym, ok := top.Data[symbol]
	if !ok {
		return nil, fmt.Errorf("腾讯K线无数据: %s", symbol)
	}
	raw, ok := sym["qfq"+qqPeriod]
	if !ok {
		raw, ok = sym[qqPeriod]
		if !ok {
			return nil, fmt.Errorf("腾讯K线字段缺失: qfq%s", qqPeriod)
		}
	}
	var rows [][]interface{}
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, fmt.Errorf("解析K线数组失败: %v", err)
	}

	result := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		if len(r) < 6 {
			continue
		}
		// 每个字段都是 string，转一道
		toStr := func(v interface{}) string {
			if s, ok := v.(string); ok {
				return s
			}
			return fmt.Sprintf("%v", v)
		}
		result = append(result, map[string]interface{}{
			"date":   toStr(r[0]),
			"open":   parseFloat(toStr(r[1])),
			"close":  parseFloat(toStr(r[2])),
			"high":   parseFloat(toStr(r[3])),
			"low":    parseFloat(toStr(r[4])),
			"volume": parseFloat(toStr(r[5])),
		})
	}
	return result, nil
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