import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
	testDir: './tests/smoke',
	fullyParallel: false,
	retries: process.env.CI ? 2 : 0,
	reporter: process.env.CI ? [['github'], ['list']] : 'list',
	use: {
		baseURL: 'http://127.0.0.1:4173',
		trace: 'retain-on-failure'
	},
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] }
		}
	],
	webServer: {
		command:
			'rm -rf .playwright-pages && mkdir -p .playwright-pages/calibrarytracker && cp -R build/. .playwright-pages/calibrarytracker && python3 -m http.server 4173 -d .playwright-pages',
		url: 'http://127.0.0.1:4173/calibrarytracker/',
		reuseExistingServer: !process.env.CI,
		timeout: 120_000
	}
});
