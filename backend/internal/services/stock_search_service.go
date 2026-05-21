package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// StockSearchResult 单条搜索结果（精简字段）
type StockSearchResult struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	PinYin string `json:"pinyin"`
	Market string `json:"market"` // "SH" / "SZ" / "BJ"
}

// StockInfo 详细信息（用于配置页一键回填）
type StockInfo struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Sector   string `json:"sector"`   // 行业细分（EM2016 最后一段，如 "纯碱"）
	Industry string `json:"industry"` // 行业全路径（"纯碱 / 化学原料 / 基础化工"）
	Market   string `json:"market"`
}

// GetStockInfo 从东方财富 F10 拉取一只 A 股的行业信息
func GetStockInfo(code string) (*StockInfo, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("股票代码必须 6 位数字")
	}
	prefix := "SZ"
	if code[0] == '6' || code[0] == '5' {
		prefix = "SH"
	} else if code[0] == '4' || code[0] == '8' || code[0] == '9' {
		prefix = "BJ"
	}

	endpoint := fmt.Sprintf("https://emweb.eastmoney.com/PC_HSF10/CompanySurvey/PageAjax?code=%s%s", prefix, code)
	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://emweb.eastmoney.com/")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Jbzl []struct {
			SecurityNameAbbr string `json:"SECURITY_NAME_ABBR"`
			EM2016           string `json:"EM2016"`
			TradeMarket      string `json:"TRADE_MARKET"`
		} `json:"jbzl"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("解析公司信息失败: %v", err)
	}
	if len(parsed.Jbzl) == 0 {
		return nil, fmt.Errorf("未获取到 %s 的公司信息", code)
	}
	d := parsed.Jbzl[0]

	// EM2016 形如 "基础化工-化学原料-纯碱"。
	// sector = 最末一级；industry = 反转后用 " / " 连接，与既有 seed 风格一致
	var sector, industry string
	if d.EM2016 != "" {
		parts := strings.Split(d.EM2016, "-")
		if n := len(parts); n > 0 {
			sector = parts[n-1]
			rev := make([]string, n)
			for i, p := range parts {
				rev[n-1-i] = p
			}
			industry = strings.Join(rev, " / ")
		}
	}

	return &StockInfo{
		Code:     code,
		Name:     d.SecurityNameAbbr,
		Sector:   sector,
		Industry: industry,
		Market:   prefix,
	}, nil
}

// SearchStocks 调用东方财富 suggest 接口，按 q（拼音首字母 / 中文名 / 股票代码）模糊匹配，仅返回 A 股
func SearchStocks(q string, limit int) ([]StockSearchResult, error) {
	if q == "" {
		return []StockSearchResult{}, nil
	}
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	endpoint := fmt.Sprintf(
		"https://searchapi.eastmoney.com/api/suggest/get?input=%s&type=14&count=%d&token=D43BF722C8E33BDC906FB84D85E326E8",
		url.QueryEscape(q), limit*2, // 多拉一些，过滤后再截断
	)

	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://www.eastmoney.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		QuotationCodeTable struct {
			Data []struct {
				Code             string `json:"Code"`
				Name             string `json:"Name"`
				PinYin           string `json:"PinYin"`
				Classify         string `json:"Classify"`
				MarketType       string `json:"MarketType"`
				SecurityTypeName string `json:"SecurityTypeName"`
			} `json:"Data"`
		} `json:"QuotationCodeTable"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("解析搜索结果失败: %v", err)
	}

	out := make([]StockSearchResult, 0, limit)
	for _, d := range parsed.QuotationCodeTable.Data {
		// 只保留 A 股（沪 A / 深 A / 北 A），过滤掉港股、美股、基金、债券、可转债等
		if d.Classify != "AStock" {
			continue
		}
		if len(d.Code) != 6 {
			continue
		}
		var market string
		switch d.MarketType {
		case "1":
			market = "SH"
		case "0":
			market = "SZ"
		case "2":
			market = "BJ"
		default:
			market = ""
		}
		out = append(out, StockSearchResult{
			Code:   d.Code,
			Name:   d.Name,
			PinYin: d.PinYin,
			Market: market,
		})
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}
