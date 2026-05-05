import { writable } from 'svelte/store';
import { browser } from '$app/environment';

import type { ToastItem, ToastTone } from '$lib/types/toast';

const DEFAULT_DURATION_MS = 4200;

function createToastStore() {
	const { subscribe, update } = writable<ToastItem[]>([]);

	function dismiss(id: string) {
		update((items) => items.filter((item) => item.id !== id));
	}

	function push(input: Omit<ToastItem, 'id'>) {
		const id = browser && typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function'
			? crypto.randomUUID()
			: `toast-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
		const toast: ToastItem = {
			id,
			durationMs: DEFAULT_DURATION_MS,
			...input
		};

		update((items) => [...items, toast]);

		const duration = toast.durationMs ?? DEFAULT_DURATION_MS;
		if (duration > 0) {
			setTimeout(() => dismiss(id), duration);
		}

		return id;
	}

	function notify(tone: ToastTone, title: string, description?: string, durationMs?: number) {
		return push({ tone, title, description, durationMs });
	}

	return {
		subscribe,
		dismiss,
		push,
		info: (title: string, description?: string, durationMs?: number) => notify('info', title, description, durationMs),
		success: (title: string, description?: string, durationMs?: number) => notify('success', title, description, durationMs),
		warning: (title: string, description?: string, durationMs?: number) => notify('warning', title, description, durationMs),
		error: (title: string, description?: string, durationMs?: number) => notify('error', title, description, durationMs)
	};
}

export const toastStore = createToastStore();
