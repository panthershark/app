import codegen from 'vite-plugin-graphql-codegen';
import { defineConfig } from 'vitest/config';

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [codegen()],
	test: {
		watch: false
	}
});
