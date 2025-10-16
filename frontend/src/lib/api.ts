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
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			...options.headers,
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

	async createMesh(name: string, description?: string) {
		return this.request('/api/meshes', {
			method: 'POST',
			body: JSON.stringify({ name, description }),
		});
	}

	async getMesh(id: number) {
		return this.request(`/api/meshes/${id}`);
	}

	async updateMesh(id: number, name?: string, description?: string) {
		return this.request(`/api/meshes/${id}`, {
			method: 'PUT',
			body: JSON.stringify({ name, description }),
		});
	}

	async deleteMesh(id: number) {
		return this.request(`/api/meshes/${id}`, {
			method: 'DELETE',
		});
	}
}

export const api = new ApiClient(API_BASE_URL);
