<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Game, type Player } from '$lib/services/api';

	let games: Game[] = [];
	let players: Player[] = [];
	let newGameName = '';
	let newPlayerName = '';
	let loading = false;
	let error = '';

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			[games, players] = await Promise.all([
				api.getGames(),
				api.getPlayers()
			]);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load data';
		} finally {
			loading = false;
		}
	}

	async function createGame() {
		if (!newGameName.trim()) return;
		
		try {
			const game = await api.createGame({
				name: newGameName,
				players: 0,
				maxPlayers: 4
			});
			games = [...games, game];
			newGameName = '';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create game';
		}
	}

	async function createPlayer() {
		if (!newPlayerName.trim()) return;
		
		try {
			const player = await api.createPlayer({
				name: newPlayerName
			});
			players = [...players, player];
			newPlayerName = '';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create player';
		}
	}
</script>

<div class="container">
	<h1>Local Games Platform</h1>

	{#if error}
		<div class="error">
			Error: {error}
			<button on:click={() => error = ''}>×</button>
		</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{/if}

	<div class="section">
		<h2>Games</h2>
		<div class="form">
			<input 
				bind:value={newGameName} 
				placeholder="Game name"
				on:keydown={(e) => e.key === 'Enter' && createGame()}
			/>
			<button on:click={createGame}>Create Game</button>
		</div>
		
		<div class="items">
			{#each games as game}
				<div class="item">
					<h3>{game.name}</h3>
					<p>Players: {game.players}/{game.maxPlayers}</p>
					<p>Created: {new Date(game.createdAt).toLocaleString()}</p>
				</div>
			{/each}
		</div>
	</div>

	<div class="section">
		<h2>Players</h2>
		<div class="form">
			<input 
				bind:value={newPlayerName} 
				placeholder="Player name"
				on:keydown={(e) => e.key === 'Enter' && createPlayer()}
			/>
			<button on:click={createPlayer}>Add Player</button>
		</div>
		
		<div class="items">
			{#each players as player}
				<div class="item">
					<h3>{player.name}</h3>
					<p>ID: {player.id}</p>
				</div>
			{/each}
		</div>
	</div>
</div>

<style>
	.container {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
	}

	.section {
		margin: 2rem 0;
	}

	.form {
		display: flex;
		gap: 1rem;
		margin: 1rem 0;
	}

	.form input {
		flex: 1;
		padding: 0.5rem;
		border: 1px solid #ccc;
		border-radius: 4px;
	}

	.form button {
		padding: 0.5rem 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	.items {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
		gap: 1rem;
	}

	.item {
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1rem;
	}

	.error {
		background: #f8d7da;
		color: #721c24;
		padding: 1rem;
		border-radius: 4px;
		margin: 1rem 0;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.loading {
		text-align: center;
		padding: 2rem;
	}
</style>
