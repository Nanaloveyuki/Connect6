<script lang="ts">
	import { onDestroy, onMount } from 'svelte';

	import { findShortestPath, getBootstrap, getUser } from '$lib/api/client';
	import ControlPanel from '$lib/components/ControlPanel.svelte';
	import GraphStage from '$lib/components/GraphStage.svelte';
	import PathSteps from '$lib/components/PathSteps.svelte';
	import SearchPanel from '$lib/components/SearchPanel.svelte';
	import SidebarCard from '$lib/components/SidebarCard.svelte';
	import { messages } from '$lib/i18n/messages';
	import { toastStore } from '$lib/stores/toast';
	import type {
		BootstrapResponse,
		ForceLayoutControlLimits,
		ForceLayoutControls,
		GraphRuntime,
		PathResponse,
		SearchState,
		SelectedPathStep,
		User
	} from '$lib/types/connect6';

	type TopTab = 'search' | 'path' | 'details' | 'settings' | 'diagnostics';
	type AppStatus = 'ready' | 'busy' | 'degraded' | 'error';

	const locale = 'en-US';
	const dictionary = messages[locale];

	const bootstrapFallback: BootstrapResponse = {
		search: {
			defaultSource: '',
			defaultTarget: '',
			sourcePlaceholder: 'github-user-a',
			targetPlaceholder: 'github-user-b',
			defaultMaxDepth: 4
		},
		controls: {
			maxDepthMin: 1,
			maxDepthMax: 6,
			defaultRepulsion: 36,
			repulsionMin: 10,
			repulsionMax: 80,
			defaultAttraction: 14,
			attractionMin: 5,
			attractionMax: 30,
			defaultNodeSize: 14,
			nodeSizeMin: 8,
			nodeSizeMax: 24
		},
		graph: {
			initialZoom: 1,
			minZoomRatio: 0.2,
			maxZoomRatio: 4,
			layoutIterations: 120,
			playbackInterval: 1200,
			statusCardMaxWidthRem: 32,
			githubCapability: {
				mode: 'limited',
				description: 'GitHub token missing. Live requests are available with stricter limits.'
			},
			relationLegend: [],
			enablePan: true,
			enableZoom: true,
			enableAnimations: true,
			enableSidebar: true
		}
	};

	let bootstrap = bootstrapFallback;
	let loading = false;
	let bootstrapping = true;
	let errorMessage = '';
	let result: PathResponse | null = null;
	let activeUser: User | null = null;
	let selectedStep: SelectedPathStep | null = null;
	let playbackTimer: ReturnType<typeof setInterval> | null = null;
	let isPlayingSteps = false;
	let activeTab: TopTab | null = 'search';
	let lastSourceToastKey = '';
	let hasShownCapabilityToast = false;
	let appStatus: AppStatus = 'busy';
	let searchState: SearchState = {
		source: '',
		target: '',
		maxDepth: bootstrapFallback.search.defaultMaxDepth
	};
	let controls: ForceLayoutControls = {
		maxDepth: bootstrapFallback.search.defaultMaxDepth,
		repulsion: bootstrapFallback.controls.defaultRepulsion,
		attraction: bootstrapFallback.controls.defaultAttraction,
		nodeSize: bootstrapFallback.controls.defaultNodeSize
	};
	let controlLimits: ForceLayoutControlLimits = {
		maxDepthMin: bootstrapFallback.controls.maxDepthMin,
		maxDepthMax: bootstrapFallback.controls.maxDepthMax,
		repulsionMin: bootstrapFallback.controls.repulsionMin,
		repulsionMax: bootstrapFallback.controls.repulsionMax,
		attractionMin: bootstrapFallback.controls.attractionMin,
		attractionMax: bootstrapFallback.controls.attractionMax,
		nodeSizeMin: bootstrapFallback.controls.nodeSizeMin,
		nodeSizeMax: bootstrapFallback.controls.nodeSizeMax
	};
	let runtime: GraphRuntime = bootstrapFallback.graph;

	const topTabs: Array<{ id: TopTab; label: string }> = [
		{ id: 'search', label: dictionary.topbar.tabs.search },
		{ id: 'path', label: dictionary.topbar.tabs.path },
		{ id: 'details', label: dictionary.topbar.tabs.details },
		{ id: 'settings', label: dictionary.topbar.tabs.settings },
		{ id: 'diagnostics', label: dictionary.topbar.tabs.diagnostics }
	];

	onMount(async () => {
		try {
			bootstrap = await getBootstrap();
			searchState = {
				source: bootstrap.search.defaultSource,
				target: bootstrap.search.defaultTarget,
				maxDepth: bootstrap.search.defaultMaxDepth
			};
			controls = {
				maxDepth: bootstrap.search.defaultMaxDepth,
				repulsion: bootstrap.controls.defaultRepulsion,
				attraction: bootstrap.controls.defaultAttraction,
				nodeSize: bootstrap.controls.defaultNodeSize
			};
			controlLimits = {
				maxDepthMin: bootstrap.controls.maxDepthMin,
				maxDepthMax: bootstrap.controls.maxDepthMax,
				repulsionMin: bootstrap.controls.repulsionMin,
				repulsionMax: bootstrap.controls.repulsionMax,
				attractionMin: bootstrap.controls.attractionMin,
				attractionMax: bootstrap.controls.attractionMax,
				nodeSizeMin: bootstrap.controls.nodeSizeMin,
				nodeSizeMax: bootstrap.controls.nodeSizeMax
			};
			runtime = bootstrap.graph;

			if (!hasShownCapabilityToast && bootstrap.graph.githubCapability.mode !== 'authenticated') {
				hasShownCapabilityToast = true;
				toastStore.warning(
					dictionary.toast.githubCapabilityTitle,
					bootstrap.graph.githubCapability.description,
					4200
				);
			}
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Failed to load bootstrap';
			toastStore.error(dictionary.toast.bootstrapErrorTitle, errorMessage);
		} finally {
			bootstrapping = false;
		}
	});

	onDestroy(() => {
		stopPlayback();
	});

	async function handleSearch() {
		loading = true;
		errorMessage = '';

		try {
			result = await findShortestPath({
				source: searchState.source,
				target: searchState.target,
				maxDepth: controls.maxDepth
			});

			activeUser = result.nodes[0] ? await getUser(result.nodes[0].id) : null;
			selectedStep = null;
			stopPlayback();
			activeTab = 'path';
		} catch (error) {
			result = null;
			activeUser = null;
			selectedStep = null;
			stopPlayback();
			errorMessage = error instanceof Error ? error.message : 'Unexpected request failure';
			toastStore.error(dictionary.toast.searchErrorTitle, errorMessage);
		} finally {
			loading = false;
		}
	}

	function handleStepSelected(event: CustomEvent<SelectedPathStep>) {
		stopPlayback();
		selectedStep = event.detail;
	}

	function stopPlayback() {
		if (playbackTimer) {
			clearInterval(playbackTimer);
			playbackTimer = null;
		}

		isPlayingSteps = false;
	}

	function handlePlaybackToggle() {
		if (!result || result.steps.length === 0) {
			return;
		}

		if (isPlayingSteps) {
			stopPlayback();
			return;
		}

		isPlayingSteps = true;
		let nextIndex = selectedStep?.index ?? 0;
		selectedStep = { index: nextIndex, ...result.steps[nextIndex] };

		playbackTimer = setInterval(() => {
			if (!result) {
				stopPlayback();
				return;
			}

			nextIndex += 1;
			if (nextIndex >= result.steps.length) {
				stopPlayback();
				return;
			}

			selectedStep = { index: nextIndex, ...result.steps[nextIndex] };
		}, runtime.playbackInterval);
	}

	function handlePlaybackRestart() {
		if (!result || result.steps.length === 0) {
			return;
		}

		selectedStep = { index: 0, ...result.steps[0] };
		stopPlayback();
	}

	async function handleNodeSelected(event: CustomEvent<{ nodeId: string }>) {
		stopPlayback();
		try {
			activeUser = await getUser(event.detail.nodeId);
			activeTab = 'details';
		} catch {
			activeUser = null;
			toastStore.warning(dictionary.toast.userErrorTitle);
		}
	}

	function notifySourceMode(path: PathResponse) {
		const key = `${path.meta.mode}:${path.meta.reason ?? ''}:${path.source}:${path.target}`;
		if (lastSourceToastKey === key) {
			return;
		}

		lastSourceToastKey = key;

		if (path.meta.mode === 'seed') {
			const description = path.meta.reason === 'github-unauthorized'
				? dictionary.toast.seedUnauthorizedBody
				: path.meta.reason === 'github-rate-limited'
					? dictionary.toast.seedRateLimitedBody
					: dictionary.toast.seedSourceBody;

			toastStore.warning(dictionary.toast.seedSourceTitle, description, 5200);
			return;
		}

		if (path.meta.mode === 'direct') {
			toastStore.info(dictionary.toast.directSourceTitle, dictionary.toast.directSourceBody, 3600);
			return;
		}

		if (path.meta.mode === 'live') {
			toastStore.success(dictionary.toast.liveSourceTitle, undefined, 2400);
		}
	}

	$: if (result) {
		notifySourceMode(result);
	}

	$: appStatus = getAppStatus();

	function toggleTab(tab: TopTab) {
		activeTab = activeTab === tab ? null : tab;
	}

	function getStatusLabel() {
		if (appStatus === 'busy') {
			return dictionary.topbar.busy;
		}

		if (appStatus === 'degraded') {
			return dictionary.topbar.degraded;
		}

		if (appStatus === 'error') {
			return dictionary.topbar.offline;
		}

		return dictionary.topbar.ready;
	}

	function getAppStatus(): AppStatus {
		if (loading || bootstrapping) {
			return 'busy';
		}

		if (errorMessage) {
			return 'error';
		}

		if (result?.meta.mode === 'seed') {
			return 'degraded';
		}

		return 'ready';
	}

	function getCapabilityLabel(mode: string) {
		return mode === 'authenticated'
			? dictionary.controls.githubCapability.authenticated
			: dictionary.controls.githubCapability.limited;
	}

	function getSourceModeLabel(mode?: string) {
		switch (mode) {
			case 'live':
				return dictionary.graph.sourceModes.live;
			case 'seed':
				return dictionary.graph.sourceModes.seed;
			case 'direct':
				return dictionary.graph.sourceModes.direct;
			default:
				return dictionary.graph.empty;
		}
	}

	function getReasonLabel(reason?: string) {
		switch (reason) {
			case 'github-following':
				return dictionary.diagnostics.reasons.githubFollowing;
			case 'github-unauthorized':
				return dictionary.diagnostics.reasons.githubUnauthorized;
			case 'github-rate-limited':
				return dictionary.diagnostics.reasons.githubRateLimited;
			case 'github-provider-unavailable':
				return dictionary.diagnostics.reasons.githubProviderUnavailable;
			case 'same-user':
				return dictionary.diagnostics.reasons.sameUser;
			case 'local-development':
				return dictionary.diagnostics.reasons.localDevelopment;
			default:
				return dictionary.diagnostics.emptyReason;
		}
	}
