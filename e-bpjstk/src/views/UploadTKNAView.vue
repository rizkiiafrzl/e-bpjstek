<template>
  <div class="upload-tk-na-page">
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
          UPLOAD TENAGA KERJA NONAKTIF
        </v-card-title>
        <v-card-text class="pa-6">
          <!-- Instructions -->
          <div class="instructions mb-4">
            <ol class="text-body-2">
              <li class="mb-2">
                Mohon diperiksa kembali bahwa tenaga kerja yang akan dinonaktifkan
                <strong class="text-red">tidak dalam status meninggal</strong> dalam bulan
                pelaporan.
              </li>
              <li class="mb-2">
                File upload adalah microsoft excel ( xls / xlsx ) dengan maksimum 10.000 baris.
              </li>
              <li class="mb-2">Susunan tombol : download template - choose file - upload.</li>
              <li class="mb-2">
                File upload harus berdasarkan dari template upah yang telah di download, nama file
                tidak boleh diubah.
              </li>
              <li class="mb-2">Format data pada file excel harus berupa text.</li>
            </ol>
          </div>

          <!-- Upload Form -->
          <v-form @submit.prevent="handleUpload" ref="uploadForm">
            <div class="upload-controls">
              <v-btn
                color="success"
                size="large"
                variant="outlined"
                prepend-icon="mdi-download"
                @click="downloadTemplate"
                class="download-btn"
              >
                DOWNLOAD TEMPLATE
              </v-btn>

              <v-file-input
                v-model="selectedFile"
                label="Choose files To Upload"
                variant="outlined"
                accept=".xls,.xlsx"
                prepend-icon="mdi-file-excel"
                :multiple="false"
                @update:model-value="onFileSelected"
                class="file-input"
              />

              <v-btn
                type="submit"
                color="success"
                size="large"
                :loading="isUploading"
                :disabled="!selectedFile"
                prepend-icon="mdi-upload"
                class="upload-btn"
              >
                UPLOAD
              </v-btn>
            </div>
          </v-form>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Breadcrumbs
const breadcrumbs = ref([
  { title: 'Data Tenaga Kerja', disabled: false, href: '/edit' },
  { title: 'Upload Tenaga Kerja Nonaktif', disabled: true },
])

// Upload form data
const uploadForm = ref(null)
const selectedFile = ref(null)
const isUploading = ref(false)

const goBack = () => {
  router.push('/dashboard')
}

const onFileSelected = (value) => {
  // Normalize: Vuetify may pass a File or an array of File
  const file = Array.isArray(value) ? value[0] : value
  if (!file) {
    selectedFile.value = null
    return
  }
  // Validate file type
  const allowedTypes = ['.xls', '.xlsx']
  const fileExtension = '.' + (file.name || '').split('.').pop().toLowerCase()
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
  selectedFile.value = file
}

const handleUpload = async () => {
  if (!selectedFile.value) {
    alert('Pilih file terlebih dahulu')
    return
  }

  isUploading.value = true

  try {
    // Upload TK NA menggunakan API upload TK NA yang benar
    const result = await api.uploadTKNA(selectedFile.value)
    // Show success message with details
    const message = `File TK Nonaktif berhasil diupload!\n\n` +
      `Total Data: ${result.totalData || 0}\n` +
      `Berhasil Dinonaktifkan: ${result.valid || 0}\n` +
      `Gagal: ${result.invalid || 0}\n` +
      `Aksi: ${result.action || 'nonaktif'}`
    alert(message)
    // Reset form
    selectedFile.value = null
  } catch (error) {
    console.error('Upload error:', error)
    alert(error?.message || 'Gagal mengupload file TK Nonaktif. Silakan coba lagi.')
  } finally {
    isUploading.value = false
  }
}

import api from '@/services/api'

const downloadTemplate = async () => {
  try {
    const base = api.baseURL || (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1')
    const url = `${base}/templates/tk_na`
    const response = await fetch(url, { method: 'GET' })
    if (!response.ok) {
      throw new Error('Gagal mengunduh template')
    }
    const blob = await response.blob()
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'template_tk_na_24101780.xlsx'
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(link.href)
  } catch (e) {
    alert('Gagal mengunduh template')
  }
}
</script>

<style scoped>
.upload-tk-na-page {
  background: #f8fafc;
  min-height: 100vh;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.v-card {
  animation: gentleSlide 0.4s ease-out;
}

@keyframes gentleSlide {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.v-card-title {
  font-weight: 600;
  letter-spacing: 0.3px;
  font-size: 1.1rem;
}

.instructions ol {
  padding-left: 20px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.instructions li {
  margin-bottom: 8px;
  font-size: 0.95rem;
  line-height: 1.4;
  animation: gentleFade 0.3s ease-out;
  animation-fill-mode: both;
}

.instructions li:nth-child(1) {
  animation-delay: 0.05s;
}
.instructions li:nth-child(2) {
  animation-delay: 0.1s;
}
.instructions li:nth-child(3) {
  animation-delay: 0.15s;
}
.instructions li:nth-child(4) {
  animation-delay: 0.2s;
}
.instructions li:nth-child(5) {
  animation-delay: 0.25s;
}

@keyframes gentleFade {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.upload-controls {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 24px;
  flex-wrap: wrap;
}

.download-btn {
  min-width: 180px;
  height: 48px;
  font-weight: 600;
  font-size: 0.9rem;
  letter-spacing: 0.5px;
  transition: all 0.2s ease;
  animation: gentleAppear 0.3s ease-out;
}

.upload-btn {
  min-width: 120px;
  height: 48px;
  font-weight: 600;
  font-size: 0.9rem;
  letter-spacing: 0.5px;
  transition: all 0.2s ease;
  animation: gentleAppear 0.3s ease-out;
  animation-delay: 0.1s;
  animation-fill-mode: both;
}

@keyframes gentleAppear {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.download-btn:hover,
.upload-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.file-input {
  flex: 1;
  min-width: 300px;
  animation: gentleAppear 0.3s ease-out;
  animation-delay: 0.05s;
  animation-fill-mode: both;
}

.v-file-input :deep(.v-field) {
  height: 48px !important;
  min-height: 48px !important;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.v-file-input :deep(.v-field__input) {
  height: 48px !important;
  min-height: 48px !important;
  padding-top: 0 !important;
  padding-bottom: 0 !important;
}

.v-file-input :deep(.v-field__outline) {
  height: 48px !important;
  min-height: 48px !important;
}

.v-file-input :deep(.v-label) {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  font-size: 0.9rem;
}

.v-file-input :deep(.v-field):hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .upload-controls {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .download-btn,
  .upload-btn {
    min-width: auto;
    width: 100%;
  }

  .file-input {
    min-width: auto;
    width: 100%;
  }
}
</style>
