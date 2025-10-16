<template>
  <div class="edit-worker-page">
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

      <!-- Edit Form -->
      <v-card elevation="2" class="form-card pa-6">
        <div class="form-header">
          <div class="title">EDIT DATA TENAGA KERJA</div>
          <div class="subtitle">Silakan edit data sesuai dengan kebutuhan!</div>
        </div>

        <v-form @submit.prevent="handleSubmit" ref="formRef" v-if="!loading">
          <v-row dense>
            <!-- Left column -->
            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.nik"
                label="No. Identitas (NIK): *"
                placeholder="16 digit"
                variant="outlined"
                :rules="[(v) => form.nationality === 'WNI' ? (!!v && /^\d{16}$/.test(v)) || 'NIK harus 16 digit' : true]"
              />
              <v-text-field
                v-model="form.nama"
                label="Nama Lengkap: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-text-field
                v-model="form.tempatLahir"
                label="Tempat Lahir: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-menu
                v-model="dobMenu"
                :close-on-content-click="false"
                transition="scale-transition"
                offset-y
                max-width="290px"
                min-width="auto"
              >
                <template #activator="{ props }">
                  <v-text-field
                    v-bind="props"
                    :model-value="displayDob"
                    label="Tanggal Lahir: *"
                    placeholder="dd-mm-yyyy"
                    variant="outlined"
                    readonly
                    :rules="[(v) => !!dob || 'Wajib diisi']"
                  />
                </template>
                <v-date-picker
                  class="compact-date-picker"
                  v-model="dob"
                  hide-actions
                  locale="en-US"
                  :first-day-of-week="0"
                  :show-adjacent-months="true"
                  elevation="0"
                  rounded
                  :title="''"
                  :width="300"
                  @update:model-value="onPickDob"
                />
              </v-menu>
              <v-text-field
                v-model="form.ibuKandung"
                label="Ibu Kandung: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-select
                v-model="form.golDarah"
                :items="['A', 'B', 'AB', 'O']"
                label="Gol Darah: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-select
                v-model="form.jenisKelamin"
                :items="[
                  { title: 'Laki-laki', value: 'L' },
                  { title: 'Perempuan', value: 'P' }
                ]"
                label="Jenis Kelamin: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-select
                v-model="form.statusKawin"
                :items="[
                  { title: 'KAWIN', value: 'KAWIN' },
                  { title: 'BELUM KAWIN', value: 'BELUM KAWIN' }
                ]"
                label="Status Kawin: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
            </v-col>

            <!-- Right column -->
            <v-col cols="12" md="6">
              <v-select
                v-model="form.statusPegawai"
                :items="[
                  { title: 'PKWT (Perjanjian Kerja Waktu Tertentu)', value: 'PKWT' },
                  { title: 'PKWTT (Perjanjian Kerja Waktu Tidak Tertentu)', value: 'PKWTT' }
                ]"
                label="Status Pegawai: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-row dense>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="tanggalAwalBekerja"
                    label="Tanggal Awal Bekerja: *"
                    type="date"
                    variant="outlined"
                    :rules="[(v) => !!v || 'Wajib diisi']"
                  />
                </v-col>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="tanggalAkhirKontrak"
                    label="Tanggal Akhir Kontrak: *"
                    type="date"
                    variant="outlined"
                    :rules="[(v) => form.statusPegawai === 'PKWT' ? (!!v || 'Wajib diisi untuk PKWT') : true]"
                  />
                </v-col>
              </v-row>
              <v-text-field
                v-model="form.noPegawai"
                label="Nomor Pegawai"
                variant="outlined"
              />
              <v-select
                v-model="form.lokasiPekerjaan"
                :items="lokasiPekerjaanItems"
                label="Lokasi Pekerjaan: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-text-field
                v-model.number="form.upah"
                label="Upah: *"
                type="number"
                variant="outlined"
                :rules="[(v) => v >= 0 || 'Upah tidak boleh negatif']"
              />
              <v-text-field
                v-model.number="form.rapel"
                label="Rapel"
                type="number"
                variant="outlined"
                :rules="[(v) => v >= 0 || 'Rapel tidak boleh negatif']"
              />
              <v-text-field
                v-model="form.alamat"
                label="Alamat: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-select
                v-model="form.kabupaten"
                :items="kabupatenItems"
                label="Kabupaten: *"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
              <v-text-field
                v-model="form.kodePos"
                label="Kode Pos: *"
                variant="outlined"
                :rules="[(v) => !v || /^\d{5}$/.test(v) || 'Kode pos 5 digit']"
              />
              <v-text-field
                v-model="form.teleponRumah"
                label="Nomor Telepon Rumah"
                variant="outlined"
              />
              <v-text-field
                v-model="form.handphone"
                label="Handphone"
                variant="outlined"
              />
              <v-text-field
                v-model="form.npwp"
                label="NPWP"
                variant="outlined"
              />
              <v-text-field
                v-model="form.email"
                label="Email"
                variant="outlined"
                type="email"
              />
              <v-text-field
                v-model="form.kpj"
                label="KPJ"
                variant="outlined"
              />
              <v-select
                v-model="form.nationality"
                :items="['WNI','WNA']"
                label="Kewarganegaraan"
                variant="outlined"
              />
              <v-text-field
                v-if="form.nationality === 'WNA'"
                v-model="form.passportNo"
                label="Nomor Paspor"
                variant="outlined"
                :rules="[(v)=> !!v || 'Nomor paspor wajib untuk WNA']"
                required
              />
              <v-menu
                v-if="form.nationality === 'WNA'"
                v-model="passportMenu"
                :close-on-content-click="false"
                transition="scale-transition"
                offset-y
                max-width="290px"
                min-width="auto"
              >
                <template #activator="{ props }">
                  <v-text-field
                    v-bind="props"
                    :model-value="displayPassportValid"
                    label="Masa Berlaku Paspor"
                    placeholder="dd-mm-yyyy"
                    variant="outlined"
                    :rules="[(v) => !!passportValidUntil || 'Wajib diisi untuk WNA']"
                    readonly
                    required
                  />
                </template>
                <v-date-picker
                  class="compact-date-picker"
                  v-model="passportValidUntil"
                  hide-actions
                  locale="en-US"
                  :first-day-of-week="0"
                  :show-adjacent-months="true"
                  elevation="0"
                  rounded
                  :title="''"
                  :width="300"
                  @update:model-value="onPickPassportValid"
                />
              </v-menu>
            </v-col>
          </v-row>

          <div class="d-flex justify-end mt-4" style="gap: 12px;">
            <v-btn 
              color="grey" 
              variant="outlined"
              @click="goBack"
            > 
              Batal 
            </v-btn>
            <v-btn 
              color="primary" 
              type="submit"
              :loading="isSubmitting"
            > 
              Simpan Perubahan 
            </v-btn>
          </div>
        </v-form>

        <!-- Loading State -->
        <div v-if="loading" class="text-center pa-8">
          <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
          <div class="mt-4">Memuat data worker...</div>
        </div>

        <!-- Error State -->
        <div v-if="error" class="text-center pa-8">
          <v-alert type="error" variant="tonal">
            {{ error }}
          </v-alert>
          <v-btn color="primary" @click="loadWorkerData" class="mt-4">
            Coba Lagi
          </v-btn>
        </div>
      </v-card>
    </v-container>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import apiService from '../services/api.js'

