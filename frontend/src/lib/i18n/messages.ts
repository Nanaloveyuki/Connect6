export const messages = {
	'en-US': {
		page: {
			title: 'Connect6',
			description: 'Explore GitHub relationship paths between two developers.'
		},
		topbar: {
			brand: 'Connect6',
			ready: 'Ready',
			busy: 'Working',
			degraded: 'Degraded',
			offline: 'Offline',
			closePanel: 'Close panel',
			tabs: {
				search: 'Search',
				path: 'Path',
				details: 'Details',
				settings: 'Settings',
				diagnostics: 'Diagnostics'
			}
		},
		toast: {
			bootstrapErrorTitle: 'Bootstrap unavailable',
			githubCapabilityTitle: 'GitHub live capability',
			searchErrorTitle: 'Path request failed',
			userErrorTitle: 'Profile unavailable',
			liveSourceTitle: 'Live GitHub data',
			seedSourceTitle: 'Fallback seed data',
			seedUnauthorizedBody: 'GitHub access is unauthorized, so local fallback data is shown.',
			seedRateLimitedBody: 'GitHub requests are rate-limited, so local fallback data is shown.',
			directSourceTitle: 'Direct match',
			seedSourceBody: 'GitHub relationships were unavailable, so local fallback data is shown.',
			directSourceBody: 'The source and target refer to the same GitHub user.'
		},
		search: {
			eyebrow: 'Search',
			title: 'Path query',
			sourceLabel: 'Source GitHub user',
			targetLabel: 'Target GitHub user',
			depthLabel: 'Maximum degree',
			submit: 'Search path',
			searching: 'Searching...'
		},
		controls: {
			eyebrow: 'Layout',
			title: 'Graph controls',
			maxDepth: 'Max depth',
			repulsion: 'Repulsion',
			attraction: 'Attraction',
			nodeSize: 'Node size',
			githubCapability: {
				authenticated: 'Authenticated',
				limited: 'Limited'
			}
		},
		diagnostics: {
			title: 'Diagnostics',
			appStatus: 'App status',
			githubCapability: 'GitHub capability',
			lastPathSource: 'Last path source',
			lastReason: 'Last degrade reason',
			emptyReason: 'No degrade reason recorded.',
			reasons: {
				githubFollowing: 'Live GitHub following graph',
				githubUnauthorized: 'GitHub access unauthorized',
				githubRateLimited: 'GitHub requests rate-limited',
				githubProviderUnavailable: 'GitHub relationship provider unavailable',
				sameUser: 'Source and target are the same user',
				localDevelopment: 'Local development fallback data'
			}
		},
		graph: {
			animationsOn: 'motion on',
			animationsOff: 'motion off',
			zoomEnabled: 'zoom enabled',
			zoomDisabled: 'zoom disabled',
			loading: 'Calculating path...',
			empty: 'Open Search to load a graph.',
			nodesLabel: 'Nodes',
			edgesLabel: 'Edges',
			degreeLabel: 'Degree',
			sourceLabel: 'Source',
			relationLabel: 'Relation',
			hoverLabel: 'Hover',
			focusLabel: 'Focus',
			focusHint: 'Selected node focused',
			stepHint: 'Selected path segment',
			sourceModes: {
				live: 'Live GitHub',
				seed: 'Fallback seed',
				direct: 'Direct match'
			}
		},
		sidebar: {
			eyebrow: 'Details',
			emptyTitle: 'No active node',
			emptyBody: 'Select a result to inspect profile details.',
			degree: 'Degree',
			nodes: 'Nodes',
			selected: 'Selected node',
			source: 'Path source',
			steps: 'Path steps',
			play: 'Play',
			pause: 'Pause',
			restart: 'Restart',
			playing: 'Playing path',
			emptySteps: 'No path steps loaded.',
			stepProgress: 'Step'
		}
	}
} as const;

export type Locale = keyof typeof messages;
export type I18nDictionary = (typeof messages)['en-US'];
