package config

import (
	"context"
	"fmt"
	"user_api/app/v1/contract"
	"user_api/utils/logger"

	gormLogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgreSql(ctx context.Context, cfg Configuration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable timezone=Asia/Jakarta", cfg.PostgreSql.Host, cfg.PostgreSql.Username, cfg.PostgreSql.Password, cfg.PostgreSql.Database, cfg.PostgreSql.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to connect database: %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(&contract.User{}); err != nil {
		logger.GetLogger(ctx).Errorf("failed to migrate database: %v", err)
		return nil, err
	}

	return db, nil
}
