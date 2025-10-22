import { defineConfig } from 'vite'
import { codecovVitePlugin } from '@codecov/vite-plugin'

export default defineConfig({
  plugins: [
    // Put the Codecov vite plugin after all other plugins
    codecovVitePlugin({
      enableBundleAnalysis: process.env.CODECOV_TOKEN !== undefined,
      bundleName: 'cli.wasm',
      uploadToken: process.env.CODECOV_TOKEN,
    }),
  ],
})
