<script lang="ts">
	import { serialStore } from '$lib/stores/serial.svelte';
	import { authStore } from '$lib/stores/auth.svelte';

	let baudRate = $state(115200);
	let connecting = $state(false);
	let testMessage = $state('');

	async function handleConnect() {
		connecting = true;
		try {
			await serialStore.connect(baudRate);
		} catch (err) {
			console.error('Connection failed:', err);
		} finally {
			connecting = false;
		}
	}

	async function handleDisconnect() {
		try {
			await serialStore.disconnect();
		} catch (err) {
			console.error('Disconnect failed:', err);
		}
	}

	async function handleSendTest() {
		if (!testMessage.trim()) return;
		try {
			await serialStore.write(testMessage + '\n');
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
				{#if !serialStore.isSupported}
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
									{#if serialStore.isConnected}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
											Connected
										</span>
										{#if serialStore.deviceInfo}
											<span class="ml-2 text-gray-500">
												VID: 0x{serialStore.deviceInfo.vendorId?.toString(16).toUpperCase() ?? '????'}
												PID: 0x{serialStore.deviceInfo.productId?.toString(16).toUpperCase() ?? '????'}
											</span>
										{/if}
									{:else}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
											Disconnected
										</span>
									{/if}
								</p>
							</div>
							<div>
								{#if serialStore.isConnected}
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
					{#if serialStore.error}
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
										<p>{serialStore.error}</p>
									</div>
								</div>
							</div>
						</div>
					{/if}

					<!-- Connection Settings -->
					{#if !serialStore.isConnected}
						<div class="bg-white border rounded-lg p-4">
							<h3 class="text-lg font-semibold text-gray-900 mb-4">Connection Settings</h3>
							<div class="space-y-4">
								<div>
									<label for="baud-rate" class="block text-sm font-medium text-gray-700 mb-1">
										Baud Rate
									</label>
									<select
										id="baud-rate"
										bind:value={baudRate}
										class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
									>
										<option value={9600}>9600</option>
										<option value={19200}>19200</option>
										<option value={38400}>38400</option>
										<option value={57600}>57600</option>
										<option value={115200}>115200 (Default for Meshtastic)</option>
										<option value={230400}>230400</option>
										<option value={460800}>460800</option>
										<option value={921600}>921600</option>
									</select>
									<p class="text-sm text-gray-500 mt-1">
										Meshtastic devices typically use 115200 baud
									</p>
								</div>
							</div>
						</div>
					{/if}

					<!-- Serial Monitor (when connected) -->
					{#if serialStore.isConnected}
						<div class="bg-white border rounded-lg p-4">
							<h3 class="text-lg font-semibold text-gray-900 mb-4">Serial Monitor</h3>

							<!-- Send Message -->
							<div class="mb-4">
								<div class="flex gap-2">
									<input
										type="text"
										bind:value={testMessage}
										placeholder="Type a message to send..."
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
									Press Enter or click Send to transmit
								</p>
							</div>

							<!-- Last Received Message -->
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">
									Last Received:
								</label>
								<div class="bg-gray-900 text-green-400 font-mono text-sm p-4 rounded-md min-h-[100px] max-h-[300px] overflow-y-auto">
									{#if serialStore.lastMessage}
										<pre class="whitespace-pre-wrap break-words">{serialStore.lastMessage}</pre>
									{:else}
										<span class="text-gray-500">No data received yet...</span>
									{/if}
								</div>
								<p class="text-sm text-gray-500 mt-1">
									Reading: {serialStore.isReading ? 'Active' : 'Inactive'}
								</p>
							</div>
						</div>
					{/if}

					<!-- Usage Instructions -->
					<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
						<h3 class="text-sm font-medium text-blue-800 mb-2">How to Use</h3>
						<ol class="list-decimal list-inside space-y-1 text-sm text-blue-700">
							<li>Connect your Meshtastic device to your computer via USB</li>
							<li>Click "Connect Device" and select your device from the browser dialog</li>
							<li>Once connected, you can send commands and view received data</li>
							<li>Use the serial monitor to test communication with your device</li>
						</ol>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>
