<template>
  <div class="add-worker-page">
    <v-container>
      <h2 class="text-h5 mb-4">Data Tenaga Kerja</h2>

      <!-- Precheck Dialog: pilihan + konten muncul di bawah tanpa next step -->
      <v-dialog v-model="dlgPrecheck" max-width="780" persistent>
        <v-card class="step-card">
          <div class="banner banner-green">
            <div class="banner-title">
              Apakah tenaga kerja sudah memiliki kartu BPJS Ketenagakerjaan ?
            </div>
          </div>
          <v-card-text class="pa-6 d-flex flex-column align-center">
            <div class="option-row">
              <v-btn
                class="option-btn sudah"
                :class="{ active: hasCard === 'sudah' }"
                @click="selectHasCard('sudah')"
                variant="outlined"
              >
                <v-icon start>mdi-account-multiple</v-icon>
                SUDAH
              </v-btn>
              <v-btn
                class="option-btn belum"
                :class="{ active: hasCard === 'belum' }"
                @click="selectHasCard('belum')"
                variant="outlined"
              >
                <v-icon start>mdi-restore-alert</v-icon>
                BELUM
              </v-btn>
            </div>
            <div v-if="precheckTouched && hasCard === null" class="error-text mt-2">
              Wajib pilih salah satu
            </div>

            <!-- KPJ input langsung jika "Sudah" -->
            <div v-if="hasCard === 'sudah'" class="w-100 d-flex flex-column align-center mt-6">
              <v-text-field
                v-model="kpj"
                class="kpj-input"
                placeholder="Input No KPJ"
                variant="outlined"
              />
              <div v-if="precheckTouched && !kpj" class="error-text mt-1">Wajib isi No KPJ</div>
            </div>

            <!-- Pilihan WNI/WNA langsung jika "Belum" -->
            <div v-if="hasCard === 'belum'" class="w-100 d-flex flex-column align-center mt-6">
              <v-radio-group v-model="nationality" class="mt-2 nationality-group">
                <v-radio label="Warga Negara Indonesia (WNI)" value="WNI" />
                <v-radio label="Warga Negara Asing (WNA)" value="WNA" />
              </v-radio-group>
              <div v-if="precheckTouched && !nationality" class="error-text mt-1">
                Wajib pilih kewarganegaraan
              </div>
            </div>
          </v-card-text>
          <v-card-actions class="justify-center pb-6">
            <v-btn variant="text" class="link-next" @click="finishPrecheck">
              <v-icon start color="success">mdi-arrow-right</v-icon>
              <span class="text-success">LANJUT</span>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <!-- Form seperti gambar -->
      <v-card elevation="2" class="form-card pa-6 mt-4">
        <div class="form-header">
          <div class="title">FORM TENAGA KERJA</div>
          <div class="subtitle">Silakan masukkan data sesuai dengan KTP!</div>
        </div>

        <v-form @submit.prevent="submit" ref="formRef">
          <v-row>
            <v-col cols="12" md="4">
              <v-text-field
                v-model="form.nik"
                label="NIK"
                placeholder="Nomor e-KTP"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
            </v-col>
            <v-col cols="12" md="4">
              <v-text-field
                v-model="form.nama"
                label="Nama Lengkap"
                placeholder="Nama Sesuai e-KTP"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
            </v-col>
            <v-col cols="12" md="4">
              <v-text-field
                v-model="form.tgl"
                label="Tanggal Lahir"
                placeholder="dd-mm-yyyy"
                variant="outlined"
                :rules="[(v) => !!v || 'Wajib diisi']"
              />
            </v-col>
          </v-row>

          <v-row class="align-center">
            <v-col cols="12" md="4">
              <TraditionalCaptcha ref="captchaRef" @verified="onCaptcha" @error="onCaptchaError" />
            </v-col>
          </v-row>

          <div class="d-flex justify-end mt-2">
            <v-btn color="success" class="submit-btn" type="submit"> Daftar </v-btn>
          </div>
        </v-form>
      </v-card>

      <!-- Success Dialog -->
      <v-dialog v-model="dlgSuccess" max-width="520">
        <v-card class="add-worker-dialog">
          <v-card-text class="text-center pa-6">
            <div class="aw-title">Data berhasil didaftarkan</div>
            <div class="aw-desc">Tenaga kerja telah ditambahkan ke daftar periode aktif.</div>
            <div class="d-flex justify-center mt-4">
              <v-btn color="success" class="aw-btn aw-btn-green" @click="goBackToEdit"
                >Kembali ke Edit</v-btn
              >
            </div>
          </v-card-text>
        </v-card>
      </v-dialog>
    </v-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TraditionalCaptcha from '../components/TraditionalCaptcha.vue'

