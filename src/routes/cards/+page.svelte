<script lang="ts">
	import { base } from '$app/paths';
	import { getAllCards, deleteCard, exportData, importData } from '$lib/db';
	import { getSystemById } from '$lib/data';
	import type { UserCard } from '$lib/types';

	let cards = $state<UserCard[]>([]);
	let showImportDialog = $state(false);
	let importError = $state('');

	$effect(() => {
		loadCards();
	});

	async function loadCards(): Promise<void> {
		cards = await getAllCards();
	}

	async function handleDelete(id: string): Promise<void> {
		await deleteCard(id);
		await loadCards();
	}

	async function handleExport(): Promise<void> {
		const data = await exportData();
		const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `calibrary-backup-${new Date().toISOString().slice(0, 10)}.json`;
		a.click();
		URL.revokeObjectURL(url);
	}

	async function handleImport(event: Event): Promise<void> {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;

		try {
			const text = await file.text();
			const data = JSON.parse(text);
			await importData(data);
			await loadCards();
			showImportDialog = false;
			importError = '';
		} catch {
			importError = 'Invalid backup file. Please select a valid JSON export.';
		}
	}
</script>

<svelte:head>
	<title>My Cards | CA Library Tracker</title>
</svelte:head>

<div class="header">
	<div>
		<h1>My Library Cards</h1>
		<p class="subtitle">
			{cards.length} card{cards.length !== 1 ? 's' : ''} tracked
		</p>
	</div>
	<div class="header-actions">
		<button class="btn btn-small" onclick={handleExport}>Export</button>
		<button class="btn btn-small" onclick={() => (showImportDialog = true)}>Import</button>
	</div>
</div>

{#if showImportDialog}
	<div class="import-dialog">
		<p>Select a previously exported JSON backup file:</p>
		<input type="file" accept=".json" onchange={handleImport} />
		{#if importError}
			<p class="error">{importError}</p>
		{/if}
		<button class="btn btn-small" onclick={() => (showImportDialog = false)}>Cancel</button>
	</div>
{/if}

{#if cards.length === 0}
	<div class="empty-state">
		<p>You have not added any library cards yet.</p>
		<p>
			<a href={`${base}/systems`}>Browse library systems</a> to start adding cards.
		</p>
	</div>
{:else}
	<div class="card-list">
		{#each cards as card (card.id)}
			{@const system = getSystemById(card.systemId)}
			<div class="card-item" class:expired={card.status === 'expired'}>
				<div class="card-info">
					<h2>
						{#if system}
							<a href={`${base}/systems/${system.id}`}>{system.name}</a>
						{:else}
							Unknown System
						{/if}
					</h2>
					<div class="card-meta">
						<span class="status-badge status-{card.status}">{card.status}</span>
						{#if card.isDigitalOnly}
							<span class="type-badge">eCard</span>
						{/if}
						{#if card.obtainedDate}
							<span>Obtained: {card.obtainedDate}</span>
						{/if}
						{#if card.expirationDate}
							<span>Expires: {card.expirationDate}</span>
						{/if}
					</div>
					{#if card.notes}
						<p class="card-notes">{card.notes}</p>
					{/if}
				</div>
				<button
					class="btn-delete"
					onclick={() => handleDelete(card.id)}
					aria-label="Delete card for {system?.name ?? 'unknown system'}"
				>
					Remove
				</button>
			</div>
		{/each}
	</div>
{/if}

<style>
	.header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: var(--space-lg);
	}

	h1 {
		font-size: 1.5rem;
		color: var(--color-primary);
	}

	.subtitle {
		color: var(--color-text-muted);
		font-size: 0.9rem;
	}

	.header-actions {
		display: flex;
		gap: var(--space-sm);
	}

	.btn {
		padding: var(--space-sm) var(--space-md);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-bg-card);
		cursor: pointer;
		font-size: 0.85rem;
		transition: border-color 0.15s;
	}

	.btn:hover {
		border-color: var(--color-primary);
	}

	.btn-small {
		padding: var(--space-xs) var(--space-sm);
		font-size: 0.8rem;
	}

	.import-dialog {
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-lg);
		margin-bottom: var(--space-lg);
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}

	.error {
		color: var(--color-danger);
		font-size: 0.85rem;
	}

	.empty-state {
		text-align: center;
		padding: var(--space-2xl);
		color: var(--color-text-muted);
	}

	.empty-state a {
		font-weight: 600;
	}

	.card-list {
		display: grid;
		gap: var(--space-sm);
	}

	.card-item {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-left: 3px solid var(--color-success);
		border-radius: var(--radius-md);
		padding: var(--space-md) var(--space-lg);
	}

	.card-item.expired {
		border-left-color: var(--color-warning);
		opacity: 0.75;
	}

	h2 {
		font-size: 1rem;
		font-weight: 600;
	}

	.card-meta {
		display: flex;
		gap: var(--space-sm);
		align-items: center;
		flex-wrap: wrap;
		margin-top: var(--space-xs);
		font-size: 0.8rem;
		color: var(--color-text-muted);
	}

	.status-badge {
		padding: 0.1rem 0.4rem;
		border-radius: var(--radius-full);
		font-size: 0.7rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.03em;
	}

	.status-active {
		background: var(--color-success);
		color: white;
	}

	.status-expired {
		background: var(--color-warning);
		color: white;
	}

	.status-lost {
		background: var(--color-danger);
		color: white;
	}

	.status-unknown {
		background: var(--color-text-muted);
		color: white;
	}

	.type-badge {
		padding: 0.1rem 0.4rem;
		border-radius: var(--radius-full);
		font-size: 0.7rem;
		font-weight: 600;
		background: var(--color-primary-light);
		color: white;
	}

	.card-notes {
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
