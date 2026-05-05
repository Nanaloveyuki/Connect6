<script lang="ts">
	import type { I18nDictionary } from '$lib/i18n/messages';
	import type { PathResponse, User } from '$lib/types/connect6';

	export let result: PathResponse | null = null;
	export let user: User | null = null;
	export let dictionary: I18nDictionary;

	function getSourceModeLabel(mode: string) {
		switch (mode) {
			case 'live':
				return dictionary.graph.sourceModes.live;
			case 'seed':
				return dictionary.graph.sourceModes.seed;
			case 'direct':
				return dictionary.graph.sourceModes.direct;
			default:
				return mode;
		}
	}
</script>

<aside class="sidebar animate-in panel-delay-3">
	<p class="eyebrow">{dictionary.sidebar.eyebrow}</p>
	{#if user}
		<p class="selected-tag">{dictionary.sidebar.selected}</p>
		<h2>{user.name || user.login}</h2>
		<p class="login">@{user.login}</p>
		{#if user.bio}
			<p class="bio">{user.bio}</p>
		{/if}
	{:else}
		<h2>{dictionary.sidebar.emptyTitle}</h2>
		<p>{dictionary.sidebar.emptyBody}</p>
	{/if}

	{#if result}
		<div class="meta-grid">
			<div>
				<span>{dictionary.sidebar.source}</span>
				<strong>{getSourceModeLabel(result.meta.mode)}</strong>
			</div>
			<div>
				<span>{dictionary.sidebar.degree}</span>
				<strong>{result.degree}</strong>
			</div>
			<div>
				<span>{dictionary.sidebar.nodes}</span>
				<strong>{result.nodes.length}</strong>
			</div>
		</div>
	{/if}
</aside>

<style>
	.sidebar {
		background: linear-gradient(180deg, rgba(255, 255, 255, 0.84), rgba(247, 251, 255, 0.72));
		backdrop-filter: blur(14px);
		border-radius: 26px;
		padding: 1.25rem;
		border: 1px solid rgba(150, 184, 214, 0.2);
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.12);
	}

	.eyebrow {
		margin: 0 0 0.35rem;
		font-size: 0.76rem;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: #67a9eb;
	}

	h2,
	p {
		margin-top: 0;
	}

	.login,
	.bio {
		color: #6f86a4;
	}

	.selected-tag {
		margin-bottom: 0.4rem;
		font-size: 0.78rem;
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: #4b94f0;
	}

	.meta-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 0.8rem;
		margin-top: 1rem;
	}

	.meta-grid div {
		padding: 0.85rem 0.9rem;
		border-radius: 16px;
		background: rgba(255, 255, 255, 0.74);
		border: 1px solid rgba(150, 184, 214, 0.14);
	}

	.meta-grid span {
		display: block;
		font-size: 0.82rem;
		color: #5f728e;
	}

	.meta-grid strong {
		display: block;
		margin-top: 0.2rem;
		font-size: 1.15rem;
	}

	.animate-in {
		opacity: 0;
		transform: translateY(14px);
		animation: rise 520ms ease forwards;
	}

	.panel-delay-3 {
		animation-delay: 210ms;
	}

	@keyframes rise {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
