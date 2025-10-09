<template>
  <div class="upload-worker-page">
    <v-container>
      <!-- Breadcrumb -->
      <div class="d-flex align-center mb-4">
        <v-btn
          variant="text"
          color="primary"
          prepend-icon="mdi-arrow-left"
          @click="goBack"
          class="mr-2"
        >
          Kembali
        </v-btn>
        <v-breadcrumbs :items="breadcrumbs" class="pa-0">
          <template v-slot:divider>
            <v-icon>mdi-chevron-right</v-icon>
          </template>
        </v-breadcrumbs>
      </div>

      <!-- Upload Section -->
      <v-card class="mb-6" elevation="2">
        <v-card-title class="bg-primary text-white">
          <v-icon left>mdi-upload</v-icon>
          UPLOAD TENAGA KERJA BARU
        </v-card-title>
        <v-card-text class="pa-6">
          <!-- Instructions -->
          <v-alert type="info" variant="tonal" class="mb-4">
            <div class="text-body-2">
              <strong>Petunjuk Upload:</strong>
              <ul class="mt-2">
                <li>
                  File upload harus berupa file Microsoft Excel (.xls / .xlsx) dengan maksimal
                  10.000 baris
                </li>
                <li>
                  Terdapat 2 pilihan upload: <strong>"Upload TK Mendaftar"</strong> untuk tenaga
                  kerja baru, dan <strong>"Upload TK Lanjutan"</strong> untuk tenaga kerja yang
                  sudah menjadi peserta dan memiliki Kartu Peserta BPJS Ketenagakerjaan (KPJ)
                </li>
                <li>
                  File yang diupload harus berdasarkan template yang disediakan, yang dapat
                  didownload dengan mengklik <strong>"Download tk mendaftar"</strong> atau
                  <strong>"Download tk lanjutan"</strong>. Nama file tidak boleh diubah
                </li>
                <li>Format data dalam file Excel harus berupa teks</li>
              </ul>
            </div>
          </v-alert>

          <!-- Upload Form -->
          <v-form @submit.prevent="handleUpload" ref="uploadForm">
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="uploadType"
                  :items="uploadTypes"
                  label="Pilihan Upload"
                  variant="outlined"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-file-input
                  v-model="selectedFile"
                  label="Choose files To Upload"
                  variant="outlined"
                  accept=".xls,.xlsx"
                  show-size
                  prepend-icon="mdi-file-excel"
                  @change="onFileSelected"
                />
              </v-col>
            </v-row>

            <v-row class="mt-4">
              <v-col cols="12" class="d-flex gap-4">
                <v-btn
                  type="submit"
                  color="primary"
                  size="large"
                  :loading="isUploading"
                  :disabled="!selectedFile || !uploadType"
                  prepend-icon="mdi-upload"
                >
                  Upload
                </v-btn>
                <v-btn
                  color="success"
                  size="large"
                  variant="outlined"
                  prepend-icon="mdi-download"
                  @click="downloadTemplate"
                >
                  Download Template {{ uploadType === 'mendaftar' ? 'Mendaftar' : 'Lanjutan' }}
                </v-btn>
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
      </v-card>

      <!-- History Section -->
      <v-card elevation="2">
        <v-card-title class="bg-info text-white">
          <v-icon left>mdi-history</v-icon>
          LIST HISTORY UPLOAD TK
        </v-card-title>
        <v-card-text class="pa-6">
          <!-- History Instructions -->
          <v-alert type="warning" variant="tonal" class="mb-4">
            <div class="text-body-2">
              Untuk melihat daftar upload yang gagal, klik tombol <strong>"Download"</strong> pada
              tabel history upload di bawah ini. Upload ulang file menggunakan form di atas.
            </div>
          </v-alert>

          <!-- History Table -->
          <v-data-table
            :headers="historyHeaders"
            :items="historyData"
            :items-per-page="10"
            class="elevation-1"
            :loading="isLoadingHistory"
          >
            <template v-slot:top>
              <div class="d-flex justify-space-between align-center pa-4">
                <v-select
                  v-model="itemsPerPage"
                  :items="[10, 25, 50]"
                  label="Show entries"
                  variant="outlined"
                  density="compact"
                  style="max-width: 150px"
                />
                <v-text-field
                  v-model="searchQuery"
                  label="Search:"
                  prepend-inner-icon="mdi-magnify"
                  variant="outlined"
                  density="compact"
                  style="max-width: 300px"
                />
              </div>
            </template>

            <template v-slot:item.action="{ item }">
              <v-btn
                size="small"
                color="primary"
                variant="outlined"
                prepend-icon="mdi-download"
                @click="downloadFailedData(item)"
              >
                Download
              </v-btn>
            </template>

            <template v-slot:no-data>
              <div class="text-center pa-4">
                <v-icon size="64" color="grey">mdi-database-off</v-icon>
                <div class="text-h6 mt-2">No data available in table</div>
              </div>
            </template>
          </v-data-table>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Breadcrumbs
