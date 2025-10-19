/*
 * Copyright (C) 2025 Michael Graff
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

// Auth store using Svelte 5 runes

import { api, type User } from '$lib/api';

class AuthStore {
	user = $state<User | null>(null);
	loading = $state(true);
	error = $state<string | null>(null);

	async init() {
		const token = api.getToken();
		if (!token) {
			this.loading = false;
			return;
		}

		try {
			this.user = await api.me();
			this.error = null;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to load user';
			api.setToken(null);
		} finally {
			this.loading = false;
		}
	}

	async login(email: string, password: string) {
		try {
			this.loading = true;
			this.error = null;
			const response = await api.login(email, password);
			this.user = response.user;
			return true;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Login failed';
			return false;
		} finally {
			this.loading = false;
		}
	}

	async register(email: string, password: string, displayName: string) {
		try {
			this.loading = true;
			this.error = null;
			const response = await api.register(email, password, displayName);
			this.user = response.user;
			return true;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Registration failed';
			return false;
		} finally {
			this.loading = false;
		}
	}

	async logout() {
		try {
			await api.logout();
		} catch (err) {
			// Ignore errors during logout
			console.error('Logout error:', err);
		} finally {
			this.user = null;
			this.error = null;
		}
	}

	get isAuthenticated() {
		return this.user !== null;
	}
}

export const authStore = new AuthStore();
