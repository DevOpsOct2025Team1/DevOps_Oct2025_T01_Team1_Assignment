import { defineConfig } from 'vitest/config'
import { playwright } from '@vitest/browser-playwright'
import react from '@vitejs/plugin-react'
import tailwindcss from "@tailwindcss/vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [tailwindcss(), react(), tsconfigPaths()],
  test: {
    setupFiles: ["tests/setup.ts"],
    globals: true,
    coverage: {
      provider: 'istanbul',
      reporter: ['text', 'cobertura'],
      reportsDirectory: './coverage',
    },
    browser: {
      enabled: true,
      headless: true,
      provider: playwright(),
      // https://vitest.dev/config/browser/playwright
      instances: [
        { browser: 'chromium' },
        { browser: 'firefox' },
      ],
    },
  },
})
