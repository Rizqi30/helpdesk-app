<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { LogOut, Ticket, Users, User, Menu, X } from 'lucide-vue-next'

const auth = useAuthStore()
const mobileMenuOpen = ref(false)
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">

          <div class="flex items-center gap-6">
            <RouterLink to="/" class="text-xl font-bold text-gray-900">
              Helpdesk
            </RouterLink>

            <div class="hidden sm:flex items-center gap-4">
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

          <div class="hidden sm:flex items-center gap-4">
            <div class="flex items-center gap-2 text-sm text-gray-600">
              <User :size="16" />
              <span>{{ auth.user?.name }}</span>
              <span
                class="px-2 py-0.5 text-xs font-medium rounded-full"
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

          <div class="flex sm:hidden">
            <button
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="p-2 text-gray-600 hover:text-gray-900 focus:outline-none"
            >
              <X v-if="mobileMenuOpen" :size="24" />
              <Menu v-else :size="24" />
            </button>
          </div>

        </div>
      </div>

      <div v-if="mobileMenuOpen" class="sm:hidden border-t border-gray-200 bg-white px-4 pt-2 pb-4 space-y-3">
        <div class="flex items-center justify-between py-2 border-b border-gray-100 text-sm text-gray-700">
          <div class="flex items-center gap-2">
            <User :size="18" />
            <span class="font-medium">{{ auth.user?.name }}</span>
          </div>
          <span
            class="px-2 py-0.5 text-xs font-medium rounded-full"
            :class="{
              'bg-red-100 text-red-700': auth.isAdmin,
              'bg-blue-100 text-blue-700': auth.isAgent,
              'bg-gray-100 text-gray-700': auth.isRequester,
            }"
          >
            {{ auth.userRole }}
          </span>
        </div>

        <RouterLink
          to="/"
          @click="mobileMenuOpen = false"
          class="flex items-center gap-2 py-2 text-sm text-gray-600 hover:text-gray-900"
        >
          <Ticket :size="18" />
          Tickets
        </RouterLink>

        <RouterLink
          v-if="auth.isAdmin"
          to="/admin/users"
          @click="mobileMenuOpen = false"
          class="flex items-center gap-2 py-2 text-sm text-gray-600 hover:text-gray-900"
        >
          <Users :size="18" />
          Users
        </RouterLink>

        <button
          @click="auth.logout()"
          class="w-full flex items-center gap-2 py-2 text-sm text-red-600 hover:bg-red-50 rounded text-left cursor-pointer"
        >
          <LogOut :size="18" />
          Logout
        </button>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <slot />
    </main>
  </div>
</template>