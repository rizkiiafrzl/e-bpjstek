package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UploadHistory menyimpan riwayat upload massal pekerja
type UploadHistory struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           uuid.UUID `json:"user_id" gorm:"type:uuid;index;not null"`
	FileName         string    `json:"file_name" gorm:"size:255;not null"`
	FileData         []byte    `json:"-" gorm:"type:bytea"` // Store Excel file as BLOB
	FileSize         int64     `json:"fileSize" gorm:"column:file_size;default:0"`
	TotalValid       int       `json:"totalValid" gorm:"column:total_valid;default:0"`
	TotalInvalid     int       `json:"totalInvalid" gorm:"column:total_invalid;default:0"`
	TotalData        int       `json:"totalData" gorm:"column:total_data;default:0"`
	ValidationStatus string    `json:"validationStatus" gorm:"column:validation_status;size:50;default:'Selesai'"`
	ValidationDate   time.Time `json:"validationDate" gorm:"column:validation_date;autoCreateTime"`
	DataSource       string    `json:"dataSource" gorm:"column:data_source;size:50;default:'Upload'"`
	Type             string    `json:"type" gorm:"column:type;size:50;default:'workers'"`
	ErrorFilePath    string    `json:"errorFilePath" gorm:"column:error_file_path;size:500"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UploadHistory) TableName() string { return "upload_histories" }

func (u *UploadHistory) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
