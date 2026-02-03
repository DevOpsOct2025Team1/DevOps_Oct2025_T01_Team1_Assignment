import { screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import Login from '../app/routes/login';
import { BrowserRouter } from 'react-router';
import { renderWithProviders } from './test-utils';

const mutateAsyncMock = vi.fn();
let loginIsPending = false;

vi.mock('../app/api/generated', () => ({
  useLogin: () => ({
    mutateAsync: mutateAsyncMock,
    isPending: loginIsPending,
  }),
}));

describe('Login Page', () => {
  beforeEach(() => {
    localStorage.clear();
    mutateAsyncMock.mockReset();
    loginIsPending = false;
  });

  it('renders login form', async () => {
    renderWithProviders(
      <BrowserRouter>
        <Login />
      </BrowserRouter>
    );

    expect(screen.getByRole('heading', { name: /login/i })).toBeDefined();
    expect(screen.getByLabelText(/username/i)).toBeDefined();
    expect(screen.getByLabelText(/password/i)).toBeDefined();
    expect(screen.getByRole('button', { name: /login/i })).toBeDefined();
  });

  it('validates empty inputs', async () => {
    const user = userEvent.setup();
    
    renderWithProviders(
      <BrowserRouter>
        <Login />
      </BrowserRouter>
    );

    const loginButton = screen.getByRole('button', { name: /login/i });
    await user.click(loginButton);

    expect(await screen.findByText(/username and password are required/i)).toBeDefined();
  });

  it('disables the form while logging in', async () => {
    loginIsPending = true;

    renderWithProviders(
      <BrowserRouter>
        <Login />
      </BrowserRouter>
    );

    const loginButton = screen.getByRole('button', { name: /logging in/i });
    expect(loginButton).toBeDisabled();
    expect(screen.getByLabelText(/username/i)).toBeDisabled();
    expect(screen.getByLabelText(/password/i)).toBeDisabled();
  });

  it('shows error message on invalid credentials', async () => {
    const user = userEvent.setup();
    mutateAsyncMock.mockRejectedValueOnce({ status: 401, message: 'Unauthorized' });
    
    renderWithProviders(
      <BrowserRouter>
        <Login />
      </BrowserRouter>
    );

    const usernameInput = screen.getByLabelText(/username/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const loginButton = screen.getByRole('button', { name: /login/i });

    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'wrongpassword');
    await user.click(loginButton);

    expect(await screen.findByText(/unauthorized/i)).toBeDefined();
  });
});
