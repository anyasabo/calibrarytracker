<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import {
		getSystemById,
		getBranchesBySystem,
		getPartnershipsForSystem,
		getReciprocatingSystemIds,
		getConsortiumOverdriveUrl
	} from '$lib/data';
	import { getCardsBySystem, addCard, getVisitsByBranch, addVisit } from '$lib/db';
	import { getReferenceData } from '$lib/data';
	import type { UserCard, BranchVisit, Branch } from '$lib/types';

	const { systems } = getReferenceData();

	const systemId = $derived($page.params.id ?? '');
	const system = $derived(getSystemById(systemId));
	const branches = $derived(getBranchesBySystem(systemId));
	const partnerships = $derived(getPartnershipsForSystem(systemId));
	const reciprocatingIds = $derived(getReciprocatingSystemIds(systemId));
	const consortiumOd = $derived(getConsortiumOverdriveUrl(systemId));

	let cards = $state<UserCard[]>([]);
	let visitedBranchIds = $state<Set<string>>(new Set());

	$effect(() => {
		loadUserData();
	});

	async function loadUserData(): Promise<void> {
		cards = await getCardsBySystem(systemId);
		const allVisits: BranchVisit[] = [];
		for (const branch of branches) {
			const branchVisits = await getVisitsByBranch(branch.id);
			allVisits.push(...branchVisits);
		}
		visitedBranchIds = new Set(allVisits.map((v) => v.branchId));
	}

	async function handleAddCard(): Promise<void> {
		const card: UserCard = {
			id: crypto.randomUUID(),
			systemId,
			cardNumber: '',
			pin: '',
			obtainedDate: new Date().toISOString().slice(0, 10),
			expirationDate: '',
			status: 'active',
			isDigitalOnly: false,
			notes: ''
		};
		await addCard(card);
		await loadUserData();
	}

	async function handleToggleVisit(branch: Branch): Promise<void> {
		if (visitedBranchIds.has(branch.id)) {
			// Already handled by the visits page — for simplicity, just note it
			return;
		}
		const visit: BranchVisit = {
			id: crypto.randomUUID(),
			branchId: branch.id,
			visitDate: new Date().toISOString().slice(0, 10),
			notes: ''
		};
		await addVisit(visit);
		await loadUserData();
	}

	const reciprocatingSystems = $derived(
		reciprocatingIds
			.map((id) => systems.find((s) => s.id === id))
			.filter((s): s is NonNullable<typeof s> => s !== undefined)
	);
</script>

<svelte:head>
	<title>{system?.name ?? 'Library System'} | CA Library Tracker</title>
</svelte:head>

