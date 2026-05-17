package models

import "time"

// Stock 股票基础信息
type Stock struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"uniqueIndex;size:10" json:"code"`
	Name        string    `gorm:"size:50" json:"name"`
	Sector      string    `gorm:"size:50" json:"sector"`
	Industry    string    `gorm:"size:100" json:"industry"`
	IsRecommend bool      `json:"isRecommend"`
	EntryPrice  float64   `json:"entryPrice"`
	HeavyPrice  float64   `json:"heavyPrice"`
	TargetPrice float64   `json:"targetPrice"`
	CoreLogic   string    `gorm:"type:text" json:"coreLogic"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Quote 实时行情
type Quote struct {
	Code      string    `gorm:"index;size:10" json:"stockCode"`
	Open      float64   `json:"open"`
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Volume    int64     `json:"volume"`
	Amount    float64   `json:"amount"`
	Change    float64   `json:"change"`
	QuoteTime time.Time `json:"quote_time"`
	PrevClose float64   `json:"prev_close"`
}

// CompanySummary 公司财务数据
type CompanySummary struct {
	ID             uint    `gorm:"primaryKey" json:"-"`
	Code           string  `gorm:"uniqueIndex;size:10" json:"code"`
	SyncedAt       string  `json:"syncedAt"`
	TotalShares    int64   `json:"totalShares"`
	NetAssets      int64   `json:"netAssets"`
	AnnualProfit   int64   `json:"annualProfit"`
	LatestDividend float64 `json:"latestDividendPer10"`
	MarketCap      float64 `json:"marketCap"`
	PE             float64 `json:"pe"`
	PB             float64 `json:"pb"`
	ROE            float64 `json:"roe"`
	DividendYield  float64 `json:"dividendYield"`
	Sector         string  `gorm:"size:100" json:"sector"`
}

// MacroCard 宏观经济卡片
type MacroCard struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Detail string `json:"detail"`
}

// CycleBar 周期条形图
type CycleBar struct {
	Label    string   `json:"label"`
	Value    int      `json:"value"`
	Gradient string   `json:"gradient"`
	Labels   []string `json:"labels"`
}

// Source 数据来源
type Source struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// CycleInsight 周期分析数据
type CycleInsight struct {
	ID           uint         `gorm:"primaryKey" json:"-"`
	Title        string       `json:"title"`
	UpdatedAt    string       `json:"updatedAt"`
	DataLayer    string       `gorm:"type:text" json:"dataLayer"`
	OpinionLayer string       `gorm:"type:text" json:"opinionLayer"`
	MacroCards   []MacroCard  `gorm:"-" json:"macroCards"`
	Bars         []CycleBar   `gorm:"-" json:"bars"`
	Conclusion   string       `gorm:"type:text" json:"conclusion"`
	Focus        string       `gorm:"type:text" json:"focus"`
	Risk         string       `gorm:"type:text" json:"risk"`
	Sources      []Source     `gorm:"-" json:"sources"`
	MacroCardsJSON string     `gorm:"type:text" json:"-"`
	BarsJSON      string      `gorm:"type:text" json:"-"`
	SourcesJSON   string      `gorm:"type:text" json:"-"`
}

// Position 持仓记录
type Position struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StockCode string    `gorm:"index;size:10" json:"stockCode"`
	Type      string    `gorm:"size:10" json:"type"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Date      time.Time `json:"date"`
	Note      string    `gorm:"type:text" json:"note"`
}

// Session 会话
type Session struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	Token     string    `gorm:"uniqueIndex" json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
