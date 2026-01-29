import { QueryClient } from '@tanstack/react-query';

export const defaultQueryOptions = {
  queries: {
    staleTime: 60_000,
    retry: 1,
    refetchOnWindowFocus: false,
  },
  mutations: {
    retry: 0,
  },
};

export const createQueryClient = () =>
  new QueryClient({
    defaultOptions: defaultQueryOptions,
  });

export const queryKeys = {
  auth: {
    login: () => ['auth', 'login'] as const,
  },
  admin: {
    createUser: () => ['admin', 'create-user'] as const,
    deleteUser: (id: string) => ['admin', 'delete-user', id] as const,
  },
};
