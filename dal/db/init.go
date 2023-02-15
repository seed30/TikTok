package db

import (
	"github.com/seed30/TikTok/pkg/consts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
	"time"
)

var DB *gorm.DB

func Init() {
	var err error

	gormlogrus := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			Colorful:      false,
			LogLevel:      logger.Info,
		},
	)

	if DB, err = gorm.Open(mysql.Open(consts.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt: true,
			Logger:      gormlogrus,
		},
	); err != nil {
		panic(err)
	}

	if err = DB.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	if err = DB.AutoMigrate(&Video{}); err != nil {
		panic(err)
	}
}
