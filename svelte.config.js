import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			fallback: '404.html',
			strict: false
		}),
		paths: {
			// Set this to your GitHub Pages repo name if deploying as a project site.
			// For a user/org site (username.github.io), leave this empty.
			// Example: base: '/calibrarytracker'
			base: ''
		}
	}
};

export default config;
