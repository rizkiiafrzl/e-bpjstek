<template>
  <div class="dashboard-wrapper">
    <v-container fluid class="dashboard-container">
      <!-- Dashboard Title -->
      <v-row class="justify-center">
        <v-col cols="12" class="text-center">
          <div class="dashboard-title">DASHBOARD PERIODE PELAPORAN BULAN JULI 2025</div>
        </v-col>
      </v-row>

      <!-- Summary Cards -->
      <v-row class="mb-4 summary-cards-row">
        <v-col
          cols="12"
          md="3"
          v-for="(card, index) in summaryCards"
          :key="index"
          class="summary-col"
        >
          <v-card elevation="3" class="summary-card" :class="`card-${index + 1}`">
            <v-card-text class="text-center summary-card-content">
              <div class="summary-value" :class="card.colorClass">{{ card.value }}</div>
              <div class="summary-label">{{ card.label }}</div>
              <v-btn
                v-if="card.button"
                color="success"
                size="small"
                variant="elevated"
                class="mt-2 detail-btn"
                @click="handleDetailClick"
              >
                {{ card.button }}
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Detail Iuran Cards -->
      <v-row class="mb-6 iuran-cards-row">
        <v-col cols="12" md="3" v-for="(iuran, index) in iuranCards" :key="index" class="iuran-col">
          <v-card elevation="3" class="iuran-card" :class="`iuran-${index + 1}`">
            <v-card-text class="iuran-card-content">
              <div class="d-flex align-center">
                <div class="iuran-icon" :class="iuran.iconClass">
                  <v-icon color="white" size="24">{{ iuran.icon }}</v-icon>
                </div>
                <div class="iuran-content">
                  <div class="iuran-label">{{ iuran.label }}</div>
                  <div class="iuran-value">{{ iuran.value }}</div>
                  <v-btn
                    v-if="iuran.button"
                    color="success"
                    size="small"
                    variant="elevated"
                    class="mt-2 kartu-btn"
                    @click="handleKartuClick"
                  >
                    {{ iuran.button }}
                  </v-btn>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Pengelolaan Periode Section -->
      <v-row class="justify-center">
        <v-col cols="12">
          <v-card elevation="3" class="management-card">
            <v-card-title class="management-title text-center">
              PENGELOLAAN PERIODE PELAPORAN BULAN JULI 2025
            </v-card-title>

            <v-card-text>
              <!-- Action Buttons -->
              <div class="action-buttons-container">
                <v-btn
                  v-for="(btn, index) in actionButtons"
                  :key="index"
                  :color="btn.color"
                  variant="elevated"
                  :prepend-icon="btn.icon"
                  class="action-btn"
                  @click="handleActionClick(btn.action)"
                >
                  {{ btn.label }}
                </v-btn>
              </div>

              <!-- Filter and Search -->
              <div class="d-flex align-center justify-center mb-6 filter-search-container">
                <div class="d-flex align-center gap-2">
                  <v-select
                    :items="filterOptions"
                    density="comfortable"
                    hide-details
                    variant="outlined"
                    class="filter-select"
                    v-model="selectedFilter"
                    @update:model-value="onFilterChange"
                    style="
                      width: 150px !important;
                      min-width: 150px !important;
                      max-width: 150px !important;
                      height: 40px !important;
                      flex: 0 0 150px !important;
                    "
                  ></v-select>
                  <v-btn
                    color="primary"
                    variant="outlined"
                    prepend-icon="mdi-refresh"
                    class="filter-btn"
                    size="default"
                    @click="loadWorkers"
                    style="
                      width: 150px !important;
                      min-width: 150px !important;
                      max-width: 150px !important;
                      height: 40px !important;
                      flex: 0 0 150px !important;
                    "
                    >Refresh TK</v-btn
                  >
                  <v-select
                    :items="['Pilih']"
                    density="comfortable"
                    hide-details
                    variant="outlined"
                    class="filter-select"
                    model-value="Pilih"
                    style="
                      width: 150px !important;
                      min-width: 150px !important;
                      max-width: 150px !important;
                      height: 40px !important;
                      flex: 0 0 150px !important;
                    "
                  ></v-select>
                  <v-text-field
                    placeholder="PENCARIAN"
                    density="comfortable"
                    hide-details
                    variant="outlined"
                    class="search-field"
                    style="
                      width: 150px !important;
                      min-width: 150px !important;
                      max-width: 150px !important;
                      height: 40px !important;
                      flex: 0 0 150px !important;
                    "
                  ></v-text-field>
                  <v-btn
                    color="primary"
                    variant="elevated"
                    prepend-icon="mdi-magnify"
                    class="search-btn"
                    size="default"
                    style="
                      width: 150px !important;
                      min-width: 150px !important;
                      max-width: 150px !important;
                      height: 40px !important;
                      flex: 0 0 150px !important;
                    "
                    >Search</v-btn
                  >
                </div>
              </div>

              <!-- Data Table -->
              <div class="table-responsive table-container">
                <v-table class="data-table" density="comfortable">
                  <thead>
                    <tr>
                      <th class="nik-column">NIK</th>
                      <th class="kpj-column">KPJ</th>
                      <th class="no-pegawai-column">No Pegawai</th>
                      <th class="nama-column">Nama</th>
                      <th class="upah-column">Upah (Rp)</th>
                      <th class="rapel-column">Rapel (Rp)</th>
                      <th class="action-column">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(row, idx) in tableData" :key="row.id || idx">
                      <td class="nik-column">{{ row.nik }}</td>
                      <td class="kpj-column">{{ row.kpj || '-' }}</td>
                      <td class="no-pegawai-column">{{ row.noPegawai || '-' }}</td>
                      <td class="nama-column">{{ row.nama }}</td>
                      <td class="upah-column">{{ row.upah }}</td>
                      <td class="rapel-column">{{ row.rapel }}</td>
                      <td class="action-column">
                        <v-btn size="small" color="primary" variant="text" icon="mdi-pencil" @click.stop.prevent="openEdit(row)"></v-btn>
                        <v-btn size="small" color="error" variant="text" icon="mdi-delete" @click.stop.prevent="confirmDelete(row)"></v-btn>
                      </td>
                    </tr>
                  </tbody>
                </v-table>
              </div>

              <!-- Edit Dialog -->
              <v-dialog v-model="editing" max-width="520">
                <v-card>
                  <v-card-title>Edit Data Karyawan</v-card-title>
                  <v-card-text>
                    <v-text-field v-model="editItem.nik" label="NIK" variant="outlined" />
                    <v-text-field v-model="editItem.kpj" label="KPJ" variant="outlined" />
                    <v-text-field v-model="editItem.noPegawai" label="No Pegawai" variant="outlined" />
                    <v-text-field v-model="editItem.nama" label="Nama" variant="outlined" />
                    <v-text-field v-model.number="editItem.upah" label="Upah" variant="outlined" type="number" />
                    <v-text-field v-model.number="editItem.rapel" label="Rapel" variant="outlined" type="number" />
                  </v-card-text>
                  <v-card-actions class="justify-end">
                    <v-btn variant="text" @click="editing = false">Batal</v-btn>
                    <v-btn color="primary" @click="saveEdit">Simpan</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <!-- Table Footer -->
              <div class="d-flex justify-space-between align-center mt-6 table-footer">
                <div class="text-caption text-grey-darken-1">Showing 10 entries</div>
                <div class="text-caption text-grey-darken-1">2018 © BPJS Ketenagakerjaan.</div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Notification Modal -->
      <v-dialog v-model="showNotificationModal" persistent max-width="500">
        <v-card>
          <v-card-text class="text-center pa-6">
            <div class="notification-icon mb-4">
              <v-icon color="primary" size="48">mdi-information</v-icon>
            </div>
            <div class="text-h6 font-weight-bold mb-3">Pemberitahuan</div>
            <div class="text-body-1 mb-4 text-left">
              Harap cek kembali! Jika Pekerja belum berhenti bekerja/PHK, pastikan Pekerja tersebut
              masih mendapatkan perlindungan jamsostek. Jika memerlukan informasi atau pertanyaan
              lebih lanjut agar menghubungi Petugas BP Jamsostek di Kantor Cabang terdaftar.
            </div>
            <div class="d-flex justify-center">
              <v-btn color="primary" variant="elevated" @click="closeNotification">OK</v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-dialog>

      <!-- Confirmation Modal -->
      <v-dialog v-model="showConfirmationModal" persistent max-width="400">
        <v-card>
          <v-card-text class="text-center pa-6">
            <div class="confirmation-icon mb-4">
              <v-icon color="primary" size="48">mdi-help-circle</v-icon>
            </div>
            <div class="text-h6 font-weight-bold mb-3">Konfirmasi</div>
            <div class="text-body-1 mb-4">Apakah PK/BU ada perubahan data?</div>
            <div class="d-flex justify-center gap-3">
              <v-btn color="primary" variant="elevated" @click="confirmYes">Ya</v-btn>
              <v-btn color="grey" variant="outlined" @click="confirmNo">Tidak</v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-dialog>

      <!-- Add Worker Flow Dialog (Step 3) -->
      <v-dialog v-model="dlgAddStep3" max-width="640">
        <v-card class="add-worker-dialog">
          <v-card-text class="text-center pa-6">
            <div class="aw-title">Pilihan tambah tenaga kerja</div>
            <div class="aw-desc">
              Pekerja Baru yang didaftarkan setelah kejadian meninggal dunia atau kecelakaan kerja,
              maka biaya obat/rawat dan manfaat JKK-JKM lainnya menjadi tanggung jawab pemberi
              kerja.
            </div>
            <v-checkbox class="mt-2" v-model="agreement" label="Saya setuju" />

            <div class="aw-divider"></div>

            <v-radio-group v-model="addType" class="aw-radio-group">
              <v-radio label="TAMBAH INDIVIDU" value="individu" />
              <v-radio label="TAMBAH MASSAL (UPLOAD)" value="massal" />
            </v-radio-group>

            <div class="d-flex justify-center gap-4 mt-4 aw-actions">
              <v-btn
                color="success"
                class="aw-btn aw-btn-green"
                :disabled="!agreement || !addType"
                @click="confirmAddFromEdit"
              >
                Pilih
              </v-btn>
              <v-btn
                color="error"
                class="aw-btn aw-btn-red"
                variant="elevated"
                @click="dlgAddStep3 = false"
              >
                Cancel
              </v-btn>
            </div>

            <v-alert
              v-if="!agreement"
              type="warning"
              variant="tonal"
              density="comfortable"
              class="mt-4"
            >
              Silakan checklist persetujuan pernyataan terlebih dahulu
            </v-alert>
          </v-card-text>
        </v-card>
      </v-dialog>
    </v-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

