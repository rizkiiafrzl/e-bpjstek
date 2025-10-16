package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Worker merepresentasikan data karyawan/peserta pada periode aktif
type Worker struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;index;not null"`
	NIK          string    `json:"nik" gorm:"size:32;index"`
	KPJ          string    `json:"kpj" gorm:"column:kpj;size:32"`
	NoPegawai    string    `json:"noPegawai" gorm:"column:no_pegawai;size:64"`
	Nama         string    `json:"nama" gorm:"size:200;not null"`
	TempatLahir  string    `json:"tempatLahir" gorm:"column:tempat_lahir;size:120"`
	DateOfBirth  time.Time `json:"dateOfBirth" gorm:"column:date_of_birth;type:date"`
	IbuKandung   string    `json:"ibuKandung" gorm:"column:ibu_kandung;size:200"`
	JenisKelamin string    `json:"jenisKelamin" gorm:"column:jenis_kelamin;size:1"`
	GolDarah     string    `json:"golDarah" gorm:"column:gol_darah;size:3"`
	StatusKawin  string    `json:"statusKawin" gorm:"column:status_kawin;size:20"`

	StatusPegawai       string     `json:"statusPegawai" gorm:"column:status_pegawai;size:16"`
	TanggalAwalBekerja  *time.Time `json:"tanggalAwalBekerja" gorm:"column:tanggal_awal_bekerja;type:date"`
	TanggalAkhirKontrak *time.Time `json:"tanggalAkhirKontrak" gorm:"column:tanggal_akhir_kontrak;type:date"`
	LokasiPekerjaan     string     `json:"lokasiPekerjaan" gorm:"column:lokasi_pekerjaan;size:120"`

	Upah               float64    `json:"upah" gorm:"default:0"`
	Rapel              float64    `json:"rapel" gorm:"default:0"`
	Nationality        string     `json:"nationality" gorm:"size:8;default:'WNI'"`
	PassportNo         string     `json:"passportNo" gorm:"column:passport_no;size:64"`
	PassportValidUntil *time.Time `json:"passportValidUntil" gorm:"column:passport_valid_until;type:date"`

	// Contact Info
	TeleponAreaRumah  string `json:"teleponAreaRumah" gorm:"column:telepon_area_rumah;size:8"`
	TeleponRumah      string `json:"teleponRumah" gorm:"column:telepon_rumah;size:32"`
	TeleponAreaKantor string `json:"teleponAreaKantor" gorm:"column:telepon_area_kantor;size:8"`
	TeleponKantor     string `json:"teleponKantor" gorm:"column:telepon_kantor;size:32"`
	TeleponExtKantor  string `json:"teleponExtKantor" gorm:"column:telepon_ext_kantor;size:16"`
	Handphone         string `json:"handphone" gorm:"column:handphone;size:32"`
	Email             string `json:"email" gorm:"column:email;size:200"`

	// Address Info
	Alamat    string `json:"alamat" gorm:"column:alamat;type:text"`
	Kabupaten string `json:"kabupaten" gorm:"column:kabupaten;size:120"`
	KodePos   string `json:"kodePos" gorm:"column:kode_pos;size:10"`

	// Additional Info
	NPWP               string    `json:"npwp" gorm:"column:npwp;size:32"`
	JenisIdentitas     string    `json:"jenisIdentitas" gorm:"column:jenis_identitas;size:16"`
	MasaLakuIdentitas  string    `json:"masaLakuIdentitas" gorm:"column:masa_laku_identitas;size:16"`
	SuratMenyuratKe    string    `json:"suratMenyuratKe" gorm:"column:surat_menyurat_ke;size:200"`
	TanggalKepesertaan string    `json:"tanggalKepesertaan" gorm:"column:tanggal_kepesertaan;size:16"`
	KodeNegara         string    `json:"kodeNegara" gorm:"column:kode_negara;size:4"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Worker) TableName() string { return "workers" }

func (w *Worker) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