const router = useRouter()

const formRef = ref(null)
const form = ref({ nik: '', nama: '', tgl: '' })
const captchaValue = ref('')
const captchaVerified = ref(false)
const captchaRef = ref(null)

// Precheck (gabung jadi satu dialog)
const dlgPrecheck = ref(false)
const hasCard = ref(null)
const kpj = ref('')
const nationality = ref(null)
const precheckTouched = ref(false)

// Success dialog
const dlgSuccess = ref(false)

onMounted(() => {
  openPrechecks()
})

const openPrechecks = () => {
  dlgPrecheck.value = true
  hasCard.value = null
  kpj.value = ''
  nationality.value = null
  precheckTouched.value = false
}

const selectHasCard = (val) => {
  hasCard.value = val
}

const finishPrecheck = () => {
  precheckTouched.value = true
  if (hasCard.value === null) return
  if (hasCard.value === 'sudah' && !kpj.value) return
  if (hasCard.value === 'belum' && !nationality.value) return
  dlgPrecheck.value = false
}

const onCaptcha = () => {
  captchaVerified.value = true
}
const onCaptchaError = () => {
  captchaVerified.value = false
}

const submit = async () => {
  const { valid } = await formRef.value.validate()
  if (!valid) return
  if (!captchaVerified.value || !captchaValue.value) return
  dlgSuccess.value = true
}

const goBackToEdit = () => {
  dlgSuccess.value = false
  // Kembali ke halaman edit; gunakan payload kosong jika tidak ada data
  router.push('/edit/{}')
}
</script>

<style scoped>
.add-worker-page {
  background: #fff;
}

/* Dialog/banner styles */
.step-card {
  border-radius: 10px;
  overflow: hidden;
}
.banner {
  color: #fff;
  padding: 16px;
  text-align: center;
}
.banner-green {
  background: #009849;
}
.banner-title {
  font-weight: 700;
}
.option-row {
  display: flex;
  gap: 16px;
  align-items: center;
  justify-content: center;
}
.option-btn {
  min-width: 140px;
  height: 40px;
  border-width: 2px !important;
  border-radius: 8px;
}
.option-btn.sudah {
  color: #009849 !important;
  border-color: #009849 !important;
}
.option-btn.sudah :deep(.v-icon) {
  color: #009849 !important;
}
.option-btn.sudah.active,
.option-btn.sudah:hover {
  background: #009849 !important;
  color: #fff !important;
}
.option-btn.belum {
  color: #ff4d4f !important;
  border-color: #ff4d4f !important;
}
.option-btn.belum :deep(.v-icon) {
  color: #ff4d4f !important;
}
.option-btn.belum.active,
.option-btn.belum:hover {
  background: #ff4d4f !important;
  color: #fff !important;
}
.primary-action {
  min-width: 160px;
}
.kpj-input {
  width: 100%;
  max-width: 520px;
}
.nationality-group :deep(.v-selection-control) {
  margin-inline-start: 12px;
}

/* Form card styles */
.form-card {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
}
.form-header .title {
  font-weight: 700;
}
.form-header .subtitle {
  color: #6b7280;
  margin-bottom: 12px;
}
.submit-btn {
  min-width: 120px;
}
.banner-blue {
  background: #2f80ed;
}
.link-next {
  font-weight: 600;
}
.error-text {
  color: #d32f2f;
  font-size: 12px;
}
</style>
