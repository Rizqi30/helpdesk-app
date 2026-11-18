import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import api from '../services/api'
import router from '../router'

interface User {
  id: number
  name: string
  email: string
  role: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || '')
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isAgent = computed(() => user.value?.role === 'agent')
  const isRequester = computed(() => user.value?.role === 'requester')

  function initialize() {
    const savedToken = localStorage.getItem('token')
    const savedUser = localStorage.getItem('user')
    if (savedToken && savedUser) {
      token.value = savedToken
      user.value = JSON.parse(savedUser)
    }
  }

  async function login(email: string, password: string) {
    loading.value = true
    try {
      const response = await api.post('/auth/login', { email, password })
      const data = response.data.data

      token.value = data.token
      user.value = data.user

      // save to localstorage
      localStorage.setItem('token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))

      return { success: true }
    } catch (error: any) {
      const message = error.response?.data?.message || 'Login failed'
      return { success: false, message }
    } finally {
      loading.value = false
    }
  }

  async function register(name: string, email: string, password: string) {
    loading.value = true
    try {
      const response = await api.post('/auth/register', { name, email, password })
      const data = response.data.data

      token.value = data.token
      user.value = data.user

      localStorage.setItem('token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))

      return { success: true }
    } catch (error: any) {
      const message = error.response?.data?.message || 'Registration failed'
      return { success: false, message }
    } finally {
      loading.value = false
    }
  }

  async function fetchProfile() {
    try {
      const response = await api.get('/me')
      user.value = response.data.data
      localStorage.setItem('user', JSON.stringify(user.value))
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push({ name: 'login' })
  }

  return {
    user,
    token,
    loading,

    isAuthenticated,
    userRole,
    isAdmin,
    isAgent,
    isRequester,

    initialize,
    login,
    register,
    fetchProfile,
    logout,
  }
})