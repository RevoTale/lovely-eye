import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import { resolve } from 'node:path';
import tailwindcss from '@tailwindcss/vite';
import { tanstackRouter } from '@tanstack/router-plugin/vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tanstackRouter({
      routesDirectory: './src/routes',
      generatedRouteTree: './src/routeTree.gen.ts',
      autoCodeSplitting: true,
    }),
    tailwindcss(),
    react(),
  ],
  // Use absolute paths from root for SPA routing compatibility
  base: '/',
  resolve: {
    alias: {
      '@': resolve(import.meta.dirname, './src'),
    },
  },
  
  css: {
    transformer: 'lightningcss',
  },
  build: {
    // Static export - output goes to dist/ folder, served by Go server
    outDir: 'dist',
    sourcemap: false,
    // Optimize for static serving
    minify: true,
    cssMinify: 'lightningcss',
    rollupOptions: {
      treeshake: true,
    },
  },
});
