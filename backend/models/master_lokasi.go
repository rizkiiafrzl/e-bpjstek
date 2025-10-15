package models

import (
	"time"
)

// MasterLokasi maps to existing table with integer PK 'id'.
type MasterLokasi struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Kode      string    `json:"kode" gorm:"column:kode;size:10;uniqueIndex;not null"`
	Nama      string    `json:"nama" gorm:"column:nama;size:120;not null"`
	Alamat    string    `json:"alamat" gorm:"column:alamat;type:text"`
	IsActive  bool      `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (MasterLokasi) TableName() string { return "master_lokasi" }
