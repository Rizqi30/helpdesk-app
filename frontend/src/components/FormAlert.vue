<script setup lang="ts">
import { CircleAlert, CircleCheck } from 'lucide-vue-next'

defineProps<{
  type: 'error' | 'success'
  message: string
}>()
</script>

<template>
  <div
    class="mb-4 p-3 text-sm rounded-md flex items-center gap-2"
    :class="{
      'bg-red-50 border border-red-200 text-red-700': type === 'error',
      'bg-green-50 border border-green-200 text-green-700': type === 'success',
    }"
  >
    <CircleAlert v-if="type === 'error'" :size="16" class="shrink-0" />
    <CircleCheck v-if="type === 'success'" :size="16" class="shrink-0" />
    <span>{{ message }}</span>
  </div>
</template>