import api from './api'
import type { Comment } from './ticket'

export interface CreateCommentPayload {
  body: string
  is_internal?: boolean
}

export const commentService = {
  async getByTicketId(ticketId: number) {
    const response = await api.get(`/tickets/${ticketId}/comments`)
    return response.data.data as Comment[]
  },

  async create(ticketId: number, payload: CreateCommentPayload) {
    const response = await api.post(`/tickets/${ticketId}/comments`, payload)
    return response.data.data as Comment
  },
}