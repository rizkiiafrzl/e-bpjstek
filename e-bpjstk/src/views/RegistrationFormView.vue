<template>
  <v-container class="py-6">
    <v-card elevation="2">
      <v-card-title class="text-h6 font-weight-bold">Formulir Pendaftaran Tenaga Kerja</v-card-title>
      <v-divider></v-divider>
      <v-card-text>
        <v-form ref="formRef" @submit.prevent="onSubmit">
          <v-row dense>
            <!-- Left column -->
            <v-col cols="12" md="6">
              <v-text-field v-model="form.nik" label="No. Identitas (NIK): *" variant="outlined" placeholder="16 digit" :rules="[rules.requiredIfWNI, rules.nikIfWNI]" />
              <v-text-field v-model="form.namaLengkap" label="Nama Lengkap: *" variant="outlined" :rules="[rules.required]" />
              <v-text-field v-model="form.tempatLahir" label="Tempat Lahir: *" variant="outlined" :rules="[rules.required]" />
              <v-text-field v-model="form.tanggalLahir" label="Tanggal Lahir: *" type="date" variant="outlined" :rules="[rules.required]" />
              <v-text-field v-model="form.kpj" label="KPJ" variant="outlined" placeholder="Nomor KPJ" />
              <v-text-field v-model="form.nomorPegawai" label="No Pegawai" variant="outlined" placeholder="Nomor Pegawai" />
              <v-text-field v-model="form.ibuKandung" label="Ibu Kandung: *" variant="outlined" :rules="[rules.required]" />
              <v-select v-model="form.golDarah" :items="golDarahItems" label="Gol Darah: *" variant="outlined" :rules="[rules.required]" />
              <v-select v-model="form.jenisKelamin" :items="jenisKelaminItems" label="Jenis Kelamin: *" variant="outlined" :rules="[rules.required]" />
              <v-select v-model="form.statusKawin" :items="statusKawinItems" label="Status Kawin: *" variant="outlined" :rules="[rules.required]" />
            </v-col>

            <!-- Right column -->
            <v-col cols="12" md="6">
              <v-select v-model="form.statusPegawai" :items="statusPegawaiItems" label="Status Pegawai: *" variant="outlined" :rules="[rules.required]" />
              <v-row dense>
                <v-col cols="12" md="6">
                  <v-text-field v-model="form.tanggalAwalBekerja" label="Tanggal Awal Bekerja: *" type="date" variant="outlined" :rules="[rules.required]" />
                </v-col>
                <v-col cols="12" md="6">
                  <v-text-field v-model="form.tanggalAkhirKontrak" label="Tanggal Akhir Kontrak: *" type="date" variant="outlined" :rules="[rules.requiredIfPKWT]" />
                </v-col>
              </v-row>
              <v-select v-model="form.lokasiPekerjaan" :items="lokasiPekerjaanItems" label="Lokasi Pekerjaan: *" variant="outlined" :rules="[rules.required]" />
              <v-text-field v-model.number="form.upah" label="Upah: *" type="number" variant="outlined" :rules="[rules.required, rules.nonNegative]" />
              <v-text-field v-model="form.alamat" label="Alamat: *" variant="outlined" :rules="[rules.required]" />
              <v-select v-model="form.kabupaten" :items="kabupatenItems" label="Kabupaten: *" variant="outlined" :rules="[rules.required]" />
              <v-text-field v-model="form.kodePos" label="Kode Pos: *" variant="outlined" :rules="[rules.required, rules.kodePos]" />
              <v-text-field v-model="form.teleponRumah" label="Nomor Telepon Rumah" variant="outlined" />
              <v-text-field v-model="form.handphone" label="Handphone" variant="outlined" />
              <v-text-field v-model="form.npwp" label="NPWP" variant="outlined" />
              <v-text-field v-model="form.email" label="Email" variant="outlined" type="email" />
            </v-col>
          </v-row>

          <div class="d-flex justify-end mt-4" style="gap: 8px;">
            <v-btn color="grey" variant="tonal" @click="onReset">Reset</v-btn>
            <v-btn color="primary" type="submit">Daftar</v-btn>
          </div>
        </v-form>
      </v-card-text>
    </v-card>

    <!-- Notifikasi Pemberitahuan -->
    <v-dialog v-model="notifOpen" persistent max-width="400">
      <v-card class="mx-auto">
        <v-card-text class="text-center pa-4">
          <div class="text-h6 font-weight-bold mb-3">Pemberitahuan</div>
          <div class="text-body-1 mb-4">
            Harap cek kembali! Jika data sudah benar, pilih "Daftar" untuk melanjutkan pendaftaran.
          </div>
          <div class="d-flex justify-center" style="gap: 12px;">
            <v-btn color="primary" @click="confirmDaftar">Daftar</v-btn>
            <v-btn color="grey" variant="outlined" @click="notifOpen = false">Tidak</v-btn>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- Popup Berhasil Didaftarkan - Simple Alert Style -->
    <v-dialog v-model="successOpen" persistent max-width="400">
      <v-card class="mx-auto">
        <v-card-text class="text-center pa-4">
          <div class="text-body-1 mb-4">Data worker berhasil diperbarui!</div>
          <v-btn color="primary" @click="goBackToEdit">OK</v-btn>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import apiService from '../services/api.js'

