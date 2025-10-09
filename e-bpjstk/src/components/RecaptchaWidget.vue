<template>
  <div class="recaptcha-widget">
    <div id="recaptcha-container" ref="container"></div>
    <div v-if="error" class="text-error text-caption mt-2">{{ error }}</div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({
  siteKey: {
    type: String,
    default: '6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI', // Test key
  },
})

const emit = defineEmits(['verified', 'expired', 'error'])

const container = ref(null)
const error = ref('')
const widgetId = ref(null)
const token = ref(null)

const loadScript = () => {
  return new Promise((resolve, reject) => {
    if (window.grecaptcha) {
      resolve()
      return
    }

    if (document.querySelector('script[src*="recaptcha"]')) {
      // Script already exists, wait for it to load
      const checkInterval = setInterval(() => {
        if (window.grecaptcha) {
          clearInterval(checkInterval)
          resolve()
        }
      }, 100)

      setTimeout(() => {
        clearInterval(checkInterval)
        if (!window.grecaptcha) {
          reject(new Error('reCAPTCHA script failed to load'))
        }
      }, 5000)
      return
    }

    const script = document.createElement('script')
    script.src = 'https://www.google.com/recaptcha/api.js?onload=onRecaptchaLoad&render=explicit'
    script.async = true
    script.defer = true
    script.onload = () => {
      console.log('reCAPTCHA script loaded')
      resolve()
    }
    script.onerror = () => {
      console.error('Failed to load reCAPTCHA script')
      reject(new Error('Failed to load reCAPTCHA script'))
    }
    document.head.appendChild(script)
  })
}

const renderWidget = () => {
  if (!window.grecaptcha || !window.grecaptcha.render) {
    error.value = 'reCAPTCHA tidak tersedia'
    return
  }

  if (!container.value) {
    error.value = 'Container tidak ditemukan'
    return
  }

  try {
    widgetId.value = window.grecaptcha.render(container.value, {
      sitekey: props.siteKey,
      theme: 'light', // 'light' or 'dark'
      size: 'normal', // 'compact', 'normal', 'invisible'
      callback: (response) => {
        console.log('reCAPTCHA verified:', response)
        token.value = response
        error.value = ''
        emit('verified', response)
      },
      'expired-callback': () => {
        console.log('reCAPTCHA expired')
        token.value = null
        emit('expired')
      },
      'error-callback': () => {
        console.log('reCAPTCHA error')
        token.value = null
        error.value = 'reCAPTCHA error. Silakan coba lagi.'
        emit('error')
      },
    })
    console.log('reCAPTCHA widget rendered with ID:', widgetId.value)
  } catch (err) {
    console.error('Error rendering reCAPTCHA:', err)
    error.value = 'Gagal memuat reCAPTCHA'
  }
}

const reset = () => {
  if (window.grecaptcha && widgetId.value !== null) {
    window.grecaptcha.reset(widgetId.value)
    token.value = null
  }
}

const execute = () => {
  if (window.grecaptcha && widgetId.value !== null) {
    window.grecaptcha.execute(widgetId.value)
  }
}

onMounted(async () => {
  try {
    await loadScript()
    setTimeout(() => {
      renderWidget()
    }, 500)
  } catch (err) {
    console.error('Failed to load reCAPTCHA:', err)
    error.value = 'Gagal memuat reCAPTCHA. Silakan refresh halaman.'
  }
})

onUnmounted(() => {
  if (window.grecaptcha && widgetId.value !== null) {
    window.grecaptcha.reset(widgetId.value)
  }
})

// Expose methods for parent component
defineExpose({
  reset,
  execute,
  token,
})
</script>

<style scoped>
.recaptcha-widget {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 16px 0;
}

.text-error {
  color: #d32f2f;
}
</style>
