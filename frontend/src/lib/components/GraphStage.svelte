<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { browser } from '$app/environment';
	import { onDestroy } from 'svelte';

	import type { I18nDictionary } from '$lib/i18n/messages';
	import type { ForceLayoutControls, GraphRuntime, PathResponse, SelectedPathStep, User } from '$lib/types/connect6';

	export let runtime: GraphRuntime;
	export let controls: ForceLayoutControls;
	export let result: PathResponse | null = null;
	export let selectedStep: SelectedPathStep | null = null;
	export let activeUser: User | null = null;
	export let loading = false;
	export let errorMessage = '';
	export let dictionary: I18nDictionary;

	const dispatch = createEventDispatcher<{
		nodeSelected: { nodeId: string };
	}>();

	let container: HTMLDivElement | null = null;
	let sigmaInstance: InstanceType<typeof import('sigma').default> | null = null;
	let hoveredNodeLabel = '';
	let hoveredEdgeRelation = '';
	let focusedNodeLabel = '';
	let selectedStepLabel = '';
	let selectedStepRelation = '';
	let selectedEdgeRelation = '';
	let sigmaModulesPromise: Promise<{
		Graph: typeof import('graphology').default;
		forceAtlas2: typeof import('graphology-layout-forceatlas2').default;
		Sigma: typeof import('sigma').default;
	}> | null = null;

	function loadSigmaModules() {
		if (!sigmaModulesPromise) {
			sigmaModulesPromise = Promise.all([
				import('graphology'),
				import('graphology-layout-forceatlas2'),
				import('sigma')
			]).then(([graphologyModule, forceAtlas2Module, sigmaModule]) => ({
				Graph: graphologyModule.default,
				forceAtlas2: forceAtlas2Module.default,
				Sigma: sigmaModule.default
			}));
		}

		return sigmaModulesPromise;
	}

	function destroyGraph() {
		sigmaInstance?.kill();
		sigmaInstance = null;
		hoveredNodeLabel = '';
		hoveredEdgeRelation = '';
		focusedNodeLabel = '';
		selectedStepLabel = '';
		selectedStepRelation = '';
		selectedEdgeRelation = '';
		if (container) {
			container.innerHTML = '';
		}
	}

	function getRelationColor(relation: string, isPathEdge: boolean) {
		const legendItem = runtime.relationLegend.find((item) => item.key === relation);
		if (legendItem) {
			return legendItem.color;
		}

		return isPathEdge ? '#5e9de0' : '#8ea6c5';
	}

	function getRelationLabel(relation: string) {
		const legendItem = runtime.relationLegend.find((item) => item.key === relation);
		return legendItem?.label ?? relation;
	}

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

	async function renderGraph() {
		if (!browser || !container || !result || loading || errorMessage) {
			destroyGraph();
			return;
		}

		destroyGraph();

		const { Graph, forceAtlas2, Sigma } = await loadSigmaModules();

		const graph = new Graph();
		const nodeCount = Math.max(result.nodes.length, 1);

		result.nodes.forEach((node, index) => {
			const angle = (Math.PI * 2 * index) / nodeCount;
			const distance = 10 + controls.repulsion * 0.15;

			graph.addNode(node.id, {
				label: node.label,
				originalLabel: node.label,
				size: controls.nodeSize,
				color: index === 0 || index === result.nodes.length - 1 ? '#1f67ff' : '#18a999',
				originalColor: index === 0 || index === result.nodes.length - 1 ? '#1f67ff' : '#18a999',
				x: Math.cos(angle) * distance,
				y: Math.sin(angle) * distance
			});
		});

		result.edges.forEach((edge, index) => {
			const key = `${edge.source}-${edge.target}-${index}`;
			const pathIndex = result.steps.findIndex((step) => step.from === edge.source && step.to === edge.target);
			const relation = edge.label ?? '';
			const edgeColor = getRelationColor(relation, pathIndex >= 0);
			graph.addEdgeWithKey(key, edge.source, edge.target, {
				color: edgeColor,
				size: pathIndex >= 0 ? 2.8 : 1.4,
				originalColor: edgeColor,
				pathIndex,
				relation
			});
		});

		forceAtlas2.assign(graph, {
			iterations: runtime.enableAnimations ? runtime.layoutIterations : 1,
			settings: {
				gravity: Math.max(0.1, controls.attraction / 12),
				scalingRatio: Math.max(1, controls.repulsion / 6),
				strongGravityMode: false,
				slowDown: 1.5
			}
		});

		const instance = new Sigma(graph, container, {
			renderLabels: true,
			labelDensity: 0.08,
			labelGridCellSize: 90,
			defaultEdgeType: 'line',
			allowInvalidContainer: false,
			enableEdgeEvents: true,
			minCameraRatio: runtime.enableZoom ? runtime.minZoomRatio : runtime.initialZoom,
			maxCameraRatio: runtime.enableZoom ? runtime.maxZoomRatio : runtime.initialZoom
		});

		sigmaInstance = instance;

		instance.on('enterNode', ({ node }) => {
			hoveredNodeLabel = String(graph.getNodeAttribute(node, 'originalLabel') ?? node);
			graph.forEachNode((key) => {
				graph.setNodeAttribute(key, 'color', key === node ? '#0f6dff' : '#b8cadf');
			});
		});

		instance.on('leaveNode', () => {
			hoveredNodeLabel = '';
			graph.forEachNode((key) => {
				graph.setNodeAttribute(key, 'color', graph.getNodeAttribute(key, 'originalColor'));
			});
		});

		instance.on('enterEdge', ({ edge }) => {
			hoveredEdgeRelation = String(graph.getEdgeAttribute(edge, 'relation') ?? '');
			graph.setEdgeAttribute(edge, 'color', '#0f6dff');
			graph.setEdgeAttribute(edge, 'size', 4.2);
		});

		instance.on('leaveEdge', ({ edge }) => {
			hoveredEdgeRelation = '';
			const originalColor = String(graph.getEdgeAttribute(edge, 'originalColor'));
			const pathIndex = Number(graph.getEdgeAttribute(edge, 'pathIndex') ?? -1);
			graph.setEdgeAttribute(edge, 'color', originalColor);
			graph.setEdgeAttribute(edge, 'size', pathIndex >= 0 ? 2.8 : 1.1);
		});

		instance.on('clickEdge', ({ edge }) => {
			selectedEdgeRelation = String(graph.getEdgeAttribute(edge, 'relation') ?? '');
			focusedNodeLabel = '';
			selectedStepLabel = '';
			selectedStepRelation = '';
		});

		instance.on('clickNode', ({ node }) => {
			selectedEdgeRelation = '';
			selectedStepLabel = '';
			focusedNodeLabel = String(graph.getNodeAttribute(node, 'originalLabel') ?? node);
			graph.forEachNode((key) => {
				const originalColor = graph.getNodeAttribute(key, 'originalColor');
				graph.setNodeAttribute(key, 'color', key === node ? '#0f6dff' : originalColor);
				graph.setNodeAttribute(key, 'size', key === node ? controls.nodeSize + 5 : controls.nodeSize);
			});

			graph.forEachEdge((edgeKey, attributes, source, target) => {
				const isConnected = source === node || target === node;
				graph.setEdgeAttribute(edgeKey, 'color', isConnected ? '#2f7df6' : attributes.originalColor);
				graph.setEdgeAttribute(edgeKey, 'size', isConnected ? 3.4 : attributes.pathIndex >= 0 ? 2.8 : 1.4);
			});

			const camera = instance.getCamera();
			if (camera) {
				camera.animate(
					{
						x: Number(graph.getNodeAttribute(node, 'x')),
						y: Number(graph.getNodeAttribute(node, 'y')),
						ratio: 0.55
					},
					{ duration: 420 }
				);
			}

			dispatch('nodeSelected', { nodeId: node });
		});

		if (selectedStep) {
			selectedStepLabel = `${selectedStep.from} -> ${selectedStep.to}`;
			selectedStepRelation = selectedStep.relation;
			graph.forEachNode((key) => {
				const originalColor = graph.getNodeAttribute(key, 'originalColor');
				const isInStep = key === selectedStep.from || key === selectedStep.to;
				graph.setNodeAttribute(key, 'color', isInStep ? '#0f6dff' : '#c4d3e4');
				graph.setNodeAttribute(key, 'size', isInStep ? controls.nodeSize + 4 : controls.nodeSize - 1);
			});

			graph.forEachEdge((edgeKey, attributes, source, target) => {
				const isSelectedEdge = source === selectedStep.from && target === selectedStep.to;
				graph.setEdgeAttribute(edgeKey, 'color', isSelectedEdge ? '#0f6dff' : attributes.originalColor);
				graph.setEdgeAttribute(edgeKey, 'size', isSelectedEdge ? 4.2 : attributes.pathIndex >= 0 ? 2.8 : 1.1);
			});

			const stepSourceX = Number(graph.getNodeAttribute(selectedStep.from, 'x'));
			const stepSourceY = Number(graph.getNodeAttribute(selectedStep.from, 'y'));
			const stepTargetX = Number(graph.getNodeAttribute(selectedStep.to, 'x'));
			const stepTargetY = Number(graph.getNodeAttribute(selectedStep.to, 'y'));
			const camera = instance.getCamera();
			camera.animate(
				{
					x: (stepSourceX + stepTargetX) / 2,
					y: (stepSourceY + stepTargetY) / 2,
					ratio: 0.65
				},
				{ duration: 420 }
			);
		}

		if (!selectedStep) {
			selectedStepLabel = '';
			selectedStepRelation = '';
		}

	}

	$: void renderGraph();

	onDestroy(() => {
		destroyGraph();
	});