const breadcrumbs = ref([
  { title: 'Data Tenaga Kerja', disabled: false, href: '/edit' },
  { title: 'Upload Tenaga Kerja Baru', disabled: true },
])

// Upload form data
const uploadForm = ref(null)
const uploadType = ref('mendaftar')
const selectedFile = ref(null)
const isUploading = ref(false)

// Upload types
const uploadTypes = ref([
  { title: 'Upload TK Mendaftar', value: 'mendaftar' },
  { title: 'Upload TK Lanjutan', value: 'lanjutan' },
])

// History table data
const historyHeaders = ref([
  { title: 'Total Valid (disimpan)', key: 'totalValid', sortable: true },
  { title: 'Total Tidak Valid', key: 'totalInvalid', sortable: true },
  { title: 'Total Data', key: 'totalData', sortable: true },
  { title: 'Status Validasi', key: 'validationStatus', sortable: true },
  { title: 'Tanggal Selesai Validasi', key: 'validationDate', sortable: true },
  { title: 'Sumber Data', key: 'dataSource', sortable: true },
  { title: 'Jenis', key: 'type', sortable: true },
  { title: 'Action', key: 'action', sortable: false },
])

const historyData = ref([])
const isLoadingHistory = ref(false)
const itemsPerPage = ref(10)
const searchQuery = ref('')

onMounted(() => {
  loadHistoryData()
})

const goBack = () => {
  router.back()
}

const onFileSelected = (file) => {
  if (file && file.length > 0) {
    const file = file[0]
    // Validate file type
    const allowedTypes = ['.xls', '.xlsx']
    const fileExtension = '.' + file.name.split('.').pop().toLowerCase()

    if (!allowedTypes.includes(fileExtension)) {
      alert('File harus berupa Excel (.xls atau .xlsx)')
      selectedFile.value = null
      return
    }

    // Validate file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      alert('Ukuran file maksimal 10MB')
      selectedFile.value = null
      return
    }
  }
}

const handleUpload = async () => {
  if (!selectedFile.value || !uploadType.value) {
    alert('Pilih file dan jenis upload terlebih dahulu')
    return
  }

  isUploading.value = true

  try {
    // Simulate upload process
    await new Promise((resolve) => setTimeout(resolve, 2000))

    // Show success message
    alert('File berhasil diupload!')

    // Reset form
    selectedFile.value = null
    uploadType.value = 'mendaftar'

    // Reload history
    loadHistoryData()
  } catch (error) {
    console.error('Upload error:', error)
    alert('Gagal mengupload file. Silakan coba lagi.')
  } finally {
    isUploading.value = false
  }
}

import api from '@/services/api'

const downloadTemplate = async () => {
  try {
    const key = uploadType.value === 'mendaftar' ? 'tk' : 'koreksi_tk'
    const filename = uploadType.value === 'mendaftar'
      ? 'template_tk_24101780.xlsx'
      : 'template_koreksi_tk_24101780.xlsx'

    const base = api.baseURL || (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1')
    const url = `${base}/templates/${key}`

    const response = await fetch(url, { method: 'GET' })
    if (!response.ok) {
      throw new Error('Gagal mengunduh template')
    }
    const blob = await response.blob()
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(link.href)
  } catch (e) {
    alert('Gagal mengunduh template')
  }
}

const downloadFailedData = (item) => {
  alert(`Downloading failed data for: ${item.validationDate}`)

  // In real implementation, this would download the failed data file
  const link = document.createElement('a')
  link.href = '#'
  link.download = `failed_data_${item.validationDate}.xlsx`
  link.click()
}

const loadHistoryData = async () => {
  isLoadingHistory.value = true

  try {
    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 1000))

    // Mock data - in real implementation, this would come from API
    historyData.value = [
      // Add sample data here if needed
    ]
  } catch (error) {
    console.error('Error loading history:', error)
  } finally {
    isLoadingHistory.value = false
  }
}
</script>

<style scoped>
.upload-worker-page {
  background: #f8fafc;
  min-height: 100vh;
}

.v-card-title {
  font-weight: 700;
  letter-spacing: 0.5px;
}

.v-alert {
  border-radius: 8px;
}

.v-data-table {
  border-radius: 8px;
}

.gap-4 {
  gap: 16px;
}
</style>