const router = useRouter()
const route = useRoute()
const workerId = route.params.id

// Breadcrumbs
const breadcrumbs = ref([
  { title: 'Data Tenaga Kerja', disabled: false, href: '/edit' },
  { title: 'Edit Tenaga Kerja', disabled: true },
])

// Form state
const formRef = ref(null)
const loading = ref(true)
const error = ref(null)
const isSubmitting = ref(false)

// Form data
const form = reactive({
  nik: '',
  nama: '',
  kpj: '',
  noPegawai: '',
  nationality: 'WNI',
  upah: 0,
  rapel: 0,
  passportNo: '',
  // Personal info
  tempatLahir: '',
  ibuKandung: '',
  jenisKelamin: '',
  golDarah: '',
  statusKawin: '',
  // Employment info
  statusPegawai: '',
  lokasiPekerjaan: '',
  // Contact & address
  alamat: '',
  kabupaten: '',
  kodePos: '',
  teleponRumah: '',
  handphone: '',
  npwp: '',
  email: '',
})

// Date fields
const dob = ref('')
const dobMenu = ref(false)
const passportValidUntil = ref('')
const passportMenu = ref(false)
const tanggalAwalBekerja = ref('')
const tanggalAkhirKontrak = ref('')

// Dropdown items
const lokasiPekerjaanItems = [
  'ACEH BARAT','ACEH BARAT DAYA','ACEH BESAR','ACEH JAYA','ACEH SELATAN','ACEH TENGAH','ACEH TENGGARA','BANDA ACEH','MEDAN','DELISERDANG','TEBING TINGGI','PADANG','BUKITTINGGI','PEKANBARU','BATAM','TANJUNGPINANG','JAMBI','PALEMBANG','BENGKULU','BANDAR LAMPUNG','JAKARTA PUSAT','JAKARTA SELATAN','JAKARTA TIMUR','JAKARTA BARAT','JAKARTA UTARA','TANGERANG','BEKASI','DEPOK','BOGOR','BANDUNG','CIREBON','CIMAHI','SEMARANG','SOLO','MAGELANG','PEKALONGAN','TEGAL','YOGYAKARTA','SURABAYA','MALANG','KEDIRI','MADIUN','GRESIK','SIDOARJO','SERANG','CILEGON','DENPASAR','MATARAM','KUPANG','PONTIANAK','BANJARMASIN','BALIKPAPAN','SAMARINDA','TANJUNGPINANG','MANADO','PALU','MAKASSAR','KENDARI','AMBON','TERNATE','JAYAPURA'
]

