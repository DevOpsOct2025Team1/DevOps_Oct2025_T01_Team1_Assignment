import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router';
import { authApi } from '../utils/api';
import type { User } from '../utils/api';
import { getStoredUser as getStoredAuthUser, isAuthenticated, isAdmin } from '../utils/auth';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Badge } from '../components/ui/badge';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '../components/ui/dialog';

interface Action {
  type: 'create' | 'delete' | 'update';
  username: string;
  details?: string;
  timestamp: Date;
}

const ACTIONS_STORAGE_KEY = 'admin_recent_actions';

function loadActions(): Action[] {
  try {
    const stored = localStorage.getItem(ACTIONS_STORAGE_KEY);
    if (stored) {
      const parsed = JSON.parse(stored);
      return parsed.map((action: any) => ({
        ...action,
        timestamp: new Date(action.timestamp),
      }));
    }
  } catch (error) {
    console.error('Failed to load actions:', error);
  }
  return [];
}

function saveActions(actions: Action[]) {
  try {
    localStorage.setItem(ACTIONS_STORAGE_KEY, JSON.stringify(actions));
  } catch (error) {
    console.error('Failed to save actions:', error);
  }
}

export default function Admin() {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [filteredUsers, setFilteredUsers] = useState<User[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [roleFilter, setRoleFilter] = useState<'all' | 'admin' | 'user'>('all');
  const [actions, setActions] = useState<Action[]>([]);
  const [loading, setLoading] = useState(true);
  
  // Add User Dialog
  const [addUserOpen, setAddUserOpen] = useState(false);
  const [newUsername, setNewUsername] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [addUserError, setAddUserError] = useState('');
  const [addUserLoading, setAddUserLoading] = useState(false);
  
  // Delete User Dialog
  const [deleteUserOpen, setDeleteUserOpen] = useState(false);
  const [userToDelete, setUserToDelete] = useState<User | null>(null);
  const [deleteUserLoading, setDeleteUserLoading] = useState(false);
  
  // Change Role Dialog
  const [changeRoleOpen, setChangeRoleOpen] = useState(false);
  const [userToUpdate, setUserToUpdate] = useState<User | null>(null);
  const [newRole, setNewRole] = useState<string>('');
  const [changeRoleLoading, setChangeRoleLoading] = useState(false);

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      const response = await authApi.listUsers();
      setUsers(response.users);
      setFilteredUsers(response.users);
    } catch (error) {
      console.error('Failed to fetch users:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!isAuthenticated()) {
      navigate('/login');
      return;
    }
    
    if (!isAdmin()) {
      navigate('/dashboard');
      return;
    }

    const storedUser = getStoredAuthUser();
    if (storedUser) {
      // Convert role to string for consistent typing
      setUser({
        ...storedUser,
        role: typeof storedUser.role === 'number' 
          ? (storedUser.role === 2 ? 'admin' : 'user')
          : storedUser.role
      });
    }
    
    // Load actions from localStorage
    setActions(loadActions());
    
    // Fetch users list
    fetchUsers();
  }, [navigate, fetchUsers]);

  useEffect(() => {
    const filtered = users.filter(u => {
      const matchesSearch = u.username.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesRole = roleFilter === 'all' || u.role === roleFilter;
      return matchesSearch && matchesRole;
    });
    setFilteredUsers(filtered);
  }, [searchQuery, roleFilter, users]);

  const addAction = (action: Action) => {
    const newActions = [action, ...actions];
    setActions(newActions);
    saveActions(newActions);
  };

  const handleAddUser = async (e: React.FormEvent) => {
    e.preventDefault();
    setAddUserError('');

    if (!newUsername.trim() || !newPassword.trim()) {
      setAddUserError('Username and password are required');
      return;
    }

    setAddUserLoading(true);
    try {
      const response = await authApi.createUser({
        username: newUsername,
        password: newPassword,
      });
      
      addAction({
        type: 'create',
        username: response.user.username,
        timestamp: new Date(),
      });
      
      setNewUsername('');
      setNewPassword('');
      setAddUserOpen(false);
      fetchUsers();
    } catch (err) {
      setAddUserError(err instanceof Error ? err.message : 'Failed to create user');
    } finally {
      setAddUserLoading(false);
    }
  };

  const handleDeleteUser = async () => {
    if (!userToDelete) return;

    console.log('Deleting user with ID:', userToDelete.id);
    setDeleteUserLoading(true);
    try {
      await authApi.deleteUser(userToDelete.id);
      
      addAction({
        type: 'delete',
        username: userToDelete.username,
        timestamp: new Date(),
      });
      
      setDeleteUserOpen(false);
      setUserToDelete(null);
      fetchUsers();
    } catch (error) {
      console.error('Failed to delete user:', error);
      console.error('User ID was:', userToDelete.id);
    } finally {
      setDeleteUserLoading(false);
    }
  };

  const handleChangeRole = async () => {
    if (!userToUpdate || !newRole) return;

    setChangeRoleLoading(true);
    try {
      await authApi.updateUserRole({
        id: userToUpdate.id,
        role: newRole,
      });
      
      addAction({
        type: 'update',
        username: userToUpdate.username,
        details: `Role changed to ${newRole}`,
        timestamp: new Date(),
      });
      
      setChangeRoleOpen(false);
      setUserToUpdate(null);
      setNewRole('');
      fetchUsers();
    } catch (error) {
      console.error('Failed to update user role:', error);
    } finally {
      setChangeRoleLoading(false);
    }
  };

  const openDeleteDialog = (user: User) => {
    setUserToDelete(user);
    setDeleteUserOpen(true);
  };

  const openChangeRoleDialog = (user: User) => {
    setUserToUpdate(user);
    setNewRole(user.role === 'admin' ? 'user' : 'admin');
    setChangeRoleOpen(true);
  };

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4">
        {/* Greeting Section */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Hello, {user.username}!</h1>
          <p className="text-gray-600 mt-2">
            Manage your users, update roles, and monitor system activity from this dashboard.
          </p>
        </div>

        <div className="bg-white rounded-lg shadow">
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-lg font-semibold text-gray-900">Users</h2>
                <p className="text-sm text-gray-500 mt-1">
                  {loading ? 'Loading...' : `${filteredUsers.length} user${filteredUsers.length !== 1 ? 's' : ''} ${searchQuery || roleFilter !== 'all' ? 'found' : 'total'}`}
                </p>
              </div>
              <Button onClick={() => setAddUserOpen(true)}>
                <svg className="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                Add User
              </Button>
            </div>
            <div className="flex gap-3 flex-1 max-w-2xl">
              <div className="relative flex-1">
                <Input
                  type="text"
                  placeholder="Search Users"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
                <svg
                  className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  />
                </svg>
              </div>
              <div className="relative min-w-[140px]">
                <select
                  value={roleFilter}
                  onChange={(e) => setRoleFilter(e.target.value as 'all' | 'admin' | 'user')}
                  className="w-full appearance-none px-4 py-2 pr-10 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors cursor-pointer"
                >
                  <option value="all">Show All</option>
                  <option value="admin">Admins</option>
                  <option value="user">Users</option>
                </select>
                <svg
                  className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-500 pointer-events-none"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
          </div>

          {/* Users List */}
          <div className="divide-y divide-gray-200">
            {loading ? (
              <div className="p-8 text-center text-gray-500">Loading users...</div>
            ) : filteredUsers.length === 0 ? (
              <div className="p-8 text-center text-gray-500">No users found</div>
            ) : (
              filteredUsers.map((u) => (
                <div key={u.id} className="p-4 flex items-center justify-between hover:bg-gray-50">
                  <div className="flex-1">
                    <span className="text-gray-900 font-medium">{u.username}</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge variant={u.role === 'admin' ? 'default' : 'secondary'}>
                      {u.role === 'admin' ? 'Admin' : 'User'}
                    </Badge>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                          <span className="sr-only">Open menu</span>
                          <svg className="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                            <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
                          </svg>
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onSelect={() => openChangeRoleDialog(u)}>
                          Change Role
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          onSelect={() => openDeleteDialog(u)}
                          className="text-red-600 hover:text-red-700"
                        >
                          Delete User
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Recent Actions */}
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
                        : action.type === 'delete'
                        ? 'bg-red-100 text-red-800'
                        : 'bg-blue-100 text-blue-800'
                    }`}>
                      {action.type.toUpperCase()}
                    </span>
                    <span className="text-gray-700">
                      {action.username}
                      {action.details && <span className="text-gray-500 text-sm ml-2">({action.details})</span>}
                    </span>
                  </div>
                  <span className="text-sm text-gray-500">
                    {action.timestamp.toLocaleString()}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Add User Dialog */}
      <Dialog open={addUserOpen} onOpenChange={setAddUserOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add New User</DialogTitle>
            <DialogDescription>
              Create a new user account with username and password.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleAddUser}>
            <div className="space-y-4">
              <div>
                <label htmlFor="new-username" className="block text-sm font-medium text-gray-700 mb-1">
                  Username
                </label>
                <Input
                  id="new-username"
                  type="text"
                  value={newUsername}
                  onChange={(e) => setNewUsername(e.target.value)}
                  disabled={addUserLoading}
                />
              </div>
              <div>
                <label htmlFor="new-password" className="block text-sm font-medium text-gray-700 mb-1">
                  Password
                </label>
                <Input
                  id="new-password"
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  disabled={addUserLoading}
                />
              </div>
              {addUserError && (
                <div className="text-red-600 text-sm bg-red-50 p-3 rounded-md">
                  {addUserError}
                </div>
              )}
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="ghost"
                onClick={() => setAddUserOpen(false)}
                disabled={addUserLoading}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={addUserLoading}>
                {addUserLoading ? 'Creating...' : 'Create User'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete User Dialog */}
      <Dialog open={deleteUserOpen} onOpenChange={setDeleteUserOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete User</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete user "{userToDelete?.username}"? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="ghost"
              onClick={() => setDeleteUserOpen(false)}
              disabled={deleteUserLoading}
            >
              Cancel
            </Button>
            <Button
              onClick={handleDeleteUser}
              disabled={deleteUserLoading}
              className="bg-red-600 hover:bg-red-700"
            >
              {deleteUserLoading ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Change Role Dialog */}
      <Dialog open={changeRoleOpen} onOpenChange={setChangeRoleOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Change User Role</DialogTitle>
            <DialogDescription>
              Change the role of user "{userToUpdate?.username}" to {newRole === 'admin' ? 'Admin' : 'User'}.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="ghost"
              onClick={() => setChangeRoleOpen(false)}
              disabled={changeRoleLoading}
            >
              Cancel
            </Button>
            <Button onClick={handleChangeRole} disabled={changeRoleLoading}>
              {changeRoleLoading ? 'Updating...' : 'Change Role'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
