package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportPeriod merekam periode pelaporan per pengguna per bulan
type ReportPeriod struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;index;not null"`
	Year   int       `json:"year" gorm:"not null"`
	Month  int       `json:"month" gorm:"not null"`
	Status string    `json:"status" gorm:"type:varchar(20);not null;default:'Draft'"`
	// Summary fields for dashboard
	TotalTk    int       `json:"totalTk" gorm:"column:total_tk;default:0"`
	TotalIuran float64   `json:"totalIuran" gorm:"column:total_iuran;default:0"`
	TotalDenda float64   `json:"totalDenda" gorm:"column:total_denda;default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (ReportPeriod) TableName() string { return "report_periods" }

func (rp *ReportPeriod) BeforeCreate(tx *gorm.DB) (err error) {
	if rp.ID == uuid.Nil {
		rp.ID = uuid.New()
	}
	return nil
}
