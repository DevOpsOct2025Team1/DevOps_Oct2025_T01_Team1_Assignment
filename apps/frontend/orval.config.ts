import { defineConfig } from 'orval';

export default defineConfig({
  api: {
    input: '../api-gateway/docs/swagger.json',
    output: {
      mode: 'single',
      target: './app/api/generated/index.ts',
      schemas: './app/api/generated/model',
      client: 'react-query',
      httpClient: 'fetch',
      baseUrl: '',
      override: {
        mutator: {
          path: './app/api/orval-client.ts',
          name: 'customFetch',
        },
      },
    },
  },
});
