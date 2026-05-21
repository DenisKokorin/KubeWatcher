import React from 'react'
import { render, screen } from '@testing-library/react'
import { describe, it, vi, expect } from 'vitest'

vi.mock('../hooks/useAuth', () => ({ useAuth: vi.fn() }))
import * as auth from '../hooks/useAuth'
import Sidebar from '../components/Sidebar'
import { MemoryRouter } from 'react-router-dom'

describe('Sidebar role behavior', () => {
  const useAuthMock = auth.useAuth as unknown as vi.Mock

  it('shows admin link for admin users', () => {
    useAuthMock.mockReturnValue({ isAdmin: () => true, user: { name: 'Admin', role: 'admin' }, logout: vi.fn() })
    render(
      <MemoryRouter>
        <Sidebar />
      </MemoryRouter>
    )
    expect(screen.getByText(/User Management/i)).toBeInTheDocument()
  })

  it('hides admin link for normal users', () => {
    useAuthMock.mockReturnValue({ isAdmin: () => false, user: { name: 'User', role: 'user' }, logout: vi.fn() })
    render(
      <MemoryRouter>
        <Sidebar />
      </MemoryRouter>
    )
    expect(screen.queryByText(/User Management/i)).not.toBeInTheDocument()
  })
})
