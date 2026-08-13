package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

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
}

func Load(path string) *AppConfig {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	cfg := new(AppConfig)
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}
	return cfg
}

func NewDatabase(cfg *AppConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", 
			cfg.Database.User, cfg.Database.Password, 
			cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return db
}