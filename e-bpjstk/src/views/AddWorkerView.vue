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
      <v-card v-if="!dlgPrecheck" elevation="2" class="form-card pa-6 mt-4">
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
                :rules="[(v) => form.nationality === 'WNI' ? (!!v && /^\d{16}$/.test(v)) || 'NIK harus 16 digit' : true]"
                :error-messages="form.nationality === 'WNI' && form.nik && !/^\d{16}$/.test(form.nik) ? ['NIK harus 16 digit'] : []"
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
                    label="Tanggal Lahir"
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
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12" md="4">
              <v-select
                v-model="form.nationality"
                :items="['WNI','WNA']"
                label="Kewarganegaraan"
                variant="outlined"
                @update:model-value="(val) => console.log('Nationality changed to:', val)"
              />
            </v-col>
            <v-col cols="12" md="4" v-if="form.nationality === 'WNA'">
              <v-text-field
                v-model="form.passportNo"
                label="Nomor Paspor"
                placeholder="Nomor paspor"
                variant="outlined"
                :rules="[(v)=> !!v || 'Nomor paspor wajib untuk WNA']"
                required
              />
            </v-col>
            <v-col cols="12" md="4" v-if="form.nationality === 'WNA'">
              <v-menu
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

          <v-row class="align-center">
            <v-col cols="12" md="4">
              <TraditionalCaptcha ref="captchaRef" @verified="onCaptcha" @error="onCaptchaError" />
            </v-col>
          </v-row>

          <div class="d-flex justify-end mt-2">
            <v-btn 
              color="success" 
              class="submit-btn" 
              type="button"
              @click="submit"
            > 
              Daftar 
            </v-btn>
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

      <!-- Consent Dialog sebelum pindah ke form lengkap -->
      <v-dialog v-model="dlgConsent" max-width="520" persistent>
        <v-card>
          <v-card-text class="text-center pa-6">
            <div class="text-h5 font-weight-bold mb-3">Persetujuan</div>
            <div>
              Dengan ini saya menyatakan bahwa data yang disampaikan merupakan data yang sebenarnya dan bersedia jika data yang didaftarkan disimpan BPJS Ketenagakerjaan sebagai data peserta.
            </div>
            <div class="d-flex justify-center mt-6" style="gap: 12px;">
              <v-btn color="success" @click="agreeConsent">Setuju</v-btn>
              <v-btn color="error" variant="outlined" @click="rejectConsent">Tidak</v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-dialog>
    </v-container>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import TraditionalCaptcha from '../components/TraditionalCaptcha.vue'
import apiService from '../services/api.js'

const router = useRouter()
const props = defineProps({ id: { type: [String, Number], required: false } })

const formRef = ref(null)
const form = ref({ nik: '', nama: '', nationality: 'WNI', passportNo: '', passportValidUntil: '', kpj: '' })
const dob = ref('') // YYYY-MM-DD
const dobMenu = ref(false)
const passportValidUntil = ref('') // YYYY-MM-DD
const passportMenu = ref(false)
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
const dlgConsent = ref(false)

onMounted(async () => {
  if (props.id && String(props.id).length > 0) {
    // Editing mode: load worker and prefill
    try {
      const data = await apiService.getWorker(String(props.id))
      form.value.nik = data.nik || ''
      form.value.nama = data.nama || ''
      form.value.nationality = data.nationality || 'WNI'
      form.value.passportNo = data.passportNo || ''
      form.value.kpj = data.kpj || ''
      hasCard.value = data.kpj ? 'sudah' : 'belum'
      kpj.value = data.kpj || ''
      if (data.dateOfBirth) {
        dob.value = data.dateOfBirth
      }
      if (data.passportValidUntil) {
        passportValidUntil.value = data.passportValidUntil
        form.value.passportValidUntil = data.passportValidUntil
      }
      dlgPrecheck.value = false
    } catch (e) {
      console.error('Gagal memuat worker', e)
    }
  } else {
    openPrechecks()
  }
})

const openPrechecks = () => {
  dlgPrecheck.value = true
  hasCard.value = null
  kpj.value = ''
  nationality.value = null
  precheckTouched.value = false
}

const selectHasCard = (val) => {
  console.log('Card selection changed to:', val)
  hasCard.value = val
}

const finishPrecheck = () => {
  console.log('Finishing precheck...')
  console.log('hasCard:', hasCard.value)
  console.log('kpj:', kpj.value)
  console.log('nationality:', nationality.value)
  
  precheckTouched.value = true
  if (hasCard.value === null) {
    console.log('No card selection made')
    return
  }
  if (hasCard.value === 'sudah' && !kpj.value) {
    console.log('KPJ required for existing card')
    return
  }
  if (hasCard.value === 'belum' && !nationality.value) {
    console.log('Nationality required for new card')
    return
  }
  
  console.log('Precheck completed, closing dialog')
  dlgPrecheck.value = false
  
  // Set form nationality based on precheck selection
  if (hasCard.value === 'belum' && nationality.value) {
    form.value.nationality = nationality.value
    console.log('Form nationality set to:', form.value.nationality)
  }
  
  // Set KPJ to form if user has card
  if (hasCard.value === 'sudah' && kpj.value) {
    form.value.kpj = kpj.value
    console.log('Form KPJ set to:', form.value.kpj)
  }
}

const onCaptcha = (val) => {
  console.log('CAPTCHA verified with value:', val)
  captchaVerified.value = true
  captchaValue.value = val || '1'
  console.log('CAPTCHA state updated - verified:', captchaVerified.value, 'value:', captchaValue.value)
}
const onCaptchaError = () => {
  console.log('CAPTCHA error occurred')
  captchaVerified.value = false
  captchaValue.value = ''
}

