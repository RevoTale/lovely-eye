import { resolve } from 'node:path';
import tailwindcss from '@tailwindcss/vite';
import { tanstackRouter } from '@tanstack/router-plugin/vite';
import react from '@vitejs/plugin-react-swc';
import { defineConfig, type Plugin } from 'vite';

const runtimeConfigPlugin = (): Plugin => ({
  name: 'lovely-eye-runtime-config',
  transformIndexHtml: {
    order: 'post',
    handler: (html) =>
      html.replace('<!-- runtime-config -->', '<script src="{{BASE_PATH}}/config.js"></script>'),
  },
});

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    runtimeConfigPlugin(),
    tanstackRouter({
      routesDirectory: './src/app/routes',
      generatedRouteTree: './src/app/route-tree.gen.ts',
      autoCodeSplitting: true,
    }),
    tailwindcss(),
    react(),
  ],
  // The deployment path is unknown at build time; <base> is injected by Go at runtime.
  base: './',
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
    rolldownOptions: {
      treeshake: true,
    },
  },
});
