<script lang="ts">
	import { meshtasticStore } from '$lib/stores/meshtastic.svelte';
	import { authStore } from '$lib/stores/auth.svelte';

	let connecting = $state(false);
	let testMessage = $state('');

	async function handleConnect() {
		connecting = true;
		try {
			await meshtasticStore.connect();
		} catch (err) {
			console.error('Connection failed:', err);
		} finally {
			connecting = false;
		}
	}

	async function handleDisconnect() {
		try {
			await meshtasticStore.disconnect();
		} catch (err) {
			console.error('Disconnect failed:', err);
		}
	}

	async function handleSendTest() {
		if (!testMessage.trim()) return;
		try {
			await meshtasticStore.sendMessage(testMessage);
			testMessage = '';
		} catch (err) {
			console.error('Send failed:', err);
		}
	}
</script>

<div class="min-h-screen bg-gray-50">
	<nav class="bg-white shadow-sm">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<a href="/" class="text-sm text-blue-600 hover:text-blue-800">← Back to Home</a>
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

	<div class="max-w-4xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
		<div class="bg-white shadow rounded-lg">
			<div class="p-6 border-b">
				<h1 class="text-3xl font-bold text-gray-900">Connect Device</h1>
				<p class="text-gray-600 mt-2">
					Connect to your Meshtastic device via USB using Web Serial API
				</p>
			</div>

			<div class="p-6 space-y-6">
				<!-- Browser Support Check -->
				{#if !meshtasticStore.isSupported}
					<div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
						<div class="flex">
							<div class="flex-shrink-0">
								<svg
									class="h-5 w-5 text-yellow-400"
									viewBox="0 0 20 20"
									fill="currentColor"
								>
									<path
										fill-rule="evenodd"
										d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
										clip-rule="evenodd"
									/>
								</svg>
							</div>
							<div class="ml-3">
								<h3 class="text-sm font-medium text-yellow-800">Browser Not Supported</h3>
								<div class="mt-2 text-sm text-yellow-700">
									<p>
										Web Serial API is not supported in your browser. Please use a Chromium-based
										browser like Chrome, Edge, or Opera.
									</p>
								</div>
							</div>
						</div>
					</div>
				{:else}
					<!-- Connection Status -->
					<div class="bg-white border rounded-lg p-4">
						<div class="flex items-center justify-between">
							<div>
								<h3 class="text-lg font-semibold text-gray-900">Connection Status</h3>
								<p class="text-sm text-gray-600 mt-1">
									{#if meshtasticStore.isConnected}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
											Connected
										</span>
									{:else if meshtasticStore.isConnecting}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
											Connecting...
										</span>
									{:else}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
											Disconnected
										</span>
									{/if}
								</p>
							</div>
							<div>
								{#if meshtasticStore.isConnected}
									<button
										onclick={handleDisconnect}
										class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 font-medium"
									>
										Disconnect
									</button>
								{:else}
									<button
										onclick={handleConnect}
										disabled={connecting}
										class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
									>
										{connecting ? 'Connecting...' : 'Connect Device'}
									</button>
								{/if}
							</div>
						</div>
					</div>

					<!-- Error Display -->
					{#if meshtasticStore.error}
						<div class="bg-red-50 border border-red-200 rounded-lg p-4">
							<div class="flex">
								<div class="flex-shrink-0">
									<svg
										class="h-5 w-5 text-red-400"
										viewBox="0 0 20 20"
										fill="currentColor"
									>
										<path
											fill-rule="evenodd"
											d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
											clip-rule="evenodd"
										/>
									</svg>
								</div>
								<div class="ml-3">
									<h3 class="text-sm font-medium text-red-800">Error</h3>
									<div class="mt-2 text-sm text-red-700">
										<p>{meshtasticStore.error}</p>
									</div>
								</div>
							</div>
						</div>
					{/if}

					<!-- Device Information (when connected) -->
					{#if meshtasticStore.isConnected && meshtasticStore.deviceInfo}
						<div class="bg-white border rounded-lg p-4">
							<h3 class="text-lg font-semibold text-gray-900 mb-4">Device Information</h3>
							<div class="grid grid-cols-2 gap-4">
								{#if meshtasticStore.deviceInfo.firmwareVersion}
									<div>
										<p class="text-sm font-medium text-gray-500">Firmware Version</p>
										<p class="text-sm text-gray-900">{meshtasticStore.deviceInfo.firmwareVersion}</p>
									</div>
								{/if}
								{#if meshtasticStore.deviceInfo.hwModel}
									<div>
										<p class="text-sm font-medium text-gray-500">Hardware Model</p>
										<p class="text-sm text-gray-900">{meshtasticStore.deviceInfo.hwModel}</p>
									</div>
								{/if}
								{#if meshtasticStore.deviceInfo.myNodeNum}
									<div>
										<p class="text-sm font-medium text-gray-500">Node Number</p>
										<p class="text-sm text-gray-900">{meshtasticStore.deviceInfo.myNodeNum}</p>
									</div>
								{/if}
								{#if meshtasticStore.deviceInfo.region}
									<div>
										<p class="text-sm font-medium text-gray-500">LoRa Region</p>
										<p class="text-sm text-gray-900">{meshtasticStore.deviceInfo.region}</p>
									</div>
								{/if}
								{#if meshtasticStore.deviceInfo.modemPreset}
									<div>
										<p class="text-sm font-medium text-gray-500">Modem Preset</p>
										<p class="text-sm text-gray-900">{meshtasticStore.deviceInfo.modemPreset}</p>
									</div>
								{/if}
								<div>
									<p class="text-sm font-medium text-gray-500">WiFi</p>
									<p class="text-sm text-gray-900">{meshtasticStore.deviceInfo.hasWifi ? 'Yes' : 'No'}</p>
								</div>
								<div>
									<p class="text-sm font-medium text-gray-500">Bluetooth</p>
									<p class="text-sm text-gray-900">{meshtasticStore.deviceInfo.hasBluetooth ? 'Yes' : 'No'}</p>
								</div>
							</div>
						</div>
					{/if}

					<!-- Messaging (when connected) -->
					{#if meshtasticStore.isConnected}
						<div class="bg-white border rounded-lg p-4">
							<h3 class="text-lg font-semibold text-gray-900 mb-4">Send Message</h3>

							<!-- Send Message -->
							<div class="mb-4">
								<div class="flex gap-2">
									<input
										type="text"
										bind:value={testMessage}
										placeholder="Type a message to send to the mesh..."
										onkeydown={(e) => e.key === 'Enter' && handleSendTest()}
										class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
									/>
									<button
										onclick={handleSendTest}
										disabled={!testMessage.trim()}
										class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
									>
										Send
									</button>
								</div>
								<p class="text-sm text-gray-500 mt-1">
									Press Enter or click Send to transmit to the mesh network
								</p>
							</div>

							<!-- Last Received Message -->
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">
									Last Received Message:
								</label>
								<div class="bg-gray-900 text-green-400 font-mono text-sm p-4 rounded-md min-h-[100px] max-h-[300px] overflow-y-auto">
									{#if meshtasticStore.lastMessage}
										<pre class="whitespace-pre-wrap break-words">{meshtasticStore.lastMessage}</pre>
									{:else}
										<span class="text-gray-500">No messages received yet...</span>
									{/if}
								</div>
							</div>
						</div>
					{/if}

					<!-- Usage Instructions -->
					<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
						<h3 class="text-sm font-medium text-blue-800 mb-2">How to Use</h3>
						<ol class="list-decimal list-inside space-y-1 text-sm text-blue-700">
							<li>Connect your Meshtastic device to your computer via USB</li>
							<li>Click "Connect Device" and select your device from the browser dialog</li>
							<li>Wait for device information to load (may take a few seconds)</li>
							<li>Send messages to the mesh network using the messaging interface</li>
							<li>Monitor received messages in the message display</li>
						</ol>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>
