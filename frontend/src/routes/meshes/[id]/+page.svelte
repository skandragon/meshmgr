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
	import { page } from '$app/stores';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte';

	let meshId = $derived(parseInt($page.params.id || '0'));
	let mesh = $state<any>(null);
	let loading = $state(true);

	// Nodes state
	let nodes = $state<any[]>([]);
	let nodesLoading = $state(false);
	let showNodeModal = $state(false);
	let expandedNodes = $state<Set<number>>(new Set());
	let nodeForm = $state({
		hardware_id: '',
		name: '',
		long_name: '',
		role: '',
		status: '',
		unmessageable: false
	});

	function toggleNodeExpansion(nodeId: number) {
		const newSet = new Set(expandedNodes);
		if (newSet.has(nodeId)) {
			newSet.delete(nodeId);
		} else {
			newSet.add(nodeId);
		}
		expandedNodes = newSet;
	}

	let error = $state('');

	async function loadMesh() {
		try {
			mesh = await api.getMesh(meshId);
		} catch (err: any) {
			console.error('Failed to load mesh:', err);
			error = 'Failed to load mesh';
		} finally {
			loading = false;
		}
	}

	async function loadNodes() {
		nodesLoading = true;
		try {
			const result = await api.listNodes(meshId);
			nodes = Array.isArray(result) ? result : [];
		} catch (err: any) {
			console.error('Failed to load nodes:', err);
			nodes = [];
		} finally {
			nodesLoading = false;
		}
	}

	async function handleCreateNode(e: Event) {
		e.preventDefault();
		error = '';
		try {
			await api.createNode(meshId, nodeForm);
			showNodeModal = false;
			nodeForm = { hardware_id: '', name: '', long_name: '', role: '', status: '', unmessageable: false };
			await loadNodes();
		} catch (err: any) {
			error = err.message || 'Failed to create node';
		}
	}

	async function handleDeleteNode(nodeId: number) {
		if (!confirm('Are you sure you want to delete this node?')) return;
		try {
			await api.deleteNode(meshId, nodeId);
			await loadNodes();
		} catch (err: any) {
			error = err.message || 'Failed to delete node';
		}
	}

	onMount(() => {
		loadMesh();
		loadNodes();
	});
</script>

