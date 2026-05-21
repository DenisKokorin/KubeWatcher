import api from './api'

export interface ExternalPost {
  id: number
  user_id: number
  title: string
  summary: string
  content: string
  source: string
  published_at: string
}

export const externalService = {
  getPosts: async (params?: { userId?: number; limit?: number; search?: string }): Promise<ExternalPost[]> => {
    const response = await api.get('/external/posts', { params })
    return response.data.data || []
  },
}
