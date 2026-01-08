import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'node:path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Use relative paths for assets - allows runtime base path configuration via Go server
  base: './',
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
      external: ['./config.js'],
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom', '@tanstack/react-router'],
          apollo: ['@apollo/client', 'graphql'],
        },
      },
    },
  },
});
