<script lang="ts">
	import { base } from '$app/paths';
	import { getAllVisits, deleteVisit } from '$lib/db';
	import { getReferenceData } from '$lib/data';
	import type { BranchVisit, Branch } from '$lib/types';

	const { branches, systems } = getReferenceData();

	let visits = $state<BranchVisit[]>([]);

	$effect(() => {
		loadVisits();
	});

	async function loadVisits(): Promise<void> {
		visits = await getAllVisits();
	}

	const branchMap = $derived(new Map(branches.map((b) => [b.id, b])));
	const systemMap = $derived(new Map(systems.map((s) => [s.id, s])));

	function getBranch(branchId: string): Branch | undefined {
		return branchMap.get(branchId);
	}

	async function handleDeleteVisit(id: string): Promise<void> {
		await deleteVisit(id);
		await loadVisits();
	}

	const sortedVisits = $derived(
		[...visits].sort((a, b) => {
			if (a.visitDate && b.visitDate) return b.visitDate.localeCompare(a.visitDate);
			if (a.visitDate) return -1;
			if (b.visitDate) return 1;
			return 0;
		})
	);
</script>

<svelte:head>
	<title>My Visits | CA Library Tracker</title>
</svelte:head>

<h1>Branch Visits</h1>
<p class="subtitle">
	{visits.length} branch{visits.length !== 1 ? 'es' : ''} visited out of {branches.length} total.
</p>

{#if visits.length === 0}
	<div class="empty-state">
		<p>You have not recorded any branch visits yet.</p>
		<p>
			<a href={`${base}/systems`}>Browse library systems</a> to find branches to visit.
		</p>
	</div>
{:else}
	<div class="visit-list">
		{#each sortedVisits as visit (visit.id)}
			{@const branch = getBranch(visit.branchId)}
			{@const system = branch ? systemMap.get(branch.systemId) : undefined}
			<div class="visit-item">
				<div class="visit-info">
					<h2>{branch?.name ?? 'Unknown Branch'}</h2>
					{#if system}
						<p class="visit-system">
							<a href={`${base}/systems/${system.id}`}>{system.name}</a>
						</p>
					{/if}
					<p class="visit-meta">
						{#if branch}
							{branch.city}, {branch.county} County
						{/if}
						{#if visit.visitDate}
							&middot; Visited {visit.visitDate}
						{/if}
					</p>
					{#if visit.notes}
						<p class="visit-notes">{visit.notes}</p>
					{/if}
				</div>
				<button
					class="btn-delete"
					onclick={() => handleDeleteVisit(visit.id)}
					aria-label="Remove visit to {branch?.name ?? 'unknown branch'}"
				>
					Remove
				</button>
			</div>
		{/each}
	</div>
{/if}

<style>
	h1 {
		font-size: 1.5rem;
		color: var(--color-primary);
		margin-bottom: var(--space-xs);
	}

	.subtitle {
		color: var(--color-text-muted);
		margin-bottom: var(--space-lg);
	}

	.empty-state {
		text-align: center;
		padding: var(--space-2xl);
		color: var(--color-text-muted);
	}

	.empty-state a {
		font-weight: 600;
	}

	.visit-list {
		display: grid;
		gap: var(--space-sm);
	}

	.visit-item {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-md) var(--space-lg);
	}

	h2 {
		font-size: 1rem;
		font-weight: 600;
	}

	.visit-system {
		font-size: 0.85rem;
		margin-top: 2px;
	}

	.visit-meta {
		font-size: 0.8rem;
		color: var(--color-text-muted);
		margin-top: var(--space-xs);
	}

	.visit-notes {
		font-size: 0.85rem;
		color: var(--color-text-muted);
		margin-top: var(--space-xs);
	}

	.btn-delete {
		background: none;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-xs) var(--space-sm);
		font-size: 0.8rem;
		cursor: pointer;
		color: var(--color-text-muted);
		transition:
			color 0.15s,
			border-color 0.15s;
		white-space: nowrap;
	}

	.btn-delete:hover {
		color: var(--color-danger);
		border-color: var(--color-danger);
	}
</style>
