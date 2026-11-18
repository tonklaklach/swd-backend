package database

import (
	"embed"
	"fmt"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed schema.sql
var schema embed.FS

func Open(dsn string) (*gorm.DB, error) {
	separator := "?"
	if strings.Contains(dsn, "?") {
		separator = "&"
	}
	db, err := gorm.Open(sqlite.Open(dsn+separator+"_foreign_keys=on&_busy_timeout=5000"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}
	pool, err := db.DB()
	if err != nil {
		return nil, err
	}
	pool.SetMaxOpenConns(1)
	data, err := schema.ReadFile("schema.sql")
	if err != nil {
		pool.Close()
		return nil, err
	}
	if err = db.Exec(string(data)).Error; err != nil {
		pool.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	return db, nil
}
func Close(db *gorm.DB) error {
	pool, err := db.DB()
	if err != nil {
		return err
	}
	return pool.Close()
}
