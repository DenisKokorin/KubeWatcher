import api from './api'

export interface Employee {
  uuid: string
  name: string
  email: string
  position: string
  team: string
  role: string
  hire_date: string
  profile_document_url?: string
}

export interface EmployeeQuery {
  search?: string
  team?: string
  position?: string
  role?: string
  sortBy?: string
  sortOrder?: string
  page?: number
  size?: number
}

export interface EmployeeList {
  employees: Employee[]
  total: number
  page: number
  size: number
}

export interface Review {
  id: string
  employee_id: string
  reviewer_id: string
  period: string
  rating: number
  comments: string
  goals: string[]
  strengths: string[]
  areas_for_improvement: string[]
  created_at: string
  updated_at: string
}

export interface EmployeeWithReviews {
  employee: Employee
  reviews: Review[]
}

export interface TeamStats {
  team: string
  employee_count: number
  average_rating: number
  reviews_count: number
}

export interface ReviewRequest {
  employee_id: string
  reviewer_id: string
  period: string
  rating: number
  comments: string
  goals: string[]
  strengths: string[]
  areas_for_improvement: string[]
}

export const reviewService = {
  getAllEmployees: async (params?: EmployeeQuery): Promise<EmployeeList> => {
    const response = await api.get('/reviews/employees', { params })
    return response.data.data || response.data
  },

  getEmployee: async (id: string): Promise<Employee> => {
    const response = await api.get(`/reviews/employees/${id}`)
    return response.data.data || response.data
  },

  getEmployeeWithReviews: async (id: string): Promise<EmployeeWithReviews> => {
    const response = await api.get(`/reviews/employees/${id}/reviews`)
    return response.data.data || response.data
  },

  createReview: async (review: ReviewRequest): Promise<{ id: string }> => {
    console.log('Creating review:', review)
    try {
      const response = await api.post('/reviews', review)
      console.log('Review created successfully:', response.data)
      return response.data.data || response.data
    } catch (error: any) {
      console.error('Failed to create review:', error.response?.status, error.response?.data)
      throw error
    }
  },

  getEmployeesByTeam: async (team: string): Promise<Employee[]> => {
    const response = await api.get(`/reviews/employees/${team}`)
    return response.data.data || response.data
  },

  getTeamStats: async (team: string): Promise<TeamStats> => {
    const response = await api.get(`/reviews/team/${team}`)
    return response.data.data || response.data
  },
}