const kabupatenItems = [
  'ACEH BARAT','ACEH BARAT DAYA','ACEH BESAR','ACEH JAYA','ACEH SELATAN','ACEH SINGKIL','ACEH TAMIANG','ACEH TENGAH','ACEH TENGGARA','ACEH TIMUR','ACEH UTARA','BENER MERIAH','PIDIE','PIDIE JAYA','SIMEULUE','KOTA BANDA ACEH','KOTA LHOKSEUMAWE','KOTA LANGSA','KOTA SABANG','KARO','DELI SERDANG','LANGKAT','ASAHAN','BATU BARA','LABUHANBATU','LABUHANBATU UTARA','LABUHANBATU SELATAN','SIMALUNGUN','SAMOSIR','TOBA','HUMBANG HASUNDUTAN','DAIRI','TAPANULI UTARA','TAPANULI TENGAH','TAPANULI SELATAN','NIAS','NIAS UTARA','NIAS SELATAN','NIAS BARAT','KOTA MEDAN','KOTA BINJAI','KOTA TEBING TINGGI','AGAM','LIMA PULUH KOTA','PADANG PARIAMAN','PASAMAN','PASAMAN BARAT','PESISIR SELATAN','SIJUNJUNG','SOLOK','TANAH DATAR','KOTA PADANG','KOTA BUKITTINGGI','KOTA PAYAKUMBUH','ROKAN HILIR','ROKAN HULU','SIAK','PELALAWAN','INDRAGIRI HULU','INDRAGIRI HILIR','KAMPAR','BENGKALIS','KEPULAUAN MERANTI','KOTA PEKANBARU','BINTAN','KARIMUN','LINGGA','NATUNA','ANAMBAS','KOTA BATAM','KOTA TANJUNGPINANG','BANGKA','BELITUNG','BANGKA BARAT','BANGKA SELATAN','BANGKA TENGAH','BELITUNG TIMUR','KOTA PANGKAL PINANG','KEPULAUAN SERIBU','KOTA ADM. JAKARTA PUSAT','KOTA ADM. JAKARTA UTARA','KOTA ADM. JAKARTA BARAT','KOTA ADM. JAKARTA SELATAN','KOTA ADM. JAKARTA TIMUR','BOGOR','SUKABUMI','CIANJUR','BANDUNG','GARUT','TASIKMALAYA','CIAMIS','KUNINGAN','CIREBON','MAJALENGKA','SUMEDANG','INDRAMAYU','SUBANG','PURWAKARTA','KARAWANG','BEKASI','KOTA BOGOR','KOTA SUKABUMI','KOTA BANDUNG','KOTA CIREBON','KOTA BEKASI','KOTA DEPOK','BANJAR','SEMARANG','KUDUS','PATi','REMBANG','BLORA','GROBOGAN','DEMAK','KENDAL','BATANG','PEKALONGAN','PEMALANG','TEGAL','BREBES','MAGELANG','TEMANGGUNG','WONOSOBO','PURWOREJO','KEBUMEN','BANYUMAS','PURBALINGGA','BANJARNEGARA','CILACAP','SLEMAN','BANTUL','GUNUNGKIDUL','KULON PROGO','KOTA YOGYAKARTA','SURABAYA','SIDOARJO','GRESIK','LAMONGAN','TUBAN','BOJONEGORO','MADIUN','NGAWI','MAGETAN','PONOROGO','PACITAN','KEDIRI','BLITAR','TULUNGAGUNG','TRENGGALEK','MALANG','PASURUAN','PROBOLINGGO','LUMAJANG','JEMBER','BONDOWOSO','SITUBONDO','BANYUWANGI','PANDEGLANG','LEBAK','SERANG','TANGERANG','KOTA SERANG','KOTA CILEGON','KOTA TANGERANG','KOTA TANGERANG SELATAN','DENPASAR','BADUNG','TABANAN','BANGLI','GIANYAR','KLUNGKUNG','BULELENG','JEMBRANA','MATARAM','LOMBOK BARAT','LOMBOK TENGAH','LOMBOK TIMUR','SUMBAWA','DOMPU','BIMA','KUPANG','TIMOR TENGAH SELATAN','TIMOR TENGAH UTARA','BELU','ALOR','ENDE','Sikka','FLORES TIMUR','SUMBA TIMUR','SUMBA BARAT','PONTIANAK','KUBU RAYA','SAMBAS','SINGKAWANG','SANGGAU','KETAPANG','SINTANG','MELAWI','KAPUAS HULU','LANDAK','BANJARMASIN','BANJARBARU','TANAH LAUT','TANAH BUMBU','KOTABARU','HULU SUNGAI SELATAN','BARITO KUALA','BALIKPAPAN','SAMARINDA','BONTANG','PASER','PENAJAM PASER UTARA','KUTAI BARAT','KUTAI TIMUR','KUTAI KARTANEGARA','MANADO','BITUNG','TOMOHON','MINAHASA','BOLAANG MONGONDOW','PALU','SIGI','DONGGALA','PARIGI MOUTONG','MAKASSAR','GOWA','TAKALAR','MAROS','BONE','WATAMPONE','PAREPARE','KENDARI','KONAWE','KONAWE SELATAN','AMBON','MALUKU TENGAH','TERNATE','TIDORE','HALMAHERA BARAT','HALMAHERA TENGAH','JAYAPURA','MERAUKE'
]

