import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router';
import { setAuth, isAuthenticated, isAdmin } from '../utils/auth';
import { useLogin } from '../api/generated';

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const loginMutation = useLogin();
  const isLoading = loginMutation.isPending;

  useEffect(() => {
    if (isAuthenticated()) {
      if (isAdmin()) {
        navigate('/admin');
      } else {
        navigate('/dashboard');
      }
    }
  }, [navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!username.trim() || !password.trim()) {
      setError('Username and password are required');
      return;
    }

    try {
      const response = await loginMutation.mutateAsync({
        data: {
          username,
          password,
        },
      });
      const authData = response.data;

      if (
        !authData ||
        !('user' in authData) ||
        !authData.user ||
        !authData.user.id ||
        !authData.user.username ||
        !authData.user.role ||
        !authData.token
      ) {
        throw new Error('Login failed. Please try again.');
      }

      console.log('Login response:', authData);
      console.log('User role:', authData.user.role);
      
      // Normalize API roles into the UI-friendly values we already use.
      const roleValue = authData.user.role as string | number;
      const userRole = typeof roleValue === 'number'
        ? (roleValue === 2 ? 'admin' : 'user')
        : roleValue.toLowerCase().includes('admin')
          ? 'admin'
          : 'user';
      
      const normalizedUser = {
        id: authData.user.id,
        username: authData.user.username,
        role: userRole
      };
      
      setAuth(normalizedUser, authData.token);
      
      if (userRole === 'admin') {
        console.log('Navigating to /admin');
        navigate('/admin');
      } else {
        console.log('Navigating to /dashboard');
        navigate('/dashboard');
      }
    } catch (err: any) {
      if (err?.status === 401) {
        setError('Incorrect Username or Password');
      } else if (err?.message) {
        setError(err.message);
      } else {
        setError('Login failed. Please try again.');
      }
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
        <h1 className="text-3xl font-bold text-gray-800 mb-6 text-center">Login</h1>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">
              Username
            </label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              disabled={isLoading}
            />
          </div>

          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              disabled={isLoading}
            />
          </div>

          {error && (
            <div className="text-red-600 text-sm bg-red-50 p-3 rounded-md">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loginMutation.isPending}
            className="w-full bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 disabled:bg-indigo-400 disabled:cursor-not-allowed transition-colors font-medium"
          >
            {loginMutation.isPending ? 'Logging in...' : 'Login'}
          </button>
        </form>
      </div>
    </div>
  );
}
