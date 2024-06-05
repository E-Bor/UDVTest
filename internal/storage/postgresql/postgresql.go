package postgresql

import (
	"StackService/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

type JSONData struct {
	ID   uint   `gorm:"primaryKey"`
	Data []byte `gorm:"type:bytea"` // Field to store raw json
}

type PostgresDB struct {
	*gorm.DB
	logger *slog.Logger
}

func InitDB(
	logger *slog.Logger,
	cfg *config.Config,
) *PostgresDB {

	logger.Info("Initializing connect PostgresDB")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Dbname,
		cfg.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to the database: %v", err)
		return nil
	}

	err = db.AutoMigrate(&JSONData{})
	if err != nil {
		logger.Error("Failed to migrate database: %v", err)
		return nil
	}

	return &PostgresDB{DB: db, logger: logger}
}

func (db *PostgresDB) InsertMessage(json []byte) error {
	jsonData := JSONData{
		Data: json,
	}
	result := db.Create(&jsonData)
	return result.Error
}

func (db *PostgresDB) PopMessage() error {
	var jsonData JSONData
	result := db.Order("id desc").First(&jsonData)
	result = db.Delete(&jsonData)
	return result.Error
}

func (db *PostgresDB) GetAllMessages() ([][]byte, error) {
	db.logger.Info("searching messages in pg")
	var jsonDataList []JSONData
	result := db.Order("id asc").Find(&jsonDataList)
	if result.Error != nil {
		db.logger.Error(result.Error.Error())
		return nil, result.Error
	}

	var messages [][]byte
	for _, jsonData := range jsonDataList {
		messages = append(messages, jsonData.Data)
	}
	return messages, nil
}
