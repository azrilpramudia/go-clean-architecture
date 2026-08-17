package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type AppConfig struct {
	Database struct {
		Host string `json:"host"`
		Port int `json:"port"`
		User string `json:"user"`
		Password string `json:"password"`
		Name string `json:"name"`
	} `json:"database"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	JWT struct {
		Secret string `json:"secret"`
		ExpiryHours int `json:"expiry_hours"`
	} `json:"jwt"`
	Notification struct {
		BaseURL string `json:"base_url"`
	} `json:"notification"`
}

func Load(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	cfg := new(AppConfig)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return cfg, nil
}

func NewDatabase(cfg *AppConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", 
			cfg.Database.User, cfg.Database.Password, 
			cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return db
}