// Modal state
const showNotificationModal = ref(false)
const showConfirmationModal = ref(false)

// Summary cards data
const summaryCards = ref([
  {
    value: '4',
    label: 'TOTAL TENAGA KERJA',
    colorClass: 'text-red',
    button: null,
  },
  {
    value: 'Rp. 22.000.000,00',
    label: 'TOTAL UPAH + RAPEL',
    colorClass: 'text-grey-darken-2',
    button: null,
  },
  {
    value: 'Rp. 118.800,00',
    label: 'TOTAL IURAN',
    colorClass: 'text-red',
    button: null,
  },
  {
    value: 'Rp. 23.544,00',
    label: 'TOTAL DENDA',
    colorClass: 'text-red',
    button: 'DETAIL',
  },
])

// Iuran cards data
const iuranCards = ref([
  {
    label: 'IURAN JKK',
    value: 'Rp 52.800,00',
    icon: 'mdi-account-group',
    iconClass: 'blue',
    button: null,
  },
  {
    label: 'IURAN JKM',
    value: 'Rp 66.000,00',
    icon: 'mdi-cog',
    iconClass: 'pink',
    button: null,
  },
  {
    label: 'IURAN JHT',
    value: 'Rp 0,00',
    icon: 'mdi-chicken',
    iconClass: 'orange',
    button: null,
  },
  {
    label: 'IURAN JP',
    value: 'Rp 0,00',
    icon: 'mdi-account-heart',
    iconClass: 'yellow',
    button: 'KARTU ANGSURAN',
  },
])

