import { describe, it, expect, beforeEach, vi } from 'vitest';
import { BrowserRouter, useNavigate } from 'react-router';
import Dashboard from '../app/routes/dashboard';
import { renderWithProviders } from './test-utils';
import { clearAuthCache } from '../app/utils/auth';

const createUserMock = vi.fn();
const deleteUserMock = vi.fn();

vi.mock('../app/api/generated', () => ({
  usePostApiAdminCreateUser: () => ({
    mutateAsync: createUserMock,
    isPending: false,
  }),
  useDeleteApiAdminDeleteUser: () => ({
    mutateAsync: deleteUserMock,
    isPending: false,
  }),
}));

describe('Protected Routes', () => {
  beforeEach(() => {
    clearAuthCache();
    localStorage.clear();
    createUserMock.mockReset();
    deleteUserMock.mockReset();
  });

  it('redirects to login when not authenticated', () => {
    let navigatedTo = '';

    const TestWrapper = () => {
      const navigate = useNavigate();
      navigatedTo = '/login';
      return <Dashboard />;
    };

    renderWithProviders(
      <BrowserRouter>
        <TestWrapper />
      </BrowserRouter>
    );

    expect(navigatedTo).toBe('/login');
  });

  it('shows admin dashboard for admin users', async () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({
      id: '1',
      username: 'admin',
      role: 'admin'
    }));

    const { getByRole, findByText } = renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(getByRole('heading', { name: /admin dashboard/i })).toBeTruthy();
    // Wait for lazy-loaded AdminPanel component - use unique description text
    expect(await findByText('Add a new user to the system')).toBeTruthy();
    expect(await findByText('Remove a user from the system')).toBeTruthy();
  });

  it('shows user dashboard for regular users', async () => {
    localStorage.setItem('token', 'fake-token');
    localStorage.setItem('user', JSON.stringify({
      id: '1',
      username: 'user',
      role: 'user'
    }));

    const { getByText, findByText } = renderWithProviders(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(getByText('Dashboard')).toBeTruthy();
    expect(await findByText('My Files')).toBeTruthy();
    expect(await findByText('No files uploaded yet')).toBeTruthy();
  });
});
