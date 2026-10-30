package databases

import (
	"chatsmith-mini/configs"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type PostgresDB = *gorm.DB

func NewPostgres(cfg configs.DatabaseConfig) (PostgresDB, error) {
	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{Logger: logger.NewGORM()})
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithoutQueryVariables())); err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}
