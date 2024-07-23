package database

import (
	"fmt"
	"os"
	"time"

	extraClausePlugin "github.com/WinterYukky/gorm-extra-clause-plugin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %s", err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %s", err))
	}
	db.Use(extraClausePlugin.New())

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	fmt.Println("DB_CONNECTED SUCCESSFULLY")
	db.Logger = logger.Default.LogMode((logger.Info))

	return db
}
