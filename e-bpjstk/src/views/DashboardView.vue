<template>
  <v-container fluid class="dashboard-container">
    <!-- Identity Bar (no app header) -->
    <v-row>
      <v-col cols="12">
        <v-alert
          color="success"
          variant="tonal"
          border="start"
          density="comfortable"
          class="identity-alert"
        >
          <div class="d-flex align-center justify-space-between flex-wrap">
            <div>
              <div class="text-h6 font-weight-bold">Selamat Datang</div>
              <div class="text-body-2">
                Anda login sebagai <strong>{{ userEmail }}</strong>
              </div>
              <div class="text-caption text-medium-emphasis">Waktu login: {{ loginTime }}</div>
            </div>
            <v-btn @click="handleLogout" color="error" variant="elevated" size="small">
              <v-icon left>mdi-logout</v-icon>
              Logout
            </v-btn>
          </div>
        </v-alert>
      </v-col>
    </v-row>

    <!-- Mutasi Data Section -->
    <v-row class="mt-4">
      <v-col cols="12">
        <v-card elevation="2" class="section-card">
          <v-card-title class="d-flex justify-space-between align-center">
            <div class="text-h6 font-weight-bold d-flex align-center">
              <v-icon left color="primary" class="mr-2">mdi-table</v-icon>
              Mutasi Data
            </div>
            <div class="text-caption text-medium-emphasis">Data Periode Pelaporan</div>
          </v-card-title>

          <v-card-text>
            <v-row>
              <!-- Summary box (full width top) -->
              <v-col cols="12" class="mb-4">
                <v-card variant="outlined" elevation="1" class="panel-card">
                  <v-table density="comfortable" class="summary-table">
                    <tbody>
                      <tr>
                        <td>Kode Tagihan</td>
                        <td class="text-right">{{ summary.kodeTagihan || '-' }}</td>
                      </tr>
                      <tr>
                        <td>Total Iuran dan Denda</td>
                        <td class="text-right">{{ formatCurrency(summary.totalIuranDanDenda) }}</td>
                      </tr>
                      <tr>
                        <td>Sisa Pembayaran Iuran Sebelumnya</td>
                        <td class="text-right">{{ formatCurrency(summary.sisaPembayaranSebelumnya) }}</td>
                      </tr>
                      <tr>
                        <td>Total Tagihan</td>
                        <td class="text-right">{{ formatCurrency(summary.totalTagihan) }}</td>
                      </tr>
                    </tbody>
                  </v-table>
                </v-card>
              </v-col>

              <!-- Controls & table (full width bottom) -->
              <v-col cols="12">
                <v-card variant="outlined" elevation="1" class="panel-card">
                  <div
                    class="d-flex justify-space-between align-center flex-wrap gap-4 mb-4 controls-row"
                  >
                    <v-btn color="primary" variant="elevated" prepend-icon="mdi-plus" @click="addPeriod">
                      Tambah Periode Pelaporan
                    </v-btn>
                    <div class="filters-group">
                      <div class="inline-control">
                        <span class="text-body-2 text-medium-emphasis">Tampil</span>
                        <v-select
                          :items="[10, 25, 50, 100]"
                          density="comfortable"
                          hide-details
                          variant="outlined"
                          style="width: 84px"
                          :model-value="pageSize"
                          @update:modelValue="onChangePageSize"
                        ></v-select>
                        <span class="text-body-2 text-medium-emphasis">entries</span>
                      </div>
                      <div class="inline-control">
                        <span class="text-body-2 text-medium-emphasis">Filter Status</span>
                        <v-select
                          class="status-select"
                          :items="['Tampil Semua', 'Draft', 'Approval', 'Finalisasi', 'Posting']"
                          density="comfortable"
                          hide-details
                          variant="outlined"
                          style="min-width: 200px"
                          :model-value="statusFilterLabel"
                          @update:modelValue="onChangeStatusFilter"
                        ></v-select>
                      </div>
                    </div>
                  </div>

                  <div class="table-responsive">
                    <v-table class="mutasi-table" density="comfortable">
                      <thead>
                        <tr>
                          <th style="width: 50px">No</th>
                          <th>Bulan Iuran</th>
                          <th>Jumlah TK</th>
                          <th>Nominal Iuran (Rp)</th>
                          <th>Nominal Denda (Rp)</th>
                          <th>Status</th>
                          <th style="width: 200px">Aksi</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(row, idx) in reportRows" :key="row.id || idx">
                          <td>{{ idx + 1 }}</td>
                          <td>{{ toMonthYear(row.year, row.month) }}</td>
                          <td>{{ row.totalTk ?? '-' }}</td>
                          <td>{{ formatCurrency(row.totalIuran) }}</td>
                          <td>{{ formatCurrency(row.totalDenda) }}</td>
                          <td>
                            <v-chip :color="row.status === 'Posting' ? 'success' : row.status === 'Draft' ? 'grey' : 'info'" size="small" variant="tonal">
                              {{ row.status || 'Draft' }}
                            </v-chip>
                          </td>
                          <td>
                            <div class="aksi-buttons">
                              <v-btn size="x-small" color="primary" variant="tonal" prepend-icon="mdi-pencil" @click="navigateToEdit(row)">Edit</v-btn>
                              <v-btn size="x-small" variant="tonal" prepend-icon="mdi-printer" @click="onPrint(row)">Cetak</v-btn>
                              <v-btn size="x-small" variant="outlined" color="error" prepend-icon="mdi-delete" @click="onDeletePeriod(row)">Hapus</v-btn>
                            </div>
                          </td>
                        </tr>
                      </tbody>
                    </v-table>
                  </div>
                  <div class="d-flex justify-space-between align-center mt-3">
                    <div class="text-caption">Total: {{ totalRows }}</div>
                    <div class="d-flex align-center" style="gap:8px">
                      <span class="text-caption">Halaman</span>
                      <v-select :items="Array.from({length: Math.max(1, Math.ceil(totalRows/pageSize))}, (_,i)=> i+1)" density="comfortable" hide-details variant="outlined" style="width:84px" :model-value="page" @update:modelValue="val => { page = Number(val)||1; loadReportPeriods() }" />
                    </div>
                  </div>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import apiService from '../services/api.js'

