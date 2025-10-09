<template>
  <div class="login-hero">
    <div class="login-hero__bg" :style="{ backgroundImage: `url(${bgImage})` }"></div>
    <v-container class="login-hero__content d-flex align-center justify-center">
      <v-card class="login-card" elevation="8">
        <v-row no-gutters>
          <!-- Left side image -->
          <v-col cols="12" md="6" class="d-none d-md-block">
            <v-img :src="sideImage" height="100%" cover></v-img>
          </v-col>

          <!-- Right side - Form -->
          <v-col cols="12" md="6" class="pa-8">
            <div class="d-flex align-center mb-6">
              <v-img :src="logoImage" max-width="120" class="me-3" alt="BPJS Logo" />
            </div>
            <h1 class="text-h5 font-weight-black mb-1">Sign In</h1>
            <p class="text-body-2 text-medium-emphasis mb-6">
              Please enter your credentials to sign in.
            </p>

            <v-form @submit.prevent="handleLogin" ref="formRef">
              <v-card-text>
                <v-text-field
                  v-model="formData.email"
                  label="Email"
                  type="email"
                  placeholder="example@mail.com"
                  :rules="emailRules"
                  :error-messages="formErrors.email ? [formErrors.email] : []"
                  variant="outlined"
                  prepend-inner-icon="mdi-email"
                  required
                />

                <v-text-field
                  v-model="formData.password"
                  label="Password"
                  :type="showPassword ? 'text' : 'password'"
                  placeholder="Enter your password"
                  :rules="passwordRules"
                  :error-messages="formErrors.password ? [formErrors.password] : []"
                  variant="outlined"
                  prepend-inner-icon="mdi-lock"
                  :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                  @click:append-inner="togglePassword"
                  required
                />

                <!-- Traditional CAPTCHA -->
                <div class="captcha-wrapper">
                  <TraditionalCaptcha
                    ref="captchaRef"
                    @verified="onCaptchaVerified"
                    @error="onCaptchaError"
                  />
                </div>

                <v-alert
                  v-if="state.error"
                  type="error"
                  variant="tonal"
                  closable
                  @click:close="state.error = null"
                >
                  {{ state.error }}
                </v-alert>
              </v-card-text>

              <v-card-actions>
                <v-btn
                  type="submit"
                  :loading="isSubmitting"
                  :disabled="isSubmitting"
                  color="primary"
                  size="large"
                  block
                >
                  <v-icon left>mdi-login</v-icon>
                  Sign In
                </v-btn>

                <v-btn
                  type="button"
                  @click="handleGoogleLogin"
                  size="large"
                  block
                  variant="outlined"
                  class="google-btn"
                  :prepend-icon="'mdi-google'"
                >
                  Sign in with Google
                </v-btn>

                <div class="text-center mt-4">
                  <span class="text-body-2 text-medium-emphasis">Don't have an account?</span>
                  <RouterLink to="/register" class="ml-2">Sign Up</RouterLink>
                </div>
              </v-card-actions>
            </v-form>
          </v-col>
        </v-row>
      </v-card>
    </v-container>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import apiService from '../services/api.js'
import TraditionalCaptcha from '../components/TraditionalCaptcha.vue'
import bgImage from '@/components/icons/background.jpg'
import sideImage from '@/components/icons/samping login.jpg'
import logoImage from '@/components/icons/Log BPJS Ketenagakerjaan.png'

const router = useRouter()
const showPassword = ref(false)
const isSubmitting = ref(false)
const state = reactive({ error: null, success: false })
const formErrors = reactive({
  email: null,
  password: null,
})
const formRef = ref(null)
const captchaToken = ref(null)
const captchaRef = ref(null)

const formData = reactive({
  email: '',
  password: '',
})

const emailRules = [
  (v) => !!v || 'Email is required',
  (v) => /.+@.+\..+/.test(v) || 'Invalid email format',
]

const passwordRules = [
  (v) => !!v || 'Password is required',
  (v) => (v && v.length >= 6) || 'Password must be at least 6 characters',
]

const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const handleLogin = async () => {
  if (!formRef.value) return
  const { valid } = await formRef.value.validate()
  if (!valid) return

  // Validate CAPTCHA
  if (!captchaToken.value) {
    state.error = 'Please complete the CAPTCHA verification.'
    return
  }

  isSubmitting.value = true
  state.error = null
  formErrors.email = null
  formErrors.password = null

  try {
    await apiService.login({
      email: formData.email,
      password: formData.password,
      captchaToken: captchaToken.value,
    })

    state.success = true
    router.push('/dashboard')
  } catch (error) {
    state.error = error.message || 'Invalid email or password. Please try again.'
    if (captchaRef.value) {
      captchaRef.value.reset()
    }
  } finally {
    isSubmitting.value = false
  }
}

const handleGoogleLogin = () => {
  const base = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
  window.location.href = `${base}/auth/google/login`
}

// CAPTCHA event handlers
const onCaptchaVerified = (token) => {
  console.log('CAPTCHA verified:', token)
  captchaToken.value = token
  state.error = null // Clear any previous errors
}

const onCaptchaError = (errorMsg) => {
  console.log('CAPTCHA error:', errorMsg)
  captchaToken.value = null
  state.error = errorMsg || 'CAPTCHA error. Please try again.'
}
</script>

<style scoped>
.login-hero {
  position: relative;
  min-height: 100vh;
  width: 100%;
  overflow: hidden;
}

.login-hero__bg {
  position: absolute;
  inset: 0;
  background: center/cover no-repeat fixed;
  filter: blur(6px);
  transform: scale(1.04);
}

.login-hero__content {
  position: relative;
  min-height: 100vh;
}

.login-card {
  position: relative;
  width: 100%;
  max-width: 960px;
  border-radius: 12px;
}

@media (max-width: 960px) {
  .login-card {
    margin: 16px;
  }
}

.google-btn :deep(.v-btn__prepend) {
  margin-inline-end: 8px;
}

.google-btn {
  justify-content: center;
  background: #ffffff !important;
  color: #374151 !important;
  border: 1px solid #d1d5db !important;
  text-transform: none;
}

.google-btn:hover {
  background: #f9fafb !important;
  border-color: #9ca3af !important;
}

.google-btn :deep(.v-icon) {
  color: #4285f4 !important;
}

.google-btn :deep(.v-btn__content) {
  width: 100%;
  justify-content: center;
}

.captcha-wrapper {
  display: flex;
  justify-content: center;
  margin: 20px 0;
}

.captcha-container {
  width: 100%;
  max-width: 300px;
  margin: 0 auto;
}

/* Improve form spacing */
.v-text-field {
  margin-bottom: 16px !important;
}

.v-card-text {
  padding: 0 !important;
}

.v-card-actions {
  padding: 0 !important;
  flex-direction: column;
  gap: 12px;
}

/* Better button styling */
.v-btn {
  height: 48px !important;
  font-weight: 500;
  letter-spacing: 0.5px;
}

/* Improve responsive design */
@media (max-width: 768px) {
  .login-card {
    margin: 8px;
    border-radius: 8px;
  }

  .pa-8 {
    padding: 24px !important;
  }

  .captcha-wrapper {
    margin: 16px 0;
  }

  .captcha-container {
    max-width: 280px;
  }
}
</style>
