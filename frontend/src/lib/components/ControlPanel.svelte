<script lang="ts">
	import type { ForceLayoutControlLimits, ForceLayoutControls } from '$lib/types/connect6';
	import type { I18nDictionary } from '$lib/i18n/messages';

	export let controls: ForceLayoutControls;
	export let limits: ForceLayoutControlLimits;
	export let dictionary: I18nDictionary;
</script>

<section class="panel animate-in panel-delay-2">
	<div class="panel-head">
		<p class="eyebrow">{dictionary.controls.eyebrow}</p>
		<h2>{dictionary.controls.title}</h2>
	</div>

	<div class="control-stack">
		<label class="metric">
			<div class="metric-head">
				<span>{dictionary.controls.maxDepth}</span>
				<strong>{controls.maxDepth}</strong>
			</div>
			<input bind:value={controls.maxDepth} min={limits.maxDepthMin} max={limits.maxDepthMax} type="range" />
		</label>
		<label class="metric">
			<div class="metric-head">
				<span>{dictionary.controls.repulsion}</span>
				<strong>{controls.repulsion}</strong>
			</div>
			<input bind:value={controls.repulsion} min={limits.repulsionMin} max={limits.repulsionMax} type="range" />
		</label>
		<label class="metric">
			<div class="metric-head">
				<span>{dictionary.controls.attraction}</span>
				<strong>{controls.attraction}</strong>
			</div>
			<input bind:value={controls.attraction} min={limits.attractionMin} max={limits.attractionMax} type="range" />
		</label>
		<label class="metric">
			<div class="metric-head">
				<span>{dictionary.controls.nodeSize}</span>
				<strong>{controls.nodeSize}</strong>
			</div>
			<input bind:value={controls.nodeSize} min={limits.nodeSizeMin} max={limits.nodeSizeMax} type="range" />
		</label>
	</div>
</section>

<style>
	.panel {
		background: linear-gradient(180deg, rgba(246, 251, 255, 0.92), rgba(234, 243, 252, 0.84));
		color: #1b3655;
		border: 1px solid rgba(150, 184, 214, 0.2);
		border-radius: 26px;
		padding: 1.35rem;
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.12);
		backdrop-filter: blur(14px);
	}

	.panel-head h2,
	.panel-head p {
		margin-top: 0;
	}

	.eyebrow {
		margin: 0 0 0.35rem;
		font-size: 0.78rem;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: #67a9eb;
	}

	.control-stack {
		display: grid;
		gap: 0.85rem;
	}

	.metric {
		display: grid;
		gap: 0.75rem;
		padding: 0.9rem 1rem;
		border-radius: 16px;
		background: rgba(255, 255, 255, 0.72);
		border: 1px solid rgba(150, 184, 214, 0.14);
	}

	.metric-head {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
	}

	.metric span {
		color: #6f86a4;
	}

	.metric strong {
		font-size: 1.05rem;
	}

	input[type='range'] {
		width: 100%;
		accent-color: #4f95ff;
	}

	.animate-in {
		opacity: 0;
		transform: translateY(14px);
		animation: rise 520ms ease forwards;
	}

	.panel-delay-2 {
		animation-delay: 140ms;
	}

	@keyframes rise {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