const router = useRouter()
const userEmail = ref('')
const loginTime = ref('')
const summary = ref({ kodeTagihan: '-', totalIuranDanDenda: 0, sisaPembayaranSebelumnya: 0, totalTagihan: 0 })

// Data periode pelaporan dari API
const reportRows = ref([])
const totalRows = ref(0)
let page = 1
const pageSize = ref(10)
const statusFilter = ref('all')
const statusFilterLabel = computed(() => {
  switch (statusFilter.value) {
    case 'Draft': return 'Draft'
    case 'Approval': return 'Approval'
    case 'Finalisasi': return 'Finalisasi'
    case 'Posting': return 'Posting'
    default: return 'Tampil Semua'
  }
})

onMounted(async () => {
  // Get user data from localStorage
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      const user = JSON.parse(userData)
      userEmail.value = user.email || user.fullName || 'user@example.com'
    } catch (error) {
      console.error('Error parsing user data:', error)
      userEmail.value = 'user@example.com'
    }
  } else {
    userEmail.value = 'user@example.com'
  }

  // Get login time
  const loginTimestamp = localStorage.getItem('loginTime')
  if (loginTimestamp) {
    loginTime.value = new Date(parseInt(loginTimestamp)).toLocaleString('id-ID')
  } else {
    loginTime.value = new Date().toLocaleString('id-ID')
  }
  await Promise.all([
    loadReportPeriods(),
    loadSummary(),
  ])
})

const loadReportPeriods = async () => {
  try {
    const res = await apiService.getReportPeriods({ page, pageSize: pageSize.value, status: statusFilter.value })
    if (res && Array.isArray(res.items)) {
      reportRows.value = res.items
      totalRows.value = res.total || res.items.length
    } else if (Array.isArray(res)) {
      reportRows.value = res
      totalRows.value = res.length
    }
  } catch (e) {
    console.error('Gagal memuat periode pelaporan', e)
  }
}

const loadSummary = async () => {
  try {
    const data = await apiService.getReportSummary()
    if (data) summary.value = data
  } catch (e) {
    console.error('Gagal memuat ringkasan', e)
  }
}

const handleLogout = async () => {
  try {
    // Call logout API
    await apiService.logout()
  } catch (error) {
    console.error('Logout error:', error)
  } finally {
    // Redirect to login page
    router.push('/login')
  }
}

// Navigate to edit page
const navigateToEdit = (row) => {
  // Show confirmation dialog
  if (confirm('Apakah Anda yakin ingin mengedit data ini?')) {
    const dataParam = encodeURIComponent(JSON.stringify(row))
    router.push(`/edit/${dataParam}`)
  }
}

