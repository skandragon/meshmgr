<!--
  ~ Copyright (C) 2025 Michael Graff
  ~
  ~ This program is free software: you can redistribute it and/or modify
  ~ it under the terms of the GNU Affero General Public License as
  ~ published by the Free Software Foundation, version 3.
  ~
  ~ This program is distributed in the hope that it will be useful,
  ~ but WITHOUT ANY WARRANTY; without even the implied warranty of
  ~ MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  ~ GNU Affero General Public License for more details.
  ~
  ~ You should have received a copy of the GNU Affero General Public License
  ~ along with this program. If not, see <http://www.gnu.org/licenses/>.
-->

<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	interface APIKey {
		id: number;
		key_name: string;
		created_at: string;
		expires_at: string | null;
		last_used_at: string | null;
	}

	let apiKeys = $state<APIKey[]>([]);
	let loading = $state(false);
	let showCreateModal = $state(false);
	let showKeyModal = $state(false);
	let keyName = $state('');
	let expiresInDays = $state<number | null>(null);
	let generatedKey = $state('');
	let error = $state('');

	async function loadAPIKeys() {
		loading = true;
		try {
			const result = await api.listAPIKeys();
			apiKeys = Array.isArray(result) ? result : [];
		} catch (err: any) {
			console.error('Failed to load API keys:', err);
			apiKeys = [];
		} finally {
			loading = false;
		}
	}

	async function handleCreateAPIKey(e: Event) {
		e.preventDefault();
		error = '';
		try {
			const expiresIn = expiresInDays ? expiresInDays * 24 * 60 * 60 : undefined;
			const result: any = await api.createAPIKey(keyName, expiresIn);
			generatedKey = result.api_key;
			showCreateModal = false;
			showKeyModal = true;
			keyName = '';
			expiresInDays = null;
			await loadAPIKeys();
		} catch (err: any) {
			error = err.message || 'Failed to create API key';
		}
	}

	async function handleDeleteAPIKey(keyId: number) {
		if (!confirm('Are you sure you want to delete this API key? This action cannot be undone.')) {
			return;
		}
		try {
			await api.deleteAPIKey(keyId);
			await loadAPIKeys();
		} catch (err: any) {
			alert(`Failed to delete API key: ${err.message}`);
		}
	}

	async function copyToClipboard() {
		try {
			await navigator.clipboard.writeText(generatedKey);
			alert('API key copied to clipboard!');
		} catch (err) {
			console.error('Failed to copy:', err);
		}
	}

	async function handleLogout() {
		await authStore.logout();
		goto('/login');
	}

	onMount(() => {
		if (!authStore.isAuthenticated) {
			goto('/login');
			return;
		}
		loadAPIKeys();
	});
</script>

