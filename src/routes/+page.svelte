<script lang="ts">
	import { base } from '$app/paths';
	import { getReferenceData } from '$lib/data';
	import { getAllCards, getAllVisits } from '$lib/db';

	const { systems, branches } = getReferenceData();

	let cardCount = $state(0);
	let visitCount = $state(0);

	$effect(() => {
		getAllCards().then((cards) => (cardCount = cards.length));
		getAllVisits().then((visits) => (visitCount = visits.length));
	});
</script>

<svelte:head>
	<title>CA Library Card Tracker</title>
</svelte:head>

<div class="hero">
	<h1>California Library Card Tracker</h1>
	<p class="subtitle">
		Track your library cards across all {systems.length} public library systems and {branches.length}
		branches in California.
	</p>
</div>

<div class="stats">
	<div class="stat-card">
		<span class="stat-number">{systems.length}</span>
		<span class="stat-label">Library Systems</span>
	</div>
	<div class="stat-card">
		<span class="stat-number">{branches.length}</span>
		<span class="stat-label">Branches</span>
	</div>
	<div class="stat-card">
		<span class="stat-number">{cardCount}</span>
		<span class="stat-label">Your Cards</span>
	</div>
	<div class="stat-card">
		<span class="stat-number">{visitCount}</span>
		<span class="stat-label">Branches Visited</span>
	</div>
</div>

<div class="actions">
	<a href={`${base}/systems`} class="btn btn-primary">Browse Libraries</a>
	<a href={`${base}/map`} class="btn btn-secondary">View Map</a>
	<a href={`${base}/cards`} class="btn btn-secondary">Manage Cards</a>
</div>

<section class="info">
	<h2>How it works</h2>
	<div class="info-grid">
		<div class="info-card">
			<h3>Get cards</h3>
			<p>
				Any California resident can get a free library card from any public library in the state.
				Each system gives you access to different ebook collections through Libby.
			</p>
		</div>
		<div class="info-card">
			<h3>Track progress</h3>
			<p>
				Keep track of which systems you have cards from, when you got them, and whether they are
				still active. Mark branches you have visited.
			</p>
		</div>
		<div class="info-card">
			<h3>Your data, your device</h3>
			<p>
				All your data stays on your device. Nothing is sent to any server. You can export your data
				as JSON anytime for backup.
			</p>
		</div>
	</div>
</section>

<style>
	.hero {
		text-align: center;
		padding: var(--space-2xl) 0 var(--space-xl);
	}

	h1 {
		font-size: 2rem;
		color: var(--color-primary);
		margin-bottom: var(--space-sm);
	}

	.subtitle {
		font-size: 1.1rem;
		color: var(--color-text-muted);
		max-width: 36rem;
		margin: 0 auto;
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
		gap: var(--space-md);
		margin: var(--space-xl) 0;
	}

	.stat-card {
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-lg);
		text-align: center;
		box-shadow: var(--shadow-sm);
	}

	.stat-number {
		display: block;
		font-size: 2rem;
		font-weight: 700;
		color: var(--color-primary);
	}

	.stat-label {
		display: block;
		font-size: 0.8rem;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
		margin-top: var(--space-xs);
	}

	.actions {
		display: flex;
		gap: var(--space-md);
		justify-content: center;
		margin: var(--space-xl) 0;
	}

	.btn {
		display: inline-block;
		padding: var(--space-sm) var(--space-lg);
		border-radius: var(--radius-md);
		font-weight: 600;
		font-size: 0.95rem;
		transition:
			background 0.15s,
			transform 0.1s;
	}

	.btn:hover {
		text-decoration: none;
		transform: translateY(-1px);
	}

	.btn-primary {
		background: var(--color-primary);
		color: white;
	}

	.btn-primary:hover {
		background: var(--color-primary-light);
	}

	.btn-secondary {
		background: var(--color-bg-card);
		color: var(--color-primary);
		border: 1px solid var(--color-border);
	}

	.btn-secondary:hover {
		border-color: var(--color-primary);
	}

	.info {
		margin-top: var(--space-2xl);
	}

	.info h2 {
		text-align: center;
		font-size: 1.4rem;
		margin-bottom: var(--space-lg);
		color: var(--color-primary);
	}

	.info-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
		gap: var(--space-lg);
	}

	.info-card {
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: var(--space-lg);
		box-shadow: var(--shadow-sm);
	}

	.info-card h3 {
		font-size: 1rem;
		margin-bottom: var(--space-sm);
		color: var(--color-primary);
	}

	.info-card p {
		font-size: 0.9rem;
		color: var(--color-text-muted);
		line-height: 1.5;
	}
</style>
