<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';

	let meshes = $state<any[]>([]);
	let loading = $state(false);
	let showCreateModal = $state(false);
	let meshName = $state('');
	let meshDescription = $state('');
	let loraRegion = $state('US');
	let modemPreset = $state('LONG_FAST');
	let frequencySlot = $state(0);
	let error = $state('');
	let initialLoad = $state(true);

	async function loadMeshes() {
		loading = true;
		try {
			const result = await api.listMeshes();
			meshes = Array.isArray(result) ? result : [];
		} catch (err: any) {
			console.error('Failed to load meshes:', err);
			meshes = [];
		} finally {
			loading = false;
		}
	}

	async function handleCreateMesh(e: Event) {
		e.preventDefault();
		error = '';
		try {
			await api.createMesh(
				meshName,
				meshDescription || undefined,
				loraRegion,
				modemPreset,
				frequencySlot
			);
			showCreateModal = false;
			meshName = '';
			meshDescription = '';
			loraRegion = 'US';
			modemPreset = 'LONG_FAST';
			frequencySlot = 0;
			await loadMeshes();
		} catch (err: any) {
			error = err.message || 'Failed to create mesh';
		}
	}

	async function handleLogout() {
		await authStore.logout();
		goto('/login');
	}

	// Load meshes when auth is ready
	$effect(() => {
		if (authStore.isAuthenticated && !authStore.loading && initialLoad) {
			initialLoad = false;
			loadMeshes();
		}
	});
</script>

<div class="min-h-screen bg-gray-50">
	<nav class="bg-white shadow-sm">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex">
					<div class="flex-shrink-0 flex items-center">
						<h1 class="text-xl font-bold text-gray-900">Meshtastic Manager</h1>
					</div>
				</div>
				<div class="flex items-center">
					{#if authStore.isAuthenticated}
						<span class="text-sm text-gray-700 mr-4">
							{authStore.user?.display_name}
						</span>
						<button
							onclick={handleLogout}
							class="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-md text-sm font-medium"
						>
							Logout
						</button>
					{:else}
						<a
							href="/login"
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
						>
							Sign in
						</a>
					{/if}
				</div>
			</div>
		</div>
	</nav>

	<div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
		{#if authStore.loading}
			<div class="text-center py-12">
				<p class="text-gray-500">Loading...</p>
			</div>
		{:else if authStore.isAuthenticated}
			<div class="px-4 py-6 sm:px-0">
				<div class="bg-white shadow rounded-lg p-6">
					<div class="flex justify-between items-center mb-6">
						<h2 class="text-2xl font-bold text-gray-900">My Meshes</h2>
						<button
							onclick={() => (showCreateModal = true)}
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
						>
							Create Mesh
						</button>
					</div>

					{#if loading}
						<p class="text-gray-500">Loading meshes...</p>
					{:else if meshes.length === 0}
						<p class="text-gray-600">No meshes yet. Create one to get started!</p>
					{:else}
						<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
							{#each meshes as mesh}
								<a
									href="/meshes/{mesh.id}"
									class="block p-4 border rounded-lg hover:shadow-md transition-shadow"
								>
									<h3 class="text-lg font-semibold text-gray-900">{mesh.name}</h3>
									{#if mesh.description}
										<p class="text-sm text-gray-600 mt-1">{mesh.description}</p>
									{/if}
									<p class="text-xs text-gray-400 mt-2">
										Created {new Date(mesh.created_at).toLocaleDateString()}
									</p>
								</a>
							{/each}
						</div>
					{/if}
				</div>
			</div>

			{#if showCreateModal}
				<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
					<div class="bg-white rounded-lg p-6 max-w-md w-full">
						<h3 class="text-lg font-bold mb-4">Create New Mesh</h3>
						<form onsubmit={handleCreateMesh}>
							<div class="mb-4">
								<label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
								<input
									type="text"
									bind:value={meshName}
									required
									class="w-full px-3 py-2 border border-gray-300 rounded-md"
								/>
							</div>
							<div class="mb-4">
								<label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
								<textarea
									bind:value={meshDescription}
									class="w-full px-3 py-2 border border-gray-300 rounded-md"
									rows="3"
								></textarea>
							</div>
							<div class="mb-4">
								<label class="block text-sm font-medium text-gray-700 mb-1">LoRa Region</label>
								<select
									bind:value={loraRegion}
									class="w-full px-3 py-2 border border-gray-300 rounded-md"
								>
									<option value="US">US</option>
									<option value="EU_868">EU 868 MHz</option>
									<option value="EU_433">EU 433 MHz</option>
									<option value="CN">China</option>
									<option value="JP">Japan</option>
									<option value="ANZ">Australia/NZ</option>
									<option value="KR">Korea</option>
									<option value="TW">Taiwan</option>
									<option value="RU">Russia</option>
									<option value="IN">India</option>
									<option value="TH">Thailand</option>
								</select>
							</div>
							<div class="mb-4">
								<label class="block text-sm font-medium text-gray-700 mb-1">Modem Preset</label>
								<select
									bind:value={modemPreset}
									class="w-full px-3 py-2 border border-gray-300 rounded-md"
								>
									<option value="SHORT_TURBO">Short Range / Turbo</option>
									<option value="SHORT_FAST">Short Range / Fast</option>
									<option value="SHORT_SLOW">Short Range / Slow</option>
									<option value="MEDIUM_FAST">Medium Range / Fast</option>
									<option value="MEDIUM_SLOW">Medium Range / Slow</option>
									<option value="LONG_FAST">Long Range / Fast</option>
									<option value="LONG_MODERATE">Long Range / Moderate</option>
									<option value="LONG_SLOW">Long Range / Slow</option>
									<option value="VERY_LONG_SLOW">Very Long Range / Slow</option>
								</select>
							</div>
							<div class="mb-4">
								<label class="block text-sm font-medium text-gray-700 mb-1">Frequency Slot (0-7)</label>
								<input
									type="number"
									bind:value={frequencySlot}
									min="0"
									max="7"
									class="w-full px-3 py-2 border border-gray-300 rounded-md"
								/>
							</div>
							{#if error}
								<p class="text-red-600 text-sm mb-4">{error}</p>
							{/if}
							<div class="flex justify-end space-x-2">
								<button
									type="button"
									onclick={() => (showCreateModal = false)}
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
		{:else}
			<div class="px-4 py-6 sm:px-0">
				<div class="bg-white shadow rounded-lg p-6 text-center">
					<h2 class="text-2xl font-bold text-gray-900 mb-4">Welcome to Meshtastic Manager</h2>
					<p class="text-gray-600 mb-6">
						Manage your Meshtastic nodes with ease. Sign in to get started.
					</p>
					<div class="space-x-4">
						<a
							href="/login"
							class="inline-block bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-md font-medium"
						>
							Sign in
						</a>
						<a
							href="/register"
							class="inline-block bg-gray-200 hover:bg-gray-300 text-gray-800 px-6 py-3 rounded-md font-medium"
						>
							Create account
						</a>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
