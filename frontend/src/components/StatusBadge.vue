<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
}>()

const statusConfig = computed(() => {
  const configs: Record<string, { label: string; class: string }> = {
    open: {
      label: 'Open',
      class: 'bg-blue-100 text-blue-700',
    },
    in_progress: {
      label: 'In Progress',
      class: 'bg-yellow-100 text-yellow-700',
    },
    resolved: {
      label: 'Resolved',
      class: 'bg-green-100 text-green-700',
    },
    closed: {
      label: 'Closed',
      class: 'bg-gray-100 text-gray-700',
    },
  }

  return configs[props.status] || { label: props.status, class: 'bg-gray-100 text-gray-700' }
})
</script>

<template>
  <span
    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
    :class="statusConfig.class"
  >
    {{ statusConfig.label }}
  </span>
</template>