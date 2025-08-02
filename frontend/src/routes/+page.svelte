<script lang="ts">
	import { onMount } from 'svelte';
	import { playerStore } from '$lib/stores/player';
	import { goto } from '$app/navigation';
	import { GameController, Users, Trophy, Sparkles } from 'lucide-svelte';

	let isLoading = true;

	onMount(() => {
		setTimeout(() => {
			isLoading = false;
		}, 1000);
	});

	function createGame() {
		goto('/games/create');
	}

	function joinGame() {
		goto('/games/join');
	}
</script>

<svelte:head>
	<title>Local Games - Multiplayer Gaming Platform</title>
	<meta name="description" content="Play multiplayer games with friends in real-time. Create or join games instantly." />
</svelte:head>

{#if isLoading}
	<div class="min-h-screen flex items-center justify-center">
		<div class="text-center">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
			<p class="text-gray-600">Loading amazing games...</p>
		</div>
	</div>
{:else}
	<div class="space-y-16">
		<!-- Hero Section -->
		<section class="text-center space-y-8">
			<div class="space-y-4">
				<h1 class="text-5xl md:text-6xl font-bold text-gradient">
					Local Games
				</h1>
				<p class="text-xl text-gray-600 max-w-2xl mx-auto">
					Play multiplayer games with friends in real-time. Create or join games instantly from any device.
				</p>
			</div>

			<div class="flex flex-col sm:flex-row gap-4 justify-center">
				<button on:click={createGame} class="btn btn-primary text-lg px-8 py-4">
					<GameController class="w-5 h-5 mr-2" />
					Create Game
				</button>
				<button on:click={joinGame} class="btn btn-secondary text-lg px-8 py-4">
					<Users class="w-5 h-5 mr-2" />
					Join Game
				</button>
			</div>
		</section>

		<!-- Features Section -->
		<section class="grid md:grid-cols-3 gap-8">
			<div class="card text-center space-y-4">
				<div class="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center mx-auto">
					<GameController class="w-6 h-6 text-primary-600" />
				</div>
				<h3 class="text-xl font-semibold">Multiple Games</h3>
				<p class="text-gray-600">
					Play various game types including the exciting Impostor game and more coming soon.
				</p>
			</div>

			<div class="card text-center space-y-4">
				<div class="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center mx-auto">
					<Users class="w-6 h-6 text-primary-600" />
				</div>
				<h3 class="text-xl font-semibold">Real-time Multiplayer</h3>
				<p class="text-gray-600">
					Connect with friends instantly through WebSocket technology for seamless gameplay.
				</p>
			</div>

			<div class="card text-center space-y-4">
				<div class="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center mx-auto">
					<Trophy class="w-6 h-6 text-primary-600" />
				</div>
				<h3 class="text-xl font-semibold">Cross-platform</h3>
				<p class="text-gray-600">
					Play on any device - desktop, tablet, or mobile with responsive design.
				</p>
			</div>
		</section>

		<!-- Game Types Section -->
		<section class="space-y-8">
			<div class="text-center">
				<h2 class="text-3xl font-bold mb-4">Available Games</h2>
				<p class="text-gray-600">Choose from our growing collection of multiplayer games</p>
			</div>

			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
				<!-- Impostor Game Card -->
				<div class="card hover:shadow-lg transition-shadow cursor-pointer" on:click={() => goto('/games/create?type=impostor')}>
					<div class="space-y-4">
						<div class="w-16 h-16 bg-red-100 rounded-lg flex items-center justify-center mx-auto">
							<Sparkles class="w-8 h-8 text-red-600" />
						</div>
						<div class="text-center">
							<h3 class="text-xl font-semibold mb-2">Impostor</h3>
							<p class="text-gray-600 mb-4">
								Find the impostor among your friends! Each player gets a card with the same answer, except one - the impostor.
							</p>
							<div class="flex justify-center space-x-2">
								<span class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm">2-8 Players</span>
								<span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm">Quick Play</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Coming Soon Cards -->
				{#each ['Word Association', 'Drawing Challenge', 'Trivia Battle'] as game}
					<div class="card opacity-60">
						<div class="space-y-4">
							<div class="w-16 h-16 bg-gray-100 rounded-lg flex items-center justify-center mx-auto">
								<Sparkles class="w-8 h-8 text-gray-400" />
							</div>
							<div class="text-center">
								<h3 class="text-xl font-semibold mb-2">{game}</h3>
								<p class="text-gray-600 mb-4">
									Coming soon! This exciting game is currently in development.
								</p>
								<span class="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm">Coming Soon</span>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</section>

		<!-- How to Play Section -->
		<section class="card">
			<div class="text-center space-y-8">
				<h2 class="text-3xl font-bold">How to Play</h2>
				<div class="grid md:grid-cols-3 gap-8">
					<div class="space-y-4">
						<div class="w-12 h-12 bg-primary-600 rounded-full flex items-center justify-center mx-auto text-white font-bold">
							1
						</div>
						<h3 class="text-lg font-semibold">Create or Join</h3>
						<p class="text-gray-600">
							Create a new game or join an existing one using a room code.
						</p>
					</div>
					<div class="space-y-4">
						<div class="w-12 h-12 bg-primary-600 rounded-full flex items-center justify-center mx-auto text-white font-bold">
							2
						</div>
						<h3 class="text-lg font-semibold">Wait for Players</h3>
						<p class="text-gray-600">
							Wait for other players to join your game room.
						</p>
					</div>
					<div class="space-y-4">
						<div class="w-12 h-12 bg-primary-600 rounded-full flex items-center justify-center mx-auto text-white font-bold">
							3
						</div>
						<h3 class="text-lg font-semibold">Start Playing</h3>
						<p class="text-gray-600">
							Once everyone is ready, start the game and have fun!
						</p>
					</div>
				</div>
			</div>
		</section>
	</div>
{/if} 