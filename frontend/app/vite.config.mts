import preact from "@preact/preset-vite";
import { loadEnv } from "vite";
import { defineConfig } from "vitest/config";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [preact()],
  server: {
    port: 8081,
  },
  css: {
    modules: {
      localsConvention: "camelCase",
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (
            id.includes("node_modules/@firebase") ||
            id.includes("node_modules/firebase")
          ) {
            return "firebase";
          }
          if (id.includes("node_modules/graphql")) {
            return "graphql";
          }
          if (id.includes("node_modules/@codemirror")) {
            return "codemirror";
          }
        },
      },
    },
  },
  test: {
    globals: true,
    environment: "happy-dom",
    watch: false,
    mockReset: true,
    env: loadEnv("test", process.cwd()),

    css: {
      modules: {
        classNameStrategy: "non-scoped",
      },
    },
  },
}));
