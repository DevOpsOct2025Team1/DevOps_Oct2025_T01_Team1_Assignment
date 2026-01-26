import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router';
import type { User } from '../utils/api';
import { getStoredUser, isAuthenticated, isAdmin } from '../utils/auth';

export default function Dashboard() {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    if (!isAuthenticated()) {
      navigate('/login');
      return;
    }
    
    if (isAdmin()) {
      navigate('/admin');
      return;
    }

    const storedUser = getStoredUser();
    setUser(storedUser);
  }, [navigate]);

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p className="text-gray-600 mt-1">Welcome, {user.username}</p>
        </div>

        <div className="grid gap-6 md:grid-cols-3">
          <div className="bg-white p-6 rounded-lg shadow">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Upload File</h2>
              <svg className="w-8 h-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
            </div>
            <p className="text-sm text-gray-600 mb-4">Upload files to the system</p>
            <button className="w-full bg-gray-300 text-gray-700 py-2 px-4 rounded-md cursor-not-allowed">
              Feature Coming Soon
            </button>
            <p className="text-xs text-gray-500 mt-2 text-center">TODO: connect API endpoint</p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Download File</h2>
              <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
            </div>
            <p className="text-sm text-gray-600 mb-4">Download your files</p>
            <button className="w-full bg-gray-300 text-gray-700 py-2 px-4 rounded-md cursor-not-allowed">
              Feature Coming Soon
            </button>
            <p className="text-xs text-gray-500 mt-2 text-center">TODO: connect API endpoint</p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Delete File</h2>
              <svg className="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </div>
            <p className="text-sm text-gray-600 mb-4">Remove files from the system</p>
            <button className="w-full bg-gray-300 text-gray-700 py-2 px-4 rounded-md cursor-not-allowed">
              Feature Coming Soon
            </button>
            <p className="text-xs text-gray-500 mt-2 text-center">TODO: connect API endpoint</p>
          </div>
        </div>

        <div className="mt-8 bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">My Files</h2>
          <div className="text-center py-8 text-gray-500">
            <svg className="w-16 h-16 mx-auto mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <p>No files uploaded yet</p>
            <p className="text-xs mt-2">TODO: connect API endpoint for file listing</p>
          </div>
        </div>

        <div className="mt-6 bg-blue-50 border border-blue-200 p-4 rounded-lg">
          <h3 className="text-sm font-semibold text-blue-900 mb-2">Note</h3>
          <p className="text-sm text-blue-800">
            File management features (upload, download, delete) are placeholders. 
            Backend API endpoints for these operations need to be implemented.
          </p>
        </div>
      </div>
    </div>
  );
}
