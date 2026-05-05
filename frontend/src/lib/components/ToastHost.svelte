<script lang="ts">
	import { fly } from 'svelte/transition';

	import { toastStore } from '$lib/stores/toast';

	function dismiss(id: string) {
		toastStore.dismiss(id);
	}
</script>

<div class="toast-layer" aria-live="polite" aria-atomic="true">
	{#each $toastStore as toast (toast.id)}
		<section class={`toast-card tone-${toast.tone}`} in:fly={{ y: -12, duration: 220 }} out:fly={{ y: -8, duration: 180 }}>
			<div class="toast-copy">
				<strong>{toast.title}</strong>
				{#if toast.description}
					<p>{toast.description}</p>
				{/if}
			</div>
			<button type="button" class="toast-dismiss" on:click={() => dismiss(toast.id)} aria-label="Dismiss notification">
				×
			</button>
		</section>
	{/each}
</div>

<style>
	.toast-layer {
		position: fixed;
		top: 1rem;
		right: 1rem;
		z-index: 50;
		display: grid;
		gap: 0.75rem;
		width: min(22rem, calc(100vw - 1.5rem));
		pointer-events: none;
	}

	.toast-card {
		pointer-events: auto;
		display: grid;
		grid-template-columns: 1fr auto;
		gap: 0.75rem;
		align-items: start;
		padding: 0.95rem 1rem;
		border-radius: 18px;
		border: 1px solid rgba(150, 184, 214, 0.22);
		background: rgba(255, 255, 255, 0.84);
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.16);
		backdrop-filter: blur(16px);
	}

	.toast-copy strong,
	.toast-copy p {
		margin: 0;
	}

	.toast-copy strong {
		display: block;
		color: #21466f;
	}

	.toast-copy p {
		margin-top: 0.2rem;
		font-size: 0.88rem;
		color: #6f86a4;
	}

	.toast-dismiss {
		border: none;
		background: transparent;
		color: #7990ad;
		font-size: 1.05rem;
		line-height: 1;
		cursor: pointer;
	}

	.tone-info {
		border-color: rgba(69, 138, 244, 0.18);
	}

	.tone-success {
		border-color: rgba(29, 187, 141, 0.2);
	}

	.tone-warning {
		border-color: rgba(243, 156, 86, 0.22);
	}

	.tone-error {
		border-color: rgba(208, 91, 112, 0.22);
	}

	@media (max-width: 860px) {
		.toast-layer {
			top: auto;
			right: 0.75rem;
			bottom: 0.75rem;
			left: 0.75rem;
			width: auto;
		}
	}
</style>
