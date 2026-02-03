import { useState } from "react";
import { useCreateUser, useDeleteUser } from "../api/generated";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./ui/card";
import { Badge } from "./ui/badge";

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
        <Card>
          <CardHeader>
            <CardTitle>Create User</CardTitle>
            <CardDescription>Add a new user to the system</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreateUser} className="space-y-4">
              <div className="space-y-2">
                <label htmlFor="create-username" className="text-sm font-medium">
                  Username
                </label>
                <Input
                  id="create-username"
                  type="text"
                  value={createUsername}
                  onChange={(e) => setCreateUsername(e.target.value)}
                  disabled={createUserMutation.isPending}
                  placeholder="Enter username"
                />
              </div>

              <div className="space-y-2">
                <label htmlFor="create-password" className="text-sm font-medium">
                  Password
                </label>
                <Input
                  id="create-password"
                  type="password"
                  value={createPassword}
                  onChange={(e) => setCreatePassword(e.target.value)}
                  disabled={createUserMutation.isPending}
                  placeholder="Enter password"
                />
              </div>

              {createError && (
                <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
                  {createError}
                </div>
              )}

              <Button
                type="submit"
                disabled={createUserMutation.isPending}
                className="w-full"
              >
                {createUserMutation.isPending ? "Creating..." : "Create User"}
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Delete User</CardTitle>
            <CardDescription>Remove a user from the system</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleDeleteUser} className="space-y-4">
              <div className="space-y-2">
                <label htmlFor="delete-id" className="text-sm font-medium">
                  User ID
                </label>
                <Input
                  id="delete-id"
                  type="text"
                  value={deleteUserId}
                  onChange={(e) => setDeleteUserId(e.target.value)}
                  disabled={deleteUserMutation.isPending}
                  placeholder="Enter user ID"
                />
              </div>

              {deleteError && (
                <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
                  {deleteError}
                </div>
              )}

              <Button
                type="submit"
                disabled={deleteUserMutation.isPending}
                variant="destructive"
                className="w-full"
              >
                {deleteUserMutation.isPending ? "Deleting..." : "Delete User"}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Recent Actions</CardTitle>
          <CardDescription>View recent user management activities</CardDescription>
        </CardHeader>
        <CardContent>
          {actions.length === 0 ? (
            <p className="text-muted-foreground text-center py-4">No actions yet</p>
          ) : (
            <div className="space-y-2">
              {actions.map((action, idx) => (
                <div key={idx} className="flex items-center justify-between py-2 px-4 bg-muted/50 rounded">
                  <div className="flex items-center space-x-3">
                    <Badge variant={action.type === "create" ? "default" : "destructive"}>
                      {action.type.toUpperCase()}
                    </Badge>
                    <span>{action.username}</span>
                  </div>
                  <span className="text-sm text-muted-foreground">
                    {action.timestamp.toLocaleTimeString()}
                  </span>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </>
  );
}