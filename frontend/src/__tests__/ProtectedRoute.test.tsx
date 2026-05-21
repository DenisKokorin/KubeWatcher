import React from 'react'
import { render, screen } from '@testing-library/react'
import { describe, it, vi, expect } from 'vitest'
import { MemoryRouter } from 'react-router-dom'

vi.mock('../hooks/useAuth', () => ({ useAuth: vi.fn() }))
import * as auth from '../hooks/useAuth'
import { ProtectedRoute } from '../components/ProtectedRoute'

describe('ProtectedRoute', () => {
  const useAuthMock = auth.useAuth as unknown as vi.Mock

  it('shows loading when auth is loading', () => {
    useAuthMock.mockReturnValue({ isAuthenticated: false, loading: true, user: null })
    render(
      <MemoryRouter>
        <ProtectedRoute><div>Secret</div></ProtectedRoute>
      </MemoryRouter>
    )
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('redirects to login when not authenticated', () => {
    useAuthMock.mockReturnValue({ isAuthenticated: false, loading: false, user: null })
    const { queryByText } = render(
      <MemoryRouter>
        <ProtectedRoute><div>Secret</div></ProtectedRoute>
      </MemoryRouter>
    )
    expect(queryByText('Secret')).not.toBeInTheDocument()
  })

  it('blocks admin required route for non-admin', () => {
    useAuthMock.mockReturnValue({ isAuthenticated: true, loading: false, user: { role: 'user' } })
    const { queryByText } = render(
      <MemoryRouter>
        <ProtectedRoute requiredRole="admin"><div>AdminOnly</div></ProtectedRoute>
      </MemoryRouter>
    )
    expect(queryByText('AdminOnly')).not.toBeInTheDocument()
  })

  it('renders children for authorized user', () => {
    useAuthMock.mockReturnValue({ isAuthenticated: true, loading: false, user: { role: 'admin' } })
    render(
      <MemoryRouter>
        <ProtectedRoute requiredRole="admin"><div>AdminOnly</div></ProtectedRoute>
      </MemoryRouter>
    )
    expect(screen.getByText('AdminOnly')).toBeInTheDocument()
  })
})