</script>

<section
	class="stage-shell animate-in panel-delay-1"
	style={`--graph-status-width-rem: ${runtime.statusCardMaxWidthRem};`}
>
	<div class:loading class="stage-canvas">
		<div class="grid-glow"></div>
		<div class="mesh"></div>
		<div bind:this={container} class="sigma-host"></div>

		{#if loading}
			<div class="state-card pulse">{dictionary.graph.loading}</div>
		{:else if errorMessage}
			<div class="state-card error">{errorMessage}</div>
		{:else if !result}
			<div class="state-card">{dictionary.graph.empty}</div>
		{/if}
	</div>

	<div class="graph-overlay graph-overlay-top">
		<div class="graph-badges">
			<span>{runtime.enableAnimations ? dictionary.graph.animationsOn : dictionary.graph.animationsOff}</span>
			<span>{runtime.enableZoom ? dictionary.graph.zoomEnabled : dictionary.graph.zoomDisabled}</span>
		</div>
		<div class="graph-status">
			{#if hoveredEdgeRelation}
				<p><strong>{dictionary.graph.relationLabel}</strong> {getRelationLabel(hoveredEdgeRelation)}</p>
			{:else if hoveredNodeLabel}
				<p><strong>{dictionary.graph.hoverLabel}</strong> {hoveredNodeLabel}</p>
			{:else if selectedEdgeRelation}
				<p><strong>{dictionary.graph.relationLabel}</strong> {getRelationLabel(selectedEdgeRelation)}</p>
			{:else if selectedStepLabel}
				<p><strong>{dictionary.graph.stepHint}</strong> {selectedStepLabel}</p>
			{:else if focusedNodeLabel}
				<p><strong>{dictionary.graph.focusLabel}</strong> {focusedNodeLabel}</p>
			{:else}
				<p>{dictionary.graph.empty}</p>
			{/if}
			{#if selectedStepRelation}
				<p class="graph-status-secondary"><strong>{dictionary.graph.relationLabel}</strong> {getRelationLabel(selectedStepRelation)}</p>
			{/if}
		</div>
	</div>

	<div class="graph-overlay graph-overlay-bottom">
		{#if result}
			<div class="graph-meta">
				<span>{dictionary.graph.sourceLabel}: {getSourceModeLabel(result.meta.mode)}</span>
				<span>{dictionary.graph.nodesLabel}: {result.nodes.length}</span>
				<span>{dictionary.graph.edgesLabel}: {result.edges.length}</span>
				<span>{dictionary.graph.degreeLabel}: {result.degree}</span>
			</div>
		{/if}
		{#if activeUser}
			<div class="active-user">
				<strong>{activeUser.name || activeUser.login}</strong>
				<span>{activeUser.login}</span>
			</div>
		{/if}
	</div>
</section>

<style>
	.stage-shell {
		position: relative;
		height: calc(100vh - 6.25rem);
		min-height: 620px;
	}

	.graph-badges {
		display: flex;
		gap: 0.6rem;
		flex-wrap: wrap;
	}

	.graph-badges span {
		padding: 0.4rem 0.7rem;
		border-radius: 999px;
		background: rgba(255, 255, 255, 0.82);
		font-size: 0.85rem;
		color: #4f6f97;
		border: 1px solid rgba(150, 184, 214, 0.22);
		backdrop-filter: blur(10px);
	}

	.stage-canvas {
		position: relative;
		height: 100%;
		min-height: 100%;
		border-radius: 30px;
		overflow: hidden;
		border: 1px solid rgba(160, 192, 223, 0.22);
		box-shadow: 0 26px 72px rgba(116, 143, 178, 0.16);
		background:
			radial-gradient(circle at 20% 20%, rgba(38, 135, 255, 0.24), transparent 28%),
			radial-gradient(circle at 80% 30%, rgba(26, 189, 180, 0.16), transparent 26%),
			linear-gradient(180deg, #f8fbff 0%, #e9f0fb 100%);
	}

	.grid-glow,
	.mesh,
	.sigma-host {
		position: absolute;
		inset: 0;
	}

	.grid-glow {
		background-image:
			linear-gradient(rgba(23, 60, 110, 0.06) 1px, transparent 1px),
			linear-gradient(90deg, rgba(23, 60, 110, 0.06) 1px, transparent 1px);
		background-size: 34px 34px;
		opacity: 0.65;
	}

	.mesh {
		background: radial-gradient(circle at center, rgba(255, 255, 255, 0.24), transparent 55%);
		animation: drift 8s ease-in-out infinite alternate;
	}

	.sigma-host {
		z-index: 0;
	}

	.state-card {
		position: relative;
		z-index: 1;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-height: 100%;
		width: 100%;
		padding: 2rem;
		font-size: 1.02rem;
		color: #29476f;
		background: linear-gradient(180deg, rgba(248, 251, 255, 0.62), rgba(233, 240, 251, 0.74));
	}

	.state-card.error {
		color: #9f2437;
	}

	.pulse {
		animation: pulse 1.4s ease-in-out infinite;
	}

	.graph-overlay {
		position: absolute;
		z-index: 2;
		left: 1.1rem;
		right: 1.1rem;
		display: flex;
		justify-content: space-between;
		gap: 1rem;
		pointer-events: none;
	}

	.graph-overlay-top {
		top: 1.1rem;
		align-items: start;
	}

	.graph-overlay-bottom {
		bottom: 1.1rem;
		align-items: end;
	}

	.graph-status,
	.graph-meta,
	.active-user {
		pointer-events: auto;
		background: rgba(255, 255, 255, 0.72);
		border: 1px solid rgba(150, 184, 214, 0.18);
		box-shadow: 0 14px 36px rgba(116, 143, 178, 0.12);
		backdrop-filter: blur(14px);
	}

	.graph-status {
		max-width: min(calc(var(--graph-status-width-rem) * 1rem), 58vw);
		padding: 0.9rem 1rem;
		border-radius: 20px;
		color: #476482;
	}

	.graph-status p {
		margin: 0;
	}

	.graph-status-secondary {
		margin-top: 0.25rem;
		font-size: 0.84rem;
		color: #6f86a4;
	}

	.graph-meta {
		display: grid;
		grid-auto-flow: column;
		gap: 0.75rem;
		padding: 0.78rem 0.95rem;
		border-radius: 18px;
	}

	.graph-meta span {
		white-space: nowrap;
		color: #536f91;
	}

	.active-user {
		display: grid;
		justify-items: end;
		padding: 0.78rem 0.95rem;
		border-radius: 18px;
	}

	.active-user span {
		font-size: 0.85rem;
		color: #6f86a4;
	}

	.animate-in {
		opacity: 0;
		transform: translateY(14px);
		animation: rise 520ms ease forwards;
	}

	.panel-delay-1 {
		animation-delay: 80ms;
	}

	.loading .state-card {
		filter: saturate(0.85);
	}

	@keyframes rise {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 0.65;
		}

		50% {
			opacity: 1;
		}
	}

	@keyframes drift {
		from {
			transform: translate3d(-1.5%, -1%, 0) scale(1);
		}
		to {
			transform: translate3d(1.5%, 1%, 0) scale(1.04);
		}
	}

	@media (max-width: 860px) {
		.stage-shell {
			height: calc(100vh - 7rem);
			min-height: 520px;
		}

		.graph-overlay,
		.graph-overlay-top,
		.graph-overlay-bottom {
			flex-direction: column;
			align-items: stretch;
		}

		.graph-status {
			max-width: none;
		}

		.graph-meta {
			grid-auto-flow: row;
			justify-items: start;
		}

		.active-user {
			justify-items: start;
		}
	}
</style>
