package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuchi/cycle-stock/internal/models"
	"github.com/yuchi/cycle-stock/internal/repository"
	"github.com/yuchi/cycle-stock/internal/services"
)

// GetStocks 获取股票列表
func GetStocks(c *gin.Context) {
	stocks, err := repository.GetStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	// 获取实时行情
	codes := make([]string, len(stocks))
	for i, s := range stocks {
		codes[i] = s.Code
	}

	client := services.NewSinaQuoteClient()
	quotes, err := client.GetQuotes(codes)
	if err != nil {
		// 行情获取失败不影响列表显示
		quotes = make(map[string]models.Quote)
	}

	// 合并数据
	result := make([]map[string]interface{}, len(stocks))
	for i, stock := range stocks {
		item := map[string]interface{}{
			"stock": stock,
		}
		if quote, ok := quotes[stock.Code]; ok {
			item["quote"] = quote
		}
		result[i] = item
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": result})
}

// GetStock 获取单只股票详情
func GetStock(c *gin.Context) {
	code := c.Param("code")
	stock, err := repository.GetStockByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"ok": false, "error": "股票不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": stock})
}

// CreateStock 创建股票
func CreateStock(c *gin.Context) {
	var stock struct {
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		Sector      string  `json:"sector"`
		Industry    string  `json:"industry"`
		IsRecommend bool    `json:"isRecommend"`
		EntryPrice  float64 `json:"entryPrice"`
		HeavyPrice  float64 `json:"heavyPrice"`
		TargetPrice float64 `json:"targetPrice"`
		CoreLogic   string  `json:"coreLogic"`
	}

	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "请求格式错误"})
		return
	}

	s := &models.Stock{
		Code:        stock.Code,
		Name:        stock.Name,
		Sector:      stock.Sector,
		Industry:    stock.Industry,
		IsRecommend: stock.IsRecommend,
		EntryPrice:  stock.EntryPrice,
		HeavyPrice:  stock.HeavyPrice,
		TargetPrice: stock.TargetPrice,
		CoreLogic:   stock.CoreLogic,
	}

	if err := repository.CreateStock(s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": s})
}

// UpdateStock 更新股票
func UpdateStock(c *gin.Context) {
	code := c.Param("code")
	stock, err := repository.GetStockByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"ok": false, "error": "股票不存在"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "请求格式错误"})
		return
	}

	if name, ok := req["name"].(string); ok {
		stock.Name = name
	}
	if sector, ok := req["sector"].(string); ok {
		stock.Sector = sector
	}
	if industry, ok := req["industry"].(string); ok {
		stock.Industry = industry
	}
	if isRecommend, ok := req["isRecommend"].(bool); ok {
		stock.IsRecommend = isRecommend
	}
	if entryPrice, ok := req["entryPrice"].(float64); ok {
		stock.EntryPrice = entryPrice
	}
	if heavyPrice, ok := req["heavyPrice"].(float64); ok {
		stock.HeavyPrice = heavyPrice
	}
	if targetPrice, ok := req["targetPrice"].(float64); ok {
		stock.TargetPrice = targetPrice
	}
	if coreLogic, ok := req["coreLogic"].(string); ok {
		stock.CoreLogic = coreLogic
	}

	if err := repository.UpdateStock(stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": stock})
}

// DeleteStock 删除股票
func DeleteStock(c *gin.Context) {
	code := c.Param("code")
	if err := repository.DeleteStock(code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GetQuotes 批量获取行情
func GetQuotes(c *gin.Context) {
	codesStr := c.Query("codes")
	if codesStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "缺少codes参数"})
		return
	}

	codes := strings.Split(codesStr, ",")
	client := services.NewSinaQuoteClient()
	quotes, err := client.GetQuotes(codes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": quotes, "errors": gin.H{}})
}

// GetCompanySummaries 获取公司财务数据
func GetCompanySummaries(c *gin.Context) {
	codesStr := c.Query("codes")
	if codesStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "缺少codes参数"})
		return
	}

	codes := strings.Split(codesStr, ",")
	summaries, err := repository.GetCompanySummaries(codes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": summaries, "errors": gin.H{}})
}

// GetChart 获取K线数据
func GetChart(c *gin.Context) {
	code := c.Param("code")
	period := c.DefaultQuery("period", "daily")

	client := services.NewSinaQuoteClient()
	data, err := client.GetKLineData(code, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": data})
}

// GetValuation 获取估值数据（PE/PB Band）
func GetValuation(c *gin.Context) {
	code := c.Param("code")

	client := services.NewSinaQuoteClient()

	// 获取当前估值
	currentVal, err := client.GetCurrentValuation(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	// 获取财务数据（季度ROE/BPS）
	financialData, err := client.GetPEPBHistory(code)
	if err != nil {
		// 财务数据获取失败不影响当前估值显示
		financialData = []map[string]interface{}{}
	}

	// 获取K线数据用于计算历史PE/PB
	klineData, err := client.GetValuationHistory(code, 365)
	if err != nil {
		// K线数据获取失败不影响当前估值显示
		klineData = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"data": map[string]interface{}{
			"current":   currentVal,
			"financial": financialData,
			"kline":     klineData,
		},
	})
}
// GetValuationBand 获取估值带数据（PE/PB Band）
func GetValuationBand(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "缺少code参数"})
		return
	}

	metric := c.Query("metric")
	if metric == "" {
		metric = "pe_ttm"
	}

	yearsStr := c.Query("years")
	years := 5
	if yearsStr != "" {
		if y, err := strconv.Atoi(yearsStr); err == nil && y > 0 && y <= 10 {
			years = y
		}
	}

	// 获取估值带数据
	data, err := services.GetValuationBand(code, metric, years)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"data": map[string]interface{}{
			"band":       data,
			"updatedAt":  data.UpdatedAt,
			"latestDate": data.Points[len(data.Points)-1].Date,
		},
	})
}
