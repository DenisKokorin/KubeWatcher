import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:7979/api/v1',
  withCredentials: true,
})

api.interceptors.request.use((config) => {
  const accessToken = localStorage.getItem('accessToken')
  if (accessToken && config.headers) {
    config.headers.Authorization = `Bearer ${accessToken}`
  }
  return config
})

export default api