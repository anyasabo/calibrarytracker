<script lang="ts">
	import { onMount } from 'svelte';
	import { getReferenceData, getReciprocatingSystemIds } from '$lib/data';
	import { getAllCards, getAllVisits, addVisit } from '$lib/db';
	import type { UserCard, BranchVisit, Branch } from '$lib/types';
	import { SvelteSet } from 'svelte/reactivity';

	const { branches } = getReferenceData();

	type FilterValue = 'all' | 'have-card' | 'reciprocal' | 'no-card';

	let cards = $state<UserCard[]>([]);
	let visits = $state<BranchVisit[]>([]);
	let filter = $state<FilterValue>('all');
	let searchQuery = $state('');
	let routeMode = $state(false);
	let routeStops = $state<Branch[]>([]);
	let mapReady = $state(false);
	let BranchMapComponent = $state<
		(typeof import('$lib/components/BranchMap.svelte'))['default'] | null
	>(null);

	const cardSystemIds = $derived(new Set(cards.map((c) => c.systemId)));

	const reciprocatingSystemIds = $derived(() => {
		const ids = new SvelteSet<string>();
		for (const sysId of cardSystemIds) {
			for (const rid of getReciprocatingSystemIds(sysId)) {
				ids.add(rid);
			}
		}
		return ids;
	});

	const visitedBranchIds = $derived(new Set(visits.map((v) => v.branchId)));

	const mappableBranches = $derived(branches.filter((b) => b.lat !== 0 || b.lng !== 0));

	$effect(() => {
		Promise.all([getAllCards(), getAllVisits()]).then(([c, v]) => {
			cards = c;
			visits = v;
			mapReady = true;
		});
	});

	onMount(async () => {
		const mod = await import('$lib/components/BranchMap.svelte');
		BranchMapComponent = mod.default;
	});

	function handleToggleRouteStop(branch: Branch): void {
		const idx = routeStops.findIndex((s) => s.id === branch.id);
		if (idx >= 0) {
			routeStops = routeStops.filter((s) => s.id !== branch.id);
		} else {
			if (routeStops.length >= 25) return;
			routeStops = [...routeStops, branch];
		}
	}

	async function handleMarkVisited(branch: Branch): Promise<void> {
		const visit: BranchVisit = {
			id: crypto.randomUUID(),
			branchId: branch.id,
			visitDate: new Date().toISOString().slice(0, 10),
			notes: ''
		};
		await addVisit(visit);
		visits = await getAllVisits();
	}

	function moveStop(idx: number, direction: -1 | 1): void {
		const newIdx = idx + direction;
		if (newIdx < 0 || newIdx >= routeStops.length) return;
		const copy = [...routeStops];
		[copy[idx], copy[newIdx]] = [copy[newIdx], copy[idx]];
		routeStops = copy;
	}

	function removeStop(idx: number): void {
		routeStops = routeStops.filter((_, i) => i !== idx);
	}

	function buildGoogleMapsUrl(): string {
		if (routeStops.length === 0) return '';
		if (routeStops.length === 1) {
			const s = routeStops[0];
			return `https://www.google.com/maps/dir/?api=1&destination=${s.lat},${s.lng}`;
		}
		const origin = routeStops[0];
		const destination = routeStops[routeStops.length - 1];
		const waypoints = routeStops
			.slice(1, -1)
			.map((s) => `${s.lat},${s.lng}`)
			.join('|');
		let url = `https://www.google.com/maps/dir/?api=1&origin=${origin.lat},${origin.lng}&destination=${destination.lat},${destination.lng}`;
		if (waypoints) url += `&waypoints=${waypoints}`;
		return url;
	}

	function buildAppleMapsUrl(): string {
		if (routeStops.length === 0) return '';
		const dest = routeStops[routeStops.length - 1];
		return `https://maps.apple.com/?daddr=${dest.lat},${dest.lng}`;
	}
</script>