// Load worker data
const loadWorkerData = async () => {
  if (!workerId) {
    error.value = 'ID worker tidak ditemukan'
    loading.value = false
    return
  }

  try {
    loading.value = true
    error.value = null
    
    const worker = await apiService.getWorker(workerId)
    
    // Populate form with worker data
    form.nik = worker.nik || ''
    form.nama = worker.nama || ''
    form.kpj = worker.kpj || ''
    form.noPegawai = worker.noPegawai || ''
    form.nationality = worker.nationality || 'WNI'
    form.upah = worker.upah || 0
    form.rapel = worker.rapel || 0
    form.passportNo = worker.passportNo || ''
    // Personal info
    form.tempatLahir = worker.tempatLahir || ''
    form.ibuKandung = worker.ibuKandung || ''
    form.jenisKelamin = worker.jenisKelamin || ''
    form.golDarah = worker.golDarah || ''
    form.statusKawin = worker.statusKawin || ''
    // Employment info
    form.statusPegawai = worker.statusPegawai || ''
    form.lokasiPekerjaan = worker.lokasiPekerjaan || ''
    // Contact & address
    form.alamat = worker.alamat || ''
    form.kabupaten = worker.kabupaten || ''
    form.kodePos = worker.kodePos || ''
    form.teleponRumah = worker.teleponRumah || ''
    form.handphone = worker.handphone || ''
    form.npwp = worker.npwp || ''
    form.email = worker.email || ''
    
    // Set dates - normalize format
    if (worker.dateOfBirth) {
      // Handle ISO format and convert to YYYY-MM-DD
      if (worker.dateOfBirth.includes('T')) {
        dob.value = worker.dateOfBirth.split('T')[0]
      } else {
        dob.value = worker.dateOfBirth
      }
    }
    if (worker.passportValidUntil) {
      // Handle ISO format and convert to YYYY-MM-DD
      if (worker.passportValidUntil.includes('T')) {
        passportValidUntil.value = worker.passportValidUntil.split('T')[0]
      } else {
        passportValidUntil.value = worker.passportValidUntil
      }
    }
    if (worker.tanggalAwalBekerja) {
      // Handle ISO format and convert to YYYY-MM-DD
      if (worker.tanggalAwalBekerja.includes('T')) {
        tanggalAwalBekerja.value = worker.tanggalAwalBekerja.split('T')[0]
      } else {
        tanggalAwalBekerja.value = worker.tanggalAwalBekerja
      }
    }
    if (worker.tanggalAkhirKontrak) {
      // Handle ISO format and convert to YYYY-MM-DD
      if (worker.tanggalAkhirKontrak.includes('T')) {
        tanggalAkhirKontrak.value = worker.tanggalAkhirKontrak.split('T')[0]
      } else {
        tanggalAkhirKontrak.value = worker.tanggalAkhirKontrak
      }
    }
    
  } catch (e) {
    console.error('Error loading worker:', e)
    error.value = e?.message || 'Gagal memuat data worker'
  } finally {
    loading.value = false
  }
}