<div class="min-h-screen bg-gray-50">
	<nav class="bg-white shadow-sm">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<a href="/" class="text-sm text-blue-600 hover:text-blue-800">← Back to Meshes</a>
				</div>
				<div class="flex items-center">
					{#if authStore.isAuthenticated}
						<span class="text-sm text-gray-700 mr-4">
							{authStore.user?.display_name}
						</span>
					{/if}
				</div>
			</div>
		</div>
	</nav>

	<div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
		{#if loading}
			<p class="text-gray-500">Loading...</p>
		{:else if mesh}
			<div class="px-4 py-6 sm:px-0">
				<div class="bg-white shadow rounded-lg">
					<div class="p-6 border-b">
						<div class="flex justify-between items-start">
							<div>
								<h1 class="text-3xl font-bold text-gray-900">{mesh.name}</h1>
								{#if mesh.description}
									<p class="text-gray-600 mt-2">{mesh.description}</p>
								{/if}
							</div>
							<a
								href="/meshes/{meshId}/edit"
								class="px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 text-sm font-medium"
							>
								Edit Mesh
							</a>
						</div>
					</div>

					<div class="p-6">
						{#if error}
							<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
								{error}
							</div>
						{/if}
							<div class="flex justify-between items-center mb-4">
								<h2 class="text-xl font-semibold">Nodes</h2>
								<button
									onclick={() => (showNodeModal = true)}
									class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm"
								>
									Add Node
								</button>
							</div>

							{#if nodesLoading}
								<p class="text-gray-500">Loading nodes...</p>
							{:else if nodes.length === 0}
								<p class="text-gray-600">No nodes yet.</p>
							{:else}
								<div class="overflow-x-auto">
									<table class="min-w-full divide-y divide-gray-200">
										<thead class="bg-gray-50">
											<tr>
												<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase w-8"></th>
												<th
													class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
													>Name</th
												>
												<th
													class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
													>Hardware ID</th
												>
												<th
													class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
													>Role</th
												>
												<th
													class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
													>Status</th
												>
												<th
													class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
													>Flags</th
												>
												<th class="px-6 py-3 text-right"></th>
											</tr>
										</thead>
										<tbody class="bg-white divide-y divide-gray-200">
											{#each nodes as node}
												<tr class="{node.pending_changes ? 'bg-yellow-50' : ''}">
													<td class="px-2 py-4 whitespace-nowrap text-sm text-gray-500">
														<button
															onclick={() => toggleNodeExpansion(node.id)}
															class="text-gray-600 hover:text-gray-900 p-1"
															aria-label="Toggle details"
														>
															{#if expandedNodes.has(node.id)}
																▼
															{:else}
																▶
															{/if}
														</button>
													</td>
													<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
														{node.name}
														<div class="text-xs text-gray-500">{node.long_name}</div>
														{#if node.pending_changes}
															<div class="text-xs text-yellow-600 font-medium">⚠ Pending changes</div>
														{/if}
													</td>
													<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
														{node.hardware_id}
													</td>
													<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
														{node.role || '-'}
													</td>
													<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
														{node.status || '-'}
													</td>
													<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
														{#if node.unmessageable}
															<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
																🔇 Unmessageable
															</span>
														{:else}
															<span class="text-gray-400">-</span>
														{/if}
													</td>
													<td class="px-6 py-4 whitespace-nowrap text-right text-sm">
														<button
															onclick={() => handleDeleteNode(node.id)}
															class="text-red-600 hover:text-red-900"
														>
															Delete
														</button>
													</td>
												</tr>
												{#if expandedNodes.has(node.id)}
													<tr class="{node.pending_changes ? 'bg-yellow-50' : 'bg-gray-50'}">
														<td colspan="7" class="px-6 py-4">
															<div class="text-sm">
																<h4 class="font-semibold mb-3 text-gray-900">Configuration State</h4>
																<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
																	<div>
																		<div class="text-xs font-medium text-gray-500 mb-2">DESIRED STATE</div>
																		<div class="space-y-2 text-xs">
																			<div class="flex justify-between">
																				<span class="font-medium">Name:</span>
																				<span class="{node.name !== node.applied_name && node.applied_name ? 'text-orange-600 font-semibold' : ''}">{node.name}</span>
																			</div>
																			<div class="flex justify-between">
																				<span class="font-medium">Long Name:</span>
																				<span class="{node.long_name !== node.applied_long_name && node.applied_long_name ? 'text-orange-600 font-semibold' : ''}">{node.long_name}</span>
																			</div>
																			<div class="flex justify-between">
																				<span class="font-medium">Role:</span>
																				<span class="{node.role !== node.applied_role && node.applied_role ? 'text-orange-600 font-semibold' : ''}">{node.role || '-'}</span>
																			</div>
																			<div class="flex justify-between">
																				<span class="font-medium">Unmessageable:</span>
																				<span class="{node.unmessageable !== node.applied_unmessageable && node.applied_unmessageable !== null ? 'text-orange-600 font-semibold' : ''}">{node.unmessageable ? 'Yes' : 'No'}</span>
																			</div>
																			{#if node.public_key}
																				<div class="flex justify-between">
																					<span class="font-medium">Public Key:</span>
																					<span class="font-mono {node.public_key !== node.applied_public_key && node.applied_public_key ? 'text-orange-600 font-semibold' : ''}">{node.public_key.substring(0, 12)}...</span>
																				</div>
																			{/if}
																		</div>
																	</div>
																	<div>
																		<div class="text-xs font-medium text-gray-500 mb-2">APPLIED STATE</div>
																		<div class="space-y-2 text-xs">
																			<div class="flex justify-between">
																				<span class="font-medium">Name:</span>
																				<span>{node.applied_name || '-'}</span>
																			</div>
																			<div class="flex justify-between">
																				<span class="font-medium">Long Name:</span>
																				<span>{node.applied_long_name || '-'}</span>
																			</div>
																			<div class="flex justify-between">
																				<span class="font-medium">Role:</span>
																				<span>{node.applied_role || '-'}</span>
																			</div>
																			<div class="flex justify-between">
																				<span class="font-medium">Unmessageable:</span>
																				<span>{node.applied_unmessageable !== null ? (node.applied_unmessageable ? 'Yes' : 'No') : '-'}</span>
																			</div>
																			{#if node.public_key}
																				<div class="flex justify-between">
																					<span class="font-medium">Public Key:</span>
																					<span class="font-mono">{node.applied_public_key ? node.applied_public_key.substring(0, 12) + '...' : '-'}</span>
																				</div>
																			{/if}
																		</div>
																		{#if node.config_applied_at}
																			<div class="mt-3 pt-3 border-t border-gray-200 text-xs text-gray-600">
																				Last applied: {new Date(node.config_applied_at).toLocaleString()}
																			</div>
																		{:else}
																			<div class="mt-3 pt-3 border-t border-gray-200 text-xs text-gray-500 italic">
																				Never applied to device
																			</div>
																		{/if}
																	</div>
																</div>
															</div>
														</td>
													</tr>
												{/if}
											{/each}
										</tbody>
									</table>
								</div>
							{/if}
					</div>
				</div>
			</div>
		{:else}
			<p class="text-red-600">Mesh not found</p>
		{/if}
	</div>
</div>

<!-- Node Modal -->
{#if showNodeModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
		<div class="bg-white rounded-lg p-6 max-w-md w-full">
			<h3 class="text-lg font-bold mb-4">Add New Node</h3>
			<form onsubmit={handleCreateNode}>
				<div class="mb-4">
					<label for="node-hardware-id" class="block text-sm font-medium text-gray-700 mb-1">Hardware ID</label>
					<input
						id="node-hardware-id"
						type="text"
						bind:value={nodeForm.hardware_id}
						required
						class="w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
				</div>
				<div class="mb-4">
					<label for="node-name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
					<input
						id="node-name"
						type="text"
						bind:value={nodeForm.name}
						required
						class="w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
				</div>
				<div class="mb-4">
					<label for="node-long-name" class="block text-sm font-medium text-gray-700 mb-1">Long Name</label>
					<input
						id="node-long-name"
						type="text"
						bind:value={nodeForm.long_name}
						required
						class="w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
				</div>
				<div class="mb-4">
					<label for="node-role" class="block text-sm font-medium text-gray-700 mb-1">Role (optional)</label>
					<select
						id="node-role"
						bind:value={nodeForm.role}
						class="w-full px-3 py-2 border border-gray-300 rounded-md"
					>
						<option value="">Select a role...</option>
						<option value="CLIENT">CLIENT - App connected or messaging device</option>
						<option value="CLIENT_MUTE">CLIENT_MUTE - Does not forward packets</option>
						<option value="CLIENT_HIDDEN">CLIENT_HIDDEN - Broadcasts only as needed</option>
						<option value="CLIENT_BASE">CLIENT_BASE - Personal base station</option>
						<option value="TRACKER">TRACKER - GPS position priority</option>
						<option value="LOST_AND_FOUND">LOST_AND_FOUND - Device recovery broadcasts</option>
						<option value="SENSOR">SENSOR - Telemetry priority</option>
						<option value="TAK">TAK - ATAK system optimized</option>
						<option value="TAK_TRACKER">TAK_TRACKER - TAK PLI broadcasts</option>
						<option value="REPEATER">REPEATER - Network coverage extension</option>
						<option value="ROUTER">ROUTER - Infrastructure node</option>
						<option value="ROUTER_LATE">ROUTER_LATE - Delayed rebroadcast router</option>
					</select>
				</div>
				<div class="mb-4">
					<label for="node-status" class="block text-sm font-medium text-gray-700 mb-1">Status (optional)</label>
					<input
						id="node-status"
						type="text"
						bind:value={nodeForm.status}
						class="w-full px-3 py-2 border border-gray-300 rounded-md"
					/>
				</div>
				<div class="mb-4">
					<label class="flex items-center">
						<input
							type="checkbox"
							bind:checked={nodeForm.unmessageable}
							class="mr-2 h-4 w-4 text-blue-600 border-gray-300 rounded"
						/>
						<span class="text-sm font-medium text-gray-700">Unmessageable (cannot receive direct messages)</span>
					</label>
				</div>
				{#if error}
					<p class="text-red-600 text-sm mb-4">{error}</p>
				{/if}
				<div class="flex justify-end space-x-2">
					<button
						type="button"
						onclick={() => (showNodeModal = false)}
						class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700"
					>
						Add Node
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