// Action buttons data
const actionButtons = ref([
  {
    label: 'TAMBAH TK',
    color: 'success',
    icon: 'mdi-plus',
    action: 'tambah-tk',
  },
  {
    label: 'UPLOAD TK NA',
    color: 'primary',
    icon: 'mdi-upload',
    action: 'upload-tk-na',
  },
  {
    label: 'UPLOAD UPAH',
    color: 'info',
    icon: 'mdi-upload',
    action: 'upload-upah',
  },
  {
    label: 'KOREKSI DATA TK MASAL',
    color: 'warning',
    icon: 'mdi-pencil',
    action: 'koreksi-data',
  },
  {
    label: 'HITUNG IURAN',
    color: 'success',
    icon: 'mdi-calculator',
    action: 'hitung-iuran',
  },
  {
    label: 'FINALISASI',
    color: 'error',
    icon: 'mdi-check',
    action: 'finalisasi',
  },
])

// Table data dari API workers
import apiService from '../services/api.js'
const tableData = ref([])
const allWorkers = ref([]) // Store all workers data
const editing = ref(false)
const editItem = ref({ id: null, nik: '', kpj: '', noPegawai: '', nama: '', upah: 0, rapel: 0 })

// Filter state - default to show all workers
const selectedFilter = ref('Semua Data')
const filterOptions = [
  { title: 'Semua Data', value: 'semua' },
  { title: 'Peserta Aktif', value: 'aktif' },
  { title: 'Peserta Non Aktif', value: 'non-aktif' },
  { title: 'Peserta Baru', value: 'baru' }
]

