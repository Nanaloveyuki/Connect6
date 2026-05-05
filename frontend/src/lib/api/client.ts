import { API_BASE_URL } from '$lib/config/env';
import type { BootstrapResponse, PathResponse, ShortestPathRequest, User } from '$lib/types/connect6';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const response = await fetch(`${API_BASE_URL}${path}`, init);
	if (!response.ok) {
		const payload = (await response.json().catch(() => null)) as { message?: string } | null;
		throw new Error(payload?.message ?? `Request failed with status ${response.status}`);
	}

	return (await response.json()) as T;
}

export function getUser(username: string): Promise<User> {
	return request<User>(`/api/v1/users/${encodeURIComponent(username)}`);
}

export function getBootstrap(): Promise<BootstrapResponse> {
	return request<BootstrapResponse>('/api/v1/bootstrap');
}

export function findShortestPath(payload: ShortestPathRequest): Promise<PathResponse> {
	return request<PathResponse>('/api/v1/graph/path', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(payload)
	});
}
