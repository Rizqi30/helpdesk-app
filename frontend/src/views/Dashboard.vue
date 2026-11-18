<script setup lang="ts">
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
</script>

<template>
  <DefaultLayout>
    <div>
      <h1 class="text-2xl font-bold text-gray-900 mb-4">Dashboard</h1>
      <p class="text-gray-600">Welcome, {{ auth.user?.name }}! You are logged in as {{ auth.userRole }}.</p>
      <p class="text-gray-400 text-sm mt-2">Ticket list will be implemented in the next step.</p>
    </div>
  </DefaultLayout>
</template>