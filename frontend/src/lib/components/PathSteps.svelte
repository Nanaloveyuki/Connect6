<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	import type { I18nDictionary } from '$lib/i18n/messages';
	import type { PathResponse, SelectedPathStep } from '$lib/types/connect6';

	export let result: PathResponse | null = null;
	export let selectedStepIndex: number | null = null;
	export let isPlaying = false;
	export let dictionary: I18nDictionary;

	const dispatch = createEventDispatcher<{
		stepSelected: SelectedPathStep;
		playToggle: void;
		restart: void;
	}>();

	function selectStep(index: number, from: string, to: string, relation: string) {
		dispatch('stepSelected', { index, from, to, relation });
	}
</script>

<section class="steps-card animate-in panel-delay-4">
	<div class="steps-head">
		<p class="eyebrow">{dictionary.sidebar.steps}</p>
		{#if result && result.steps.length > 0}
			<div class="playback-tools">
				<button type="button" on:click={() => dispatch('playToggle')}>{isPlaying ? dictionary.sidebar.pause : dictionary.sidebar.play}</button>
				<button type="button" on:click={() => dispatch('restart')}>{dictionary.sidebar.restart}</button>
			</div>
		{/if}
	</div>

	{#if isPlaying}
		<p class="playback-note">{dictionary.sidebar.playing}</p>
	{/if}

	{#if result && result.steps.length > 0}
		<ol>
			{#each result.steps as step, index}
				<li>
					<button
						class="step-button"
						class:selected={selectedStepIndex === index}
						type="button"
						on:click={() => selectStep(index, step.from, step.to, step.relation)}
					>
						<strong>{step.from}</strong>
						<span>{step.relation}</span>
						<strong>{step.to}</strong>
					</button>
				</li>
			{/each}
		</ol>
	{:else}
		<p class="empty">{dictionary.sidebar.emptySteps}</p>
	{/if}
</section>

<style>
	.steps-card {
		background: linear-gradient(180deg, rgba(255, 255, 255, 0.84), rgba(247, 251, 255, 0.72));
		backdrop-filter: blur(14px);
		border-radius: 26px;
		padding: 1.25rem;
		border: 1px solid rgba(150, 184, 214, 0.2);
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.12);
	}

	.eyebrow {
		margin: 0 0 0.5rem;
		font-size: 0.76rem;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: #67a9eb;
	}

	.steps-head {
		display: flex;
		justify-content: space-between;
		gap: 0.75rem;
		align-items: center;
	}

	.playback-tools {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.playback-tools button {
		padding: 0.42rem 0.75rem;
		border-radius: 999px;
		border: 1px solid rgba(150, 184, 214, 0.24);
		background: rgba(255, 255, 255, 0.72);
		cursor: pointer;
	}

	.playback-note {
		margin: 0 0 0.7rem;
		font-size: 0.82rem;
		color: #7e93ae;
	}

	ol {
		margin: 0;
		padding-left: 1.15rem;
		display: grid;
		gap: 0.75rem;
	}

	li {
		color: #5f728e;
		list-style: decimal;
	}

	.step-button {
		width: 100%;
		text-align: left;
		padding: 0.7rem 0.8rem;
		border-radius: 14px;
		border: 1px solid transparent;
		background: rgba(255, 255, 255, 0.58);
		color: inherit;
		cursor: pointer;
		transition: border-color 150ms ease, background 150ms ease, transform 150ms ease;
	}

	.step-button:hover {
		border-color: rgba(94, 157, 224, 0.24);
		background: rgba(255, 255, 255, 0.82);
		transform: translateY(-1px);
	}

	.step-button.selected {
		border-color: rgba(47, 125, 246, 0.28);
		background: rgba(237, 246, 255, 0.92);
	}

	li strong {
		color: #18314f;
	}

	li span {
		margin: 0 0.4rem;
	}

	.empty {
		margin: 0;
		color: #7e93ae;
	}

	.animate-in {
		opacity: 0;
		transform: translateY(14px);
		animation: rise 520ms ease forwards;
	}

	.panel-delay-4 {
		animation-delay: 260ms;
	}

	@keyframes rise {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@media (max-width: 860px) {
		.steps-head {
			align-items: start;
			flex-direction: column;
		}
	}
</style>
