<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';

  let playerId: string | null = null;
  let playerName: string | null = null;

  onMount(() => {
    // Get player info from localStorage
    const storedPlayerId = localStorage.getItem('playerId');
    const storedPlayerName = localStorage.getItem('playerName');
    
    if (storedPlayerId && storedPlayerName) {
      playerId = storedPlayerId;
      playerName = storedPlayerName;
    }
  });

  function logout() {
    localStorage.removeItem('playerId');
    localStorage.removeItem('playerName');
    playerId = null;
    playerName = null;
    window.location.href = '/';
  }
</script>

<svelte:head>
  <title>Local Games - Multiplayer Gaming Platform</title>
  <meta name="description" content="Play multiplayer games with friends in real-time" />
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="">
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-100">
  <!-- Navigation -->
  <nav class="bg-white shadow-lg border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <a href="/" class="flex items-center space-x-2">
            <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-lg">G</span>
            </div>
            <span class="text-xl font-bold text-gray-900">Local Games</span>
          </a>
        </div>
        
        <div class="flex items-center space-x-4">
          {#if playerId && playerName}
            <span class="text-gray-700">Welcome, {playerName}!</span>
            <button on:click={logout} class="btn btn-secondary">
              Logout
            </button>
          {:else}
            <a href="/login" class="btn btn-primary">
              Login
            </a>
          {/if}
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
      <div class="text-center text-gray-600">
        <p>&copy; 2024 Local Games. Built with SvelteKit and Go.</p>
      </div>
    </div>
  </footer>
</div>