{#if !system}
	<p>Library system not found.</p>
{:else}
	<div class="breadcrumb">
		<a href={`${base}/systems`}>Libraries</a> &rsaquo; {system.name}
	</div>

	<h1>{system.name}</h1>

	<div class="system-info">
		<div class="info-section">
			<h2>Contact</h2>
			<p>{system.address}</p>
			<p>{system.city}, {system.county} County</p>
			{#if system.phone}
				<p>Phone: {system.phone}</p>
			{/if}
			{#if system.website}
				<p><a href={system.website} target="_blank" rel="noopener">{system.website}</a></p>
			{/if}
			{#if system.email}
				<p><a href="mailto:{system.email}">{system.email}</a></p>
			{/if}
		</div>

		{#if system.digitalAccess.hasOverdrive || system.digitalAccess.offersEcard || consortiumOd}
			<div class="info-section">
				<h2>Digital Access</h2>
				{#if system.digitalAccess.hasOverdrive}
					<p>
						Libby/OverDrive:
						<a href={system.digitalAccess.overdriveUrl} target="_blank" rel="noopener">
							Browse collection
						</a>
					</p>
				{:else if consortiumOd}
					<p>
						Libby/OverDrive (via {consortiumOd.coopName}):
						<a href={consortiumOd.url} target="_blank" rel="noopener"> Browse collection </a>
					</p>
				{/if}
				{#if system.digitalAccess.offersEcard}
					<p>
						Online eCard:
						<a href={system.digitalAccess.ecardUrl} target="_blank" rel="noopener"> Get eCard </a>
					</p>
					{#if system.digitalAccess.ecardNotes}
						<p class="note">{system.digitalAccess.ecardNotes}</p>
					{/if}
				{/if}
			</div>
		{/if}
	</div>

	<!-- Cards section -->
	<section class="section">
		<div class="section-header">
			<h2>Your Cards ({cards.length})</h2>
			<button class="btn btn-primary-small" onclick={handleAddCard}>Add Card</button>
		</div>
		{#if cards.length === 0}
			<p class="muted">No cards for this system yet.</p>
		{:else}
			{#each cards as card (card.id)}
				<div class="card-badge">
					<span class="status-{card.status}">{card.status}</span>
					{#if card.isDigitalOnly}eCard{/if}
					{#if card.obtainedDate}— obtained {card.obtainedDate}{/if}
				</div>
			{/each}
		{/if}
	</section>

	<!-- Branches section -->
	<section class="section">
		<h2>Branches ({branches.length})</h2>
		<p class="muted">
			{visitedBranchIds.size} of {branches.length} visited
		</p>
		<div class="branch-list">
			{#each branches as branch (branch.id)}
				{@const visited = visitedBranchIds.has(branch.id)}
				<div class="branch-item" class:visited>
					<button
						class="visit-toggle"
						onclick={() => handleToggleVisit(branch)}
						aria-label="{visited ? 'Visited' : 'Mark as visited'}: {branch.name}"
						disabled={visited}
					>
						{visited ? '✓' : '○'}
					</button>
					<div class="branch-info">
						<strong>{branch.name}</strong>
						<span class="branch-meta">
							{branch.address}, {branch.city}
							{branch.zipCode}
						</span>
					</div>
				</div>
			{/each}
		</div>
	</section>

	<!-- Reciprocity section -->
	{#if reciprocatingSystems.length > 0}
		<section class="section">
			<h2>Reciprocal Borrowing</h2>
			<p class="muted">Your card from {system.name} may also work at these systems:</p>
			{#if partnerships.length > 0}
				{#each partnerships as partnership (partnership.id)}
					<div class="partnership">
						<h3>{partnership.name}</h3>
						<p class="note">{partnership.description}</p>
					</div>
				{/each}
			{/if}
			<ul class="reciprocal-list">
				{#each reciprocatingSystems as rSystem (rSystem.id)}
					<li>
						<a href={`${base}/systems/${rSystem.id}`}>{rSystem.name}</a>
						<span class="muted">— {rSystem.city}, {rSystem.county} County</span>
					</li>
				{/each}
			</ul>
		</section>
	{/if}
{/if}

<style>
	.breadcrumb {
		font-size: 0.85rem;
		color: var(--color-text-muted);
		margin-bottom: var(--space-md);
	}

	h1 {
		font-size: 1.5rem;
		color: var(--color-primary);
		margin-bottom: var(--space-lg);
	}

	.system-info {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: var(--space-lg);
		margin-bottom: var(--space-xl);
	}

	.info-section {
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-lg);
	}

	.info-section h2 {
		font-size: 0.9rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: var(--color-text-muted);
		margin-bottom: var(--space-sm);
	}

	.info-section p {
		font-size: 0.9rem;
		margin-top: var(--space-xs);
	}

	.section {
		margin-bottom: var(--space-xl);
	}

	.section h2 {
		font-size: 1.1rem;
		color: var(--color-primary);
		margin-bottom: var(--space-sm);
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--space-sm);
	}

	.btn-primary-small {
		background: var(--color-primary);
		color: white;
		border: none;
		padding: var(--space-xs) var(--space-md);
		border-radius: var(--radius-md);
		font-size: 0.85rem;
		cursor: pointer;
		font-weight: 600;
	}

	.btn-primary-small:hover {
		background: var(--color-primary-light);
	}

	.card-badge {
		display: inline-flex;
		gap: var(--space-sm);
		align-items: center;
		padding: var(--space-xs) var(--space-sm);
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: 0.85rem;
		margin-top: var(--space-xs);
	}

	.muted {
		color: var(--color-text-muted);
		font-size: 0.9rem;
	}

	.note {
		color: var(--color-text-muted);
		font-size: 0.85rem;
		font-style: italic;
	}

	.branch-list {
		display: grid;
		gap: 2px;
		margin-top: var(--space-sm);
	}

	.branch-item {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		padding: var(--space-sm) var(--space-md);
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
	}

	.branch-item.visited {
		border-left: 3px solid var(--color-success);
	}

	.visit-toggle {
		background: none;
		border: none;
		font-size: 1.1rem;
		cursor: pointer;
		width: 1.5rem;
		text-align: center;
		color: var(--color-text-muted);
	}

	.branch-item.visited .visit-toggle {
		color: var(--color-success);
	}

	.branch-info {
		display: flex;
		flex-direction: column;
	}

	.branch-info strong {
		font-size: 0.9rem;
	}

	.branch-meta {
		font-size: 0.8rem;
		color: var(--color-text-muted);
	}

	.partnership {
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-md);
		margin: var(--space-sm) 0;
	}

	.partnership h3 {
		font-size: 0.95rem;
		margin-bottom: var(--space-xs);
	}

	.reciprocal-list {
		list-style: none;
		margin-top: var(--space-sm);
	}

	.reciprocal-list li {
		padding: var(--space-xs) 0;
		font-size: 0.9rem;
		border-bottom: 1px solid var(--color-border);
	}

	.reciprocal-list li:last-child {
		border-bottom: none;
	}

	.status-active {
		color: var(--color-success);
		font-weight: 600;
	}

	.status-expired {
		color: var(--color-warning);
		font-weight: 600;
	}

	.status-lost {
		color: var(--color-danger);
		font-weight: 600;
	}

	.status-unknown {
		color: var(--color-text-muted);
		font-weight: 600;
	}
</style>