const onPickDob = (val) => {
  // Normalisasi ke format YYYY-MM-DD
  try {
    const dt = typeof val === 'string' ? new Date(val) : val
    if (dt instanceof Date && !isNaN(dt)) {
      const y = dt.getFullYear()
      const m = String(dt.getMonth() + 1).padStart(2, '0')
      const d = String(dt.getDate()).padStart(2, '0')
      dob.value = `${y}-${m}-${d}`
    }
  } catch {
    // ignore parse error and keep original value
  }
  dobMenu.value = false
}

const onPickPassportValid = (val) => {
  // Normalisasi ke format YYYY-MM-DD
  try {
    const dt = typeof val === 'string' ? new Date(val) : val
    if (dt instanceof Date && !isNaN(dt)) {
      const y = dt.getFullYear()
      const m = String(dt.getMonth() + 1).padStart(2, '0')
      const d = String(dt.getDate()).padStart(2, '0')
      passportValidUntil.value = `${y}-${m}-${d}`
      form.value.passportValidUntil = passportValidUntil.value
    }
  } catch {
    // ignore parse error and keep original value
  }
  passportMenu.value = false
}

const displayDob = computed(() => {
  if (!dob.value) return ''
  const v = dob.value
  if (typeof v === 'string' && v.includes('-')) {
    const [y, m, d] = v.split('-')
    return `${d}-${m}-${y}`
  }
  if (v instanceof Date && !isNaN(v)) {
    const y = v.getFullYear()
    const m = String(v.getMonth() + 1).padStart(2, '0')
    const d = String(v.getDate()).padStart(2, '0')
    return `${d}-${m}-${y}`
  }
  return ''
})

const displayPassportValid = computed(() => {
  if (!passportValidUntil.value) return ''
  const v = passportValidUntil.value
  if (typeof v === 'string' && v.includes('-')) {
    const [y, m, d] = v.split('-')
    return `${d}-${m}-${y}`
  }
  if (v instanceof Date && !isNaN(v)) {
    const y = v.getFullYear()
    const m = String(v.getMonth() + 1).padStart(2, '0')
    const d = String(v.getDate()).padStart(2, '0')
    return `${d}-${m}-${y}`
  }
  return ''
})

const submit = async () => {
  console.log('Submit button clicked')
  console.log('Form data:', form.value)
  console.log('DOB:', dob.value)
  console.log('CAPTCHA verified:', captchaVerified.value)
  console.log('CAPTCHA value:', captchaValue.value)
  
  // Manual validation untuk NIK WNI
  if (form.value.nationality === 'WNI') {
    if (!form.value.nik || !/^\d{16}$/.test(form.value.nik)) {
      console.log('NIK validation failed for WNI')
      alert('NIK harus 16 digit untuk WNI')
      return
    }
  }
  
  // Manual validation untuk WNA
  if (form.value.nationality === 'WNA') {
    if (!form.value.passportNo || form.value.passportNo.trim() === '') {
      console.log('Passport number validation failed for WNA')
      alert('Nomor paspor wajib diisi untuk WNA')
      return
    }
    if (!passportValidUntil.value) {
      console.log('Passport valid until validation failed for WNA')
      alert('Masa berlaku paspor wajib diisi untuk WNA')
      return
    }
  }
  
  // Manual validation untuk nama
  if (!form.value.nama || form.value.nama.trim() === '') {
    console.log('Nama validation failed')
    alert('Nama lengkap wajib diisi')
    return
  }
  
  // Manual validation untuk tanggal lahir
  if (!dob.value) {
    console.log('Date of birth not filled')
    alert('Tanggal lahir wajib diisi')
    return
  }
  
  // Manual validation untuk CAPTCHA
  if (!captchaVerified.value || !captchaValue.value) {
    console.log('CAPTCHA not verified')
    alert('Mohon verifikasi CAPTCHA terlebih dahulu')
    return
  }
  
  console.log('All validations passed, showing consent dialog...')
  dlgConsent.value = true
}

const goBackToEdit = () => {
  dlgSuccess.value = false
  router.push('/dashboard')
}

// Consent handlers
const agreeConsent = async () => {
  // Kumpulkan data minimal lalu simpan ke backend sebelum navigate ke halaman edit detail
  const payload = {
    nik: form.value.nationality === 'WNI' ? form.value.nik : '',
    nama: form.value.nama,
    noPegawai: '',
    kpj: form.value.kpj || '',
    dateOfBirth: dob.value || '',
    upah: 0,
    rapel: 0,
    nationality: form.value.nationality,
    passportNo: form.value.nationality === 'WNA' ? form.value.passportNo : '',
    passportValidUntil: form.value.nationality === 'WNA' ? passportValidUntil.value : '',
  }
  
  // Debug logging
  console.log('=== DEBUG PAYLOAD ===')
  console.log('form.value.kpj:', form.value.kpj)
  console.log('kpj.value:', kpj.value)
  console.log('hasCard.value:', hasCard.value)
  console.log('Full payload:', payload)
  console.log('====================')
  try {
    const created = await apiService.createWorker(payload)
    dlgConsent.value = false
    const newId = created?.id || created?.worker?.id
    if (newId) {
      router.push(`/tenaga/form-lanjutan/${newId}`)
    } else {
      const dataParam = encodeURIComponent(JSON.stringify(payload))
      router.push(`/edit/${dataParam}`)
    }
  } catch (e) {
    alert(e?.message || 'Gagal menyimpan data tenaga kerja')
  }
}

const rejectConsent = () => {
  dlgConsent.value = false
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

