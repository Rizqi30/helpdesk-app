import api from './api'

export interface Ticket {
  id: number
  title: string
  description: string
  category: string
  priority: string
  status: string
  requester_id: number
  assignee_id: number | null
  created_at: string
  updated_at: string
  requester: {
    id: number
    name: string
    email: string
    role: string
  }
  assignee?: {
    id: number
    name: string
    email: string
    role: string
  }
  comments?: Comment[]
}

export interface Comment {
  id: number
  ticket_id: number
  author_id: number
  body: string
  is_internal: boolean
  created_at: string
  author: {
    id: number
    name: string
    email: string
    role: string
  }
}

export interface TicketListResult {
  tickets: Ticket[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export interface TicketFilter {
  status?: string
  priority?: string
  search?: string
  sort_by?: string
  sort_order?: string
  page?: number
  limit?: number
}

export interface CreateTicketPayload {
  title: string
  description: string
  category?: string
  priority?: string
}

export interface UpdateTicketPayload {
  title?: string
  description?: string
  category?: string
  priority?: string
  status?: string
  assignee_id?: number
}

export const ticketService = {
  async getAll(filter: TicketFilter = {}) {
    const params = new URLSearchParams()

    if (filter.status) params.append('status', filter.status)
    if (filter.priority) params.append('priority', filter.priority)
    if (filter.search) params.append('search', filter.search)
    if (filter.sort_by) params.append('sort_by', filter.sort_by)
    if (filter.sort_order) params.append('sort_order', filter.sort_order)
    if (filter.page) params.append('page', filter.page.toString())
    if (filter.limit) params.append('limit', filter.limit.toString())

    const response = await api.get('/tickets', { params })
    return response.data.data as TicketListResult
  },

  async getById(id: number) {
    const response = await api.get(`/tickets/${id}`)
    return response.data.data as Ticket
  },

  async create(payload: CreateTicketPayload) {
    const response = await api.post('/tickets', payload)
    return response.data.data as Ticket
  },

  async update(id: number, payload: UpdateTicketPayload) {
    const response = await api.patch(`/tickets/${id}`, payload)
    return response.data.data as Ticket
  },
}