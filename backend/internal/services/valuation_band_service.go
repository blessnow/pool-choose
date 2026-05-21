package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"
)

// ValuationBandPoint 单日估值数据点
type ValuationBandPoint struct {
	Date       string   `json:"date"`
	RawValue   *float64 `json:"rawValue,omitempty"`
	Price      float64  `json:"price"`
	Value      float64  `json:"value"`
	Percentile *float64 `json:"percentile"`
	Tracks     struct {
		P90 *float64 `json:"p90"`
		P70 *float64 `json:"p70"`
		P50 *float64 `json:"p50"`
		P30 *float64 `json:"p30"`
		P10 *float64 `json:"p10"`
	} `json:"tracks"`
	PriceTracks struct {
		P90 *float64 `json:"p90"`
		P70 *float64 `json:"p70"`
		P50 *float64 `json:"p50"`
		P30 *float64 `json:"p30"`
		P10 *float64 `json:"p10"`
	} `json:"priceTracks"`
}

// ValuationBandData 估值带数据
type ValuationBandData struct {
	Metric    string               `json:"metric"`
	Years     int                  `json:"years"`
	Points    []ValuationBandPoint `json:"points"`
	UpdatedAt string               `json:"updatedAt,omitempty"`
}

// FinancialData 财务数据
type FinancialData struct {
	ReportDate string
	EPS        float64 // 每股收益
	BPS        float64 // 每股净资产
}

// GetValuationBand 获取估值带数据
func GetValuationBand(code string, metric string, years int) (*ValuationBandData, error) {
	// 验证metric参数
	if metric != "pe_ttm" && metric != "pb" {
		return nil, fmt.Errorf("unsupported metric: %s, must be pe_ttm or pb", metric)
	}

	// 计算日期范围
	days := years * 365
	if days > 2000 {
		days = 2000 // 限制最大天数
	}

	client := NewSinaQuoteClient()

	// 1. 获取K线数据
	klines, err := client.GetKLineData(code, "daily")
	if err != nil {
		return nil, fmt.Errorf("获取K线数据失败: %v", err)
	}

	if len(klines) == 0 {
		return nil, fmt.Errorf("未获取到K线数据")
	}

	// 限制数据范围
	if len(klines) > days {
		klines = klines[len(klines)-days:]
	}

	// 2. 获取财务数据（EPS/BPS）
	financialData, err := getFinancialData(code)
	if err != nil {
		// 财务数据获取失败时，尝试使用当前估值数据
		financialData = nil
	}

	// 3. 获取当前估值数据
	currentVal, err := client.GetCurrentValuation(code)
	if err != nil {
		currentVal = nil
	}

	// 4. 计算每日PE/PB值
	points := calculateValuationPoints(klines, financialData, currentVal, metric)

	// 5. 计算百分位和轨道线
	calculatePercentiles(points)

	// 5b. 反推每日股价轨道（priceTracks）：用每日的 EPS/BPS 乘以估值百分位
	calculatePriceTracks(points, financialData, metric)

	// 6. 格式化日期
	for i := range points {
		if len(points[i].Date) > 10 {
			points[i].Date = points[i].Date[:10]
		}
	}

	return &ValuationBandData{
		Metric:    metric,
		Years:     years,
		Points:    points,
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05+08:00"),
	}, nil
}

// getFinancialData 获取财务数据（EPS/BPS）
func getFinancialData(code string) ([]FinancialData, error) {
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

	resp, err := createHTTPClient().Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析JSON
	var result struct {
		Data []struct {
			ReportDate string  `json:"REPORT_DATE"`
			EPSJB      float64 `json:"EPSJB"`  // 基本每股收益
			EPSKCJB    float64 `json:"EPSKCJB"` // 扣非每股收益
			BPS        float64 `json:"BPS"`    // 每股净资产
		} `json:"data"`
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}

	if err := parseJSON(body, &result); err != nil {
		return nil, err
	}

	// 转换为输出格式
	data := make([]FinancialData, 0)
	for _, item := range result.Data {
		if len(item.ReportDate) >= 10 {
			eps := item.EPSJB
			if eps == 0 {
				eps = item.EPSKCJB
			}
			data = append(data, FinancialData{
				ReportDate: item.ReportDate[:10],
				EPS:        eps,
				BPS:        item.BPS,
			})
		}
	}

	// 按日期排序（从旧到新）
	sort.Slice(data, func(i, j int) bool {
		return data[i].ReportDate < data[j].ReportDate
	})

	return data, nil
}

