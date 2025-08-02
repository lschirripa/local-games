import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';

export interface SocketMessage {
	type: string;
	data: any;
	gameId?: string;
	playerId?: string;
}

export interface SocketState {
	isConnected: boolean;
	isConnecting: boolean;
	error: string | null;
}

export const socketStore: Writable<SocketState> = writable({
	isConnected: false,
	isConnecting: false,
	error: null
});

export function setSocketState(state: Partial<SocketState>) {
	socketStore.update(current => ({ ...current, ...state }));
}

export function setConnected(connected: boolean) {
	socketStore.update(state => ({ ...state, isConnected: connected }));
}

export function setConnecting(connecting: boolean) {
	socketStore.update(state => ({ ...state, isConnecting: connecting }));
}

export function setError(error: string | null) {
	socketStore.update(state => ({ ...state, error }));
} 