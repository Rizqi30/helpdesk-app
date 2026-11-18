<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import FormInput from '@/components/FormInput.vue'
import FormAlert from '@/components/FormAlert.vue'

const router = useRouter()
const auth = useAuthStore()

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const errors = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const serverError = ref('')

function validate(): boolean {
  let valid = true
  errors.name = ''
  errors.email = ''
  errors.password = ''
  errors.confirmPassword = ''

  if (!form.name) {
    errors.name = 'Name is required'
    valid = false
  } else if (form.name.length < 2) {
    errors.name = 'Name must be at least 2 characters'
    valid = false
  }

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
  } else if (form.password.length < 6) {
    errors.password = 'Password must be at least 6 characters'
    valid = false
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = 'Please confirm your password'
    valid = false
  } else if (form.password !== form.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match'
    valid = false
  }

  return valid
}

async function handleRegister() {
  serverError.value = ''

  if (!validate()) return

  const result = await auth.register(form.name, form.email, form.password)
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
        Create Account
      </h1>

      <FormAlert v-if="serverError" type="error" :message="serverError" />

      <form @submit.prevent="handleRegister" class="space-y-4">
        <FormInput
          v-model="form.name"
          label="Full Name"
          placeholder="John Doe"
          :error="errors.name"
          :required="true"
        />

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
          placeholder="Minimum 6 characters"
          :error="errors.password"
          :required="true"
        />

        <FormInput
          v-model="form.confirmPassword"
          label="Confirm Password"
          type="password"
          placeholder="Re-enter your password"
          :error="errors.confirmPassword"
          :required="true"
        />

        <button
          type="submit"
          :disabled="auth.loading"
          class="w-full py-2 px-4 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
        >
          {{ auth.loading ? 'Creating account...' : 'Register' }}
        </button>
      </form>

      <p class="mt-4 text-center text-sm text-gray-600">
        Already have an account?
        <RouterLink to="/login" class="text-blue-600 hover:underline">Login</RouterLink>
      </p>
    </div>
  </div>
</template>