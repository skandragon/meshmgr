<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte';

	let meshId = $derived(parseInt($page.params.id || '0'));
	let mesh = $state<any>(null);
	let loading = $state(true);
	let activeSection = $state('general');
	let saving = $state(false);
	let error = $state('');
	let successMessage = $state('');

	// General form
	let generalForm = $state({
		name: '',
		description: ''
	});

	// LoRa form
	let loraForm = $state({
		lora_region: '',
		modem_preset: '',
		frequency_slot: 0
	});

	// LoRa configuration
	interface RegionInfo {
		code: string;
		name: string;
	}

	interface PresetInfo {
		code: string;
		name: string;
	}

	interface PresetSlots {
		LongFast: number;
		LongSlow: number;
		LongMod: number;
		MediumFast: number;
		MediumSlow: number;
		ShortFast: number;
		ShortSlow: number;
		ShortTurbo: number;
	}

	interface LoRaConfig {
		regions: RegionInfo[];
		presets: PresetInfo[];
		slots: Record<string, PresetSlots>;
	}

	let loraConfig = $state<LoRaConfig | null>(null);

	// Admin Keys
	let adminKeys = $state<any[]>([]);
	let keysLoading = $state(false);
	let showKeyModal = $state(false);
	let keyForm = $state({ public_key: '', key_name: '' });

	// Check if a preset is available for the selected region
	function isPresetAvailable(presetCode: string): boolean {
		if (!loraConfig || !loraForm.lora_region) return true;
		const maxSlot = loraConfig.slots[loraForm.lora_region]?.[presetCode as keyof PresetSlots];
		return maxSlot !== undefined && maxSlot > 0;
	}

	// Computed slot range for current selection
	const maxRadioSlot = $derived(
		loraConfig?.slots[loraForm.lora_region]?.[loraForm.modem_preset as keyof PresetSlots] ?? 319
	);
	const slotMax = $derived(maxRadioSlot + 1);
	const slotLabel = $derived.by(() => {
		if (maxRadioSlot === 0) {
			return `Frequency Slot (0 = hash default only for ${loraForm.lora_region}/${loraForm.modem_preset})`;
		} else if (maxRadioSlot === 1) {
			return `Frequency Slot (0 = hash default, 1 available for ${loraForm.lora_region}/${loraForm.modem_preset})`;
		} else {
			return `Frequency Slot (0 = hash default, 1-${maxRadioSlot + 1} available for ${loraForm.lora_region}/${loraForm.modem_preset})`;
		}
	});

	async function loadMesh() {
		try {
			mesh = await api.getMesh(meshId);
			// Populate forms with current mesh data
			generalForm = {
				name: mesh.name,
				description: mesh.description || ''
			};
			loraForm = {
				lora_region: mesh.lora_region || 'US',
				modem_preset: mesh.modem_preset || 'LongFast',
				frequency_slot: mesh.frequency_slot || 0
			};
		} catch (err: any) {
			console.error('Failed to load mesh:', err);
			error = 'Failed to load mesh';
		} finally {
			loading = false;
		}
	}

	async function loadLoRaConfig() {
		try {
			const response = await fetch('/api/lora-config');
			if (!response.ok) {
				throw new Error('Failed to fetch LoRa config');
			}
			loraConfig = await response.json();
		} catch (err) {
			console.error('Failed to load LoRa config:', err);
		}
	}

	async function handleSaveGeneral(e: Event) {
		e.preventDefault();
		saving = true;
		error = '';
		successMessage = '';
		try {
			await api.updateMesh(meshId, {
				name: generalForm.name,
				description: generalForm.description || undefined
			});
			successMessage = 'General settings saved successfully';
			await loadMesh();
		} catch (err: any) {
			error = err.message || 'Failed to save general settings';
		} finally {
			saving = false;
		}
	}

	async function handleSaveLoRa(e: Event) {
		e.preventDefault();
		saving = true;
		error = '';
		successMessage = '';
		try {
			await api.updateMesh(meshId, {
				lora_region: loraForm.lora_region,
				modem_preset: loraForm.modem_preset,
				frequency_slot: loraForm.frequency_slot
			});
			successMessage = 'LoRa settings saved successfully';
			await loadMesh();
		} catch (err: any) {
			error = err.message || 'Failed to save LoRa settings';
		} finally {
			saving = false;
		}
	}

	async function loadAdminKeys() {
		keysLoading = true;
		try {
			const result = await api.listAdminKeys(meshId);
			adminKeys = Array.isArray(result) ? result : [];
		} catch (err: any) {
			console.error('Failed to load admin keys:', err);
			adminKeys = [];
		} finally {
			keysLoading = false;
		}
	}

	async function handleCreateKey(e: Event) {
		e.preventDefault();
		error = '';
		successMessage = '';
		try {
			await api.createAdminKey(meshId, keyForm.public_key, keyForm.key_name || undefined);
			showKeyModal = false;
			keyForm = { public_key: '', key_name: '' };
			successMessage = 'Admin key added successfully';
			await loadAdminKeys();
		} catch (err: any) {
			error = err.message || 'Failed to create admin key';
		}
	}

	async function handleDeleteKey(keyId: number) {
		if (!confirm('Are you sure you want to delete this admin key?')) return;
		error = '';
		successMessage = '';
		try {
			await api.deleteAdminKey(meshId, keyId);
			successMessage = 'Admin key deleted successfully';
			await loadAdminKeys();
		} catch (err: any) {
			error = err.message || 'Failed to delete admin key';
		}
	}

	function changeSection(section: string) {
		activeSection = section;
		error = '';
		successMessage = '';
		if (section === 'keys' && adminKeys.length === 0) loadAdminKeys();
	}

	onMount(() => {
		loadMesh();
		loadLoRaConfig();
	});
