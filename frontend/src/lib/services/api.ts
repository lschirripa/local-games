import type { Player } from '$lib/stores/player';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class ApiService {
	private baseUrl: string;

	constructor() {
		this.baseUrl = API_BASE_URL;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<T> {
		const url = `${this.baseUrl}${endpoint}`;
		const config: RequestInit = {
			headers: {
				'Content-Type': 'application/json',
				...options.headers
			},
			...options
		};

		try {
			const response = await fetch(url, config);
			
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			return await response.json();
		} catch (error) {
			console.error('API request failed:', error);
			throw error;
		}
	}

	// Authentication
	async createSession(sessionId: string): Promise<Player> {
		return this.request<Player>('/api/auth/session', {
			method: 'POST',
			body: JSON.stringify({ sessionId })
		});
	}

	async deleteSession(): Promise<void> {
		return this.request<void>('/api/auth/session', {
			method: 'DELETE'
		});
	}

	async getCurrentPlayer(): Promise<Player> {
		return this.request<Player>('/api/auth/me');
	}

	// Games
	async getGames(): Promise<any[]> {
		return this.request<any[]>('/api/games');
	}

	async createGame(gameData: {
		gameType: string;
		maxPlayers?: number;
		gameConfig?: any;
	}): Promise<any> {
		return this.request<any>('/api/games', {
			method: 'POST',
			body: JSON.stringify(gameData)
		});
	}

	async getGame(gameId: string): Promise<any> {
		return this.request<any>(`/api/games/${gameId}`);
	}

	async updateGame(gameId: string, updates: any): Promise<any> {
		return this.request<any>(`/api/games/${gameId}`, {
			method: 'PUT',
			body: JSON.stringify(updates)
		});
	}

	async deleteGame(gameId: string): Promise<void> {
		return this.request<void>(`/api/games/${gameId}`, {
			method: 'DELETE'
		});
	}

	async joinGame(gameId: string, playerId: string): Promise<any> {
		return this.request<any>(`/api/games/${gameId}/join`, {
			method: 'POST',
			body: JSON.stringify({ playerId })
		});
	}

	async leaveGame(gameId: string, playerId: string): Promise<void> {
		return this.request<void>(`/api/games/${gameId}/leave`, {
			method: 'POST',
			body: JSON.stringify({ playerId })
		});
	}

	// Players
	async getPlayer(playerId: string): Promise<Player> {
		return this.request<Player>(`/api/players/${playerId}`);
	}

	async updatePlayer(playerId: string, updates: Partial<Player>): Promise<Player> {
		return this.request<Player>(`/api/players/${playerId}`, {
			method: 'PUT',
			body: JSON.stringify(updates)
		});
	}

	// Health check
	async healthCheck(): Promise<{ status: string; message: string }> {
		return this.request<{ status: string; message: string }>('/health');
	}
}

export const apiService = new ApiService(); 