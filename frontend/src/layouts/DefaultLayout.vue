<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { LogOut, Ticket, Users, User } from 'lucide-vue-next'

const auth = useAuthStore()
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navbar -->
    <nav class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">
          <!-- Logo / Title -->
          <div class="flex items-center gap-6">
            <RouterLink to="/" class="text-xl font-bold text-gray-900">
              Helpdesk
            </RouterLink>

            <!-- Nav Links -->
            <div class="flex items-center gap-4">
              <RouterLink
                to="/"
                class="flex items-center gap-1.5 text-sm text-gray-600 hover:text-gray-900"
              >
                <Ticket :size="16" />
                Tickets
              </RouterLink>

              <RouterLink
                v-if="auth.isAdmin"
                to="/admin/users"
                class="flex items-center gap-1.5 text-sm text-gray-600 hover:text-gray-900"
              >
                <Users :size="16" />
                Users
              </RouterLink>
            </div>
          </div>

          <!-- User Info & Logout -->
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2 text-sm text-gray-600">
              <User :size="16" />
              <span>{{ auth.user?.name }}</span>
              <span class="px-2 py-0.5 text-xs font-medium rounded-full"
                :class="{
                  'bg-red-100 text-red-700': auth.isAdmin,
                  'bg-blue-100 text-blue-700': auth.isAgent,
                  'bg-gray-100 text-gray-700': auth.isRequester,
                }"
              >
                {{ auth.userRole }}
              </span>
            </div>

            <button
              @click="auth.logout()"
              class="flex items-center gap-1.5 text-sm text-gray-500 hover:text-red-600 cursor-pointer"
            >
              <LogOut :size="16" />
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Page Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <slot />
    </main>
  </div>
</template>