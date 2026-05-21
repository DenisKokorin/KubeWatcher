import React, { createContext, useContext, useEffect, useState } from 'react'
import api from '../services/api'

interface User {
  uuid: string
  name: string
  email: string
  position: string
  team: string
  role: string
}

type AuthContextValue = {
  isAuthenticated: boolean
  user: User | null
  loading: boolean
  login: (email: string, password: string) => Promise<any>
  register: (userData: {
    name: string
    email: string
    position: string
    team: string
    password: string
  }) => Promise<any>
  logout: () => Promise<void>
  isAdmin: () => boolean
  isUser: () => boolean
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    try {
      const response = await api.get('/auth/me')
      setUser(response.data)
      setIsAuthenticated(true)
    } catch (error) {
      setIsAuthenticated(false)
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  const login = async (email: string, password: string) => {
    const response = await api.post('/auth/login', {
      email,
      password,
    })
    const accessToken = response.data?.AccessToken || response.data?.accessToken
    if (accessToken) {
      localStorage.setItem('accessToken', accessToken)
    }
    setIsAuthenticated(true)
    try {
      const userResponse = await api.get('/auth/me')
      setUser(userResponse.data)
    } catch (userError) {
      console.warn('Failed to fetch user data after login:', userError)
    }
    return response.data
  }

  const register = async (userData: {
    name: string
    email: string
    position: string
    team: string
    password: string
  }) => {
    const response = await api.post('/auth/register', userData)
    return response.data
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout', {})
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      localStorage.removeItem('accessToken')
      setIsAuthenticated(false)
      setUser(null)
    }
  }

  const isAdmin = () => user?.role === 'admin'
  const isUser = () => user?.role === 'user'

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        user,
        loading,
        login,
        register,
        logout,
        isAdmin,
        isUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = (): AuthContextValue => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
