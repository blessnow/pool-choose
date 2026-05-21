package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/yuchi/cycle-stock/internal/models"
)

// macroCache 15 分钟内存缓存：CPI/PPI 月度发布，无需高频拉取
var macroCache struct {
	sync.Mutex
	data    []models.MacroCard
	fetched time.Time
}

const macroCacheTTL = 15 * time.Minute

// FetchMacroCards 返回最新一期 CPI/PPI 同环比数据，封装成 4 张 MacroCard
func FetchMacroCards() ([]models.MacroCard, error) {
	macroCache.Lock()
	if time.Since(macroCache.fetched) < macroCacheTTL && len(macroCache.data) > 0 {
		data := macroCache.data
		macroCache.Unlock()
		return data, nil
	}
	macroCache.Unlock()

	cpi, errCPI := fetchEastMoneyCPI()
	ppi, errPPI := fetchEastMoneyPPI()

	if errCPI != nil && errPPI != nil {
		return nil, fmt.Errorf("拉取宏观数据失败: cpi=%v ppi=%v", errCPI, errPPI)
	}

	cards := []models.MacroCard{
		buildCard("CPI 同比", cpi, "same"),
		buildCard("CPI 环比", cpi, "seq"),
		buildCard("PPI 同比", ppi, "same"),
		buildCard("PPI 累计同比", ppi, "accum"),
	}

	macroCache.Lock()
	macroCache.data = cards
	macroCache.fetched = time.Now()
	macroCache.Unlock()

	return cards, nil
}

type macroPoint struct {
	reportDate string
	same       float64 // 同比 %
	sequential float64 // 环比 %
	accumulate float64 // 累计同比 %
}

func buildCard(label string, p *macroPoint, kind string) models.MacroCard {
	if p == nil {
		return models.MacroCard{Label: label, Value: "--", Detail: "数据暂不可用"}
	}
	var v float64
	switch kind {
	case "same":
		v = p.same
	case "seq":
		v = p.sequential
	case "accum":
		v = p.accumulate
	}
	sign := ""
	if v > 0 {
		sign = "+"
	}
	return models.MacroCard{
		Label:  label,
		Value:  fmt.Sprintf("%s%.2f%%", sign, v),
		Detail: fmt.Sprintf("%s 国家统计局", formatReportMonth(p.reportDate)),
	}
}

func formatReportMonth(d string) string {
	if len(d) >= 7 {
		// "2026-04-30" -> "2026年04月"
		return d[:4] + "年" + d[5:7] + "月"
	}
	return d
}

// fetchEastMoneyJSON 通用东方财富 datacenter 拉取
func fetchEastMoneyJSON(reportName string, dest interface{}) error {
	url := fmt.Sprintf(
		"https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=%s&columns=ALL&pageSize=1&sortColumns=REPORT_DATE&sortTypes=-1",
		reportName,
	)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://data.eastmoney.com/")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}

// fetchEastMoneyCPI: RPT_ECONOMY_CPI 表含 NATIONAL_SAME（同比%）与 NATIONAL_SEQUENTIAL（环比%）
func fetchEastMoneyCPI() (*macroPoint, error) {
	var parsed struct {
		Result struct {
			Data []struct {
				ReportDate         string  `json:"REPORT_DATE"`
				NationalSame       float64 `json:"NATIONAL_SAME"`
				NationalSequential float64 `json:"NATIONAL_SEQUENTIAL"`
				NationalAccumulate float64 `json:"NATIONAL_ACCUMULATE"`
			} `json:"data"`
		} `json:"result"`
	}
	if err := fetchEastMoneyJSON("RPT_ECONOMY_CPI", &parsed); err != nil {
		return nil, fmt.Errorf("解析 CPI 失败: %v", err)
	}
	if len(parsed.Result.Data) == 0 {
		return nil, fmt.Errorf("CPI 未返回数据")
	}
	d := parsed.Result.Data[0]
	// NATIONAL_ACCUMULATE 是定基（100=持平），转成同比差值
	accum := d.NationalAccumulate - 100
	return &macroPoint{
		reportDate: d.ReportDate,
		same:       d.NationalSame,
		sequential: d.NationalSequential,
		accumulate: accum,
	}, nil
}

// fetchEastMoneyPPI: RPT_ECONOMY_PPI 表只有 BASE_SAME（同比%）与 BASE_ACCUMULATE（累计定基 100=持平）
// 该表没有 PPI 环比字段，环比留空（呈现层用累计同比替代）
func fetchEastMoneyPPI() (*macroPoint, error) {
	var parsed struct {
		Result struct {
			Data []struct {
				ReportDate      string  `json:"REPORT_DATE"`
				BaseSame        float64 `json:"BASE_SAME"`
				BaseAccumulate  float64 `json:"BASE_ACCUMULATE"`
			} `json:"data"`
		} `json:"result"`
	}
	if err := fetchEastMoneyJSON("RPT_ECONOMY_PPI", &parsed); err != nil {
		return nil, fmt.Errorf("解析 PPI 失败: %v", err)
	}
	if len(parsed.Result.Data) == 0 {
		return nil, fmt.Errorf("PPI 未返回数据")
	}
	d := parsed.Result.Data[0]
	accum := d.BaseAccumulate - 100
	return &macroPoint{
		reportDate: d.ReportDate,
		same:       d.BaseSame,
		accumulate: accum,
	}, nil
}
