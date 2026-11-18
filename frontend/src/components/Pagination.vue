<script setup lang="ts">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

defineProps<{
  page: number
  totalPages: number
  total: number
}>()

defineEmits<{
  'update:page': [value: number]
}>()
</script>

<template>
  <div class="flex items-center justify-between pt-4 border-t border-gray-200">
    <p class="text-sm text-gray-600">
      Total: <span class="font-medium">{{ total }}</span> tiket
    </p>

    <div class="flex items-center gap-2">
      <button
        :disabled="page <= 1"
        @click="$emit('update:page', page - 1)"
        class="p-2 rounded-md border border-gray-300 text-gray-600 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <ChevronLeft :size="16" />
      </button>

      <span class="text-sm text-gray-600">
        Halaman {{ page }} dari {{ totalPages }}
      </span>

      <button
        :disabled="page >= totalPages"
        @click="$emit('update:page', page + 1)"
        class="p-2 rounded-md border border-gray-300 text-gray-600 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <ChevronRight :size="16" />
      </button>
    </div>
  </div>
</template>