// Tambah periode pelaporan (default: bulan berjalan), cegah duplikat per bulan
const addPeriod = async () => {
  try {
    await apiService.createReportPeriod({})
    await loadReportPeriods()
    alert('Periode pelaporan berhasil ditambahkan')
  } catch (e) {
    alert(e?.message || 'Gagal menambah periode pelaporan')
  }
}

const onChangeStatusFilter = async (label) => {
  const map = { 'Tampil Semua': 'all', 'Draft': 'Draft', 'Approval': 'Approval', 'Finalisasi': 'Finalisasi', 'Posting': 'Posting' }
  statusFilter.value = map[label] || 'all'
  page = 1
  await loadReportPeriods()
}

const onChangePageSize = async (val) => {
  pageSize.value = Number(val) || 10
  page = 1
  await loadReportPeriods()
}

const onDeletePeriod = async (row) => {
  if (!confirm('Hapus periode ini? Hanya status Draft bisa dihapus.')) return
  try {
    await apiService.deleteReportPeriod(row.id)
    await loadReportPeriods()
  } catch (e) {
    alert(e?.message || 'Gagal menghapus periode')
  }
}

const onCalculate = async (row) => {
  try {
    await apiService.calculateReportPeriod(row.id)
    await loadReportPeriods()
  } catch (e) {
    alert(e?.message || 'Gagal menghitung iuran')
  }
}

const onFinalize = async (row) => {
  if (!confirm('Finalisasi periode ini?')) return
  try {
    await apiService.finalizeReportPeriod(row.id)
    await loadReportPeriods()
  } catch (e) {
    alert(e?.message || 'Gagal finalisasi')
  }
}

const onPrint = async (row) => {
  try {
    const res = await apiService.printReportPeriod(row.id)
    alert('Print-ready: ' + (res?.message || 'OK'))
  } catch (e) {
    alert(e?.message || 'Gagal cetak')
  }
}

// Util
const toMonthYear = (year, month) => {
  const m = String(month).padStart(2, '0')
  return `${m}/${year}`
}

const formatCurrency = (n) => {
  if (n == null) return '0,00'
  try {
    return new Intl.NumberFormat('id-ID', { style: 'decimal', minimumFractionDigits: 2 }).format(n)
  } catch {
    return String(n)
  }
}
</script>

<style scoped>
.dashboard-container {
  background: #f8fafc;
  min-height: 100vh;
}

.identity-alert {
  border-left-width: 6px !important;
  border-radius: 10px;
  padding: 12px 16px !important;
}

.summary-table tbody tr td:first-child {
  color: #334155;
}
.summary-table tbody tr td:last-child {
  font-weight: 600;
  color: #111827;
}

.mutasi-table thead th {
  background: #009849;
  color: #ffffff;
  font-weight: 700;
}

.mutasi-table tbody tr:nth-child(even) {
  background: #f6fff9;
}

.section-card {
  border-radius: 12px;
}
.panel-card {
  padding: 16px;
  border-radius: 10px;
}

.filters-group {
  display: flex;
  align-items: center;
  gap: 16px;
}
.inline-control {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.status-select {
  min-width: 220px;
}

.aksi-buttons {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* Join left summary box and right table so terlihat menyatu */
/* reset penggabungan agar terpisah */

.dashboard-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  color: white !important;
}

.dashboard-header .v-card-title h1 {
  color: white !important;
}

.dashboard-header .v-card-title p {
  color: rgba(255, 255, 255, 0.8) !important;
}

.welcome-card {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%) !important;
  color: white !important;
}

.welcome-card .v-card-title {
  color: white !important;
}

.feature-card {
  transition: all 0.3s ease;
  cursor: pointer;
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15) !important;
}

/* Mobile responsive */
@media (max-width: 600px) {
  .dashboard-container {
    padding: 8px !important;
  }

  .table-responsive {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }

  .mutasi-table thead th,
  .mutasi-table tbody td {
    white-space: nowrap;
  }

  .controls-row {
    flex-direction: column;
    align-items: stretch;
    gap: 8px !important;
  }

  .filters-group {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .feature-card .v-card-text {
    padding: 16px !important;
  }

  .feature-card .v-avatar {
    width: 48px !important;
    height: 48px !important;
  }
}

/* Tablet responsive */
@media (min-width: 600px) and (max-width: 960px) {
  .dashboard-container {
    padding: 16px !important;
  }
}
</style>
