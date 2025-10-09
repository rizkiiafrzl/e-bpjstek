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
                        <td class="text-right">424101780000</td>
                      </tr>
                      <tr>
                        <td>Total Iuran dan Denda</td>
                        <td class="text-right">Rp0</td>
                      </tr>
                      <tr>
                        <td>Sisa Pembayaran Iuran Sebelumnya</td>
                        <td class="text-right">Rp0</td>
                      </tr>
                      <tr>
                        <td>Total Tagihan</td>
                        <td class="text-right">Rp0</td>
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
                    <v-btn color="primary" variant="elevated" prepend-icon="mdi-plus">
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
                          model-value="10"
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
                          model-value="Tampil Semua"
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
                          <th style="width: 160px">Aksi</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(row, idx) in demoRows" :key="idx">
                          <td>{{ idx + 1 }}</td>
                          <td>{{ row.bulan }}</td>
                          <td>{{ row.tk }}</td>
                          <td>{{ row.iuran }}</td>
                          <td>{{ row.denda }}</td>
                          <td>
                            <v-chip
                              :color="
                                row.status === 'Posting'
                                  ? 'success'
                                  : row.status === 'Draft'
                                    ? 'grey'
                                    : 'info'
                              "
                              size="small"
                              variant="tonal"
                            >
                              {{ row.status }}
                            </v-chip>
                          </td>
                          <td>
                            <div class="aksi-buttons">
                              <v-btn
                                size="x-small"
                                color="primary"
                                variant="tonal"
                                prepend-icon="mdi-pencil"
                                @click="navigateToEdit(row)"
                                >Edit</v-btn
                              >
                              <v-btn size="x-small" variant="tonal" prepend-icon="mdi-printer"
                                >Cetak</v-btn
                              >
                              <v-btn
                                size="x-small"
                                variant="outlined"
                                color="error"
                                prepend-icon="mdi-delete"
                                >Hapus</v-btn
                              >
                            </div>
                          </td>
                        </tr>
                      </tbody>
                    </v-table>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiService from '../services/api.js'

const router = useRouter()
const userEmail = ref('')
const loginTime = ref('')

// Demo rows untuk tabel mutasi (statik; nanti bisa diganti API)
const demoRows = ref([
  { bulan: '06/2025', tk: 4, iuran: '118.800,00', denda: '0,00', status: 'Draft' },
  { bulan: '05/2025', tk: 4, iuran: '118.800,00', denda: '0,00', status: 'Draft' },
  { bulan: '04/2025', tk: 4, iuran: '21.600,00', denda: '0,00', status: 'Draft' },
  { bulan: '03/2025', tk: 3, iuran: '40.500,00', denda: '0,00', status: 'Posting' },
])

onMounted(() => {
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
})

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