// Date picker handlers
const onPickDob = (val) => {
  try {
    const dt = typeof val === 'string' ? new Date(val) : val
    if (dt instanceof Date && !isNaN(dt)) {
      const y = dt.getFullYear()
      const m = String(dt.getMonth() + 1).padStart(2, '0')
      const d = String(dt.getDate()).padStart(2, '0')
      dob.value = `${y}-${m}-${d}`
    }
  } catch {
    // ignore parse error
  }
  dobMenu.value = false
}

const onPickPassportValid = (val) => {
  try {
    const dt = typeof val === 'string' ? new Date(val) : val
    if (dt instanceof Date && !isNaN(dt)) {
      const y = dt.getFullYear()
      const m = String(dt.getMonth() + 1).padStart(2, '0')
      const d = String(dt.getDate()).padStart(2, '0')
      passportValidUntil.value = `${y}-${m}-${d}`
    }
  } catch {
    // ignore parse error
  }
  passportMenu.value = false
}

// Display computed properties
const displayDob = computed(() => {
  if (!dob.value) return ''
  const v = dob.value
  
  // Handle different date formats
  if (typeof v === 'string') {
    // Handle ISO format like "2025-10-01T00:00:00Z"
    if (v.includes('T')) {
      const datePart = v.split('T')[0]
      const [y, m, d] = datePart.split('-')
      return `${d}-${m}-${y}`
    }
    // Handle YYYY-MM-DD format
    if (v.includes('-') && v.length === 10) {
      const [y, m, d] = v.split('-')
      return `${d}-${m}-${y}`
    }
  }
  return ''
})

