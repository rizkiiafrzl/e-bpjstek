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
            <h1 class="text-h5 font-weight-black mb-1">Create Account</h1>
            <p class="text-body-2 text-medium-emphasis mb-6">
              Please enter your details to create an account.
            </p>

            <v-form @submit.prevent="handleRegister" ref="formRef">
              <v-card-text class="pa-0">
                <v-text-field
                  v-model="form.fullName"
                  label="Full Name"
                  placeholder="Enter your full name"
                  :rules="fullNameRules"
                  :error-messages="formErrors.fullName ? [formErrors.fullName] : []"
                  variant="outlined"
                  class="mb-4"
                  prepend-inner-icon="mdi-account"
                  required
                />

                <v-text-field
                  v-model="form.email"
                  label="Email"
                  type="email"
                  placeholder="example@mail.com"
                  :rules="emailRules"
                  :error-messages="formErrors.email ? [formErrors.email] : []"
                  variant="outlined"
                  class="mb-4"
                  prepend-inner-icon="mdi-email"
                  required
                />

                <v-text-field
                  v-model="form.password"
                  label="Password"
                  :type="showPassword ? 'text' : 'password'"
                  placeholder="Enter your password"
                  :rules="passwordRules"
                  :error-messages="formErrors.password ? [formErrors.password] : []"
                  variant="outlined"
                  class="mb-4"
                  prepend-inner-icon="mdi-lock"
                  :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                  @click:append-inner="togglePassword"
                  required
                />

                <v-text-field
                  v-model="form.confirmPassword"
                  label="Confirm Password"
                  :type="showConfirmPassword ? 'text' : 'password'"
                  placeholder="Confirm your password"
                  :rules="confirmPasswordRules"
                  :error-messages="formErrors.confirmPassword ? [formErrors.confirmPassword] : []"
                  variant="outlined"
                  class="mb-4"
                  prepend-inner-icon="mdi-lock-check"
                  :append-inner-icon="showConfirmPassword ? 'mdi-eye-off' : 'mdi-eye'"
                  @click:append-inner="toggleConfirmPassword"
                  required
                />

                <!-- Traditional CAPTCHA -->
                <div class="mb-4 d-flex justify-center">
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
                  class="mb-4"
                  closable
                  @click:close="state.error = null"
                >
                  {{ state.error }}
                </v-alert>
              </v-card-text>

              <v-card-actions class="pa-0 flex-column">
                <v-btn
                  type="submit"
                  :loading="isSubmitting"
                  :disabled="isSubmitting"
                  color="primary"
                  size="large"
                  block
                  class="mb-4"
                >
                  <v-icon left>mdi-account-plus</v-icon>
                  Create Account
                </v-btn>

                <v-btn
                  type="button"
                  @click="handleGoogleRegister"
                  size="large"
                  block
                  variant="outlined"
                  class="google-btn"
                  :prepend-icon="'mdi-google'"
                >
                  Sign up with Google
                </v-btn>

                <div class="text-center mt-4">
                  <span class="text-body-2 text-medium-emphasis">Already have an account?</span>
                  <RouterLink to="/login" class="ml-2">Sign In</RouterLink>
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
const showConfirmPassword = ref(false)
const isSubmitting = ref(false)
const state = reactive({ error: null, success: false })
const formErrors = reactive({
  fullName: null,
  email: null,
  password: null,
  confirmPassword: null,
})
const formRef = ref(null)
const captchaToken = ref(null)
const captchaRef = ref(null)

const form = reactive({
  fullName: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const fullNameRules = [
  (v) => !!v || 'Full name is required',
  (v) => (v && v.length >= 2) || 'Full name must be at least 2 characters',
]

const emailRules = [
  (v) => !!v || 'Email is required',
  (v) => /.+@.+\..+/.test(v) || 'Invalid email format',
]

const passwordRules = [
  (v) => !!v || 'Password is required',
  (v) => (v && v.length >= 8) || 'Password must be at least 8 characters',
]

const confirmPasswordRules = [
  (v) => !!v || 'Please confirm your password',
  (v) => v === form.password || 'Passwords do not match',
]

const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const toggleConfirmPassword = () => {
  showConfirmPassword.value = !showConfirmPassword.value
}

const handleRegister = async () => {
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
  formErrors.fullName = null
  formErrors.email = null
  formErrors.password = null
  formErrors.confirmPassword = null

  try {
    await apiService.register({
      full_name: form.fullName,
      email: form.email,
      password: form.password,
      captchaToken: captchaToken.value,
    })

    state.success = true
    router.push('/login')

    form.fullName = ''
    form.email = ''
    form.password = ''
    form.confirmPassword = ''
    if (captchaRef.value) {
      captchaRef.value.reset()
    }
  } catch (error) {
    state.error = error.message || 'An error occurred during registration. Please try again.'
    if (captchaRef.value) {
      captchaRef.value.reset()
    }
  } finally {
    isSubmitting.value = false
  }
}

const handleGoogleRegister = () => {
  console.log('Google register clicked')
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
</style>
