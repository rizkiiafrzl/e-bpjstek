<template>
  <AddWorkerView :id="idValue" />
  <!-- Reuse the same form as AddWorker, but in edit mode via id prop -->
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AddWorkerView from './AddWorkerView.vue'

const route = useRoute()
// Ensure we pass a plain string value, not a Ref, to the child component
const idValue = computed(() => String(route.params.id || ''))
</script>

<style scoped>
</style>