// Load workers data
const loadWorkers = async () => {
  try {
    const rows = await apiService.getWorkers()
    allWorkers.value = Array.isArray(rows) ? rows : []
    applyFilter()
  } catch (e) {
    console.error('Gagal memuat workers', e)
  }
}

// Apply filter based on selected option
const applyFilter = () => {
  let filteredWorkers = [...allWorkers.value]
  
  switch (selectedFilter.value) {
    case 'semua':
      // Semua Data: tampilkan semua workers
      filteredWorkers = allWorkers.value
      break
    case 'aktif':
      // Peserta Aktif: workers dengan upah > 0 dan tidak ada tanggal akhir kontrak
      filteredWorkers = allWorkers.value.filter(worker => 
        (worker.upah > 0) && 
        (!worker.tanggalAkhirKontrak || worker.tanggalAkhirKontrak === null)
      )
      break
    case 'non-aktif':
      // Peserta Non Aktif: workers dengan tanggal akhir kontrak yang sudah lewat
      filteredWorkers = allWorkers.value.filter(worker => {
        if (!worker.tanggalAkhirKontrak) return false
        const endDate = new Date(worker.tanggalAkhirKontrak)
        const today = new Date()
        return endDate < today
      })
      break
    case 'baru': {
      // Peserta Baru: workers yang baru dibuat dalam 30 hari terakhir
      const thirtyDaysAgo = new Date()
      thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30)
      filteredWorkers = allWorkers.value.filter(worker => {
        const createdDate = new Date(worker.createdAt)
        return createdDate > thirtyDaysAgo
      })
      break
    }
    default:
      filteredWorkers = allWorkers.value
  }
  
  // Format data for table display
  tableData.value = filteredWorkers.map((r) => ({
    id: r.id,
    nik: r.nik,
    kpj: r.kpj,
    noPegawai: r.noPegawai,
    nama: r.nama,
    upah: new Intl.NumberFormat('id-ID', { style: 'decimal', minimumFractionDigits: 2 }).format(r.upah || 0),
    rapel: new Intl.NumberFormat('id-ID', { style: 'decimal', minimumFractionDigits: 2 }).format(r.rapel || 0),
  }))
  
  // Update summary cards
  updateSummaryCards(filteredWorkers)
}