</script>

<svelte:head>
	<title>{dictionary.page.title}</title>
	<meta name="description" content={dictionary.page.description} />
</svelte:head>

<main class="page-shell">
	<section class="topbar animate-in">
		<div class="brand-block">
			<p class="brand-mark">{dictionary.topbar.brand}</p>
			<span class:busy={appStatus === 'busy'} class:degraded={appStatus === 'degraded'} class:error={appStatus === 'error'} class="status-pill">
				{getStatusLabel()}
			</span>
		</div>

		<nav class="tab-strip" aria-label="Primary panels">
			{#each topTabs as tab}
				<button
					class:active={activeTab === tab.id}
					type="button"
					on:click={() => toggleTab(tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</nav>
	</section>

	<section class="graph-layout">
		<GraphStage
			runtime={runtime}
			{controls}
			{result}
			{selectedStep}
			{activeUser}
			loading={loading || bootstrapping}
			errorMessage={errorMessage}
			{dictionary}
			on:nodeSelected={handleNodeSelected}
		/>

		<div class="panel-layer" aria-live="polite">
			<section class="tab-panel animate-in" data-open={activeTab === 'search'} hidden={activeTab !== 'search'}>
				<div class="panel-header">
					<h2>{dictionary.topbar.tabs.search}</h2>
					<button type="button" class="close-button" on:click={() => (activeTab = null)}>
						{dictionary.topbar.closePanel}
					</button>
				</div>
				<SearchPanel {bootstrap} state={searchState} {loading} {dictionary} onSubmit={handleSearch} />
			</section>

			<section class="tab-panel animate-in" data-open={activeTab === 'path'} hidden={activeTab !== 'path'}>
				<div class="panel-header">
					<h2>{dictionary.topbar.tabs.path}</h2>
					<button type="button" class="close-button" on:click={() => (activeTab = null)}>
						{dictionary.topbar.closePanel}
					</button>
				</div>
				{#if result && result.steps.length > 0}
					<div class="path-progress">
						<span>{dictionary.sidebar.stepProgress}</span>
						<strong>{(selectedStep?.index ?? 0) + 1} / {result.steps.length}</strong>
					</div>
				{/if}
				<PathSteps
					{result}
					selectedStepIndex={selectedStep?.index ?? null}
					isPlaying={isPlayingSteps}
					{dictionary}
					on:stepSelected={handleStepSelected}
					on:playToggle={handlePlaybackToggle}
					on:restart={handlePlaybackRestart}
				/>
			</section>

			<section class="tab-panel animate-in" data-open={activeTab === 'details'} hidden={activeTab !== 'details'}>
				<div class="panel-header">
					<h2>{dictionary.topbar.tabs.details}</h2>
					<button type="button" class="close-button" on:click={() => (activeTab = null)}>
						{dictionary.topbar.closePanel}
					</button>
				</div>
				<SidebarCard user={activeUser} {result} {dictionary} />
			</section>

			<section class="tab-panel animate-in" data-open={activeTab === 'settings'} hidden={activeTab !== 'settings'}>
				<div class="panel-header">
					<h2>{dictionary.topbar.tabs.settings}</h2>
					<button type="button" class="close-button" on:click={() => (activeTab = null)}>
						{dictionary.topbar.closePanel}
					</button>
				</div>
				<ControlPanel {controls} limits={controlLimits} {dictionary} />
				<div class="capability-card">
					<span class:limited={runtime.githubCapability.mode !== 'authenticated'} class="capability-pill">
						{getCapabilityLabel(runtime.githubCapability.mode)}
					</span>
					<p>{runtime.githubCapability.description}</p>
				</div>
				{#if runtime.relationLegend.length > 0}
					<div class="legend-card">
						{#each runtime.relationLegend as item}
							<div class="legend-row">
								<span class="legend-swatch" style={`--legend-color: ${item.color};`}></span>
								<div>
									<strong>{item.label}</strong>
									<p>{item.key}</p>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</section>

			<section class="tab-panel animate-in" data-open={activeTab === 'diagnostics'} hidden={activeTab !== 'diagnostics'}>
				<div class="panel-header">
					<h2>{dictionary.diagnostics.title}</h2>
					<button type="button" class="close-button" on:click={() => (activeTab = null)}>
						{dictionary.topbar.closePanel}
					</button>
				</div>
				<div class="diagnostics-card">
					<div class="diagnostic-row">
						<span>{dictionary.diagnostics.appStatus}</span>
						<strong>{getStatusLabel()}</strong>
					</div>
					<div class="diagnostic-row">
						<span>{dictionary.diagnostics.githubCapability}</span>
						<strong>{getCapabilityLabel(runtime.githubCapability.mode)}</strong>
					</div>
					<div class="diagnostic-row">
						<span>{dictionary.diagnostics.lastPathSource}</span>
						<strong>{getSourceModeLabel(result?.meta.mode)}</strong>
					</div>
					<div class="diagnostic-row diagnostic-row-stack">
						<span>{dictionary.diagnostics.lastReason}</span>
						<strong>{getReasonLabel(result?.meta.reason)}</strong>
					</div>
				</div>
			</section>
		</div>
	</section>
</main>

<style>
	.page-shell {
		min-height: 100vh;
		padding: 0.9rem;
		display: grid;
		grid-template-rows: auto 1fr;
		gap: 0.9rem;
	}

	.topbar {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
		padding: 0.8rem 1rem;
		border-radius: 28px;
		background: rgba(255, 255, 255, 0.56);
		border: 1px solid rgba(255, 255, 255, 0.62);
		box-shadow: 0 20px 60px rgba(114, 140, 173, 0.14);
		backdrop-filter: blur(18px);
	}

	.brand-block {
		display: flex;
		align-items: center;
		gap: 0.8rem;
	}

	.brand-mark {
		margin: 0;
		font-size: 1.05rem;
		font-weight: 800;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: #2f7df6;
	}

	.status-pill {
		padding: 0.4rem 0.7rem;
		border-radius: 999px;
		background: rgba(236, 250, 244, 0.82);
		border: 1px solid rgba(29, 187, 141, 0.18);
		color: #17936f;
		font-size: 0.78rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
	}

	.status-pill.busy {
		background: rgba(240, 246, 255, 0.88);
		border-color: rgba(69, 138, 244, 0.18);
		color: #2f7df6;
	}

	.status-pill.degraded {
		background: rgba(255, 246, 237, 0.92);
		border-color: rgba(243, 156, 86, 0.2);
		color: #c67730;
	}

	.status-pill.error {
		background: rgba(255, 241, 243, 0.88);
		border-color: rgba(208, 91, 112, 0.18);
		color: #c04d64;
	}

	.tab-strip {
		display: flex;
		gap: 0.55rem;
		flex-wrap: wrap;
	}

	.tab-strip button {
		padding: 0.6rem 0.9rem;
		border-radius: 999px;
		border: 1px solid rgba(150, 184, 214, 0.16);
		background: rgba(255, 255, 255, 0.54);
		color: #5c7698;
		cursor: pointer;
		transition: background 160ms ease, border-color 160ms ease, color 160ms ease, transform 160ms ease;
	}

	.tab-strip button:hover,
	.tab-strip button.active {
		background: rgba(255, 255, 255, 0.9);
		border-color: rgba(69, 138, 244, 0.2);
		color: #235ec7;
		transform: translateY(-1px);
	}

	.graph-layout {
		position: relative;
		min-height: 0;
	}

	.panel-layer {
		position: absolute;
		z-index: 5;
		top: 1.25rem;
		left: 1.25rem;
		width: min(26rem, calc(100vw - 2.5rem));
		pointer-events: none;
	}

	.tab-panel {
		pointer-events: auto;
		display: grid;
		gap: 0.85rem;
	}

	.panel-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.75rem;
		padding: 0.3rem 0.2rem 0 0.2rem;
	}

	.panel-header h2 {
		margin: 0;
		font-size: 1rem;
		color: #2c4c74;
	}

	.close-button {
		padding: 0.45rem 0.7rem;
		border-radius: 999px;
		border: 1px solid rgba(150, 184, 214, 0.16);
		background: rgba(255, 255, 255, 0.72);
		color: #5f728e;
		cursor: pointer;
	}

	.path-progress {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.85rem 1rem;
		border-radius: 18px;
		background: rgba(255, 255, 255, 0.74);
		border: 1px solid rgba(150, 184, 214, 0.18);
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.1);
		backdrop-filter: blur(14px);
		color: #546f90;
	}

	.path-progress strong {
		color: #1f416a;
	}

	.capability-card {
		display: grid;
		gap: 0.65rem;
		padding: 1rem;
		border-radius: 18px;
		background: rgba(255, 255, 255, 0.74);
		border: 1px solid rgba(150, 184, 214, 0.18);
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.1);
		backdrop-filter: blur(14px);
	}

	.capability-card p {
		margin: 0;
		color: #6f86a4;
	}

	.capability-pill {
		justify-self: start;
		padding: 0.38rem 0.68rem;
		border-radius: 999px;
		background: rgba(236, 250, 244, 0.82);
		border: 1px solid rgba(29, 187, 141, 0.18);
		color: #17936f;
		font-size: 0.78rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
	}

	.capability-pill.limited {
		background: rgba(255, 246, 237, 0.92);
		border-color: rgba(243, 156, 86, 0.2);
		color: #c67730;
	}

	.legend-card {
		display: grid;
		gap: 0.75rem;
		padding: 1rem;
		border-radius: 18px;
		background: rgba(255, 255, 255, 0.74);
		border: 1px solid rgba(150, 184, 214, 0.18);
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.1);
		backdrop-filter: blur(14px);
	}

	.legend-row {
		display: grid;
		grid-template-columns: auto 1fr;
		gap: 0.75rem;
		align-items: center;
	}

	.legend-row strong,
	.legend-row p {
		margin: 0;
	}

	.legend-row p {
		font-size: 0.82rem;
		color: #7187a3;
	}

	.legend-swatch {
		width: 0.85rem;
		height: 0.85rem;
		border-radius: 999px;
		background: var(--legend-color);
		box-shadow: 0 0 0 6px color-mix(in srgb, var(--legend-color) 18%, transparent);
	}

	.diagnostics-card {
		display: grid;
		gap: 0.75rem;
		padding: 1rem;
		border-radius: 18px;
		background: rgba(255, 255, 255, 0.74);
		border: 1px solid rgba(150, 184, 214, 0.18);
		box-shadow: 0 20px 44px rgba(116, 143, 178, 0.1);
		backdrop-filter: blur(14px);
	}

	.diagnostic-row {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
		align-items: baseline;
	}

	.diagnostic-row span,
	.diagnostic-row strong {
		margin: 0;
	}

	.diagnostic-row span {
		color: #6f86a4;
		font-size: 0.88rem;
	}

	.diagnostic-row strong {
		text-align: right;
		color: #21466f;
	}

	.diagnostic-row-stack {
		align-items: start;
		flex-direction: column;
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

	@media (max-width: 860px) {
		.page-shell {
			padding: 0.65rem;
		}

		.topbar {
			align-items: start;
			flex-direction: column;
		}

		.brand-block {
			width: 100%;
			justify-content: space-between;
		}

		.panel-layer {
			left: 0.8rem;
			right: 0.8rem;
			width: auto;
			top: 0.8rem;
		}
	}
</style>
