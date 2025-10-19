/*
 * Copyright (C) 2025 Michael Graff
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

import { test, expect } from '@playwright/test';

test.describe('Mesh Management', () => {
	let userEmail: string;
	let userPassword: string;

	async function createMesh(page: any, name: string, description?: string, loraRegion?: string, modemPreset?: string, frequencySlot?: number) {
		await page.getByRole('button', { name: 'Create Mesh' }).click();
		await expect(page.getByRole('heading', { name: 'Create New Mesh' })).toBeVisible();

		const modal = page.locator('div.fixed').filter({ hasText: 'Create New Mesh' });
		await modal.locator('label:has-text("Name") + input').fill(name);
		if (description) {
			await modal.locator('textarea').first().fill(description);
		}
		if (loraRegion) {
			await modal.locator('label:has-text("LoRa Region") + select').selectOption(loraRegion);
		}
		if (modemPreset) {
			await modal.locator('label:has-text("Modem Preset") + select').selectOption(modemPreset);
		}
		if (frequencySlot !== undefined) {
			await modal.locator('label:has-text("Frequency Slot") + input').fill(frequencySlot.toString());
		}

		await page.getByRole('button', { name: 'Create', exact: true }).click();
		await expect(page.getByText(name).first()).toBeVisible({ timeout: 10000 });
	}

	test.beforeEach(async ({ page }) => {
		// Register a new user for each test with unique timestamp + random string
		userEmail = `test-${Date.now()}-${Math.random().toString(36).substring(7)}@example.com`;
		userPassword = 'testpassword123';

		await page.goto('/register');
		await page.fill('input[name="displayName"]', 'Test User');
		await page.fill('input[name="email"]', userEmail);
		await page.fill('input[name="password"]', userPassword);

		// Submit and wait for redirect
		await page.click('button[type="submit"]');
		await page.waitForURL('/', { timeout: 10000 });
	});

	test('should create a new mesh', async ({ page }) => {
		// Click create mesh button
		await page.getByRole('button', { name: 'Create Mesh' }).click();

		// Wait for modal to appear
		await expect(page.getByRole('heading', { name: 'Create New Mesh' })).toBeVisible();

		// Fill in mesh details
		const nameInput = page.locator('input[type="text"]').first();
		await nameInput.fill('Test Mesh');
		const descTextarea = page.locator('textarea').first();
		await descTextarea.fill('This is a test mesh');

		// Submit form
		await page.getByRole('button', { name: 'Create', exact: true }).click();

		// Wait for modal to close and mesh to appear
		await expect(page.getByText('Test Mesh').first()).toBeVisible({ timeout: 10000 });
	});

	test('should navigate to mesh detail page', async ({ page }) => {
		// Create a mesh first
		await createMesh(page, 'Detail Test Mesh');

		// Click on the mesh to view details
		await page.getByText('Detail Test Mesh').first().click();

		// Should be on mesh detail page
		await expect(page).toHaveURL(/\/meshes\/\d+/);
		await expect(page.locator('h1')).toContainText('Detail Test Mesh');
	});

	test('should add a node to mesh', async ({ page }) => {
		// Create mesh
		await createMesh(page, 'Node Test Mesh');

		// Navigate to mesh detail
		await page.getByText('Node Test Mesh').first().click();

		// Should be on Nodes tab by default
		await expect(page.getByText('No nodes yet.')).toBeVisible({ timeout: 10000 });

		// Click Add Node
		await page.getByRole('button', { name: 'Add Node' }).click();

		// Wait for modal to appear and fill in node details
		await expect(page.getByRole('heading', { name: 'Add New Node' })).toBeVisible();
		const modal = page.locator('div.fixed').filter({ hasText: 'Add New Node' });
		await modal.locator('label:has-text("Hardware ID") + input').fill('!abc123');
		await modal.locator('label:has-text("Name") + input').first().fill('test-node');
		await modal.locator('label:has-text("Long Name") + input').fill('Test Node One');
		await modal.locator('label:has-text("Role") + select').selectOption('CLIENT');

		// Submit
		await modal.getByRole('button', { name: 'Add Node' }).click();

		// Verify node appears in list
		await expect(page.getByText('test-node')).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('!abc123')).toBeVisible();
	});

	test('should delete a node', async ({ page }) => {
		// Create mesh and navigate to it
		await createMesh(page, 'Delete Node Mesh');
		await page.getByText('Delete Node Mesh').first().click();

		// Add node
		await page.getByRole('button', { name: 'Add Node' }).click();
		await expect(page.getByRole('heading', { name: 'Add New Node' })).toBeVisible();
		const modal = page.locator('div.fixed').filter({ hasText: 'Add New Node' });
		await modal.locator('label:has-text("Hardware ID") + input').fill('!xyz789');
		await modal.locator('label:has-text("Name") + input').first().fill('delete-me');
		await modal.locator('label:has-text("Long Name") + input').fill('Delete Me Node');
		await modal.getByRole('button', { name: 'Add Node' }).click();

		// Verify node exists
		await expect(page.getByText('delete-me')).toBeVisible({ timeout: 10000 });

		// Delete the node
		page.on('dialog', dialog => dialog.accept());
		await page.getByRole('button', { name: 'Delete' }).first().click();

		// Verify node is gone
		await expect(page.getByText('No nodes yet.')).toBeVisible({ timeout: 10000 });
	});

	test('should add admin keys (up to 3)', async ({ page }) => {
		// Create mesh
		await createMesh(page, 'Keys Test Mesh');
		await page.getByText('Keys Test Mesh').first().click();

		// Switch to Admin Keys tab
		await page.getByRole('button', { name: 'Admin Keys' }).click();
		await expect(page.getByText('No admin keys yet.')).toBeVisible({ timeout: 10000 });

		// Add first key
		await page.getByRole('button', { name: 'Add Key' }).click();
		await expect(page.getByRole('heading', { name: 'Add Admin Key' })).toBeVisible();
		let modal = page.locator('div.fixed').filter({ hasText: 'Add Admin Key' });
		await modal.locator('textarea').fill('ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...');
		await modal.locator('label:has-text("Key Name") + input').fill('Key 1');
		await modal.getByRole('button', { name: 'Add Key' }).click();
		await expect(page.getByText('Key 1')).toBeVisible({ timeout: 10000 });

		// Add second key
		await page.getByRole('button', { name: 'Add Key' }).click();
		await expect(page.getByRole('heading', { name: 'Add Admin Key' })).toBeVisible();
		modal = page.locator('div.fixed').filter({ hasText: 'Add Admin Key' });
		await modal.locator('textarea').fill('ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD...');
		await modal.locator('label:has-text("Key Name") + input').fill('Key 2');
		await modal.getByRole('button', { name: 'Add Key' }).click();
		await expect(page.getByText('Key 2')).toBeVisible({ timeout: 10000 });

		// Add third key
		await page.getByRole('button', { name: 'Add Key' }).click();
		await expect(page.getByRole('heading', { name: 'Add Admin Key' })).toBeVisible();
		modal = page.locator('div.fixed').filter({ hasText: 'Add Admin Key' });
		await modal.locator('textarea').fill('ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQE...');
		await modal.locator('label:has-text("Key Name") + input').fill('Key 3');
		await modal.getByRole('button', { name: 'Add Key' }).click();
		await expect(page.getByText('Key 3')).toBeVisible({ timeout: 10000 });

		// Verify "Add Key" button is disabled (max 3 keys)
		await expect(page.getByRole('button', { name: 'Add Key' })).toBeDisabled();
	});

	test('should delete an admin key', async ({ page }) => {
		// Create mesh and navigate to it
		await createMesh(page, 'Delete Key Mesh');
		await page.getByText('Delete Key Mesh').first().click();

		// Switch to Admin Keys tab
		await page.getByRole('button', { name: 'Admin Keys' }).click();

		// Add a key
		await page.getByRole('button', { name: 'Add Key' }).click();
		await expect(page.getByRole('heading', { name: 'Add Admin Key' })).toBeVisible();
		const modal = page.locator('div.fixed').filter({ hasText: 'Add Admin Key' });
		await modal.locator('textarea').fill('ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQF...');
		await modal.locator('label:has-text("Key Name") + input').fill('To Delete');
		await modal.getByRole('button', { name: 'Add Key' }).click();
		await expect(page.getByText('To Delete')).toBeVisible({ timeout: 10000 });

		// Delete the key
		page.on('dialog', dialog => dialog.accept());
		await page.getByRole('button', { name: 'Delete' }).first().click();

		// Verify key is gone
		await expect(page.getByText('No admin keys yet.')).toBeVisible({ timeout: 10000 });
		await expect(page.getByRole('button', { name: 'Add Key' })).not.toBeDisabled();
	});

	test('should create mesh with LoRa configuration', async ({ page }) => {
		// Click create mesh button
		await page.getByRole('button', { name: 'Create Mesh' }).click();
		await expect(page.getByRole('heading', { name: 'Create New Mesh' })).toBeVisible();

		// Verify LoRa config fields exist
		const modal = page.locator('div.fixed').filter({ hasText: 'Create New Mesh' });
		await expect(modal.locator('label:has-text("LoRa Region")')).toBeVisible();
		await expect(modal.locator('label:has-text("Modem Preset")')).toBeVisible();
		await expect(modal.locator('label:has-text("Frequency Slot")')).toBeVisible();

		// Fill in mesh with LoRa config - use US/ShortFast (defaults, always available)
		await modal.locator('label:has-text("Name") + input').fill('LoRa Test Mesh');

		// Wait a moment for LoRa config to load (simple timeout approach)
		await page.waitForTimeout(1000);

		const regionSelect = modal.locator('label:has-text("LoRa Region") + select');
		await regionSelect.selectOption('US');

		const presetSelect = modal.locator('label:has-text("Modem Preset") + select');
		await presetSelect.selectOption('ShortFast');

		await modal.locator('label:has-text("Frequency Slot") + input').fill('5');

		// Submit form
		await page.getByRole('button', { name: 'Create', exact: true }).click();

		// Verify mesh appears
		await expect(page.getByText('LoRa Test Mesh')).toBeVisible({ timeout: 10000 });
	});

	test('should add node with unmessageable flag', async ({ page }) => {
		// Create mesh and navigate to it
		await createMesh(page, 'Unmessageable Test Mesh');
		await page.getByText('Unmessageable Test Mesh').first().click();

		// Add node with unmessageable flag
		await page.getByRole('button', { name: 'Add Node' }).click();
		await expect(page.getByRole('heading', { name: 'Add New Node' })).toBeVisible();

		const modal = page.locator('div.fixed').filter({ hasText: 'Add New Node' });
		await modal.locator('label:has-text("Hardware ID") + input').fill('!test123');
		await modal.locator('label:has-text("Name") + input').first().fill('silent-node');
		await modal.locator('label:has-text("Long Name") + input').fill('Silent Node');

		// Check unmessageable checkbox
		await modal.locator('label:has-text("Unmessageable") input[type="checkbox"]').check();

		// Submit
		await modal.getByRole('button', { name: 'Add Node' }).click();

		// Verify node appears
		await expect(page.getByText('silent-node')).toBeVisible({ timeout: 10000 });

		// Expand the node to verify unmessageable state
		const expandButton = page.locator('button[aria-label="Toggle details"]').first();
		await expandButton.click();
		await expect(page.getByText('Configuration State')).toBeVisible();
		await expect(page.getByText('DESIRED STATE')).toBeVisible();
	});

	test('should expand node to show state comparison', async ({ page }) => {
		// Create mesh and add a node
		await createMesh(page, 'State Expansion Test');
		await page.getByText('State Expansion Test').first().click();

		// Add node
		await page.getByRole('button', { name: 'Add Node' }).click();
		await expect(page.getByRole('heading', { name: 'Add New Node' })).toBeVisible();
		const modal = page.locator('div.fixed').filter({ hasText: 'Add New Node' });
		await modal.locator('label:has-text("Hardware ID") + input').fill('!expand123');
		await modal.locator('label:has-text("Name") + input').first().fill('expand-test');
		await modal.locator('label:has-text("Long Name") + input').fill('Expansion Test Node');
		await modal.locator('label:has-text("Role") + select').selectOption('CLIENT');
		await modal.getByRole('button', { name: 'Add Node' }).click();

		// Wait for node to appear
		await expect(page.getByText('expand-test')).toBeVisible({ timeout: 10000 });

		// Find and click the expand button (▶ arrow)
		const expandButton = page.locator('button[aria-label="Toggle details"]').first();
		await expandButton.click();

		// Verify state comparison section appears
		await expect(page.getByText('Configuration State')).toBeVisible();
		await expect(page.getByText('DESIRED STATE')).toBeVisible();
		await expect(page.getByText('APPLIED STATE')).toBeVisible();
		await expect(page.getByText('Never applied to device')).toBeVisible();

		// Verify the expansion arrow changed to ▼
		await expect(expandButton).toContainText('▼');

		// Click again to collapse
		await expandButton.click();
		await expect(page.getByText('Configuration State')).not.toBeVisible();
	});
});
