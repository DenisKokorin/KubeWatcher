import React from 'react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, vi, expect } from 'vitest'

const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

vi.mock('../hooks/useAuth', () => ({ useAuth: vi.fn() }))
import * as auth from '../hooks/useAuth'
import Login from '../pages/Login'

describe('Login page', () => {
  const useAuthMock = auth.useAuth as unknown as vi.Mock

  it('submits login and navigates on success', async () => {
    const loginMock = vi.fn(() => Promise.resolve({}))
    useAuthMock.mockReturnValue({ login: loginMock, register: vi.fn() })
    render(<Login />)

    await userEvent.type(screen.getByLabelText(/Email/i), 'test@example.com')
    await userEvent.type(screen.getByLabelText(/Password/i), 'secret')
    await userEvent.click(screen.getByRole('button', { name: /Sign In/i }))

    expect(loginMock).toHaveBeenCalledWith('test@example.com', 'secret')
    // navigate called with /dashboard
    expect(mockNavigate).toHaveBeenCalledWith('/dashboard')
  })

  it('shows error when login fails', async () => {
    const loginMock = vi.fn(() => Promise.reject({ response: { data: { error: 'Bad creds' } } }))
    useAuthMock.mockReturnValue({ login: loginMock, register: vi.fn() })
    render(<Login />)

    await userEvent.type(screen.getByLabelText(/Email/i), 'bad@example.com')
    await userEvent.type(screen.getByLabelText(/Password/i), 'wrong')
    await userEvent.click(screen.getByRole('button', { name: /Sign In/i }))

    expect(await screen.findByText(/Bad creds/)).toBeInTheDocument()
  })
})