const route = useRoute()
const router = useRouter()
const workerId = route.params.id

const formRef = ref(null)
const workerNationality = ref('WNI')
const notifOpen = ref(false)
const successOpen = ref(false)
const workerData = ref(null) // Store worker data for later use

const form = reactive({
  nik: '',
  namaLengkap: '',
  kpj: '',
  tempatLahir: '',
  tanggalLahir: '',
  ibuKandung: '',
  golDarah: null,
  jenisKelamin: null,
  statusKawin: null,
  statusPegawai: null,
  tanggalAwalBekerja: '',
  tanggalAkhirKontrak: '',
  nomorPegawai: '',
  lokasiPekerjaan: null,
  upah: null,
  alamat: '',
  kabupaten: null,
  kodePos: '',
  teleponRumah: '',
  handphone: '',
  npwp: '',
  email: '',
})

const statusPegawaiItems = [
  { title: '-Pilih-', value: null },
  { title: 'PKWT (Perjanjian Kerja Waktu Tertentu)', value: 'PKWT' },
  { title: 'PKWTT (Perjanjian Kerja Waktu Tidak Tertentu)', value: 'PKWTT' },
]
const golDarahItems = ['A', 'B', 'AB', 'O']
const jenisKelaminItems = [
  { title: 'Laki-laki', value: 'L' },
  { title: 'Perempuan', value: 'P' },
]
const statusKawinItems = [
  { title: '-Pilih-', value: null },
  { title: 'KAWIN', value: 'KAWIN' },
  { title: 'BELUM KAWIN', value: 'BELUM KAWIN' },
]