// Update summary cards based on filtered data
const updateSummaryCards = (workers) => {
  const totalWorkers = workers.length
  const totalUpah = workers.reduce((sum, worker) => sum + (worker.upah || 0) + (worker.rapel || 0), 0)
  const totalIuran = totalWorkers * 29700 // Rp 29.700 per worker
  const totalDenda = totalIuran * 0.02 // 2% denda
  
  summaryCards.value[0].value = totalWorkers.toString()
  summaryCards.value[1].value = new Intl.NumberFormat('id-ID', { 
    style: 'currency', 
    currency: 'IDR',
    minimumFractionDigits: 2 
  }).format(totalUpah)
  summaryCards.value[2].value = new Intl.NumberFormat('id-ID', { 
    style: 'currency', 
    currency: 'IDR',
    minimumFractionDigits: 2 
  }).format(totalIuran)
  summaryCards.value[3].value = new Intl.NumberFormat('id-ID', { 
    style: 'currency', 
    currency: 'IDR',
    minimumFractionDigits: 2 
  }).format(totalDenda)
}

// Handle filter change
const onFilterChange = (value) => {
  selectedFilter.value = value
  applyFilter()
}

onMounted(async () => {
  // Load workers data immediately
  await loadWorkers()
  // Show notification modal after data is loaded
  showNotificationModal.value = true
})

const openEdit = (row) => {
  router.push(`/tenaga/edit/${row.id}`)
}

const saveEdit = async () => {
  try {
    await apiService.updateWorker(editItem.value.id, {
      nik: editItem.value.nik,
      kpj: editItem.value.kpj,
      noPegawai: editItem.value.noPegawai,
      nama: editItem.value.nama,
      upah: editItem.value.upah,
      rapel: editItem.value.rapel,
    })
    await loadWorkers()
    editing.value = false
  } catch (e) {
    alert(e?.message || 'Gagal menyimpan perubahan')
  }
}

const confirmDelete = async (row) => {
  if (!confirm('Hapus data karyawan ini?')) return
  try {
    await apiService.deleteWorker(row.id)
    await loadWorkers() // Reload data and apply current filter
  } catch (e) {
    alert(e?.message || 'Gagal menghapus data')
  }
}

const handleDetailClick = () => {
  console.log('Detail clicked')
}

const handleKartuClick = () => {
  console.log('Kartu Angsuran clicked')
}

const handleActionClick = (action) => {
  console.log('Action clicked:', action)
  if (action === 'tambah-tk') {
    openAddWorkerFlow()
  } else if (action === 'upload-tk-na') {
    router.push('/tenaga/upload-na')
  } else if (action === 'upload-upah') {
    router.push('/tenaga/upload-upah')
  } else if (action === 'koreksi-data') {
    router.push('/tenaga/koreksi-data')
  }
}

// Add worker popup flow (Step 3 only here)
const dlgAddStep3 = ref(false)
const agreement = ref(false)
const addType = ref(null)

const router = useRouter()

const openAddWorkerFlow = () => {
  agreement.value = false
  addType.value = null
  dlgAddStep3.value = true
}

const confirmAddFromEdit = () => {
  if (!agreement.value || !addType.value) return
  dlgAddStep3.value = false
  window.scrollTo(0, 0)
  if (addType.value === 'individu') {
    router.push('/tenaga/tambah')
  } else {
    // massal upload - redirect to upload page
    router.push('/tenaga/upload')
  }
}

const closeNotification = () => {
  showNotificationModal.value = false
  // Show confirmation modal after notification is closed
  showConfirmationModal.value = true
}

const confirmYes = () => {
  showConfirmationModal.value = false
  console.log('User confirmed: Yes')
}

const confirmNo = () => {
  showConfirmationModal.value = false
  console.log('User confirmed: No')
}
</script>

<style scoped>
/* Global overrides for uniform sizing */
.filter-search-container * {
  box-sizing: border-box !important;
}

/* Nuclear option - force all elements to same size */
.filter-search-container .v-select,
.filter-search-container .v-text-field,
.filter-search-container .v-btn,
.filter-search-container .v-field,
.filter-search-container .v-field__input,
.filter-search-container .v-btn__content,
.filter-search-container .v-input,
.filter-search-container .v-input__control,
.filter-search-container .v-input__details {
  width: 150px !important;
  min-width: 150px !important;
  max-width: 150px !important;
  height: 40px !important;
  flex: 0 0 150px !important;
  flex-shrink: 0 !important;
  flex-grow: 0 !important;
  flex-basis: 150px !important;
}
.dashboard-wrapper {
  min-height: 100vh;
  background: #f8fafc;
}

