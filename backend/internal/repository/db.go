package repository

import (
	"github.com/yuchi/cycle-stock/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(path string) error {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db

	// 自动迁移
	err = DB.AutoMigrate(
		&models.Stock{},
		&models.CompanySummary{},
		&models.CycleInsight{},
		&models.Position{},
		&models.Session{},
	)
	if err != nil {
		return err
	}

	seedStocks()
	return nil
}

func GetStocks() ([]models.Stock, error) {
	var stocks []models.Stock
	err := DB.Find(&stocks).Error
	return stocks, err
}

func GetStockByCode(code string) (*models.Stock, error) {
	var stock models.Stock
	err := DB.Where("code = ?", code).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func CreateStock(stock *models.Stock) error {
	return DB.Create(stock).Error
}

func UpdateStock(stock *models.Stock) error {
	return DB.Save(stock).Error
}

func DeleteStock(code string) error {
	return DB.Where("code = ?", code).Delete(&models.Stock{}).Error
}

func GetCompanySummaries(codes []string) (map[string]models.CompanySummary, error) {
	var summaries []models.CompanySummary
	err := DB.Where("code IN ?", codes).Find(&summaries).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]models.CompanySummary)
	for _, s := range summaries {
		result[s.Code] = s
	}
	return result, nil
}

func SaveCompanySummary(summary *models.CompanySummary) error {
	return DB.Save(summary).Error
}

func GetCycleInsight() (*models.CycleInsight, error) {
	var insight models.CycleInsight
	err := DB.First(&insight).Error
	if err != nil {
		return nil, err
	}
	return &insight, nil
}

func SaveCycleInsight(insight *models.CycleInsight) error {
	return DB.Save(insight).Error
}

func GetPositions(code string) ([]models.Position, error) {
	var positions []models.Position
	if code != "" {
		err := DB.Where("stock_code = ?", code).Find(&positions).Error
		return positions, err
	}
	err := DB.Find(&positions).Error
	return positions, err
}

func CreatePosition(position *models.Position) error {
	return DB.Create(position).Error
}

func UpdatePosition(position *models.Position) error {
	return DB.Save(position).Error
}

func DeletePosition(id uint) error {
	return DB.Delete(&models.Position{}, id).Error
}

func CreateSession(session *models.Session) error {
	return DB.Create(session).Error
}

func GetSession(token string) (*models.Session, error) {
	var session models.Session
	err := DB.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteSession(token string) error {
	return DB.Where("token = ?", token).Delete(&models.Session{}).Error
}