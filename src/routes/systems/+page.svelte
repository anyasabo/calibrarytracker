<script lang="ts">
	import { base } from '$app/paths';
	import { getReferenceData } from '$lib/data';
	import { getAllCards } from '$lib/db';
	import type { UserCard } from '$lib/types';

	const { systems } = getReferenceData();

	let searchQuery = $state('');
	let cards = $state<UserCard[]>([]);

	$effect(() => {
		getAllCards().then((c) => (cards = c));
	});

	const cardSystemIds = $derived(new Set(cards.map((c) => c.systemId)));

	const filteredSystems = $derived(
		systems
			.filter((s) => {
				const q = searchQuery.toLowerCase();
				if (!q) return true;
				return (
					s.name.toLowerCase().includes(q) ||
					s.county.toLowerCase().includes(q) ||
					s.city.toLowerCase().includes(q)
				);
			})
			.sort((a, b) => a.name.localeCompare(b.name))
	);
</script>

<svelte:head>
	<title>Library Systems | CA Library Tracker</title>
</svelte:head>

<h1>Library Systems</h1>
<p class="subtitle">
	{systems.length} public library systems in California. Search by name, city, or county.
</p>

<div class="search">
	<input
		type="search"
		placeholder="Search libraries..."
		bind:value={searchQuery}
		aria-label="Search library systems"
	/>
	<span class="result-count">{filteredSystems.length} results</span>
</div>

<div class="system-list">
	{#each filteredSystems as system (system.id)}
		<a
			href={`${base}/systems/${system.id}`}
			class="system-card"
			class:has-card={cardSystemIds.has(system.id)}
		>
			<div class="system-header">
				<h2>{system.name}</h2>
				{#if cardSystemIds.has(system.id)}
					<span class="badge">Card</span>
				{/if}
			</div>
			<p class="system-meta">{system.city}, {system.county} County</p>
			{#if system.website}
				<p class="system-link">{system.website}</p>
			{/if}
		</a>
	{/each}
</div>

{#if filteredSystems.length === 0}
	<p class="no-results">No library systems match your search.</p>
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

	.search {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		margin-bottom: var(--space-lg);
	}

	input[type='search'] {
		flex: 1;
		padding: var(--space-sm) var(--space-md);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: 0.95rem;
		outline: none;
		transition: border-color 0.15s;
	}

	input[type='search']:focus {
		border-color: var(--color-border-focus);
		box-shadow: 0 0 0 2px rgba(30, 58, 95, 0.15);
	}

	.result-count {
		font-size: 0.8rem;
		color: var(--color-text-muted);
		white-space: nowrap;
	}

	.system-list {
		display: grid;
		gap: var(--space-sm);
	}

	.system-card {
		display: block;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-md) var(--space-lg);
		transition:
			border-color 0.15s,
			box-shadow 0.15s;
		color: var(--color-text);
	}

	.system-card:hover {
		border-color: var(--color-primary);
		box-shadow: var(--shadow-sm);
		text-decoration: none;
	}

	.system-card.has-card {
		border-left: 3px solid var(--color-success);
	}

	.system-header {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
	}

	h2 {
		font-size: 1rem;
		font-weight: 600;
	}

	.badge {
		font-size: 0.7rem;
		background: var(--color-success);
		color: white;
		padding: 0.1rem 0.4rem;
		border-radius: var(--radius-full);
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.03em;
	}

	.system-meta {
		font-size: 0.85rem;
		color: var(--color-text-muted);
		margin-top: var(--space-xs);
	}

	.system-link {
		font-size: 0.8rem;
		color: var(--color-primary-light);
		margin-top: var(--space-xs);
	}

	.no-results {
		text-align: center;
		color: var(--color-text-muted);
		padding: var(--space-2xl);
	}
</style>