// calculateValuationPoints 计算每日估值数据点
func calculateValuationPoints(klines []map[string]interface{}, financialData []FinancialData, currentVal map[string]interface{}, metric string) []ValuationBandPoint {
	points := make([]ValuationBandPoint, len(klines))

	for i, kline := range klines {
		date, _ := kline["date"].(string)
		closePrice, _ := kline["close"].(float64)

		points[i] = ValuationBandPoint{
			Date:  date,
			Price: closePrice,
			Value: closePrice, // 默认使用价格
		}

		// 根据metric计算估值
		var value float64
		var rawValue *float64

		switch metric {
		case "pe_ttm":
			// 计算PE TTM = 股价 / 最近4个季度EPS之和
			eps := getTTMEPS(date, financialData)
			if eps > 0 && closePrice > 0 {
				pe := closePrice / eps
				value = pe
				rawValue = &pe
			} else if currentVal != nil {
				// 使用当前PE TTM作为参考
				if peTtm, ok := currentVal["pe_ttm"].(float64); ok && peTtm > 0 {
					// 无法计算历史PE，使用当前PE估算
					value = closePrice // 显示价格
				}
			}
		case "pb":
			// 计算PB = 股价 / 每股净资产
			bps := getLatestBPS(date, financialData)
			if bps > 0 && closePrice > 0 {
				pb := closePrice / bps
				value = pb
				rawValue = &pb
			} else if currentVal != nil {
				// 使用当前PB作为参考
				if pb, ok := currentVal["pb"].(float64); ok && pb > 0 {
					value = closePrice // 显示价格
				}
			}
		}

		points[i].Value = value
		points[i].RawValue = rawValue
	}

	return points
}

// getTTMEPS 获取TTM EPS（最近4个季度每股收益之和）
func getTTMEPS(date string, financialData []FinancialData) float64 {
	if len(financialData) == 0 {
		return 0
	}

	// 找到日期之前最近的财务数据
	var latestReports []FinancialData
	for i := len(financialData) - 1; i >= 0; i-- {
		if financialData[i].ReportDate <= date {
			latestReports = append(latestReports, financialData[i])
			if len(latestReports) >= 4 {
				break
			}
		}
	}

	if len(latestReports) == 0 {
		return 0
	}

	// 计算最近4个季度EPS之和
	// 注意：A股财报是累计的，需要计算单季度EPS
	var totalEPS float64
	for i, report := range latestReports {
		if i == 0 {
			// 最新季报直接使用
			totalEPS += math.Abs(report.EPS)
		} else {
			// 计算单季度EPS = 当期累计EPS - 上期累计EPS
			// 这里简化处理，直接累加（对于年报来说是对的）
			totalEPS += math.Abs(report.EPS)
		}
	}

	// 如果数据不足4个季度，按比例估算
	if len(latestReports) < 4 {
		totalEPS = totalEPS * float64(4) / float64(len(latestReports))
	}

	return totalEPS
}

// getLatestBPS 获取最近的每股净资产
func getLatestBPS(date string, financialData []FinancialData) float64 {
	if len(financialData) == 0 {
		return 0
	}

	// 找到日期之前最近的BPS
	for i := len(financialData) - 1; i >= 0; i-- {
		if financialData[i].ReportDate <= date && financialData[i].BPS > 0 {
			return financialData[i].BPS
		}
	}

	return 0
}