// Placeholder lists resembling long options (expanded)
const lokasiPekerjaanItems = [
  '-Pilih-','ACEH BARAT','ACEH BARAT DAYA','ACEH BESAR','ACEH JAYA','ACEH SELATAN','ACEH TENGAH','ACEH TENGGARA','BANDA ACEH','MEDAN','DELISERDANG','TEBING TINGGI','PADANG','BUKITTINGGI','PEKANBARU','BATAM','TANJUNGPINANG','JAMBI','PALEMBANG','BENGKULU','BANDAR LAMPUNG','JAKARTA PUSAT','JAKARTA SELATAN','JAKARTA TIMUR','JAKARTA BARAT','JAKARTA UTARA','TANGERANG','BEKASI','DEPOK','BOGOR','BANDUNG','CIREBON','CIMAHI','SEMARANG','SOLO','MAGELANG','PEKALONGAN','TEGAL','YOGYAKARTA','SURABAYA','MALANG','KEDIRI','MADIUN','GRESIK','SIDOARJO','SERANG','CILEGON','DENPASAR','MATARAM','KUPANG','PONTIANAK','BANJARMASIN','BALIKPAPAN','SAMARINDA','TANJUNGPINANG','MANADO','PALU','MAKASSAR','KENDARI','AMBON','TERNATE','JAYAPURA'
]
const kabupatenItems = [
  '-Pilih-','ACEH BARAT','ACEH BARAT DAYA','ACEH BESAR','ACEH JAYA','ACEH SELATAN','ACEH SINGKIL','ACEH TAMIANG','ACEH TENGAH','ACEH TENGGARA','ACEH TIMUR','ACEH UTARA','BENER MERIAH','PIDIE','PIDIE JAYA','SIMEULUE','KOTA BANDA ACEH','KOTA LHOKSEUMAWE','KOTA LANGSA','KOTA SABANG','KARO','DELI SERDANG','LANGKAT','ASAHAN','BATU BARA','LABUHANBATU','LABUHANBATU UTARA','LABUHANBATU SELATAN','SIMALUNGUN','SAMOSIR','TOBA','HUMBANG HASUNDUTAN','DAIRI','TAPANULI UTARA','TAPANULI TENGAH','TAPANULI SELATAN','NIAS','NIAS UTARA','NIAS SELATAN','NIAS BARAT','KOTA MEDAN','KOTA BINJAI','KOTA TEBING TINGGI','AGAM','LIMA PULUH KOTA','PADANG PARIAMAN','PASAMAN','PASAMAN BARAT','PESISIR SELATAN','SIJUNJUNG','SOLOK','TANAH DATAR','KOTA PADANG','KOTA BUKITTINGGI','KOTA PAYAKUMBUH','ROKAN HILIR','ROKAN HULU','SIAK','PELALAWAN','INDRAGIRI HULU','INDRAGIRI HILIR','KAMPAR','BENGKALIS','KEPULAUAN MERANTI','KOTA PEKANBARU','BINTAN','KARIMUN','LINGGA','NATUNA','ANAMBAS','KOTA BATAM','KOTA TANJUNGPINANG','BANGKA','BELITUNG','BANGKA BARAT','BANGKA SELATAN','BANGKA TENGAH','BELITUNG TIMUR','KOTA PANGKAL PINANG','KEPULAUAN SERIBU','KOTA ADM. JAKARTA PUSAT','KOTA ADM. JAKARTA UTARA','KOTA ADM. JAKARTA BARAT','KOTA ADM. JAKARTA SELATAN','KOTA ADM. JAKARTA TIMUR','BOGOR','SUKABUMI','CIANJUR','BANDUNG','GARUT','TASIKMALAYA','CIAMIS','KUNINGAN','CIREBON','MAJALENGKA','SUMEDANG','INDRAMAYU','SUBANG','PURWAKARTA','KARAWANG','BEKASI','KOTA BOGOR','KOTA SUKABUMI','KOTA BANDUNG','KOTA CIREBON','KOTA BEKASI','KOTA DEPOK','BANJAR','SEMARANG','KUDUS','PATi','REMBANG','BLORA','GROBOGAN','DEMAK','KENDAL','BATANG','PEKALONGAN','PEMALANG','TEGAL','BREBES','MAGELANG','TEMANGGUNG','WONOSOBO','PURWOREJO','KEBUMEN','BANYUMAS','PURBALINGGA','BANJARNEGARA','CILACAP','SLEMAN','BANTUL','GUNUNGKIDUL','KULON PROGO','KOTA YOGYAKARTA','SURABAYA','SIDOARJO','GRESIK','LAMONGAN','TUBAN','BOJONEGORO','MADIUN','NGAWI','MAGETAN','PONOROGO','PACITAN','KEDIRI','BLITAR','TULUNGAGUNG','TRENGGALEK','MALANG','PASURUAN','PROBOLINGGO','LUMAJANG','JEMBER','BONDOWOSO','SITUBONDO','BANYUWANGI','PANDEGLANG','LEBAK','SERANG','TANGERANG','KOTA SERANG','KOTA CILEGON','KOTA TANGERANG','KOTA TANGERANG SELATAN','DENPASAR','BADUNG','TABANAN','BANGLI','GIANYAR','KLUNGKUNG','BULELENG','JEMBRANA','MATARAM','LOMBOK BARAT','LOMBOK TENGAH','LOMBOK TIMUR','SUMBAWA','DOMPU','BIMA','KUPANG','TIMOR TENGAH SELATAN','TIMOR TENGAH UTARA','BELU','ALOR','ENDE','Sikka','FLORES TIMUR','SUMBA TIMUR','SUMBA BARAT','PONTIANAK','KUBU RAYA','SAMBAS','SINGKAWANG','SANGGAU','KETAPANG','SINTANG','MELAWI','KAPUAS HULU','LANDAK','BANJARMASIN','BANJARBARU','TANAH LAUT','TANAH BUMBU','KOTABARU','HULU SUNGAI SELATAN','BARITO KUALA','BALIKPAPAN','SAMARINDA','BONTANG','PASER','PENAJAM PASER UTARA','KUTAI BARAT','KUTAI TIMUR','KUTAI KARTANEGARA','MANADO','BITUNG','TOMOHON','MINAHASA','BOLAANG MONGONDOW','PALU','SIGI','DONGGALA','PARIGI MOUTONG','MAKASSAR','GOWA','TAKALAR','MAROS','BONE','WATAMPONE','PAREPARE','KENDARI','KONAWE','KONAWE SELATAN','AMBON','MALUKU TENGAH','TERNATE','TIDORE','HALMAHERA BARAT','HALMAHERA TENGAH','JAYAPURA','MERAUKE'
]

