<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import FormInput from '@/components/FormInput.vue'
import FormAlert from '@/components/FormAlert.vue'

const router = useRouter()
const auth = useAuthStore()

// Form data
const form = reactive({
  email: '',
  password: '',
})

const errors = reactive({
  email: '',
  password: '',
})

const serverError = ref('')

function validate(): boolean {
  let valid = true
  errors.email = ''
  errors.password = ''

  if (!form.email) {
    errors.email = 'Email is required'
    valid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = 'Please enter a valid email address'
    valid = false
  }

  if (!form.password) {
    errors.password = 'Password is required'
    valid = false
  }

  return valid
}

async function handleLogin() {
  serverError.value = ''

  if (!validate()) return

  const result = await auth.login(form.email, form.password)
  if (result.success) {
    router.push({ name: 'dashboard' })
  } else {
    serverError.value = result.message
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div class="w-full max-w-md p-8 bg-white rounded-lg shadow-md">
      <h1 class="text-2xl font-bold text-center text-gray-900 mb-6">
        Helpdesk Login
      </h1>

      <FormAlert v-if="serverError" type="error" :message="serverError" />

      <form @submit.prevent="handleLogin" class="space-y-4">
        <FormInput
          v-model="form.email"
          label="Email"
          type="email"
          placeholder="email@example.com"
          :error="errors.email"
          :required="true"
        />

        <FormInput
          v-model="form.password"
          label="Password"
          type="password"
          placeholder="••••••••"
          :error="errors.password"
          :required="true"
        />

        <button
          type="submit"
          :disabled="auth.loading"
          class="w-full py-2 px-4 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
        >
          {{ auth.loading ? 'Logging in...' : 'Login' }}
        </button>
      </form>

      <p class="mt-4 text-center text-sm text-gray-600">
        Don't have an account?
        <RouterLink to="/register" class="text-blue-600 hover:underline">Register</RouterLink>
      </p>
    </div>
  </div>
</template>