<div class="min-h-screen bg-gray-50">
	<nav class="bg-white shadow-sm">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex">
					<div class="flex-shrink-0 flex items-center">
						<a href="/" class="text-xl font-bold text-gray-900">Meshtastic Manager</a>
					</div>
					<div class="hidden sm:ml-6 sm:flex sm:space-x-8">
						<a
							href="/"
							class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 hover:text-gray-900"
						>
							Meshes
						</a>
						<a
							href="/api-keys"
							class="inline-flex items-center px-1 pt-1 border-b-2 border-blue-500 text-sm font-medium text-gray-900"
						>
							API Keys
						</a>
					</div>
				</div>
				<div class="flex items-center gap-2">
					{#if authStore.isAuthenticated}
						<span class="text-sm text-gray-700 mr-2">
							{authStore.user?.display_name}
						</span>
						<button
							onclick={handleLogout}
							class="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-md text-sm font-medium"
						>
							Logout
						</button>
					{/if}
				</div>
			</div>
		</div>
	</nav>

	<div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
		<div class="px-4 py-6 sm:px-0">
			<div class="bg-white shadow rounded-lg p-6">
				<div class="flex justify-between items-center mb-6">
					<div>
						<h2 class="text-2xl font-bold text-gray-900">API Keys</h2>
						<p class="mt-1 text-sm text-gray-500">
							Manage API keys for CLI tool access
						</p>
					</div>
					<button
						onclick={() => (showCreateModal = true)}
						class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
					>
						Create API Key
					</button>
				</div>

				{#if loading}
					<p class="text-gray-500">Loading API keys...</p>
				{:else if apiKeys.length === 0}
					<p class="text-gray-600">No API keys yet. Create one to use with the CLI tool.</p>
				{:else}
					<div class="overflow-x-auto">
						<table class="min-w-full divide-y divide-gray-200">
							<thead class="bg-gray-50">
								<tr>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Name
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Key ID
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Created
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Expires
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Last Used
									</th>
									<th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
										Actions
									</th>
								</tr>
							</thead>
							<tbody class="bg-white divide-y divide-gray-200">
								{#each apiKeys as key}
									<tr>
										<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
											{key.key_name}
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-600">
											{key.id}
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
											{new Date(key.created_at).toLocaleDateString()}
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
											{key.expires_at ? new Date(key.expires_at).toLocaleDateString() : 'Never'}
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
											{key.last_used_at ? new Date(key.last_used_at).toLocaleString() : 'Never'}
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
											<button
												onclick={() => handleDeleteAPIKey(key.id)}
												class="text-red-600 hover:text-red-900"
											>
												Delete
											</button>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>

{#if showCreateModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg p-6 max-w-md w-full">
			<h3 class="text-lg font-bold mb-4">Create New API Key</h3>
			<form onsubmit={handleCreateAPIKey}>
				<div class="mb-4">
					<label for="key-name" class="block text-sm font-medium text-gray-700 mb-1">
						Name
					</label>
					<input
						id="key-name"
						type="text"
						bind:value={keyName}
						required
						placeholder="CLI Upload Key"
						class="w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
				</div>
				<div class="mb-4">
					<label for="expires-in" class="block text-sm font-medium text-gray-700 mb-1">
						Expires In (days)
					</label>
					<input
						id="expires-in"
						type="number"
						bind:value={expiresInDays}
						min="1"
						placeholder="Leave empty for no expiration"
						class="w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
					<p class="mt-1 text-xs text-gray-500">
						Leave empty for an API key that never expires
					</p>
				</div>
				{#if error}
					<p class="text-red-600 text-sm mb-4">{error}</p>
				{/if}
				<div class="flex justify-end space-x-2">
					<button
						type="button"
						onclick={() => {
							showCreateModal = false;
							error = '';
						}}
						class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700"
					>
						Create
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

{#if showKeyModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg p-6 max-w-2xl w-full">
			<h3 class="text-lg font-bold mb-4 text-red-600">⚠️ Save Your API Key</h3>
			<p class="text-sm text-gray-600 mb-4">
				This is the only time you will be able to see this API key. Save it now!
			</p>
			<div class="mb-4 p-4 bg-gray-50 rounded-md border border-gray-300 font-mono text-sm break-all">
				{generatedKey}
			</div>
			<div class="bg-blue-50 border-l-4 border-blue-400 p-4 mb-4">
				<p class="text-sm text-blue-700">
					<strong>Usage:</strong> Set this as an environment variable or use with the CLI tool:
				</p>
				<pre class="mt-2 text-xs bg-gray-800 text-green-400 p-2 rounded overflow-x-auto">export MESHMANAGER_API_KEY="{generatedKey}"</pre>
			</div>
			<div class="flex justify-end space-x-2">
				<button
					onclick={copyToClipboard}
					class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300"
				>
					Copy to Clipboard
				</button>
				<button
					onclick={() => {
						showKeyModal = false;
						generatedKey = '';
					}}
					class="px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700"
				>
					I've Saved It
				</button>
			</div>
		</div>
	</div>
{/if}
