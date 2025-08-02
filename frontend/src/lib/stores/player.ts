import { writable } from 'svelte/store';

export interface Player {
	id: string;
	username: string;
	sessionId: string;
	createdAt: string;
	lastSeen: string;
}

export const playerStore = writable<Player | null>(null);

export function setPlayer(player: Player) {
	playerStore.set(player);
}

export function clearPlayer() {
	playerStore.set(null);
}

export function updatePlayer(updates: Partial<Player>) {
	playerStore.update(player => {
		if (!player) return null;
		return { ...player, ...updates };
	});
} 