<svelte:head>
	<title>Map | CA Library Tracker</title>
	<link
		rel="stylesheet"
		href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"
		integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
		crossorigin=""
	/>
	<link
		rel="stylesheet"
		href="https://unpkg.com/leaflet.markercluster@1.5.3/dist/MarkerCluster.css"
		crossorigin=""
	/>
	<link
		rel="stylesheet"
		href="https://unpkg.com/leaflet.markercluster@1.5.3/dist/MarkerCluster.Default.css"
		crossorigin=""
	/>
</svelte:head>

<div class="map-page">
	<div class="controls">
		<div class="controls-top">
			<input
				type="search"
				placeholder="Search branches..."
				bind:value={searchQuery}
				aria-label="Search branches"
			/>
			<select bind:value={filter} aria-label="Filter branches">
				<option value="all">All branches</option>
				<option value="have-card">Have card</option>
				<option value="reciprocal">Reciprocal access</option>
				<option value="no-card">No card</option>
			</select>
		</div>

		<div class="legend">
			<span class="legend-item"><span class="dot dot-green"></span> Have card</span>
			<span class="legend-item"><span class="dot dot-blue"></span> Reciprocal</span>
			<span class="legend-item"><span class="dot dot-gray"></span> No card</span>
			<span class="legend-count">{mappableBranches.length} branches</span>
		</div>

		<div class="route-toggle">
			<button
				class="btn-route"
				class:active={routeMode}
				onclick={() => {
					routeMode = !routeMode;
					if (!routeMode) routeStops = [];
				}}
			>
				{routeMode ? 'Exit route mode' : 'Plan route'}
			</button>
		</div>
	</div>

	{#if routeMode && routeStops.length > 0}
		<div class="route-panel">
			<h3>Route ({routeStops.length}/25 stops)</h3>
			<ol class="route-list">
				{#each routeStops as stop, idx (stop.id)}
					<li>
						<span class="stop-name">{stop.name}</span>
						<span class="stop-actions">
							<button onclick={() => moveStop(idx, -1)} disabled={idx === 0} aria-label="Move up"
								>&#8593;</button
							>
							<button
								onclick={() => moveStop(idx, 1)}
								disabled={idx === routeStops.length - 1}
								aria-label="Move down">&#8595;</button
							>
							<button onclick={() => removeStop(idx)} aria-label="Remove">&#215;</button>
						</span>
					</li>
				{/each}
			</ol>
			<div class="route-actions">
				<a href={buildGoogleMapsUrl()} target="_blank" rel="noopener" class="btn-nav">
					Open in Google Maps
				</a>
				<a href={buildAppleMapsUrl()} target="_blank" rel="noopener" class="btn-nav btn-nav-apple">
					Open in Apple Maps
				</a>
				<button class="btn-clear" onclick={() => (routeStops = [])}>Clear route</button>
			</div>
			{#if routeStops.length >= 25}
				<p class="route-warn">Google Maps supports up to 25 waypoints per route.</p>
			{/if}
		</div>
	{/if}

	<div class="map-wrapper">
		{#if mapReady && BranchMapComponent}
			<BranchMapComponent
				branches={mappableBranches}
				{cardSystemIds}
				reciprocatingSystemIds={reciprocatingSystemIds()}
				{visitedBranchIds}
				{filter}
				{searchQuery}
				{routeMode}
				{routeStops}
				onToggleRouteStop={handleToggleRouteStop}
				onMarkVisited={handleMarkVisited}
			/>
		{:else if mapReady}
			<div class="map-loading">Loading map...</div>
		{/if}
	</div>
</div>

<style>
	.map-page {
		display: flex;
		flex-direction: column;
		height: calc(100vh - var(--header-height));
		margin: calc(-1 * var(--space-xl)) calc(-1 * var(--space-lg));
		position: relative;
	}

	.controls {
		position: absolute;
		top: var(--space-md);
		left: var(--space-md);
		z-index: 1000;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-sm);
		box-shadow: var(--shadow-md);
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		max-width: 320px;
	}

	.controls-top {
		display: flex;
		gap: var(--space-xs);
	}

	.controls input[type='search'] {
		flex: 1;
		min-width: 0;
		padding: var(--space-xs) var(--space-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		font-size: 0.85rem;
		outline: none;
	}

	.controls input[type='search']:focus {
		border-color: var(--color-border-focus);
	}

	.controls select {
		padding: var(--space-xs) var(--space-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		font-size: 0.85rem;
		background: white;
		cursor: pointer;
	}

	.legend {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-sm);
		align-items: center;
		font-size: 0.75rem;
		color: var(--color-text-muted);
	}

	.legend-item {
		display: flex;
		align-items: center;
		gap: 3px;
	}

	.legend-count {
		margin-left: auto;
		font-weight: 500;
	}

	.dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		display: inline-block;
	}

	.dot-green {
		background: #22c55e;
	}
	.dot-blue {
		background: #3b82f6;
	}
	.dot-gray {
		background: #9ca3af;
	}

	.route-toggle {
		display: flex;
	}

	.btn-route {
		width: 100%;
		padding: var(--space-xs) var(--space-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		background: var(--color-bg-card);
		font-size: 0.85rem;
		font-weight: 500;
		cursor: pointer;
	}

	.btn-route:hover {
		border-color: var(--color-primary);
	}

	.btn-route.active {
		background: var(--color-primary);
		color: white;
		border-color: var(--color-primary);
	}

	.route-panel {
		position: absolute;
		top: var(--space-md);
		right: var(--space-md);
		z-index: 1000;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-md);
		box-shadow: var(--shadow-md);
		max-width: 280px;
		max-height: calc(100vh - var(--header-height) - 2rem);
		overflow-y: auto;
	}

	.route-panel h3 {
		font-size: 0.9rem;
		color: var(--color-primary);
		margin-bottom: var(--space-sm);
	}

	.route-list {
		list-style: none;
		padding: 0;
		margin: 0;
	}

	.route-list li {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-xs);
		padding: var(--space-xs) 0;
		border-bottom: 1px solid var(--color-border);
		font-size: 0.8rem;
	}

	.route-list li:last-child {
		border-bottom: none;
	}

	.stop-name {
		flex: 1;
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.stop-actions {
		display: flex;
		gap: 2px;
	}

	.stop-actions button {
		background: none;
		border: 1px solid var(--color-border);
		border-radius: 3px;
		width: 22px;
		height: 22px;
		font-size: 0.75rem;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0;
	}

	.stop-actions button:hover {
		background: var(--color-bg);
	}
	.stop-actions button:disabled {
		opacity: 0.3;
		cursor: default;
	}

	.route-actions {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		margin-top: var(--space-sm);
	}

	.btn-nav {
		display: block;
		text-align: center;
		padding: var(--space-xs) var(--space-sm);
		background: var(--color-primary);
		color: white;
		border-radius: var(--radius-sm);
		font-size: 0.85rem;
		font-weight: 500;
	}

	.btn-nav:hover {
		background: var(--color-primary-light);
		text-decoration: none;
	}

	.btn-nav-apple {
		background: #1d1d1f;
	}
	.btn-nav-apple:hover {
		background: #333;
	}

	.btn-clear {
		background: none;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		padding: var(--space-xs);
		font-size: 0.8rem;
		cursor: pointer;
		color: var(--color-text-muted);
	}

	.btn-clear:hover {
		border-color: var(--color-danger);
		color: var(--color-danger);
	}

	.route-warn {
		font-size: 0.75rem;
		color: var(--color-warning);
		margin-top: var(--space-xs);
	}

	.map-wrapper {
		flex: 1;
		min-height: 0;
	}

	.map-loading {
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-muted);
		font-size: 0.95rem;
	}

	@media (max-width: 640px) {
		.controls {
			left: var(--space-xs);
			right: var(--space-xs);
			max-width: none;
		}

		.route-panel {
			right: var(--space-xs);
			left: var(--space-xs);
			max-width: none;
			top: auto;
			bottom: var(--space-xs);
			max-height: 40vh;
		}

		.controls-top {
			flex-direction: column;
		}
	}
</style>
