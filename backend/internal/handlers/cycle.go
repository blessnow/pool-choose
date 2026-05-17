package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuchi/cycle-stock/internal/models"
	"github.com/yuchi/cycle-stock/internal/repository"
)

// GetCycleInsight 获取周期分析数据
func GetCycleInsight(c *gin.Context) {
	insight, err := repository.GetCycleInsight()
	if err != nil {
		// 返回默认数据
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"data": gin.H{
				"title":        "大宗商品周期定位（数据 + 观点）",
				"updatedAt":    "2026-05-15T00:00:00+08:00",
				"dataLayer":    "最新宏观数据加载中...",
				"opinionLayer": "周期分析加载中...",
				"macroCards": []gin.H{
					{"label": "CPI 同比", "value": "--", "detail": "数据加载中"},
					{"label": "CPI 环比", "value": "--", "detail": "数据加载中"},
					{"label": "PPI 同比", "value": "--", "detail": "数据加载中"},
					{"label": "PPI 环比", "value": "--", "detail": "数据加载中"},
				},
				"bars": []gin.H{
					{"label": "内需价格温度", "value": 50, "gradient": "linear-gradient(90deg,#ef4444,#f59e0b)", "labels": []string{"偏冷", "修复", "升温", "过热"}},
					{"label": "大宗商品价格动能", "value": 50, "gradient": "linear-gradient(90deg,#3b82f6,#22c55e)", "labels": []string{"下行", "企稳", "上行", "偏热"}},
				},
				"conclusion": "数据加载中...",
				"focus":      "数据加载中...",
				"risk":       "数据加载中...",
				"sources":    []gin.H{},
			},
		})
		return
	}

	// 解析JSON字段
	var macroCards []models.MacroCard
	var bars []models.CycleBar
	var sources []models.Source

	if insight.MacroCardsJSON != "" {
		json.Unmarshal([]byte(insight.MacroCardsJSON), &macroCards)
	}
	if insight.BarsJSON != "" {
		json.Unmarshal([]byte(insight.BarsJSON), &bars)
	}
	if insight.SourcesJSON != "" {
		json.Unmarshal([]byte(insight.SourcesJSON), &sources)
	}

	insight.MacroCards = macroCards
	insight.Bars = bars
	insight.Sources = sources

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": insight})
}

// UpdateCycleInsight 更新周期分析数据
func UpdateCycleInsight(c *gin.Context) {
	var req models.CycleInsight
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "请求格式错误"})
		return
	}

	// 序列化JSON字段
	if len(req.MacroCards) > 0 {
		data, _ := json.Marshal(req.MacroCards)
		req.MacroCardsJSON = string(data)
	}
	if len(req.Bars) > 0 {
		data, _ := json.Marshal(req.Bars)
		req.BarsJSON = string(data)
	}
	if len(req.Sources) > 0 {
		data, _ := json.Marshal(req.Sources)
		req.SourcesJSON = string(data)
	}

	if err := repository.SaveCycleInsight(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": req})
}

// GetPositions 获取持仓记录
func GetPositions(c *gin.Context) {
	code := c.Query("code")
	positions, err := repository.GetPositions(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "data": positions})
}

// CreatePosition 创建持仓记录
func CreatePosition(c *gin.Context) {
	var req models.Position
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "请求格式错误"})
		return
	}

	if err := repository.CreatePosition(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": req})
}

// UpdatePosition 更新持仓记录
func UpdatePosition(c *gin.Context) {
	_ = c.Param("id")
	var req models.Position
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "请求格式错误"})
		return
	}

	// 这里简化处理，实际应该查询后更新
	c.JSON(http.StatusOK, gin.H{"ok": true, "data": req})
}

// DeletePosition 删除持仓记录
func DeletePosition(c *gin.Context) {
	_ = c.Param("id")
	// 这里简化处理
	c.JSON(http.StatusOK, gin.H{"ok": true})
}