import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, vi, expect } from 'vitest'

vi.mock('../services/api', () => ({ default: { get: vi.fn(), post: vi.fn() } }))
import api from '../services/api'
import { AuthProvider, useAuth } from '../hooks/useAuth'

function Consumer() {
  const { isAuthenticated, loading } = useAuth()
  return <div>{loading ? 'loading' : isAuthenticated ? 'auth' : 'noauth'}</div>
}

describe('AuthProvider', () => {
  it('handles failed /auth/me by marking unauthenticated', async () => {
    ;(api.get as unknown as vi.Mock).mockRejectedValueOnce({ response: { status: 401 } })
    render(
      <AuthProvider>
        <Consumer />
      </AuthProvider>
    )

    await waitFor(() => expect(screen.getByText('noauth')).toBeInTheDocument())
  })

  it('handles successful /auth/me and marks authenticated', async () => {
    ;(api.get as unknown as vi.Mock).mockResolvedValueOnce({ data: { email: 'a@b.com' } })
    render(
      <AuthProvider>
        <Consumer />
      </AuthProvider>
    )

    await waitFor(() => expect(screen.getByText('auth')).toBeInTheDocument())
  })
})
