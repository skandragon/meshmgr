import { test, expect } from '@playwright/test';

test.describe('Mesh Management', () => {
	let userEmail: string;
	let userPassword: string;

	test.beforeEach(async ({ page }) => {
		// Register a new user for each test
		userEmail = `test-${Date.now()}@example.com`;
		userPassword = 'testpassword123';

		await page.goto('/register');
		await page.fill('input[name="displayName"]', 'Test User');
		await page.fill('input[name="email"]', userEmail);
		await page.fill('input[name="password"]', userPassword);
		await page.click('button[type="submit"]');

		// Wait for redirect to dashboard
		await expect(page).toHaveURL('/', { timeout: 10000 });
	});

	test('should create a new mesh', async ({ page }) => {
		// Click create mesh button
		await page.click('text=Create Mesh');

		// Wait for modal to appear
		await page.waitForSelector('text=Create New Mesh');

		// Fill in mesh details
		await page.locator('label:has-text("Name") + input').fill('Test Mesh');
		await page.locator('label:has-text("Description") + textarea').fill('This is a test mesh');

		// Submit form
		await page.click('button[type="submit"]:has-text("Create")');

		// Wait for modal to close and mesh to appear
		await expect(page.locator('text=Test Mesh').first()).toBeVisible({ timeout: 10000 });
	});

	test('should navigate to mesh detail page', async ({ page }) => {
		// Create a mesh first
		await page.click('text=Create Mesh');
		await page.waitForSelector('text=Create New Mesh');
		await page.locator('label:has-text("Name") + input').fill('Detail Test Mesh');
		await page.click('button[type="submit"]:has-text("Create")');

		// Click on the mesh to view details
		await page.click('text=Detail Test Mesh');

		// Should be on mesh detail page
		await expect(page).toHaveURL(/\/meshes\/\d+/);
		await expect(page.locator('h1')).toContainText('Detail Test Mesh');
	});

	test('should add a node to mesh', async ({ page }) => {
		// Create mesh
		await page.click('text=Create Mesh');
		await page.waitForSelector('text=Create New Mesh');
		await page.locator('label:has-text("Name") + input').fill('Node Test Mesh');
		await page.click('button[type="submit"]:has-text("Create")');

		// Navigate to mesh detail
		await page.click('text=Node Test Mesh');

		// Should be on Nodes tab by default
		await expect(page.locator('text=No nodes yet')).toBeVisible();

		// Click Add Node
		await page.click('text=Add Node');

		// Fill in node details
		await page.fill('input[placeholder*="Hardware ID" i], label:has-text("Hardware ID") + input', '!abc123');
		await page.fill('label:has-text("Name") + input', 'test-node');
		await page.fill('label:has-text("Long Name") + input', 'Test Node One');
		await page.fill('label:has-text("Role") + input', 'CLIENT');

		// Submit
		await page.click('button[type="submit"]:has-text("Add Node")');

		// Verify node appears in list
		await expect(page.locator('text=test-node')).toBeVisible();
		await expect(page.locator('text=!abc123')).toBeVisible();
	});

	test('should delete a node', async ({ page }) => {
		// Create mesh and add node
		await page.click('text=Create Mesh');
		await page.waitForSelector('text=Create New Mesh');
		await page.locator('label:has-text("Name") + input').fill('Delete Node Mesh');
		await page.click('button[type="submit"]:has-text("Create")');
		await page.click('text=Delete Node Mesh');

		// Add node
		await page.click('text=Add Node');
		await page.fill('label:has-text("Hardware ID") + input', '!xyz789');
		await page.fill('label:has-text("Name") + input', 'delete-me');
		await page.fill('label:has-text("Long Name") + input', 'Delete Me Node');
		await page.click('button[type="submit"]:has-text("Add Node")');

		// Verify node exists
		await expect(page.locator('text=delete-me')).toBeVisible();

		// Delete the node
		page.on('dialog', dialog => dialog.accept());
		await page.click('text=Delete >> nth=0');

		// Verify node is gone
		await expect(page.locator('text=No nodes yet')).toBeVisible();
	});

	test('should add admin keys (up to 3)', async ({ page }) => {
		// Create mesh
		await page.click('text=Create Mesh');
		await page.waitForSelector('text=Create New Mesh');
		await page.locator('label:has-text("Name") + input').fill('Keys Test Mesh');
		await page.click('button[type="submit"]:has-text("Create")');
		await page.click('text=Keys Test Mesh');

		// Switch to Admin Keys tab
		await page.click('text=Admin Keys');
		await expect(page.locator('text=No admin keys yet')).toBeVisible();

		// Add first key
		await page.click('text=Add Key');
		await page.fill('textarea', 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...');
		await page.fill('label:has-text("Key Name") + input', 'Key 1');
		await page.click('button[type="submit"]:has-text("Add Key")');
		await expect(page.locator('text=Key 1')).toBeVisible();

		// Add second key
		await page.click('text=Add Key');
		await page.fill('textarea', 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD...');
		await page.fill('label:has-text("Key Name") + input', 'Key 2');
		await page.click('button[type="submit"]:has-text("Add Key")');
		await expect(page.locator('text=Key 2')).toBeVisible();

		// Add third key
		await page.click('text=Add Key');
		await page.fill('textarea', 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQE...');
		await page.fill('label:has-text("Key Name") + input', 'Key 3');
		await page.click('button[type="submit"]:has-text("Add Key")');
		await expect(page.locator('text=Key 3')).toBeVisible();

		// Verify "Add Key" button is disabled (max 3 keys)
		await expect(page.locator('button:has-text("Add Key")')).toBeDisabled();
	});

	test('should delete an admin key', async ({ page }) => {
		// Create mesh and add a key
		await page.click('text=Create Mesh');
		await page.waitForSelector('text=Create New Mesh');
		await page.locator('label:has-text("Name") + input').fill('Delete Key Mesh');
		await page.click('button[type="submit"]:has-text("Create")');
		await page.click('text=Delete Key Mesh');

		// Switch to Admin Keys tab
		await page.click('text=Admin Keys');

		// Add a key
		await page.click('text=Add Key');
		await page.fill('textarea', 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQF...');
		await page.fill('label:has-text("Key Name") + input', 'To Delete');
		await page.click('button[type="submit"]:has-text("Add Key")');

		// Delete the key
		page.on('dialog', dialog => dialog.accept());
		await page.click('text=Delete >> nth=0');

		// Verify key is gone
		await expect(page.locator('text=No admin keys yet')).toBeVisible();
		await expect(page.locator('button:has-text("Add Key")')).not.toBeDisabled();
	});

	test('should grant and revoke mesh access', async ({ page }) => {
		// Create a second user first
		const secondUserEmail = `viewer-${Date.now()}@example.com`;

		// Open new context to create second user
		const context = page.context();
		const secondPage = await context.newPage();
		await secondPage.goto('/register');
		await secondPage.fill('input[name="displayName"]', 'Viewer User');
		await secondPage.fill('input[name="email"]', secondUserEmail);
		await secondPage.fill('input[name="password"]', 'password123');
		await secondPage.click('button[type="submit"]');
		await secondPage.close();

		// Back to main page - create mesh
		await page.click('text=Create Mesh');
		await page.waitForSelector('text=Create New Mesh');
		await page.locator('label:has-text("Name") + input').fill('Access Test Mesh');
		await page.click('button[type="submit"]:has-text("Create")');
		await page.click('text=Access Test Mesh');

		// Switch to Access tab
		await page.click('text=Access >> nth=1'); // nth=1 to avoid clicking the nav item
		await expect(page.locator('text=No shared access yet')).toBeVisible();

		// Grant access to second user
		await page.click('text=Grant Access');
		await page.fill('input[type="email"]', secondUserEmail);
		await page.selectOption('select', 'viewer');
		await page.click('button[type="submit"]:has-text("Grant Access")');

		// Verify access is granted
		await expect(page.locator(`text=${secondUserEmail}`)).toBeVisible();
		await expect(page.locator('text=viewer')).toBeVisible();

		// Revoke access
		page.on('dialog', dialog => dialog.accept());
		await page.click('text=Revoke');

		// Verify access is revoked
		await expect(page.locator('text=No shared access yet')).toBeVisible();
	});
});
