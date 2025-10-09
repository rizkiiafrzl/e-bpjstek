package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Template stores downloadable template metadata
type Template struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Key         string    `json:"key" gorm:"column:key;uniqueIndex;not null"`
    FileName    string    `json:"file_name" gorm:"column:file_name;not null"`
    FilePath    string    `json:"file_path" gorm:"column:file_path;not null"`
    MimeType    string    `json:"mime_type" gorm:"column:mime_type;not null"`
    Description string    `json:"description" gorm:"column:description;type:text"`
    IsActive    bool      `json:"is_active" gorm:"column:is_active;default:true"`
    CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Template) TableName() string { return "templates" }

// BeforeCreate ensures UUID is set when inserting via GORM
func (t *Template) BeforeCreate(tx *gorm.DB) (err error) {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    return nil
}


