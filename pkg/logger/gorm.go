package logger

import (
	"log"
	"os"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

const slowQueryThreshold = time.Second

func NewGORM() gormlogger.Interface {
	return gormlogger.New(
		log.New(os.Stderr, "", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             slowQueryThreshold,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gormlogger.Warn,
		},
	)
}
