import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router';
import { authApi } from '../utils/api';
import type { User } from '../utils/api';
import { getStoredUser, isAuthenticated, isAdmin } from '../utils/auth';

interface Action {
  type: 'create' | 'delete';
  username: string;
  timestamp: Date;
}

export default function Admin() {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [actions, setActions] = useState<Action[]>([]);
  
  const [createUsername, setCreateUsername] = useState('');
  const [createPassword, setCreatePassword] = useState('');
  const [createError, setCreateError] = useState('');
  const [createLoading, setCreateLoading] = useState(false);
  
  const [deleteUserId, setDeleteUserId] = useState('');
  const [deleteError, setDeleteError] = useState('');
  const [deleteLoading, setDeleteLoading] = useState(false);

  useEffect(() => {
    if (!isAuthenticated()) {
      navigate('/login');
      return;
    }
    
    if (!isAdmin()) {
      navigate('/dashboard');
      return;
    }

    const storedUser = getStoredUser();
    setUser(storedUser);
  }, [navigate]);

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreateError('');

    if (!createUsername.trim() || !createPassword.trim()) {
      setCreateError('Username and password are required');
      return;
    }

    setCreateLoading(true);
    try {
      const response = await authApi.createUser({
        username: createUsername,
        password: createPassword,
      });
      
      setActions([
        {
          type: 'create',
          username: response.user.username,
          timestamp: new Date(),
        },
        ...actions,
      ]);
      
      setCreateUsername('');
      setCreatePassword('');
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : 'Failed to create user');
    } finally {
      setCreateLoading(false);
    }
  };

  const handleDeleteUser = async (e: React.FormEvent) => {
    e.preventDefault();
    setDeleteError('');

    if (!deleteUserId.trim()) {
      setDeleteError('User ID is required');
      return;
    }

    setDeleteLoading(true);
    try {
      await authApi.deleteUser(deleteUserId);
      
      setActions([
        {
          type: 'delete',
          username: `User ${deleteUserId}`,
          timestamp: new Date(),
        },
        ...actions,
      ]);
      
      setDeleteUserId('');
    } catch (err) {
      setDeleteError(err instanceof Error ? err.message : 'Failed to delete user');
    } finally {
      setDeleteLoading(false);
    }
  };

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
          <p className="text-gray-600 mt-1">Welcome, {user.username}</p>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          <div className="bg-white p-6 rounded-lg shadow">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Create User</h2>
            <form onSubmit={handleCreateUser} className="space-y-4">
              <div>
                <label htmlFor="create-username" className="block text-sm font-medium text-gray-700 mb-1">
                  Username
                </label>
                <input
                  id="create-username"
                  type="text"
                  value={createUsername}
                  onChange={(e) => setCreateUsername(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  disabled={createLoading}
                />
              </div>

              <div>
                <label htmlFor="create-password" className="block text-sm font-medium text-gray-700 mb-1">
                  Password
                </label>
                <input
                  id="create-password"
                  type="password"
                  value={createPassword}
                  onChange={(e) => setCreatePassword(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  disabled={createLoading}
                />
              </div>

              {createError && (
                <div className="text-red-600 text-sm bg-red-50 p-3 rounded-md">
                  {createError}
                </div>
              )}

              <button
                type="submit"
                disabled={createLoading}
                className="w-full bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 disabled:bg-indigo-400 disabled:cursor-not-allowed transition-colors font-medium"
              >
                {createLoading ? 'Creating...' : 'Create User'}
              </button>
            </form>
          </div>

          <div className="bg-white p-6 rounded-lg shadow">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Delete User</h2>
            <form onSubmit={handleDeleteUser} className="space-y-4">
              <div>
                <label htmlFor="delete-id" className="block text-sm font-medium text-gray-700 mb-1">
                  User ID
                </label>
                <input
                  id="delete-id"
                  type="text"
                  value={deleteUserId}
                  onChange={(e) => setDeleteUserId(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-red-500 focus:border-transparent"
                  disabled={deleteLoading}
                />
              </div>

              {deleteError && (
                <div className="text-red-600 text-sm bg-red-50 p-3 rounded-md">
                  {deleteError}
                </div>
              )}

              <button
                type="submit"
                disabled={deleteLoading}
                className="w-full bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700 disabled:bg-red-400 disabled:cursor-not-allowed transition-colors font-medium"
              >
                {deleteLoading ? 'Deleting...' : 'Delete User'}
              </button>
            </form>
          </div>
        </div>

        <div className="mt-8 bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Recent Actions</h2>
          {actions.length === 0 ? (
            <p className="text-gray-500 text-center py-4">No actions yet</p>
          ) : (
            <div className="space-y-2">
              {actions.map((action, idx) => (
                <div key={idx} className="flex items-center justify-between py-2 px-4 bg-gray-50 rounded">
                  <div className="flex items-center space-x-3">
                    <span className={`px-2 py-1 text-xs font-medium rounded ${
                      action.type === 'create' 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {action.type.toUpperCase()}
                    </span>
                    <span className="text-gray-700">{action.username}</span>
                  </div>
                  <span className="text-sm text-gray-500">
                    {action.timestamp.toLocaleTimeString()}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
