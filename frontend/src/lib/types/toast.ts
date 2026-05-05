export type ToastTone = 'info' | 'success' | 'warning' | 'error';

export type ToastItem = {
	id: string;
	title: string;
	description?: string;
	tone: ToastTone;
	durationMs?: number;
};
