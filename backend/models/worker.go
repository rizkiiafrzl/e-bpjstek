package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Worker merepresentasikan data karyawan/peserta pada periode aktif
type Worker struct {
	ID                 uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID             uuid.UUID  `json:"user_id" gorm:"type:uuid;index;not null"`
	NIK                string     `json:"nik" gorm:"size:32;index"`
	KPJ                string     `json:"kpj" gorm:"size:32"`
	NoPegawai          string     `json:"noPegawai" gorm:"column:no_pegawai;size:64"`
	Nama               string     `json:"nama" gorm:"size:200;not null"`
	DateOfBirth        time.Time  `json:"dateOfBirth" gorm:"column:date_of_birth;type:date"`
	Upah               float64    `json:"upah" gorm:"default:0"`
	Rapel              float64    `json:"rapel" gorm:"default:0"`
	Nationality        string     `json:"nationality" gorm:"size:8;default:'WNI'"`
	PassportNo         string     `json:"passportNo" gorm:"column:passport_no;size:64"`
	PassportValidUntil *time.Time `json:"passportValidUntil" gorm:"column:passport_valid_until;type:date"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	User User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Worker) TableName() string { return "workers" }

func (w *Worker) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
