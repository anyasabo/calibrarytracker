<script lang="ts">
	import { onMount } from 'svelte';
	import type { Branch } from '$lib/types';
	import { getSystemById } from '$lib/data';

	type MarkerStatus = 'have-card' | 'reciprocal' | 'no-card';

	interface Props {
		branches: Branch[];
		cardSystemIds: Set<string>;
		reciprocatingSystemIds: Set<string>;
		visitedBranchIds: Set<string>;
		filter: 'all' | MarkerStatus;
		searchQuery: string;
		routeMode: boolean;
		routeStops: Branch[];
		onToggleRouteStop: (branch: Branch) => void;
		onMarkVisited: (branch: Branch) => void;
	}

	let {
		branches,
		cardSystemIds,
		reciprocatingSystemIds,
		visitedBranchIds,
		filter,
		searchQuery,
		routeMode,
		routeStops,
		onToggleRouteStop,
		onMarkVisited
	}: Props = $props();

	let mapContainer: HTMLDivElement;
	let map: import('leaflet').Map | undefined;
	let clusterGroup: import('leaflet.markercluster').MarkerClusterGroup | undefined;

	const CA_CENTER: [number, number] = [37.2, -119.5];
	const CA_ZOOM = 6;

	const COLORS: Record<MarkerStatus, string> = {
		'have-card': '#22c55e',
		reciprocal: '#3b82f6',
		'no-card': '#9ca3af'
	};

	function getStatus(branch: Branch): MarkerStatus {
		if (cardSystemIds.has(branch.systemId)) return 'have-card';
		if (reciprocatingSystemIds.has(branch.systemId)) return 'reciprocal';
		return 'no-card';
	}

	function matchesSearch(branch: Branch, query: string): boolean {
		if (!query) return true;
		const q = query.toLowerCase();
		const system = getSystemById(branch.systemId);
		return (
			branch.name.toLowerCase().includes(q) ||
			branch.city.toLowerCase().includes(q) ||
			branch.county.toLowerCase().includes(q) ||
			(system?.name.toLowerCase().includes(q) ?? false)
		);
	}

	function buildPopupHtml(branch: Branch): string {
		const system = getSystemById(branch.systemId);
		const status = getStatus(branch);
		const visited = visitedBranchIds.has(branch.id);
		const isRouteStop = routeStops.some((s) => s.id === branch.id);

		const statusLabel =
			status === 'have-card' ? 'You have a card' : status === 'reciprocal' ? 'Reciprocal access' : '';

		const branchStatusBadge =
			branch.status !== 'open'
				? `<span class="popup-badge popup-badge-${branch.status}">${branch.status}</span>`
				: '';

		const googleUrl = `https://www.google.com/maps/dir/?api=1&destination=${branch.lat},${branch.lng}`;
		const appleUrl = `https://maps.apple.com/?daddr=${branch.lat},${branch.lng}`;

		return `
			<div class="branch-popup">
				<strong>${branch.name}</strong>
				${branchStatusBadge}
				<div class="popup-type">${branch.outletType}</div>
				${system ? `<a href="/systems/${system.id}" class="popup-system">${system.name}</a>` : ''}
				<div class="popup-address">${branch.address}, ${branch.city} ${branch.zipCode}</div>
				${branch.phone ? `<div class="popup-phone">${branch.phone}</div>` : ''}
				${statusLabel ? `<div class="popup-card-status popup-card-${status}">${statusLabel}</div>` : ''}
				${visited ? '<div class="popup-visited">Visited</div>' : ''}
				<div class="popup-nav">
					<a href="${googleUrl}" target="_blank" rel="noopener">Google Maps</a>
					<a href="${appleUrl}" target="_blank" rel="noopener">Apple Maps</a>
				</div>
				${
					routeMode
						? `<button class="popup-btn" data-action="route" data-branch-id="${branch.id}">
							${isRouteStop ? 'Remove from route' : 'Add to route'}
						</button>`
						: ''
				}
				${
					!visited
						? `<button class="popup-btn popup-btn-visit" data-action="visit" data-branch-id="${branch.id}">
							Mark visited
						</button>`
						: ''
				}
			</div>
		`;
	}

	function getFilteredBranches(): Branch[] {
		return branches.filter((b) => {
			if (b.lat === 0 && b.lng === 0) return false;
			if (!matchesSearch(b, searchQuery)) return false;
			if (filter === 'all') return true;
			return getStatus(b) === filter;
		});
	}

	function updateMarkers(L: typeof import('leaflet')): void {
		if (!map || !clusterGroup) return;

		clusterGroup.clearLayers();

		const filtered = getFilteredBranches();
		const markers: import('leaflet').CircleMarker[] = [];

		for (const branch of filtered) {
			const status = getStatus(branch);
			const visited = visitedBranchIds.has(branch.id);

			const marker = L.circleMarker([branch.lat, branch.lng], {
				radius: 7,
				fillColor: COLORS[status],
				color: visited ? '#000' : COLORS[status],
				weight: visited ? 2 : 1,
				opacity: 1,
				fillOpacity: visited ? 1.0 : 0.7
			});

			marker.bindPopup(() => buildPopupHtml(branch), { maxWidth: 280 });

			markers.push(marker);
		}

		clusterGroup.addLayers(markers);
	}

	onMount(() => {
		initMap();
		return () => {
			map?.remove();
		};
	});

	async function initMap(): Promise<void> {
		const L = await import('leaflet');
		await import('leaflet.markercluster');

		map = L.map(mapContainer, {
			center: CA_CENTER,
			zoom: CA_ZOOM,
			zoomControl: true
		});

		L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
			attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
			maxZoom: 19
		}).addTo(map);

		// eslint-disable-next-line @typescript-eslint/no-explicit-any -- leaflet.markercluster augments L at runtime
		const cg = new (L as any).MarkerClusterGroup({
			maxClusterRadius: 50,
			chunkedLoading: true,
			showCoverageOnHover: false,
			iconCreateFunction: (cluster: import('leaflet.markercluster').MarkerCluster) => {
				const count = cluster.getChildCount();
				let size: string;
				if (count < 10) size = 'small';
				else if (count < 100) size = 'medium';
				else size = 'large';
				return L.divIcon({
					html: `<div><span>${count}</span></div>`,
					className: `marker-cluster marker-cluster-${size}`,
					iconSize: L.point(40, 40)
				});
			}
		}) as import('leaflet.markercluster').MarkerClusterGroup;

		clusterGroup = cg;
		map.addLayer(cg as unknown as import('leaflet').Layer);

		updateMarkers(L);

		mapContainer.addEventListener('click', (e) => {
			const target = e.target as HTMLElement;
			const btn = target.closest('[data-action]') as HTMLElement | null;
			if (!btn) return;
			const action = btn.dataset.action;
			const branchId = btn.dataset.branchId;
			if (!branchId) return;

			const branch = branches.find((b) => b.id === branchId);
			if (!branch) return;

			if (action === 'route') {
				onToggleRouteStop(branch);
				map?.closePopup();
			} else if (action === 'visit') {
				onMarkVisited(branch);
				map?.closePopup();
			}
		});
	}

	$effect(() => {
		// Re-read reactive deps to trigger on changes
		void filter;
		void searchQuery;
		void cardSystemIds;
		void reciprocatingSystemIds;
		void visitedBranchIds;
		void routeMode;
		void routeStops;

		if (map) {
			import('leaflet').then((L) => updateMarkers(L));
		}
	});

	export function getVisibleCount(): number {
		return getFilteredBranches().length;
	}
