<template>
  <div class="traditional-captcha">
    <div class="captcha-container">
      <!-- CAPTCHA Image -->
      <div class="captcha-image-container">
        <canvas ref="captchaCanvas" class="captcha-canvas"></canvas>
        <button type="button" @click="generateCaptcha" class="refresh-btn" title="Refresh CAPTCHA">
          <v-icon>mdi-refresh</v-icon>
        </button>
      </div>

      <!-- Input Field -->
      <v-text-field
        v-model="userInput"
        label="Enter the numbers you see"
        placeholder="Type the numbers"
        variant="outlined"
        class="captcha-input"
        :error-messages="error ? [error] : []"
        @input="onInput"
        @keyup.enter="$emit('submit')"
      />
    </div>

    <!-- Error Message -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'

const emit = defineEmits(['verified', 'error', 'submit'])

const captchaCanvas = ref(null)
const userInput = ref('')
const captchaAnswer = ref('')
const error = ref('')

// Generate random 4-digit number
const generateRandomNumber = () => {
  return Math.floor(1000 + Math.random() * 9000).toString()
}

// Draw distorted CAPTCHA on canvas
const drawCaptcha = (ctx, text) => {
  const canvas = captchaCanvas.value
  if (!canvas) return

  // Clear canvas
  ctx.clearRect(0, 0, canvas.width, canvas.height)

  // Set canvas background
  ctx.fillStyle = '#f0f0f0'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // Add noise lines
  for (let i = 0; i < 20; i++) {
    ctx.strokeStyle = `rgba(${Math.random() * 255}, ${Math.random() * 255}, ${Math.random() * 255}, 0.3)`
    ctx.lineWidth = Math.random() * 2
    ctx.beginPath()
    ctx.moveTo(Math.random() * canvas.width, Math.random() * canvas.height)
    ctx.lineTo(Math.random() * canvas.width, Math.random() * canvas.height)
    ctx.stroke()
  }

  // Draw text with distortion
  ctx.font = 'bold 32px Arial'
  ctx.fillStyle = '#333'

  for (let i = 0; i < text.length; i++) {
    const x = 15 + i * 30
    const y = 40 + Math.sin(i) * 5 // Wave distortion
    const rotation = (Math.random() - 0.5) * 0.3 // Random rotation

    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(rotation)
    ctx.fillText(text[i], 0, 0)
    ctx.restore()
  }

  // Add more noise dots
  for (let i = 0; i < 50; i++) {
    ctx.fillStyle = `rgba(${Math.random() * 255}, ${Math.random() * 255}, ${Math.random() * 255}, 0.5)`
    ctx.beginPath()
    ctx.arc(
      Math.random() * canvas.width,
      Math.random() * canvas.height,
      Math.random() * 2,
      0,
      2 * Math.PI,
    )
    ctx.fill()
  }
}

// Generate new CAPTCHA
const generateCaptcha = () => {
  captchaAnswer.value = generateRandomNumber()
  userInput.value = ''
  error.value = ''

  console.log('Generated CAPTCHA answer:', captchaAnswer.value)

  nextTick(() => {
    const canvas = captchaCanvas.value
    if (canvas) {
      // Ensure canvas size is set
      canvas.width = 180
      canvas.height = 48
      const ctx = canvas.getContext('2d')
      drawCaptcha(ctx, captchaAnswer.value)
      console.log('CAPTCHA drawn on canvas')
    }
  })

  // Removed emit('error', null) to avoid triggering generic error upstream
}

// Handle input
const onInput = () => {
  console.log('User input:', userInput.value)
  console.log('Expected answer:', captchaAnswer.value)

  if (userInput.value === captchaAnswer.value) {
    error.value = ''
    emit('verified', userInput.value)
    console.log('CAPTCHA verified successfully!')
  } else if (
    userInput.value.length === captchaAnswer.value.length &&
    userInput.value !== captchaAnswer.value
  ) {
    error.value = 'Incorrect CAPTCHA. Please try again.'
    emit('error', 'Incorrect CAPTCHA')
    console.log('CAPTCHA failed - wrong answer')
  } else {
    // Clear error when user is still typing
    error.value = ''
  }
}

// Initialize CAPTCHA
onMounted(() => {
  // Set canvas size
  const canvas = captchaCanvas.value
  if (canvas) {
    canvas.width = 180
    canvas.height = 48
  }
  generateCaptcha()
})

// Expose methods
defineExpose({
  generateCaptcha,
  reset: () => {
    userInput.value = ''
    error.value = ''
    generateCaptcha()
  },
  isValid: () => {
    const isValid = userInput.value === captchaAnswer.value
    console.log('CAPTCHA validation check:', {
      userInput: userInput.value,
      expected: captchaAnswer.value,
      isValid: isValid,
    })
    return isValid
  },
  getValue: () => userInput.value,
  getAnswer: () => captchaAnswer.value, // For debugging
})
</script>

<style scoped>
.traditional-captcha {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.captcha-container {
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  gap: 12px;
  justify-content: center;
}

.captcha-image-container {
  position: relative;
  display: flex;
  align-items: center;
  height: 48px;
}

.captcha-canvas {
  width: 180px;
  height: 48px;
  border: 2px solid #ddd;
  border-radius: 4px;
  background: #f9f9f9;
  cursor: pointer;
  display: block; /* avoid baseline offset */
}

.refresh-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  background: rgba(255, 255, 255, 0.8);
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.2s;
}

.refresh-btn:hover {
  background: rgba(255, 255, 255, 1);
}

.captcha-input {
  width: 180px;
  max-width: 180px;
  min-width: 180px;
  align-self: center;
}

/* Force v-text-field control to 48px height to match canvas */
.captcha-input :deep(.v-input__control) {
  margin-top: 0;
  min-height: 48px;
}

.captcha-input :deep(.v-field) {
  height: 48px;
  min-height: 48px;
}

.captcha-input :deep(.v-field__input) {
  min-height: 48px;
  height: 48px;
  padding-top: 0;
  padding-bottom: 0;
}

.captcha-input :deep(.v-field__outline) {
  top: 0;
  bottom: 0;
}

.captcha-input :deep(.v-label) {
  top: 12px;
  font-size: 12px;
}

.captcha-input :deep(.v-label--active) {
  top: -8px;
  font-size: 12px;
}

.captcha-input :deep(input) {
  text-align: center;
  letter-spacing: 2px;
  font-weight: 600;
  padding-top: 12px;
  padding-bottom: 12px;
}

.error-message {
  color: #d32f2f;
  font-size: 12px;
  text-align: center;
}
</style>