.header-bar {
  position: sticky;
  top: 0;
  z-index: 1000;
}

.nav-btn {
  margin: 0 4px;
  transition: all 0.3s ease;
}

.nav-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.dashboard-container {
  padding: 32px;
  animation: fadeInUp 0.6s ease-out;
}

.header-section {
  animation: slideInLeft 0.8s ease-out;
}

.breadcrumb {
  animation: slideInRight 0.8s ease-out;
}

.dashboard-title {
  font-size: 2rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: 2rem;
  color: #424242;
  animation: fadeInDown 1s ease-out;
}

/* Summary Cards */
.summary-cards-row {
  animation: fadeInUp 1.2s ease-out;
}

.summary-col {
  animation: slideInUp 1.4s ease-out;
}

.summary-card {
  border-radius: 12px;
  background: white;
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  height: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.summary-card-content {
  padding: 16px !important;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.summary-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent);
  transition: left 0.5s;
}

.summary-card:hover::before {
  left: 100%;
}

.summary-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

.summary-value {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  transition: all 0.3s ease;
}

.summary-label {
  font-size: 0.875rem;
  color: #666;
  font-weight: 500;
}

.detail-btn {
  transition: all 0.3s ease;
}

.detail-btn:hover {
  transform: scale(1.05);
}

/* Iuran Cards */
.iuran-cards-row {
  animation: fadeInUp 1.6s ease-out;
}

.iuran-col {
  animation: slideInUp 1.8s ease-out;
}

.iuran-card {
  border-radius: 12px;
  background: white;
  transition: all 0.3s ease;
  cursor: pointer;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.iuran-card-content {
  padding: 16px !important;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.iuran-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}

.iuran-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  transition: all 0.3s ease;
}

.iuran-card:hover .iuran-icon {
  transform: scale(1.1);
}

.iuran-icon.blue {
  background: #2196f3;
}

.iuran-icon.pink {
  background: #e91e63;
}

.iuran-icon.orange {
  background: #ff9800;
}

.iuran-icon.yellow {
  background: #ffc107;
}

.iuran-content {
  flex: 1;
}

.iuran-label {
  font-size: 0.875rem;
  color: #666;
  font-weight: 500;
}

.iuran-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: #333;
}

.kartu-btn {
  transition: all 0.3s ease;
}

.kartu-btn:hover {
  transform: scale(1.05);
}

/* Management Card */
.management-card {
  border-radius: 12px;
  background: white;
  animation: fadeInUp 2s ease-out;
}

.management-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #424242;
  padding: 32px 32px 24px 32px;
  text-align: center;
}

.action-buttons-container {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  margin-bottom: 32px;
  padding: 0 32px;
  justify-content: center;
}

