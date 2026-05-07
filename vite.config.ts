import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	build: {
		rollupOptions: {
			output: {
				manualChunks(id) {
					if (id.includes('/data/branches.json')) return 'ref-branches';
					if (id.includes('/data/systems.json')) return 'ref-systems';
					if (id.includes('/data/cooperatives.json')) return 'ref-cooperatives';
					if (id.includes('/data/partnerships.json')) return 'ref-partnerships';
					if (id.includes('leaflet.markercluster')) return 'leaflet-markercluster';
					if (id.includes('/leaflet/')) return 'leaflet';
				}
			}
		}
	}
});
