import { useState } from "react";
import { useCreateUser, useDeleteUser } from "../api/generated";

interface Action {
  type: "create" | "delete";
  username: string;
  timestamp: Date;
}

export default function AdminPanel() {
  const [actions, setActions] = useState<Action[]>([]);
  const createUserMutation = useCreateUser();
  const deleteUserMutation = useDeleteUser();

  const [createUsername, setCreateUsername] = useState("");
  const [createPassword, setCreatePassword] = useState("");
  const [createError, setCreateError] = useState("");

  const [deleteUserId, setDeleteUserId] = useState("");
  const [deleteError, setDeleteError] = useState("");

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreateError("");

    if (!createUsername.trim() || !createPassword.trim()) {
      setCreateError("Username and password are required");
      return;
    }

    try {
      const response = await createUserMutation.mutateAsync({
        data: {
          username: createUsername,
          password: createPassword,
        },
      });
      const createdUser =
        response.data && "user" in response.data ? response.data.user : undefined;

      setActions([
        {
          type: "create",
          username: createdUser?.username || createUsername,
          timestamp: new Date(),
        },
        ...actions,
      ]);

      setCreateUsername("");
      setCreatePassword("");
    } catch (err: any) {
      setCreateError(err?.message || "Failed to create user");
    }
  };

  const handleDeleteUser = async (e: React.FormEvent) => {
    e.preventDefault();
    setDeleteError("");

    if (!deleteUserId.trim()) {
      setDeleteError("User ID is required");
      return;
    }

    try {
      await deleteUserMutation.mutateAsync({
        data: {
          id: deleteUserId,
        },
      });

      setActions([
        {
          type: "delete",
          username: `User ${deleteUserId}`,
          timestamp: new Date(),
        },
        ...actions,
      ]);

      setDeleteUserId("");
    } catch (err: any) {
      setDeleteError(err?.message || "Failed to delete user");
    }
  };

  return (
    <>
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
                disabled={createUserMutation.isPending}
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
                disabled={createUserMutation.isPending}
              />
            </div>

            {createError && (
              <div className="text-red-600 text-sm bg-red-50 p-3 rounded-md">
                {createError}
              </div>
            )}

            <button
              type="submit"
              disabled={createUserMutation.isPending}
              className="w-full bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 disabled:bg-indigo-400 disabled:cursor-not-allowed transition-colors font-medium"
            >
              {createUserMutation.isPending ? "Creating..." : "Create User"}
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
                disabled={deleteUserMutation.isPending}
              />
            </div>

            {deleteError && (
              <div className="text-red-600 text-sm bg-red-50 p-3 rounded-md">
                {deleteError}
              </div>
            )}

            <button
              type="submit"
              disabled={deleteUserMutation.isPending}
              className="w-full bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700 disabled:bg-red-400 disabled:cursor-not-allowed transition-colors font-medium"
            >
              {deleteUserMutation.isPending ? "Deleting..." : "Delete User"}
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
                    action.type === "create"
                      ? "bg-green-100 text-green-800"
                      : "bg-red-100 text-red-800"
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
    </>
  );
}