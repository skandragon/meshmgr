<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	async function handleLogout() {
		await authStore.logout();
		goto('/login');
	}
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
					<h2 class="text-2xl font-bold text-gray-900 mb-4">Welcome back!</h2>
					<p class="text-gray-600">
						You are logged in as <span class="font-semibold">{authStore.user?.email}</span>
					</p>
					<div class="mt-6">
						<p class="text-gray-700">Your meshes will appear here soon.</p>
					</div>
				</div>
			</div>
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
