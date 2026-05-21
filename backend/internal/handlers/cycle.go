package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuchi/cycle-stock/internal/models"
	"github.com/yuchi/cycle-stock/internal/repository"
	"github.com/yuchi/cycle-stock/internal/services"
)

// GetCycleInsight 获取周期分析数据
func GetCycleInsight(c *gin.Context) {
	// 拉取最新 CPI/PPI，失败时回退到 nil，由下方逻辑兜底
	liveMacro, macroErr := services.FetchMacroCards()

	insight, err := repository.GetCycleInsight()
	if err != nil {
		// DB 无数据时的默认响应：宏观卡片优先用实时数据，其它字段用占位
		macroCards := liveMacro
		if macroErr != nil || len(macroCards) == 0 {
			macroCards = []models.MacroCard{
				{Label: "CPI 同比", Value: "--", Detail: "数据加载中"},
				{Label: "CPI 环比", Value: "--", Detail: "数据加载中"},
				{Label: "PPI 同比", Value: "--", Detail: "数据加载中"},
				{Label: "PPI 环比", Value: "--", Detail: "数据加载中"},
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"data": gin.H{
				"title":        "大宗商品周期定位（数据 + 观点）",
				"updatedAt":    "2026-05-15T00:00:00+08:00",
				"dataLayer":    "最新宏观数据加载中...",
				"opinionLayer": "周期分析加载中...",
				"macroCards":   macroCards,
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

	// 实时宏观数据覆盖 DB 中的旧 macroCards；拉取失败才回退到 DB
	if macroErr == nil && len(liveMacro) > 0 {
		macroCards = liveMacro
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