const API_BASE = '/api';

export interface Game {
	id: string;
	name: string;
	players: number;
	maxPlayers: number;
	createdAt: string;
}

export interface Player {
	id: string;
	name: string;
}

class ApiService {
	private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
		const url = `${API_BASE}${endpoint}`;
		const response = await fetch(url, {
			headers: {
				'Content-Type': 'application/json',
				...options?.headers
			},
			...options
		});

		if (!response.ok) {
			throw new Error(`API Error: ${response.status} ${response.statusText}`);
		}

		return response.json();
	}

	// Game methods
	async getGames(): Promise<Game[]> {
		return this.request<Game[]>('/games');
	}

	async createGame(game: Omit<Game, 'id' | 'createdAt'>): Promise<Game> {
		return this.request<Game>('/games', {
			method: 'POST',
			body: JSON.stringify(game)
		});
	}

	async getGame(id: string): Promise<Game> {
		return this.request<Game>(`/games/${id}`);
	}

	// Player methods
	async getPlayers(): Promise<Player[]> {
		return this.request<Player[]>('/players');
	}

	async createPlayer(player: Omit<Player, 'id'>): Promise<Player> {
		return this.request<Player>('/players', {
			method: 'POST',
			body: JSON.stringify(player)
		});
	}
}

export const api = new ApiService();