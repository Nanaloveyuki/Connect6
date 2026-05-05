<script lang="ts">
	import type { BootstrapResponse, SearchState } from '$lib/types/connect6';
	import type { I18nDictionary } from '$lib/i18n/messages';

	export let bootstrap: BootstrapResponse;
	export let state: SearchState;
	export let loading = false;
	export let dictionary: I18nDictionary;
	export let onSubmit: () => Promise<void> | void;
</script>

<section class="panel animate-in">
	<p class="eyebrow">{dictionary.search.eyebrow}</p>
	<h2>{dictionary.search.title}</h2>
	<form on:submit|preventDefault={onSubmit}>
		<label>
			<span>{dictionary.search.sourceLabel}</span>
			<input bind:value={state.source} placeholder={bootstrap.search.sourcePlaceholder} />
		</label>

		<label>
			<span>{dictionary.search.targetLabel}</span>
			<input bind:value={state.target} placeholder={bootstrap.search.targetPlaceholder} />
		</label>

		<label>
			<span>{dictionary.search.depthLabel}</span>
			<input
				bind:value={state.maxDepth}
				min={bootstrap.controls.maxDepthMin}
				max={bootstrap.controls.maxDepthMax}
				type="number"
			/>
		</label>

		<button disabled={loading} type="submit">
			{#if loading}{dictionary.search.searching}{:else}{dictionary.search.submit}{/if}
		</button>
	</form>
</section>

<style>
	.panel {
		background: linear-gradient(180deg, rgba(255, 255, 255, 0.88), rgba(255, 255, 255, 0.72));
		border: 1px solid rgba(150, 184, 214, 0.2);
		border-radius: 26px;
		padding: 1.5rem;
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.12);
		backdrop-filter: blur(14px);
	}

	.eyebrow {
		margin: 0 0 0.35rem;
		font-size: 0.76rem;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: #67a9eb;
	}

	h2 {
		margin-top: 0;
	}

	form {
		display: grid;
		gap: 1rem;
	}

	label {
		display: grid;
		gap: 0.45rem;
		font-weight: 600;
		color: #516b8f;
	}

	input {
		width: 100%;
		padding: 0.92rem 1rem;
		border-radius: 16px;
		border: 1px solid rgba(146, 182, 219, 0.26);
		background: rgba(255, 255, 255, 0.86);
		color: #18314f;
		outline: none;
		transition: border-color 160ms ease, box-shadow 160ms ease, transform 160ms ease;
	}

	input:focus {
		border-color: rgba(69, 168, 255, 0.6);
		box-shadow: 0 0 0 4px rgba(69, 168, 255, 0.12);
		transform: translateY(-1px);
	}

	button {
		border: none;
		border-radius: 16px;
		padding: 0.95rem 1.1rem;
		background: linear-gradient(135deg, #5ab8ff 0%, #458af4 100%);
		color: white;
		font-weight: 700;
		cursor: pointer;
		box-shadow: 0 14px 28px rgba(69, 138, 244, 0.24);
	}

	button:disabled {
		opacity: 0.74;
		cursor: wait;
	}

	.animate-in {
		opacity: 0;
		transform: translateY(14px);
		animation: rise 520ms ease forwards;
	}

	@keyframes rise {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