</script>

<div class="map-container" bind:this={mapContainer}></div>

<style>
	.map-container {
		width: 100%;
		height: 100%;
	}

	:global(.branch-popup) {
		font-family: var(--font-sans, system-ui, sans-serif);
		font-size: 0.85rem;
		line-height: 1.4;
	}

	:global(.branch-popup strong) {
		font-size: 0.95rem;
		display: block;
		margin-bottom: 2px;
	}

	:global(.popup-type) {
		font-size: 0.75rem;
		color: #6b7280;
		text-transform: capitalize;
		margin-bottom: 4px;
	}

	:global(.popup-system) {
		font-size: 0.8rem;
		display: block;
		margin-bottom: 4px;
	}

	:global(.popup-address),
	:global(.popup-phone) {
		font-size: 0.8rem;
		color: #374151;
		margin-bottom: 2px;
	}

	:global(.popup-card-status) {
		font-size: 0.75rem;
		font-weight: 600;
		margin-top: 4px;
		padding: 2px 6px;
		border-radius: 4px;
		display: inline-block;
	}

	:global(.popup-card-have-card) {
		background: #dcfce7;
		color: #166534;
	}

	:global(.popup-card-reciprocal) {
		background: #dbeafe;
		color: #1e40af;
	}

	:global(.popup-visited) {
		font-size: 0.75rem;
		font-weight: 600;
		color: #166534;
		margin-top: 2px;
	}

	:global(.popup-badge) {
		font-size: 0.7rem;
		padding: 1px 5px;
		border-radius: 3px;
		font-weight: 600;
		margin-left: 4px;
	}

	:global(.popup-badge-closed) {
		background: #fee2e2;
		color: #991b1b;
	}

	:global(.popup-badge-temporarily-closed) {
		background: #fef3c7;
		color: #92400e;
	}

	:global(.popup-nav) {
		display: flex;
		gap: 8px;
		margin-top: 6px;
		padding-top: 6px;
		border-top: 1px solid #e5e7eb;
	}

	:global(.popup-nav a) {
		font-size: 0.8rem;
		font-weight: 500;
	}

	:global(.popup-btn) {
		display: block;
		width: 100%;
		margin-top: 6px;
		padding: 4px 8px;
		font-size: 0.8rem;
		font-weight: 500;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		background: #f9fafb;
		cursor: pointer;
		text-align: center;
	}

	:global(.popup-btn:hover) {
		background: #f3f4f6;
		border-color: #9ca3af;
	}

	:global(.popup-btn-visit) {
		background: #dcfce7;
		border-color: #86efac;
		color: #166534;
	}

	:global(.popup-btn-visit:hover) {
		background: #bbf7d0;
	}

	:global(.marker-cluster) {
		background-clip: padding-box;
		border-radius: 20px;
	}

	:global(.marker-cluster div) {
		width: 30px;
		height: 30px;
		margin-left: 5px;
		margin-top: 5px;
		text-align: center;
		border-radius: 15px;
		font-size: 12px;
		font-weight: 700;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #fff;
	}

	:global(.marker-cluster-small) {
		background: rgba(156, 163, 175, 0.4);
	}
	:global(.marker-cluster-small div) {
		background: rgba(156, 163, 175, 0.8);
	}

	:global(.marker-cluster-medium) {
		background: rgba(59, 130, 246, 0.4);
	}
	:global(.marker-cluster-medium div) {
		background: rgba(59, 130, 246, 0.8);
	}

	:global(.marker-cluster-large) {
		background: rgba(34, 197, 94, 0.4);
	}
	:global(.marker-cluster-large div) {
		background: rgba(34, 197, 94, 0.8);
	}
</style>
