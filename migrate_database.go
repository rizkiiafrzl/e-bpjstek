package main

import (
	"farm-management-backend/database"
	"farm-management-backend/models"
	"fmt"
	"log"
)

func main() {
	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migration
	fmt.Println("🔄 Running database migration...")
	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Check if all required fields exist
	fmt.Println("🔍 Checking database structure...")

	// Test create a sample worker to verify all fields work
	sampleWorker := models.Worker{
		// Basic Info
		NIK:       "3175051405900001",
		KPJ:       "12345678901234567890",
		NoPegawai: "TEST001",
		Nama:      "Test Worker",

		// Personal Info
		TempatLahir:  "Jakarta",
		IbuKandung:   "Test Ibu",
		JenisKelamin: "Laki-laki",
		GolDarah:     "O",
		StatusKawin:  "Belum Kawin",

		// Employment Info
		StatusPegawai:   "Tetap",
		LokasiPekerjaan: "Test Location",

		// Contact Info (NEW FIELDS)
		TeleponAreaRumah:  "021",
		TeleponRumah:      "7654321",
		TeleponAreaKantor: "021",
		TeleponKantor:     "021-1234",
		TeleponExtKantor:  "1234",
		Handphone:         "081234567890",
		Email:             "test@example.com",

		// Address Info
		Alamat:    "Test Address",
		Kabupaten: "Test Kabupaten",
		KodePos:   "10130",

		// Additional Info (NEW FIELDS)
		NPWP:               "012345678901234",
		JenisIdentitas:     "KTP",
		MasaLakuIdentitas:  "2016-05-14",
		SuratMenyuratKe:    "Test Surat",
		TanggalKepesertaan: "2020-07-01",
		KodeNegara:         "ID",

		// Financial Info
		Upah:        5000000,
		Rapel:       0,
		Nationality: "WNI",
	}

	// Try to create the worker (this will test if all fields are properly mapped)
	result := db.Create(&sampleWorker)
	if result.Error != nil {
		fmt.Printf("❌ Error creating sample worker: %v\n", result.Error)
		fmt.Println("This indicates that some fields are missing in the database schema.")
	} else {
		fmt.Println("✅ Sample worker created successfully!")
		fmt.Println("All required fields are properly mapped in the database.")

		// Clean up test data
		db.Delete(&sampleWorker)
		fmt.Println("🧹 Test data cleaned up.")
	}

	fmt.Println("✅ Database migration and verification completed!")
}







