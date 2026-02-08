import { describe, it, expect, beforeEach, vi } from 'vitest';
import { BrowserRouter } from 'react-router';
import Dashboard from '../app/routes/dashboard';
import { renderWithProviders } from './test-utils';
import { clearAuthCache } from '../app/utils/auth';

const createUserMock = vi.fn(() => ({
  mutateAsync: vi.fn(),
  isPending: false,
}));

const deleteUserMock = vi.fn(() => ({
  mutateAsync: vi.fn(),
  isPending: false,
}));

const getUsersMock = vi.fn(() => ({
  data: [],
  isLoading: false,
}));

vi.mock('../app/api/generated', () => ({
  usePostApiAdminCreateUser: () => createUserMock(),
  useDeleteApiAdminDeleteUser: () => deleteUserMock(),
  useGetApiAdminListUsers: () => getUsersMock(),
  getGetApiAdminListUsersQueryKey: () => ['users'],
}));

describe('Protected Routes', () => {
  beforeEach(() => {
    vi.resetModules();
    clearAuthCache();
    localStorage.clear();
    createUserMock.mockClear();
    deleteUserMock.mockClear();
    getUsersMock.mockClear();
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

    expect(await findByText(/Admin/i)).toBeTruthy();
    expect(await findByText('What would you like to manage today?')).toBeTruthy();
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
