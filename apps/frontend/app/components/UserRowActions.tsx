import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useDeleteApiAdminDeleteUser, getGetApiAdminListUsersQueryKey } from "../api/generated";
import { Button } from "./ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "./ui/alert-dialog";
import { MoreVertical, Trash2 } from "lucide-react";

interface User {
  id: string;
  username: string;
  role: string;
}

interface UserRowActionsProps {
  user: User;
}

export function UserRowActions({ user }: UserRowActionsProps) {
  const queryClient = useQueryClient();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  const deleteUserMutation = useDeleteApiAdminDeleteUser({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({queryKey: getGetApiAdminListUsersQueryKey()}).then();
        setShowDeleteDialog(false);
      },
    },
  });

  const handleDelete = async () => {
    try {
      await deleteUserMutation.mutateAsync({
        data: { id: user.id },
      });
    } catch (err) {
      console.error("Failed to delete user", err);
    }
  };

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <MoreVertical className="h-4 w-4" />
            <span className="sr-only">Open menu</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem
            className="text-red-600 focus:text-red-600 cursor-pointer"
            onSelect={() => setShowDeleteDialog(true)}
          >
            <Trash2 className="h-4 w-4" />
            Delete User
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the user
              <span className="font-semibold"> {user.username}</span> and remove their
              data from our servers.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              className="bg-red-600 hover:bg-red-700 focus:ring-red-600"
            >
              {deleteUserMutation.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
