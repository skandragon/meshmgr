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

// API client for communicating with the Go backend

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface User {
	id: number;
	email: string;
	display_name: string;
	created_at: string;
	updated_at: string;
}

export interface AuthResponse {
	token: string;
	user: User;
}

export interface ErrorResponse {
	error: string;
}

class ApiClient {
	private baseUrl: string;
	private token: string | null = null;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
		// Load token from localStorage if available
		if (typeof window !== 'undefined') {
			this.token = localStorage.getItem('auth_token');
		}
	}

	setToken(token: string | null) {
		this.token = token;
		if (typeof window !== 'undefined') {
			if (token) {
				localStorage.setItem('auth_token', token);
			} else {
				localStorage.removeItem('auth_token');
			}
		}
	}

	getToken(): string | null {
		return this.token;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<T> {
		const headers: Record<string, string> = {
			'Content-Type': 'application/json',
			...(options.headers as Record<string, string>),
		};

		if (this.token) {
			headers['Authorization'] = `Bearer ${this.token}`;
		}

		const response = await fetch(`${this.baseUrl}${endpoint}`, {
			...options,
			headers,
		});

		if (!response.ok) {
			const error: ErrorResponse = await response.json().catch(() => ({
				error: 'An unexpected error occurred',
			}));
			throw new Error(error.error);
		}

		return response.json();
	}

	async register(email: string, password: string, displayName: string): Promise<AuthResponse> {
		const response = await this.request<AuthResponse>('/api/auth/register', {
			method: 'POST',
			body: JSON.stringify({
				email,
				password,
				display_name: displayName,
			}),
		});
		this.setToken(response.token);
		return response;
	}

	async login(email: string, password: string): Promise<AuthResponse> {
		const response = await this.request<AuthResponse>('/api/auth/login', {
			method: 'POST',
			body: JSON.stringify({
				email,
				password,
			}),
		});
		this.setToken(response.token);
		return response;
	}

	async logout(): Promise<void> {
		await this.request('/api/auth/logout', {
			method: 'POST',
		});
		this.setToken(null);
	}

	async me(): Promise<User> {
		return this.request<User>('/api/auth/me');
	}

	async listMeshes() {
		return this.request('/api/meshes');
	}

	async createMesh(
		name: string,
		description?: string,
		loraRegion?: string,
		modemPreset?: string,
		frequencySlot?: number
	) {
		return this.request('/api/meshes', {
			method: 'POST',
			body: JSON.stringify({
				name,
				description,
				lora_region: loraRegion,
				modem_preset: modemPreset,
				frequency_slot: frequencySlot
			}),
		});
	}

	async getMesh(id: number) {
		return this.request(`/api/meshes/${id}`);
	}

	async updateMesh(
		id: number,
		data: {
			name?: string;
			description?: string;
			lora_region?: string;
			modem_preset?: string;
			frequency_slot?: number;
		}
	) {
		return this.request(`/api/meshes/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data),
		});
	}

	async deleteMesh(id: number) {
		return this.request(`/api/meshes/${id}`, {
			method: 'DELETE',
		});
	}

	// Mesh Access Management
	async listMeshAccess(meshId: number) {
		return this.request(`/api/meshes/${meshId}/access`);
	}

	async grantMeshAccess(meshId: number, userEmail: string, accessLevel: string) {
		return this.request(`/api/meshes/${meshId}/access`, {
			method: 'POST',
			body: JSON.stringify({ user_email: userEmail, access_level: accessLevel }),
		});
	}

	async updateMeshAccess(meshId: number, userId: number, accessLevel: string) {
		return this.request(`/api/meshes/${meshId}/access/${userId}`, {
			method: 'PUT',
			body: JSON.stringify({ access_level: accessLevel }),
		});
	}

	async revokeMeshAccess(meshId: number, userId: number) {
		return this.request(`/api/meshes/${meshId}/access/${userId}`, {
			method: 'DELETE',
		});
	}

	// Admin Keys
	async listAdminKeys(meshId: number) {
		return this.request(`/api/meshes/${meshId}/admin-keys`);
	}

	async getAdminKey(meshId: number, keyId: number) {
		return this.request(`/api/meshes/${meshId}/admin-keys/${keyId}`);
	}

	async createAdminKey(meshId: number, publicKey: string, keyName?: string) {
		return this.request(`/api/meshes/${meshId}/admin-keys`, {
			method: 'POST',
			body: JSON.stringify({ public_key: publicKey, key_name: keyName }),
		});
	}

	async deleteAdminKey(meshId: number, keyId: number) {
		return this.request(`/api/meshes/${meshId}/admin-keys/${keyId}`, {
			method: 'DELETE',
		});
	}

	// Nodes
	async listNodes(meshId: number) {
		return this.request(`/api/meshes/${meshId}/nodes`);
	}

	async getNode(meshId: number, nodeId: number) {
		return this.request(`/api/meshes/${meshId}/nodes/${nodeId}`);
	}

	async createNode(meshId: number, data: {
		hardware_id: string;
		name: string;
		long_name: string;
		role?: string;
		public_key?: string;
		private_key?: string;
		status?: string;
		unmessageable?: boolean;
	}) {
		return this.request(`/api/meshes/${meshId}/nodes`, {
			method: 'POST',
			body: JSON.stringify(data),
		});
	}

	async updateNode(meshId: number, nodeId: number, data: {
		name?: string;
		long_name?: string;
		role?: string;
		public_key?: string;
		private_key?: string;
		status?: string;
		unmessageable?: boolean;
		pending_changes?: boolean;
	}) {
		return this.request(`/api/meshes/${meshId}/nodes/${nodeId}`, {
			method: 'PUT',
			body: JSON.stringify(data),
		});
	}

	async updateNodeStatus(meshId: number, nodeId: number, status: string) {
		return this.request(`/api/meshes/${meshId}/nodes/${nodeId}/status`, {
			method: 'PATCH',
			body: JSON.stringify({ status }),
		});
	}

	async deleteNode(meshId: number, nodeId: number) {
		return this.request(`/api/meshes/${meshId}/nodes/${nodeId}`, {
			method: 'DELETE',
		});
	}

	// User API Keys
	async listAPIKeys() {
		return this.request('/api/user/api-keys');
	}

	async createAPIKey(keyName: string, expiresIn?: number) {
		return this.request('/api/user/api-keys', {
			method: 'POST',
			body: JSON.stringify({
				key_name: keyName,
				expires_in: expiresIn
			}),
		});
	}

	async deleteAPIKey(keyId: number) {
		return this.request(`/api/user/api-keys/${keyId}`, {
			method: 'DELETE',
		});
	}
}

export const api = new ApiClient(API_BASE_URL);
