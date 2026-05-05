export type User = {
	login: string;
	name: string;
	avatarUrl?: string;
	profileUrl?: string;
	bio?: string;
};

export type GraphNode = {
	id: string;
	label: string;
	avatarUrl?: string;
	profileUrl?: string;
};

export type GraphEdge = {
	source: string;
	target: string;
	label?: string;
};

export type RelationLegendItem = {
	key: string;
	label: string;
	color: string;
};

export type PathStep = {
	from: string;
	to: string;
	relation: string;
};

export type PathResponse = {
	source: string;
	target: string;
	degree: number;
	meta: {
		mode: string;
		reason?: string;
	};
	nodes: GraphNode[];
	edges: GraphEdge[];
	steps: PathStep[];
};

export type ShortestPathRequest = {
	source: string;
	target: string;
	maxDepth?: number;
};

export type BootstrapResponse = {
	search: {
		defaultSource: string;
		defaultTarget: string;
		sourcePlaceholder: string;
		targetPlaceholder: string;
		defaultMaxDepth: number;
	};
	controls: {
		maxDepthMin: number;
		maxDepthMax: number;
		defaultRepulsion: number;
		repulsionMin: number;
		repulsionMax: number;
		defaultAttraction: number;
		attractionMin: number;
		attractionMax: number;
		defaultNodeSize: number;
		nodeSizeMin: number;
		nodeSizeMax: number;
	};
	graph: {
		initialZoom: number;
		minZoomRatio: number;
		maxZoomRatio: number;
		layoutIterations: number;
		playbackInterval: number;
		statusCardMaxWidthRem: number;
		githubCapability: {
			mode: string;
			description: string;
		};
		relationLegend: RelationLegendItem[];
		enablePan: boolean;
		enableZoom: boolean;
		enableAnimations: boolean;
		enableSidebar: boolean;
	};
};

export type SearchState = {
	source: string;
	target: string;
	maxDepth: number;
};

export type ForceLayoutControls = {
	maxDepth: number;
	repulsion: number;
	attraction: number;
	nodeSize: number;
};

export type ForceLayoutControlLimits = {
	maxDepthMin: number;
	maxDepthMax: number;
	repulsionMin: number;
	repulsionMax: number;
	attractionMin: number;
	attractionMax: number;
	nodeSizeMin: number;
	nodeSizeMax: number;
};

export type GraphRuntime = {
	initialZoom: number;
	minZoomRatio: number;
	maxZoomRatio: number;
	layoutIterations: number;
	playbackInterval: number;
	statusCardMaxWidthRem: number;
	githubCapability: {
		mode: string;
		description: string;
	};
	relationLegend: RelationLegendItem[];
	enablePan: boolean;
	enableZoom: boolean;
	enableAnimations: boolean;
	enableSidebar: boolean;
};

export type SelectedPathStep = {
	index: number;
	from: string;
	to: string;
	relation: string;
};