</script>

<div class="min-h-screen bg-gray-50">
	<!-- Navigation -->
	<nav class="bg-white shadow-sm">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<a href="/meshes/{meshId}" class="text-sm text-blue-600 hover:text-blue-800">
						← Back to Mesh Details
					</a>
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

	{#if loading}
		<div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
			<p class="text-gray-500">Loading...</p>
		</div>
	{:else if mesh}
		<div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
			<!-- Page Header -->
			<div class="mb-6">
				<h1 class="text-3xl font-bold text-gray-900">Edit Mesh: {mesh.name}</h1>
				<p class="text-gray-600 mt-1">Update mesh settings and configuration</p>
			</div>

			<!-- Success/Error Messages -->
			{#if successMessage}
				<div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-6">
					{successMessage}
				</div>
			{/if}
			{#if error}
				<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
					{error}
				</div>
			{/if}

			<!-- Two Column Layout -->
			<div class="grid grid-cols-12 gap-6">
				<!-- Sidebar Navigation -->
				<div class="col-span-12 md:col-span-3">
					<nav class="bg-white shadow rounded-lg p-4 space-y-1">
						<button
							onclick={() => changeSection('general')}
							class="w-full text-left px-4 py-2 rounded-md text-sm font-medium {activeSection ===
							'general'
								? 'bg-blue-50 text-blue-700'
								: 'text-gray-600 hover:bg-gray-50'}"
						>
							General
						</button>
						<button
							onclick={() => changeSection('lora')}
							class="w-full text-left px-4 py-2 rounded-md text-sm font-medium {activeSection ===
							'lora'
								? 'bg-blue-50 text-blue-700'
								: 'text-gray-600 hover:bg-gray-50'}"
						>
							LoRa Configuration
						</button>
						<button
							onclick={() => changeSection('keys')}
							class="w-full text-left px-4 py-2 rounded-md text-sm font-medium {activeSection ===
							'keys'
								? 'bg-blue-50 text-blue-700'
								: 'text-gray-600 hover:bg-gray-50'}"
						>
							Admin Keys
						</button>
					</nav>
				</div>

				<!-- Main Content Area -->
				<div class="col-span-12 md:col-span-9">
					<div class="bg-white shadow rounded-lg p-6">
						<!-- General Section -->
						{#if activeSection === 'general'}
							<div>
								<h2 class="text-2xl font-bold text-gray-900 mb-4">General Settings</h2>
								<p class="text-gray-600 mb-6">
									Basic information about your mesh network
								</p>

								<form onsubmit={handleSaveGeneral}>
									<div class="mb-6">
										<label
											for="mesh-name"
											class="block text-sm font-medium text-gray-700 mb-2"
										>
											Mesh Name <span class="text-red-500">*</span>
										</label>
										<input
											id="mesh-name"
											type="text"
											bind:value={generalForm.name}
											required
											class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
											placeholder="Enter mesh name"
										/>
									</div>

									<div class="mb-6">
										<label
											for="mesh-description"
											class="block text-sm font-medium text-gray-700 mb-2"
										>
											Description
										</label>
										<textarea
											id="mesh-description"
											bind:value={generalForm.description}
											rows="4"
											class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
											placeholder="Enter mesh description (optional)"
										></textarea>
									</div>

									<div class="flex justify-end">
										<button
											type="submit"
											disabled={saving}
											class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
										>
											{saving ? 'Saving...' : 'Save Changes'}
										</button>
									</div>
								</form>
							</div>
						{/if}

						<!-- LoRa Section -->
						{#if activeSection === 'lora'}
							<div>
								<h2 class="text-2xl font-bold text-gray-900 mb-4">LoRa Configuration</h2>
								<p class="text-gray-600 mb-6">
									Configure radio settings for your mesh network
								</p>

								<form onsubmit={handleSaveLoRa}>
									<div class="mb-6">
										<label
											for="lora-region"
											class="block text-sm font-medium text-gray-700 mb-2"
										>
											LoRa Region
										</label>
										<select
											id="lora-region"
											bind:value={loraForm.lora_region}
											class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
										>
											{#if loraConfig}
												{#each loraConfig.regions as region}
													<option value={region.code}>{region.name}</option>
												{/each}
											{:else}
												<option value="US">Loading...</option>
											{/if}
										</select>
										<p class="text-sm text-gray-500 mt-1">
											Select the regulatory region for your mesh network
										</p>
									</div>

									<div class="mb-6">
										<label
											for="modem-preset"
											class="block text-sm font-medium text-gray-700 mb-2"
										>
											Modem Preset
										</label>
										<select
											id="modem-preset"
											bind:value={loraForm.modem_preset}
											class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
										>
											{#if loraConfig}
												{#each loraConfig.presets as preset}
													<option
														value={preset.code}
														disabled={!isPresetAvailable(preset.code)}
													>
														{preset.name}{!isPresetAvailable(preset.code)
															? ' (unavailable for this region)'
															: ''}
													</option>
												{/each}
											{:else}
												<option value="LongFast">Loading...</option>
											{/if}
										</select>
										<p class="text-sm text-gray-500 mt-1">
											Choose the modem preset that balances range and speed for your needs
										</p>
									</div>

									<div class="mb-6">
										<label
											for="frequency-slot"
											class="block text-sm font-medium text-gray-700 mb-2"
										>
											{slotLabel}
										</label>
										<input
											id="frequency-slot"
											type="number"
											bind:value={loraForm.frequency_slot}
											min="0"
											max={slotMax}
											class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
										/>
										<p class="text-sm text-gray-500 mt-1">
											Use 0 for automatic hash-based frequency selection, or specify a slot number
											for fixed frequency
										</p>
									</div>

									<div class="flex justify-end">
										<button
											type="submit"
											disabled={saving}
											class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
										>
											{saving ? 'Saving...' : 'Save Changes'}
										</button>
									</div>
								</form>
							</div>
						{/if}

						<!-- Admin Keys Section -->
						{#if activeSection === 'keys'}
							<div>
								<div class="flex justify-between items-center mb-4">
									<div>
										<h2 class="text-2xl font-bold text-gray-900">Admin Keys</h2>
										<p class="text-gray-600 mt-1">Manage admin keys for device administration (maximum 3)</p>
									</div>
									<button
										onclick={() => (showKeyModal = true)}
										disabled={adminKeys.length >= 3}
										class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm disabled:opacity-50 disabled:cursor-not-allowed"
									>
										Add Key
									</button>
								</div>

								{#if keysLoading}
									<p class="text-gray-500">Loading admin keys...</p>
								{:else if adminKeys.length === 0}
									<div class="text-center py-12 border-2 border-dashed border-gray-300 rounded-lg">
										<p class="text-gray-600">No admin keys configured</p>
										<p class="text-sm text-gray-500 mt-2">
											Add admin keys to enable secure device administration
										</p>
									</div>
								{:else}
									<div class="space-y-4">
										{#each adminKeys as key}
											<div class="border rounded-lg p-4">
												<div class="flex justify-between items-start">
													<div class="flex-1">
														<h3 class="font-medium text-gray-900">
															{key.key_name || 'Unnamed Key'}
														</h3>
														<p class="text-sm text-gray-500 font-mono mt-2 break-all">
															{key.public_key}
														</p>
														<p class="text-xs text-gray-400 mt-2">
															Added {new Date(key.created_at).toLocaleDateString()}
														</p>
													</div>
													<button
														onclick={() => handleDeleteKey(key.id)}
														class="text-red-600 hover:text-red-900 ml-4 text-sm"
													>
														Delete
													</button>
												</div>
											</div>
										{/each}
									</div>
								{/if}
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	{:else}
		<div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
			<p class="text-red-600">Mesh not found</p>
		</div>
	{/if}
</div>

<!-- Admin Key Modal -->
{#if showKeyModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg p-6 max-w-md w-full">
			<h3 class="text-lg font-bold mb-4">Add Admin Key</h3>
			<form onsubmit={handleCreateKey}>
				<div class="mb-4">
					<label
						for="key-public-key"
						class="block text-sm font-medium text-gray-700 mb-1"
					>
						Public Key <span class="text-red-500">*</span>
					</label>
					<textarea
						id="key-public-key"
						bind:value={keyForm.public_key}
						required
						rows="4"
						class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm focus:ring-blue-500 focus:border-blue-500"
						placeholder="Enter admin public key..."
					></textarea>
				</div>
				<div class="mb-4">
					<label for="key-name" class="block text-sm font-medium text-gray-700 mb-1">
						Key Name
					</label>
					<input
						id="key-name"
						type="text"
						bind:value={keyForm.key_name}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
						placeholder="My Laptop (optional)"
					/>
				</div>
				<div class="flex justify-end space-x-2">
					<button
						type="button"
						onclick={() => (showKeyModal = false)}
						class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700"
					>
						Add Key
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
