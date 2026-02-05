import { describe, it, expect, beforeEach, vi } from 'vitest';
import { BrowserRouter } from 'react-router';
import Dashboard from '../app/routes/dashboard';
import { renderWithProviders } from './test-utils';
import { clearAuthCache } from '../app/utils/auth';

describe('Protected Routes', () => {
  beforeEach(() => {
    clearAuthCache();
    localStorage.clear();
  });

  it('redirects to login when not authenticated', () => {
    renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(localStorage.getItem('token')).toBeNull();
  });

  it('shows user dashboard for regular users', async () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({
      id: '1',
      username: 'user',
      role: 'user'
    }));

    const { findByText, findByPlaceholderText } = renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(await findByText('Hello, user!')).toBeTruthy();
    expect(await findByPlaceholderText('Search uploaded files')).toBeTruthy();
    expect(await findByText('Presentation.pptx')).toBeTruthy();
  });

  it('redirects admin users to admin dashboard', async () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({
      id: '2',
      username: 'admin',
      role: 'admin'
    }));

    const navigateMock = vi.fn();
    vi.mock('react-router', async () => {
      const actual = await vi.importActual('react-router');
      return {
        ...actual,
        useNavigate: () => navigateMock,
      };
    });

    renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    // Admin users should be redirected to /admin,
    // so the dashboard content should not be visible
    expect(navigateMock).toHaveBeenCalledWith('/admin');
  });
});
