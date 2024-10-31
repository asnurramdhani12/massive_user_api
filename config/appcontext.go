package config

import (
	"context"
	"user_api/utils/logger"

	"gorm.io/gorm"
)

type appContext struct {
	cfg Configuration
	db  *gorm.DB
}

var appcontext appContext

func NewAppContext(ctx context.Context) error {
	logger.New(ctx)

	cfg, err := NewConfig(ctx)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to load config: %v", err)
		return err
	}

	db, err := NewPostgreSql(ctx, *cfg)
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to connect database: %v", err)
		return err
	}

	appcontext = appContext{
		cfg: *cfg,
		db:  db,
	}

	return nil
}

func Config() Configuration {
	return appcontext.cfg
}

func DB() *gorm.DB {
	return appcontext.db
}