.action-btn {
  transition: all 0.3s ease;
  animation: fadeInUp 2.2s ease-out;
  min-width: 180px;
  height: 40px;
  font-size: 0.875rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.action-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.filter-search-container {
  padding: 0 32px;
  flex-wrap: nowrap !important;
  gap: 8px !important;
  justify-content: center !important;
}

.filter-search-container .d-flex {
  gap: 8px !important;
  align-items: center !important;
  flex-wrap: nowrap !important;
  justify-content: center !important;
}

/* Force uniform sizes for all filter controls */
.filter-select,
.filter-btn,
.search-field,
.search-btn {
  width: 150px !important;
  min-width: 150px !important;
  max-width: 150px !important;
  height: 40px !important;
  flex-shrink: 0 !important;
  box-sizing: border-box !important;
}

/* Override Vuetify specific classes */
.v-select.filter-select,
.v-btn.filter-btn,
.v-text-field.search-field,
.v-btn.search-btn {
  width: 150px !important;
  min-width: 150px !important;
  max-width: 150px !important;
  height: 40px !important;
}

/* Override Vuetify input field */
.v-text-field.search-field .v-field {
  width: 150px !important;
  min-width: 150px !important;
  max-width: 150px !important;
}

/* Override Vuetify button */
.v-btn.filter-btn,
.v-btn.search-btn {
  width: 150px !important;
  min-width: 150px !important;
  max-width: 150px !important;
  padding: 0 16px !important;
}

/* Force all controls to have identical dimensions */
.filter-search-container .v-select,
.filter-search-container .v-text-field,
.filter-search-container .v-btn {
  width: 150px !important;
  min-width: 150px !important;
  max-width: 150px !important;
  height: 40px !important;
  flex: 0 0 150px !important;
}

.filter-btn,
.search-btn {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Additional specificity for Vuetify components */
.filter-search-container .v-field,
.filter-search-container .v-field__input,
.filter-search-container .v-btn__content {
  width: 100% !important;
  max-width: 100% !important;
}

.table-container {
  padding: 0 32px;
}

.table-footer {
  padding: 0 32px;
}

/* Data Table */
.data-table {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  animation: fadeInUp 2.4s ease-out;
  table-layout: fixed;
}

.data-table thead th {
  background: #f5f5f5;
  color: #424242;
  font-weight: 600;
  padding: 16px 12px;
  text-align: left;
}

.data-table tbody tr {
  transition: all 0.3s ease;
}

.data-table tbody tr:hover {
  background: #f8f9fa;
}

.data-table tbody tr:nth-child(even) {
  background: #fafafa;
}

.data-table td {
  padding: 12px;
  border-bottom: 1px solid #e0e0e0;
}

/* Data Table - Column Widths */
.nik-column {
  width: 18%;
  min-width: 140px;
}

.kpj-column {
  width: 10%;
  min-width: 80px;
}

.no-pegawai-column {
  width: 10%;
  min-width: 80px;
}

.nama-column {
  width: 30%;
  min-width: 180px;
}

.upah-column {
  width: 15%;
  min-width: 120px;
}

.rapel-column {
  width: 12%;
  min-width: 100px;
}

.action-column {
  width: 5%;
  min-width: 60px;
  text-align: center;
}

/* Modal Animations */
.notification-icon,
.confirmation-icon {
  display: flex;
  justify-content: center;
  animation: pulse 2s infinite;
}

/* Keyframe Animations */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideInLeft {
  from {
    opacity: 0;
    transform: translateX(-30px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateX(30px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes pulse {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
  100% {
    transform: scale(1);
  }
}

/* Mobile responsive */
@media (max-width: 600px) {
  .dashboard-container {
    padding: 16px;
  }

  .dashboard-title {
    font-size: 1.5rem;
  }

  .action-buttons-container {
    flex-direction: column;
    align-items: stretch;
    padding: 0 16px;
  }

  .action-btn {
    margin-bottom: 8px;
    min-width: 100%;
    width: 100%;
  }

  .filter-search-container {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    padding: 0 16px;
  }

  .filter-select,
  .filter-btn,
  .search-field,
  .search-btn {
    width: 100%;
    min-width: 100%;
    max-width: 100%;
    height: 44px;
  }

  .summary-value {
    font-size: 2rem;
  }

  .iuran-icon {
    width: 40px;
    height: 40px;
    margin-right: 12px;
  }
}
.elevated-btn {
  min-width: 120px;
  font-weight: 600;
}
.add-worker-dialog {
  border-radius: 12px;
}
.aw-title {
  font-size: 26px;
  font-weight: 800;
  color: #3a3a3a;
  margin-bottom: 8px;
  text-transform: none;
}
.aw-desc {
  max-width: 520px;
  margin: 0 auto 8px auto;
  color: #6b7280;
}
.aw-divider {
  width: 100%;
  height: 1px;
  background: #e5e7eb;
  margin: 12px 0 8px 0;
}
.aw-radio-group :deep(.v-label) {
  font-weight: 700;
  letter-spacing: 0.2px;
}
.aw-actions .aw-btn {
  min-width: 140px;
  height: 48px;
  font-weight: 700;
  border-radius: 8px;
}
.aw-btn-green {
  background: #2ecc71 !important;
  color: #fff !important;
}
.aw-btn-red {
  background: #f26a6a !important;
  color: #fff !important;
}
</style>
