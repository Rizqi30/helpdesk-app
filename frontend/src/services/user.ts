import api from './api'

export interface AdminStats {
  total_tickets: number
  open_tickets: number
  in_progress_tickets: number
  resolved_tickets: number
  closed_tickets: number
  overdue_tickets: number
  total_users: number
}

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

  async getStats() {
    const response = await api.get('/users/stats')
    return response.data.data as AdminStats
  },
}