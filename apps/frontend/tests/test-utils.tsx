import type { ReactElement } from 'react';
import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

export const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
      mutations: {
        retry: false,
      },
    },
  });

export const renderWithProviders = (ui: ReactElement) => {
  const queryClient = createTestQueryClient();
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
};
