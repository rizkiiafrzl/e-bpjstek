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
              <div class="d-flex gap-2">
                <v-btn
                  size="small"
                  color="success"
                  variant="outlined"
                  prepend-icon="mdi-download"
                  @click="downloadFailedData(item)"
                >
                  Download
                </v-btn>
                <v-btn
                  size="small"
                  color="error"
                  variant="outlined"
                  prepend-icon="mdi-delete"
                  @click="confirmDeleteHistory(item)"
                >
                  Hapus
                </v-btn>
              </div>
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

      <!-- Delete Confirmation Dialog -->
      <v-dialog v-model="deleteDialog" max-width="500">
        <v-card>
          <v-card-title class="text-h5">
            <v-icon left color="error">mdi-alert-circle</v-icon>
            Konfirmasi Hapus
          </v-card-title>
          <v-card-text>
            <p>Apakah Anda yakin ingin menghapus file upload ini?</p>
            <p class="text-body-2 text-medium-emphasis">
              File: <strong>{{ selectedItem?.fileName || 'N/A' }}</strong><br>
              Tanggal: <strong>{{ selectedItem?.validationDate || 'N/A' }}</strong>
            </p>
            <v-alert type="warning" variant="tonal" class="mt-3">
              <div class="text-body-2">
                <strong>Peringatan:</strong> Tindakan ini tidak dapat dibatalkan. 
                File dan data terkait akan dihapus secara permanen.
              </div>
            </v-alert>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="grey"
              variant="text"
              @click="deleteDialog = false"
            >
              Batal
            </v-btn>
            <v-btn
              color="error"
              variant="flat"
              :loading="isDeleting"
              @click="deleteHistoryItem"
            >
              <v-icon left>mdi-delete</v-icon>
              Hapus
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
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

// Delete functionality
const deleteDialog = ref(false)
const selectedItem = ref(null)
const isDeleting = ref(false)

onMounted(() => {
  loadHistoryData()
})

const goBack = () => {
  router.push('/dashboard')
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
    let result
    // Upload TK Mendaftar → /workers/upload-tk; Upload TK Lanjutan → /koreksi-data-tk
    if (uploadType.value === 'lanjutan') {
      result = await api.koreksiTK(selectedFile.value)
    } else {
      result = await api.uploadTK(selectedFile.value)
    }
    
    // Show success message with details
    const base = `File berhasil diupload!` +
      `\n\nTotal Data: ${result.totalData || 0}` +
      `\nValid: ${result.valid || 0}` +
      `\nInvalid: ${result.invalid || 0}`
    // Tampilkan detail error bila ada
    const details = Array.isArray(result.errors) && result.errors.length
      ? `\n\nDetail Error (maks 20):\n- ` + result.errors.slice(0, 20).join(`\n- `)
      : ''
    alert(base + details)
    selectedFile.value = null
    uploadType.value = 'mendaftar'
    await loadHistoryData()
  } catch (error) {
    console.error('Upload error:', error)
    alert(error?.message || 'Gagal mengupload file. Silakan coba lagi.')
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
    const list = await api.getUploadHistory()
    historyData.value = Array.isArray(list) ? list : []
  } catch (error) {
    console.error('Error loading history:', error)
  } finally {
    isLoadingHistory.value = false
  }
}

const confirmDeleteHistory = (item) => {
  selectedItem.value = item
  deleteDialog.value = true
}

const deleteHistoryItem = async () => {
  if (!selectedItem.value) return
  
  isDeleting.value = true
  try {
    console.log('Deleting upload history with ID:', selectedItem.value.id)
    await api.deleteUploadHistory(selectedItem.value.id)
    alert('File berhasil dihapus!')
    deleteDialog.value = false
    selectedItem.value = null
    await loadHistoryData() // Reload history data
  } catch (error) {
    console.error('Error deleting history:', error)
    console.error('Error details:', {
      message: error?.message,
      stack: error?.stack,
      selectedItem: selectedItem.value
    })
    
    let errorMessage = 'Gagal menghapus file. Silakan coba lagi.'
    
    if (error?.message === 'Failed to fetch') {
      errorMessage = 'Tidak dapat terhubung ke server. Pastikan backend sedang berjalan di http://localhost:8080'
    } else if (error?.message?.includes('404')) {
      errorMessage = 'File tidak ditemukan atau sudah dihapus.'
    } else if (error?.message?.includes('401')) {
      errorMessage = 'Sesi login telah berakhir. Silakan login ulang.'
    } else if (error?.message) {
      errorMessage = error.message
    }
    
    alert(errorMessage)
  } finally {
    isDeleting.value = false
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

.gap-2 {
  gap: 8px;
}
</style>