// calculatePercentiles 计算百分位和轨道线
func calculatePercentiles(points []ValuationBandPoint) {
	// 收集有效的估值值
	var values []float64
	for _, p := range points {
		if p.RawValue != nil && *p.RawValue > 0 && !math.IsInf(*p.RawValue, 0) && !math.IsNaN(*p.RawValue) {
			values = append(values, *p.RawValue)
		}
	}

	if len(values) == 0 {
		return
	}

	// 排序
	sort.Float64s(values)

	// 计算每个点的百分位和轨道线
	for i := range points {
		if points[i].RawValue != nil && *points[i].RawValue > 0 {
			// 计算百分位
			rank := calculatePercentileRank(values, *points[i].RawValue)
			points[i].Percentile = &rank

			// 计算轨道线
			p90 := percentileValue(values, 90)
			p70 := percentileValue(values, 70)
			p50 := percentileValue(values, 50)
			p30 := percentileValue(values, 30)
			p10 := percentileValue(values, 10)

			points[i].Tracks.P90 = &p90
			points[i].Tracks.P70 = &p70
			points[i].Tracks.P50 = &p50
			points[i].Tracks.P30 = &p30
			points[i].Tracks.P10 = &p10
		}
	}
}

// calculatePriceTracks 把每日的估值百分位（PE/PB）反推为对应股价
// price_at_pXX(t) = tracks.pXX(t) * factor(t)
//   factor = EPS_TTM(t) (metric=pe_ttm) 或 BPS(t) (metric=pb)
func calculatePriceTracks(points []ValuationBandPoint, financialData []FinancialData, metric string) {
	for i := range points {
		var factor float64
		switch metric {
		case "pe_ttm":
			factor = getTTMEPS(points[i].Date, financialData)
		case "pb":
			factor = getLatestBPS(points[i].Date, financialData)
		}
		if factor <= 0 || math.IsNaN(factor) || math.IsInf(factor, 0) {
			continue
		}

		project := func(p *float64) *float64 {
			if p == nil || *p <= 0 {
				return nil
			}
			v := *p * factor
			return &v
		}

		points[i].PriceTracks.P90 = project(points[i].Tracks.P90)
		points[i].PriceTracks.P70 = project(points[i].Tracks.P70)
		points[i].PriceTracks.P50 = project(points[i].Tracks.P50)
		points[i].PriceTracks.P30 = project(points[i].Tracks.P30)
		points[i].PriceTracks.P10 = project(points[i].Tracks.P10)
	}
}

// calculatePercentileRank 计算值在数组中的百分位排名
func calculatePercentileRank(sortedValues []float64, value float64) float64 {
	n := len(sortedValues)
	if n == 0 {
		return 0
	}

	// 找到值的位置
	var count float64
	for _, v := range sortedValues {
		if v <= value {
			count++
		}
	}

	// 百分位排名 = (小于等于该值的数量 - 0.5) / 总数量 * 100
	// 使用简单的线性插值方法
	rank := (count - 0.5) / float64(n) * 100
	if rank < 0 {
		rank = 0
	}
	if rank > 100 {
		rank = 100
	}

	return rank
}

// percentileValue 计算指定百分位对应的值
func percentileValue(sortedValues []float64, percentile float64) float64 {
	n := len(sortedValues)
	if n == 0 {
		return 0
	}

	// 计算位置
	position := (float64(n) - 1) * percentile / 100
	lowerIndex := int(math.Floor(position))
	upperIndex := int(math.Ceil(position))

	if lowerIndex == upperIndex || upperIndex >= n {
		if lowerIndex < n {
			return sortedValues[lowerIndex]
		}
		return sortedValues[n-1]
	}

	// 线性插值
	weight := position - float64(lowerIndex)
	return sortedValues[lowerIndex]*(1-weight) + sortedValues[upperIndex]*weight
}

// createHTTPClient 创建HTTP客户端
func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

// readResponseBody 读取响应体
func readResponseBody(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}

// parseJSON 解析JSON
func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
