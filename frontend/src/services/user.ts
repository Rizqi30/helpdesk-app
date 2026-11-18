import api from './api'

export interface User {
  id: number
  name: string
  email: string
  role: string
  created_at: string
}

export const userService = {
  async getAll() {
    const response = await api.get('/users')
    return response.data.data as User[]
  },

  async updateRole(userId: number, role: string) {
    const response = await api.patch(`/users/${userId}/role`, { role })
    return response.data.data as User
  },
}