const rules = {
  required: v => !!v || 'Wajib diisi',
  requiredIfWNI: v => (!!v || workerNationality.value !== 'WNI') || 'Wajib diisi',
  nikIfWNI: v => (workerNationality.value === 'WNI' ? (/^\d{16}$/.test(v) ? true : 'NIK harus 16 digit') : true),
  kodePos: v => (!v || /^\d{5}$/.test(v) ? true : 'Kode pos 5 digit'),
  nonNegative: v => (v == null || Number(v) >= 0 ? true : 'Tidak boleh negatif'),
  requiredIfPKWT: v => (form.statusPegawai === 'PKWT' ? (!!v || 'Wajib diisi untuk PKWT') : true),
}

onMounted(async () => {
  if (!workerId) return
  try {
    const w = await apiService.getWorker(workerId)
    workerData.value = w // Store worker data
    form.nik = w?.nik || ''
    form.namaLengkap = w?.nama || ''
    form.kpj = w?.kpj || ''
    form.nomorPegawai = w?.noPegawai || ''
    form.upah = w?.upah ?? null
    // Prefill from AddWorker step
    workerNationality.value = w?.nationality || 'WNI'
    if (w?.dateOfBirth) {
      form.tanggalLahir = String(w.dateOfBirth).substring(0, 10)
    }
  } catch (e) {
    console.error('Gagal memuat data worker', e)
  }
})

const onSubmit = async () => {
  if (!formRef.value) return
  const { valid } = await formRef.value.validate()
  if (!valid) return
  notifOpen.value = true
}

const confirmDaftar = async () => {
  try {
    if (workerId) {
      await apiService.updateWorker(workerId, {
        nik: form.nik,
        nama: form.namaLengkap,
        noPegawai: form.nomorPegawai,
        kpj: form.kpj || '', // Use KPJ from form
        upah: Number(form.upah || 0),
        // registration extras
        tempatLahir: form.tempatLahir,
        dateOfBirth: form.tanggalLahir,
        ibuKandung: form.ibuKandung,
        jenisKelamin: form.jenisKelamin,
        golDarah: form.golDarah,
        statusKawin: form.statusKawin,
        statusPegawai: form.statusPegawai,
        tanggalAwalBekerja: form.tanggalAwalBekerja,
        tanggalAkhirKontrak: form.tanggalAkhirKontrak,
        lokasiPekerjaan: form.lokasiPekerjaan,
        alamat: form.alamat,
        kabupaten: form.kabupaten,
        kodePos: form.kodePos,
        teleponRumah: form.teleponRumah,
        handphone: form.handphone,
        npwp: form.npwp,
        email: form.email,
      })
    }
    notifOpen.value = false
    successOpen.value = true
  } catch (e) {
    notifOpen.value = false
    alert(e?.message || 'Gagal menyimpan data')
  }
}

const onReset = () => {
  Object.keys(form).forEach(k => (form[k] = Array.isArray(form[k]) ? [] : ''))
  form.golDarah = null
  form.jenisKelamin = null
  form.statusKawin = null
  form.statusPegawai = null
  form.lokasiPekerjaan = null
  form.kabupaten = null
  form.upah = null
}

const goBackToEdit = () => {
  successOpen.value = false
  router.push('/dashboard')
}
</script>

<style scoped>
/* Menyerupai layout dua kolom seperti screenshot */
</style>
