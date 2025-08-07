import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { garmApi } from '../api/client.js';

interface AuthState {
	isAuthenticated: boolean;
	user: string | null;
	loading: boolean;
}

const initialState: AuthState = {
	isAuthenticated: false,
	user: null,
	loading: true
};

export const authStore = writable<AuthState>(initialState);

// Cookie utilities
function setCookie(name: string, value: string, days: number = 7): void {
	if (!browser) return;
	
	const expires = new Date();
	expires.setTime(expires.getTime() + (days * 24 * 60 * 60 * 1000));
	document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;SameSite=Lax`;
	console.log('Set cookie:', name, value.substring(0, 20) + '...');
}

function getCookie(name: string): string | null {
	if (!browser) return null;
	
	const nameEQ = name + "=";
	const ca = document.cookie.split(';');
	for (let i = 0; i < ca.length; i++) {
		let c = ca[i];
		while (c.charAt(0) === ' ') c = c.substring(1, c.length);
		if (c.indexOf(nameEQ) === 0) {
			const value = c.substring(nameEQ.length, c.length);
			console.log('Found cookie:', name, value.substring(0, 20) + '...');
			return value;
		}
	}
	console.log('Cookie not found:', name, 'Available cookies:', document.cookie);
	return null;
}

function deleteCookie(name: string): void {
	if (!browser) return;
	document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:01 GMT;path=/`;
}

// Auth utilities
export const auth = {
	async login(username: string, password: string): Promise<void> {
		try {
			authStore.update(state => ({ ...state, loading: true }));
			
			const response = await garmApi.login({ username, password });
			
			// Store JWT token in cookies for server authentication and set it in the API client
			if (browser) {
				setCookie('garm_token', response.token);
				setCookie('garm_user', username);
			}
			
			// Set the token in the API client for future requests
			garmApi.setToken(response.token);
			
			authStore.set({
				isAuthenticated: true,
				user: username,
				loading: false
			});
		} catch (error) {
			authStore.update(state => ({ ...state, loading: false }));
			throw error;
		}
	},

	logout(): void {
		if (browser) {
			deleteCookie('garm_token');
			deleteCookie('garm_user');
		}
		
		authStore.set({
			isAuthenticated: false,
			user: null,
			loading: false
		});
	},

	init(): void {
		console.log('Auth init called, browser:', browser);
		if (browser) {
			const token = getCookie('garm_token');
			const user = getCookie('garm_user');
			
			console.log('Auth init - token:', token ? 'present' : 'missing', 'user:', user || 'missing');
			
			if (token && user) {
				// Set the token in the API client for future requests
				garmApi.setToken(token);
				
				// Optimistically set authenticated state
				console.log('Setting authenticated state');
				authStore.set({
					isAuthenticated: true,
					user,
					loading: false
				});
			} else {
				console.log('No token or user, setting unauthenticated');
				authStore.update(state => ({ ...state, loading: false }));
			}
		} else {
			console.log('Not in browser, setting loading false');
			authStore.update(state => ({ ...state, loading: false }));
		}
	},

	// Check if token is still valid by making a test API call
	async checkAuth(): Promise<boolean> {
		try {
			await garmApi.getControllerInfo();
			return true;
		} catch {
			// Token is invalid, logout
			auth.logout();
			return false;
		}
	}
};