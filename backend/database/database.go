package database

import (
	"farm-management-backend/config"
	"farm-management-backend/models"
	"fmt"
	"log"
	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

// Connect establishes database connection
func Connect() (*DB, error) {
	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5433")
	user := config.GetEnv("DB_USER", "postgres")
	password := config.GetEnv("DB_PASSWORD", "password")
	dbname := config.GetEnv("DB_NAME", "e-bpjstk")
	sslmode := config.GetEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return &DB{db}, nil
}

// Migrate runs database migrations
func Migrate(db *DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Template{},
		&models.ReportPeriod{},
		&models.Worker{},
		&models.UploadHistory{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

// SeedTemplates inserts initial template metadata if not present
func SeedTemplates(db *DB) error {
	type pair struct{ key, file string }
	files := []pair{
		{"upah", "templates/template_upah_24101780.xlsx"},
		{"tk", "templates/template_tk_24101780.xlsx"},
		{"tk_na", "templates/template_tk_na_24101780.xlsx"},
		{"koreksi_tk", "templates/template_koreksi_tk_24101780.xlsx"},
	}
	for _, p := range files {
		var count int64
		db.Model(&models.Template{}).Where("key = ?", p.key).Count(&count)
		if count == 0 {
			fn := filepath.Base(p.file)
			tmpl := models.Template{
				Key:      p.key,
				FileName: fn,
				FilePath: filepath.ToSlash(p.file),
				MimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				IsActive: true,
			}
			if err := db.Create(&tmpl).Error; err != nil {
				return fmt.Errorf("failed to seed template %s: %w", p.key, err)
			}
		}
	}
	return nil
}
