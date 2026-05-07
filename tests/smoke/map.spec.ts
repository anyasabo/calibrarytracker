import { expect, test } from '@playwright/test';

const PAGES_BASE_PATH = '/calibrarytracker';

test('map route renders markers without runtime errors', async ({ page }) => {
	const runtimeErrors: string[] = [];

	page.on('pageerror', (error) => runtimeErrors.push(`pageerror: ${error.message}`));
	page.on('console', (message) => {
		if (message.type() === 'error') {
			runtimeErrors.push(`console: ${message.text()}`);
		}
	});

	await page.goto(`${PAGES_BASE_PATH}/`);
	await page.locator('header nav a[href="/calibrarytracker/map"]').click();

	await expect(page.locator('.map-container')).toBeVisible();
	await expect.poll(async () => await page.locator('.marker-cluster').count()).toBeGreaterThan(0);
	expect(runtimeErrors).toEqual([]);
});

test('navbar links are generated with the Pages project base path', async ({ page }) => {
	await page.goto(`${PAGES_BASE_PATH}/`);
	await expect(page.locator('header nav a[href="/calibrarytracker/systems"]')).toBeVisible();
	await expect(page.locator('header nav a[href="/calibrarytracker/map"]')).toBeVisible();
	await expect(page.locator('header nav a[href="/calibrarytracker/cards"]')).toBeVisible();
	await expect(page.locator('header nav a[href="/calibrarytracker/visits"]')).toBeVisible();
});
