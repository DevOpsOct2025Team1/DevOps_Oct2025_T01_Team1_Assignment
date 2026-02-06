import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { UseMutationOptions, UseQueryOptions } from "@tanstack/react-query";
import { customFetch } from "./orval-client";
import type { InternalHandlersUserResponse } from "./generated/model";

// Types for the missing endpoints (not in swagger.json)
export interface ListUsersResponse {
  users: InternalHandlersUserResponse[];
}

export interface UpdateUserRoleRequest {
  id: string;
  role: string;
}

export interface UpdateUserRoleResponse {
  user: InternalHandlersUserResponse;
}

// List users query (GET /api/admin)
const listUsers = async (): Promise<ListUsersResponse> => {
  const response = await customFetch<{ data: ListUsersResponse }>(
    "/api/admin",
    { method: "GET" }
  );
  return response.data;
};

export const useListUsers = (
  options?: Omit<
    UseQueryOptions<ListUsersResponse, Error>,
    "queryKey" | "queryFn"
  >
) => {
  return useQuery<ListUsersResponse, Error>({
    queryKey: ["admin", "users"],
    queryFn: listUsers,
    ...options,
  });
};

// Update user role mutation (POST /api/admin/update_user_role)
const updateUserRole = async (
  data: UpdateUserRoleRequest
): Promise<UpdateUserRoleResponse> => {
  const response = await customFetch<{ data: UpdateUserRoleResponse }>(
    "/api/admin/update_user_role",
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    }
  );
  return response.data;
};

export const useUpdateUserRole = (
  options?: UseMutationOptions<
    UpdateUserRoleResponse,
    Error,
    UpdateUserRoleRequest
  >
) => {
  const queryClient = useQueryClient();

  return useMutation<UpdateUserRoleResponse, Error, UpdateUserRoleRequest>({
    mutationFn: updateUserRole,
    onSuccess: () => {
      // Invalidate users list to refetch after role update
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
    ...options,
  });
};
