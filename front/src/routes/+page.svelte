<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let playerId: string | null = null;
  let playerName: string | null = null;
  let showLoginModal = false;
  let showCreateGameModal = false;
  let games: any[] = [];
  let loading = true;

  // Form data
  let loginForm = {
    name: ''
  };

  let createGameForm = {
    name: '',
    type: 'impostor',
    maxPlayers: 8,
    minPlayers: 3,
    categories: ['football_players', 'movies', 'animals'],
    rounds: 3,
    timePerRound: 60
  };

  onMount(() => {
    // Check if player is logged in
    const storedPlayerId = localStorage.getItem('playerId');
    const storedPlayerName = localStorage.getItem('playerName');
    
    if (storedPlayerId && storedPlayerName) {
      playerId = storedPlayerId;
      playerName = storedPlayerName;
      loadGames();
    } else {
      loading = false;
    }
  });

  async function createPlayer() {
    if (!loginForm.name.trim()) return;

    try {
      const response = await fetch('http://localhost:8080/api/v1/players', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: loginForm.name }),
      });

      if (response.ok) {
        const data = await response.json();
        playerId = data.player.id;
        playerName = data.player.name;
        
        // Store in localStorage
        localStorage.setItem('playerId', playerId);
        localStorage.setItem('playerName', playerName);
        
        showLoginModal = false;
        loginForm.name = '';
        loadGames();
      } else {
        const error = await response.json();
        alert('Error creating player: ' + error.error);
      }
    } catch (error) {
      console.error('Error creating player:', error);
      alert('Failed to create player. Please try again.');
    }
  }

  async function loadGames() {
    try {
      const response = await fetch('http://localhost:8080/api/v1/games');
      if (response.ok) {
        const data = await response.json();
        games = data.games || [];
      }
    } catch (error) {
      console.error('Error loading games:', error);
    } finally {
      loading = false;
    }
  }

  async function createGame() {
    if (!createGameForm.name.trim()) return;

    try {
      const response = await fetch('http://localhost:8080/api/v1/games', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Player-ID': playerId,
        },
        body: JSON.stringify({
          name: createGameForm.name,
          type: createGameForm.type,
          max_players: createGameForm.maxPlayers,
          min_players: createGameForm.minPlayers,
          settings: {
            categories: createGameForm.categories,
            rounds: createGameForm.rounds,
            time_per_round: createGameForm.timePerRound,
            voting_enabled: true,
            auto_start: false
          }
        }),
      });

      if (response.ok) {
        const data = await response.json();
        showCreateGameModal = false;
        createGameForm.name = '';
        loadGames();
        // Navigate to the new game
        goto(`/game/${data.game.id}`);
      } else {
        const error = await response.json();
        alert('Error creating game: ' + error.error);
      }
    } catch (error) {
      console.error('Error creating game:', error);
      alert('Failed to create game. Please try again.');
    }
  }

  async function joinGame(gameId: string) {
    try {
      const response = await fetch(`http://localhost:8080/api/v1/games/${gameId}/join`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ player_id: playerId }),
      });

      if (response.ok) {
        goto(`/game/${gameId}`);
      } else {
        const error = await response.json();
        alert('Error joining game: ' + error.error);
      }
    } catch (error) {
      console.error('Error joining game:', error);
      alert('Failed to join game. Please try again.');
    }
  }

  function openLoginModal() {
    showLoginModal = true;
  }

  function closeLoginModal() {
    showLoginModal = false;
    loginForm.name = '';
  }

  function openCreateGameModal() {
    if (!playerId) {
      openLoginModal();
      return;
    }
    showCreateGameModal = true;
  }

  function closeCreateGameModal() {
    showCreateGameModal = false;
    createGameForm.name = '';
  }
</script>

<svelte:head>
  <title>Local Games - Home</title>
</svelte:head>

