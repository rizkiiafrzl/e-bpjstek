package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Struct untuk response API
type Worker struct {
	ID                 string    `json:"id"`
	NoPegawai          string    `json:"noPegawai"`
	Nama               string    `json:"nama"`
	NIK                string    `json:"nik"`
	TeleponAreaRumah   string    `json:"teleponAreaRumah"`
	TeleponRumah       string    `json:"teleponRumah"`
	TeleponAreaKantor  string    `json:"teleponAreaKantor"`
	TeleponKantor      string    `json:"teleponKantor"`
	TeleponExtKantor   string    `json:"teleponExtKantor"`
	Handphone          string    `json:"handphone"`
	Email              string    `json:"email"`
	TempatLahir        string    `json:"tempatLahir"`
	JenisKelamin       string    `json:"jenisKelamin"`
	GolDarah           string    `json:"golDarah"`
	StatusPegawai      string    `json:"statusPegawai"`
	Upah               float64   `json:"upah"`
	NPWP               string    `json:"npwp"`
	JenisIdentitas     string    `json:"jenisIdentitas"`
	TanggalKepesertaan string    `json:"tanggalKepesertaan"`
	KodeNegara         string    `json:"kodeNegara"`
	CreatedAt          time.Time `json:"created_at"`
}

type UploadHistory struct {
	ID               string    `json:"id"`
	FileName         string    `json:"fileName"`
	TotalData        int       `json:"totalData"`
	TotalValid       int       `json:"totalValid"`
	TotalInvalid     int       `json:"totalInvalid"`
	ValidationStatus string    `json:"validationStatus"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"created_at"`
}

func main() {
	baseURL := "http://localhost:8080/api/v1"
	token := "YOUR_TOKEN_HERE" // Ganti dengan token yang valid

	fmt.Println("🔍 Checking uploaded data...")

	// 1. Cek data workers terbaru
	fmt.Println("\n=== DATA WORKERS TERBARU ===")
	workers, err := getWorkers(baseURL, token)
	if err != nil {
		fmt.Printf("❌ Error getting workers: %v\n", err)
		return
	}

	fmt.Printf("📊 Total workers: %d\n", len(workers))
	fmt.Println("\n📋 Data terbaru (5 record):")
	for i, worker := range workers {
		if i >= 5 {
			break
		}
		fmt.Printf("\n%d. %s (%s)\n", i+1, worker.Nama, worker.NoPegawai)
		fmt.Printf("   NIK: %s\n", worker.NIK)
		fmt.Printf("   Email: %s\n", worker.Email)
		fmt.Printf("   Telepon Area Rumah: %s\n", worker.TeleponAreaRumah)
		fmt.Printf("   Telepon Kantor: %s\n", worker.TeleponKantor)
		fmt.Printf("   Jenis Identitas: %s\n", worker.JenisIdentitas)
		fmt.Printf("   Tanggal Kepesertaan: %s\n", worker.TanggalKepesertaan)
		fmt.Printf("   Kode Negara: %s\n", worker.KodeNegara)
		fmt.Printf("   Upah: %.0f\n", worker.Upah)
		fmt.Printf("   Created: %s\n", worker.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// 2. Cek upload history
	fmt.Println("\n=== UPLOAD HISTORY ===")
	histories, err := getUploadHistory(baseURL, token)
	if err != nil {
		fmt.Printf("❌ Error getting upload history: %v\n", err)
		return
	}

	fmt.Printf("📊 Total uploads: %d\n", len(histories))
	fmt.Println("\n📋 Upload history terbaru:")
	for i, history := range histories {
		if i >= 5 {
			break
		}
		fmt.Printf("\n%d. %s\n", i+1, history.FileName)
		fmt.Printf("   Type: %s\n", history.Type)
		fmt.Printf("   Total Data: %d\n", history.TotalData)
		fmt.Printf("   Valid: %d\n", history.TotalValid)
		fmt.Printf("   Invalid: %d\n", history.TotalInvalid)
		fmt.Printf("   Status: %s\n", history.ValidationStatus)
		fmt.Printf("   Uploaded: %s\n", history.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// 3. Cek field mapping
	fmt.Println("\n=== FIELD MAPPING CHECK ===")
	checkFieldMapping(workers)

	fmt.Println("\n✅ Data check completed!")
}

func getWorkers(baseURL, token string) ([]Worker, error) {
	req, err := http.NewRequest("GET", baseURL+"/workers", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var workers []Worker
	err = json.Unmarshal(body, &workers)
	return workers, err
}

func getUploadHistory(baseURL, token string) ([]UploadHistory, error) {
	req, err := http.NewRequest("GET", baseURL+"/workers/upload-history", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var histories []UploadHistory
	err = json.Unmarshal(body, &histories)
	return histories, err
}

func checkFieldMapping(workers []Worker) {
	fmt.Println("🔍 Checking field mapping from Excel to Database:")

	// Cek field yang seharusnya ada dari Excel upload
	fieldChecks := map[string]int{
		"TeleponAreaRumah":   0,
		"TeleponAreaKantor":  0,
		"TeleponKantor":      0,
		"TeleponExtKantor":   0,
		"JenisIdentitas":     0,
		"TanggalKepesertaan": 0,
		"KodeNegara":         0,
	}

	for _, worker := range workers {
		if worker.TeleponAreaRumah != "" {
			fieldChecks["TeleponAreaRumah"]++
		}
		if worker.TeleponAreaKantor != "" {
			fieldChecks["TeleponAreaKantor"]++
		}
		if worker.TeleponKantor != "" {
			fieldChecks["TeleponKantor"]++
		}
		if worker.TeleponExtKantor != "" {
			fieldChecks["TeleponExtKantor"]++
		}
		if worker.JenisIdentitas != "" {
			fieldChecks["JenisIdentitas"]++
		}
		if worker.TanggalKepesertaan != "" {
			fieldChecks["TanggalKepesertaan"]++
		}
		if worker.KodeNegara != "" {
			fieldChecks["KodeNegara"]++
		}
	}

	for field, count := range fieldChecks {
		if count > 0 {
			fmt.Printf("✅ %s: %d records\n", field, count)
		} else {
			fmt.Printf("❌ %s: No data found\n", field)
		}
	}
}







