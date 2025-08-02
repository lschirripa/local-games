<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { playerStore } from '$lib/stores/player';
	import { socketStore } from '$lib/stores/socket';
	import { apiService } from '$lib/services/api';
	import { v4 as uuidv4 } from 'uuid';

	onMount(async () => {
		// Initialize player session
		let sessionId = localStorage.getItem('sessionId');
		if (!sessionId) {
			sessionId = uuidv4();
			localStorage.setItem('sessionId', sessionId);
		}

		try {
			const player = await apiService.createSession(sessionId);
			playerStore.set(player);
		} catch (error) {
			console.error('Failed to create session:', error);
		}
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50">
	<!-- Navigation -->
	<nav class="bg-white shadow-sm border-b border-gray-200">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<a href="/" class="flex items-center space-x-2">
						<div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
							<span class="text-white font-bold text-lg">🎮</span>
						</div>
						<span class="text-xl font-bold text-gradient">Local Games</span>
					</a>
				</div>

				<div class="flex items-center space-x-4">
					{#if $playerStore}
						<span class="text-sm text-gray-600">
							Welcome, {$playerStore.username || 'Player'}
						</span>
					{/if}
					<a href="/games" class="btn btn-primary">Play Games</a>
				</div>
			</div>
		</div>
	</nav>

	<!-- Main Content -->
	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<slot />
	</main>

	<!-- Footer -->
	<footer class="bg-white border-t border-gray-200 mt-auto">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="text-center text-sm text-gray-500">
				<p>&copy; 2024 Local Games. Built with Svelte and Go.</p>
			</div>
		</div>
	</footer>
</div>

<style>
	:global(html) {
		height: 100%;
	}

	:global(body) {
		height: 100%;
	}
</style> 