<div class="space-y-8">
  <!-- Hero Section -->
  <div class="text-center space-y-6">
    <h1 class="text-4xl md:text-6xl font-bold text-gray-900">
      Welcome to <span class="text-primary-600">Local Games</span>
    </h1>
    <p class="text-xl text-gray-600 max-w-2xl mx-auto">
      Play multiplayer games with friends in real-time. Create or join games and have fun together!
    </p>
    
    {#if !playerId}
      <div class="flex justify-center space-x-4">
        <button on:click={openLoginModal} class="btn btn-primary text-lg px-8 py-3">
          Get Started
        </button>
      </div>
    {:else}
      <div class="flex justify-center space-x-4">
        <button on:click={openCreateGameModal} class="btn btn-primary text-lg px-8 py-3">
          Create Game
        </button>
        <button on:click={loadGames} class="btn btn-secondary text-lg px-8 py-3">
          Refresh Games
        </button>
      </div>
    {/if}
  </div>

  <!-- Games Section -->
  {#if playerId}
    <div class="space-y-6">
      <h2 class="text-2xl font-bold text-gray-900">Available Games</h2>
      
      {#if loading}
        <div class="flex justify-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        </div>
      {:else if games.length === 0}
        <div class="text-center py-12">
          <div class="text-gray-400 text-6xl mb-4">🎮</div>
          <h3 class="text-xl font-semibold text-gray-600 mb-2">No games available</h3>
          <p class="text-gray-500">Be the first to create a game!</p>
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {#each games as game}
            <div class="card hover:shadow-xl transition-shadow duration-200">
              <div class="flex justify-between items-start mb-4">
                <h3 class="text-lg font-semibold text-gray-900">{game.name}</h3>
                <span class="px-2 py-1 text-xs font-medium rounded-full {game.status === 'waiting' ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'}">
                  {game.status}
                </span>
              </div>
              
              <div class="space-y-2 text-sm text-gray-600 mb-4">
                <p><strong>Type:</strong> {game.type}</p>
                <p><strong>Players:</strong> {game.players?.length || 0}/{game.max_players}</p>
                <p><strong>Created:</strong> {new Date(game.created_at).toLocaleDateString()}</p>
              </div>
              
              <button 
                on:click={() => joinGame(game.id)}
                class="w-full btn btn-primary"
                disabled={game.status !== 'waiting'}
              >
                {game.status === 'waiting' ? 'Join Game' : 'Game in Progress'}
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {/if}

  <!-- Features Section -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-16">
    <div class="text-center">
      <div class="text-4xl mb-4">🎯</div>
      <h3 class="text-xl font-semibold mb-2">Impostor Game</h3>
      <p class="text-gray-600">Find the impostor among your friends. Everyone gets the same word except one!</p>
    </div>
    
    <div class="text-center">
      <div class="text-4xl mb-4">⚡</div>
      <h3 class="text-xl font-semibold mb-2">Real-time</h3>
      <p class="text-gray-600">Play with friends in real-time using WebSocket connections.</p>
    </div>
    
    <div class="text-center">
      <div class="text-4xl mb-4">📱</div>
      <h3 class="text-xl font-semibold mb-2">Responsive</h3>
      <p class="text-gray-600">Works perfectly on desktop, tablet, and mobile devices.</p>
    </div>
  </div>
</div>

<!-- Login Modal -->
{#if showLoginModal}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-xl p-8 max-w-md w-full mx-4">
      <h2 class="text-2xl font-bold mb-6">Join the Fun!</h2>
      
      <div class="space-y-4">
        <div>
          <label for="playerName" class="block text-sm font-medium text-gray-700 mb-2">
            Your Name
          </label>
          <input
            id="playerName"
            type="text"
            bind:value={loginForm.name}
            placeholder="Enter your name"
            class="input"
            maxlength="50"
          />
        </div>
        
        <div class="flex space-x-3 pt-4">
          <button on:click={createPlayer} class="btn btn-primary flex-1">
            Create Player
          </button>
          <button on:click={closeLoginModal} class="btn btn-secondary flex-1">
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<!-- Create Game Modal -->
{#if showCreateGameModal}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-xl p-8 max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
      <h2 class="text-2xl font-bold mb-6">Create New Game</h2>
      
      <div class="space-y-4">
        <div>
          <label for="gameName" class="block text-sm font-medium text-gray-700 mb-2">
            Game Name
          </label>
          <input
            id="gameName"
            type="text"
            bind:value={createGameForm.name}
            placeholder="Enter game name"
            class="input"
            maxlength="100"
          />
        </div>
        
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="maxPlayers" class="block text-sm font-medium text-gray-700 mb-2">
              Max Players
            </label>
            <input
              id="maxPlayers"
              type="number"
              bind:value={createGameForm.maxPlayers}
              min="3"
              max="20"
              class="input"
            />
          </div>
          
          <div>
            <label for="minPlayers" class="block text-sm font-medium text-gray-700 mb-2">
              Min Players
            </label>
            <input
              id="minPlayers"
              type="number"
              bind:value={createGameForm.minPlayers}
              min="2"
              max="10"
              class="input"
            />
          </div>
        </div>
        
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="rounds" class="block text-sm font-medium text-gray-700 mb-2">
              Rounds
            </label>
            <input
              id="rounds"
              type="number"
              bind:value={createGameForm.rounds}
              min="1"
              max="10"
              class="input"
            />
          </div>
          
          <div>
            <label for="timePerRound" class="block text-sm font-medium text-gray-700 mb-2">
              Time per Round (seconds)
            </label>
            <input
              id="timePerRound"
              type="number"
              bind:value={createGameForm.timePerRound}
              min="30"
              max="300"
              class="input"
            />
          </div>
        </div>
        
        <div class="flex space-x-3 pt-4">
          <button on:click={createGame} class="btn btn-primary flex-1">
            Create Game
          </button>
          <button on:click={closeCreateGameModal} class="btn btn-secondary flex-1">
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
