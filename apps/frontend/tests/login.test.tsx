import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, beforeEach } from 'vitest';
import Login from '../app/routes/login';
import { BrowserRouter } from 'react-router';

describe('Login Page', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('renders login form', async () => {
    render(
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
    
    render(
      <BrowserRouter>
        <Login />
      </BrowserRouter>
    );

    const loginButton = screen.getByRole('button', { name: /login/i });
    await user.click(loginButton);

    expect(await screen.findByText(/username and password are required/i)).toBeDefined();
  });

  it('shows error message on invalid credentials', async () => {
    const user = userEvent.setup();
    
    render(
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

    expect(await screen.findByText(/login failed/i)).toBeDefined();
  });
});