const displayPassportValid = computed(() => {
  if (!passportValidUntil.value) return ''
  const v = passportValidUntil.value
  
  // Handle different date formats
  if (typeof v === 'string') {
    // Handle ISO format like "2025-10-01T00:00:00Z"
    if (v.includes('T')) {
      const datePart = v.split('T')[0]
      const [y, m, d] = datePart.split('-')
      return `${d}-${m}-${y}`
    }
    // Handle YYYY-MM-DD format
    if (v.includes('-') && v.length === 10) {
      const [y, m, d] = v.split('-')
      return `${d}-${m}-${y}`
    }
  }
  return ''
})

// Form submission
const handleSubmit = async () => {
  if (!formRef.value) return
  
  const { valid } = await formRef.value.validate()
  if (!valid) return

  // Manual validation
  if (form.nationality === 'WNI' && (!form.nik || !/^\d{16}$/.test(form.nik))) {
    alert('NIK harus 16 digit untuk WNI')
    return
  }

  if (form.nationality === 'WNA' && !form.passportNo) {
    alert('Nomor paspor wajib untuk WNA')
    return
  }

  if (!form.nama.trim()) {
    alert('Nama lengkap wajib diisi')
    return
  }

  isSubmitting.value = true

  try {
    await apiService.updateWorker(workerId, {
      nik: form.nik,
      nama: form.nama,
      kpj: form.kpj,
      noPegawai: form.noPegawai,
      nationality: form.nationality,
      upah: form.upah,
      rapel: form.rapel,
      passportNo: form.passportNo,
      dateOfBirth: dob.value,
      passportValidUntil: passportValidUntil.value,
      // Personal info
      tempatLahir: form.tempatLahir,
      ibuKandung: form.ibuKandung,
      jenisKelamin: form.jenisKelamin,
      golDarah: form.golDarah,
      statusKawin: form.statusKawin,
      // Employment info
      statusPegawai: form.statusPegawai,
      lokasiPekerjaan: form.lokasiPekerjaan,
      tanggalAwalBekerja: tanggalAwalBekerja.value,
      tanggalAkhirKontrak: tanggalAkhirKontrak.value,
      // Contact & address
      alamat: form.alamat,
      kabupaten: form.kabupaten,
      kodePos: form.kodePos,
      teleponRumah: form.teleponRumah,
      handphone: form.handphone,
      npwp: form.npwp,
      email: form.email,
    })

    alert('Data worker berhasil diperbarui!')
    goBack()
  } catch (e) {
    console.error('Error updating worker:', e)
    alert(e?.message || 'Gagal memperbarui data worker')
  } finally {
    isSubmitting.value = false
  }
}

// Navigation
const goBack = () => {
  router.push('/dashboard')
}

// Load data on mount
onMounted(() => {
  loadWorkerData()
})
</script>

<style scoped>
.edit-worker-page {
  background: #f8fafc;
  min-height: 100vh;
}

.form-card {
  border-radius: 12px;
}

.form-header .title {
  font-weight: 700;
  font-size: 1.5rem;
  color: #1f2937;
}

.form-header .subtitle {
  color: #6b7280;
  margin-bottom: 24px;
}

/* Compact date-picker look */
.compact-date-picker :deep(.v-date-picker-title) {
  display: none;
}

.compact-date-picker :deep(.v-picker-title) {
  display: none !important;
}

.compact-date-picker :deep(.v-date-picker) {
  width: 300px !important;
}

.compact-date-picker :deep(.v-date-picker-header__content) {
  font-weight: 600;
}

.compact-date-picker :deep(.v-date-picker-month__day) {
  height: 32px;
  width: 32px;
}
</style>