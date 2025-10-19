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

	// Admin Keys - 3 fixed slots
	let keySlots = $state<Array<{ id: number | null; public_key: string; key_name: string; saving: boolean }>>([
		{ id: null, public_key: '', key_name: '', saving: false },
		{ id: null, public_key: '', key_name: '', saving: false },
		{ id: null, public_key: '', key_name: '', saving: false }
	]);
	let keysLoading = $state(false);

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
			const keys = Array.isArray(result) ? result : [];
			// Map keys to slots (max 3)
			keySlots = [
				keys[0] ? { id: keys[0].id, public_key: keys[0].public_key, key_name: keys[0].key_name || '', saving: false } : { id: null, public_key: '', key_name: '', saving: false },
				keys[1] ? { id: keys[1].id, public_key: keys[1].public_key, key_name: keys[1].key_name || '', saving: false } : { id: null, public_key: '', key_name: '', saving: false },
				keys[2] ? { id: keys[2].id, public_key: keys[2].public_key, key_name: keys[2].key_name || '', saving: false } : { id: null, public_key: '', key_name: '', saving: false }
			];
		} catch (err: any) {
			console.error('Failed to load admin keys:', err);
		} finally {
			keysLoading = false;
		}
	}

	async function handleSaveKey(index: number) {
		const slot = keySlots[index];
		const keyLabel = index === 0 ? 'Primary' : index === 1 ? 'Secondary' : 'Tertiary';

		keySlots[index].saving = true;
		error = '';
		successMessage = '';

		try {
			// If key is empty and slot has an ID, delete it
			if (!slot.public_key.trim() && slot.id) {
				await api.deleteAdminKey(meshId, slot.id);
				keySlots[index] = { id: null, public_key: '', key_name: '', saving: false };
				successMessage = `${keyLabel} key removed successfully`;
				return;
			}

			// If key is empty and no ID, just clear the slot
			if (!slot.public_key.trim()) {
				keySlots[index] = { id: null, public_key: '', key_name: '', saving: false };
				return;
			}

			// Save or update the key
			if (slot.id) {
				// Update existing key (delete and recreate)
				await api.deleteAdminKey(meshId, slot.id);
			}
			const newKey = await api.createAdminKey(meshId, slot.public_key, slot.key_name || undefined);
			keySlots[index] = { id: newKey.id, public_key: slot.public_key, key_name: slot.key_name, saving: false };
			successMessage = `${keyLabel} key saved successfully`;
		} catch (err: any) {
			error = err.message || `Failed to save ${keyLabel.toLowerCase()} key`;
		} finally {
			keySlots[index].saving = false;
		}
	}

	async function handleDeleteMesh() {
		const confirmed = confirm(
			'Are you sure you want to delete this mesh?\n\n' +
			'WARNING: This will permanently delete the mesh and ALL associated nodes. ' +
			'This action cannot be undone.'
		);

		if (!confirmed) return;

		const verification = prompt(
			'This is your last chance to cancel.\n\n' +
			'Type DELETE to confirm deletion of this mesh and all its nodes:'
		);

		if (verification !== 'DELETE') {
			if (verification !== null) {
				alert('Deletion cancelled. You must type DELETE exactly to confirm.');
			}
			return;
		}

		try {
			await api.deleteMesh(meshId);
			goto('/');
		} catch (err: any) {
			error = err.message || 'Failed to delete mesh';
		}
	}

	function changeSection(section: string) {
		activeSection = section;
		error = '';
		successMessage = '';
		if (section === 'keys') loadAdminKeys();
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
						<div class="pt-2 mt-2 border-t border-gray-200">
							<button
								onclick={() => changeSection('danger')}
								class="w-full text-left px-4 py-2 rounded-md text-sm font-medium {activeSection ===
								'danger'
									? 'bg-red-50 text-red-700'
									: 'text-red-600 hover:bg-red-50'}"
							>
								Danger Zone
							</button>
						</div>
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
								<h2 class="text-2xl font-bold text-gray-900 mb-2">Admin Keys</h2>
								<p class="text-gray-600 mb-6">Configure up to 3 admin keys for device administration</p>

								{#if keysLoading}
									<p class="text-gray-500">Loading admin keys...</p>
								{:else}
									<div class="space-y-4">
										{#each keySlots as slot, index}
											{@const keyLabel = index === 0 ? 'Primary' : index === 1 ? 'Secondary' : 'Tertiary'}
											<div class="border rounded-lg p-4">
												<div class="flex items-start gap-4">
													<div class="flex-1 space-y-3">
														<div>
															<label for="key-name-{index}" class="block text-sm font-medium text-gray-700 mb-1">
																{keyLabel} Key Name (optional)
															</label>
															<input
																id="key-name-{index}"
																type="text"
																bind:value={slot.key_name}
																class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
																placeholder="e.g., Main Admin Key"
															/>
														</div>
														<div>
															<label for="key-value-{index}" class="block text-sm font-medium text-gray-700 mb-1">
																Public Key
															</label>
															<input
																id="key-value-{index}"
																type="text"
																bind:value={slot.public_key}
																class="w-full px-3 py-2 border border-gray-300 rounded-md font-mono text-sm"
																placeholder="Enter admin public key..."
															/>
														</div>
													</div>
													<div class="pt-6">
														<button
															onclick={() => handleSaveKey(index)}
															disabled={slot.saving}
															class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap"
														>
															{slot.saving ? 'Saving...' : 'Save'}
														</button>
													</div>
												</div>
											</div>
										{/each}
									</div>
								{/if}
							</div>
						{/if}

						<!-- Danger Zone Section -->
						{#if activeSection === 'danger'}
							<div>
								<h2 class="text-2xl font-bold text-red-700 mb-2">Danger Zone</h2>
								<p class="text-gray-600 mb-6">Irreversible and destructive actions</p>

								<div class="border-2 border-red-200 bg-red-50 rounded-lg p-6">
									<div class="flex items-start gap-4">
										<div class="flex-1">
											<h3 class="text-lg font-semibold text-gray-900 mb-2">Delete this mesh</h3>
											<p class="text-sm text-gray-700 mb-2">
												Once you delete a mesh, there is no going back.
											</p>
											<p class="text-sm text-red-600 font-medium">
												⚠️ This will permanently delete all nodes associated with this mesh.
											</p>
										</div>
										<div>
											<button
												onclick={handleDeleteMesh}
												class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 font-medium whitespace-nowrap"
											>
												Delete Mesh
											</button>
										</div>
									</div>
								</div>
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
