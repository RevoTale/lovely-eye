import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'node:path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Use absolute paths from root for SPA routing compatibility
  base: '/',
  resolve: {
    alias: {
      '@': resolve(import.meta.dirname, './src'),
    },
  },
  build: {
    // Static export - output goes to dist/ folder, served by Go server
    outDir: 'dist',
    sourcemap: false,
    // Optimize for static serving
    minify: 'esbuild',
    rollupOptions: {
      // config.js is runtime-injected, not bundled
      //external: ['./config.js']
    },
  },
});
