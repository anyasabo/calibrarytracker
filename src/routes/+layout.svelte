<script lang="ts">
	import '../app.css';
	import { base } from '$app/paths';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div class="app">
	<header>
		<nav>
			<a href={`${base}/`} class="logo">CA Library Tracker</a>
			<ul>
				<li><a href={`${base}/systems`}>Libraries</a></li>
				<li><a href={`${base}/map`}>Map</a></li>
				<li><a href={`${base}/cards`}>My Cards</a></li>
				<li><a href={`${base}/visits`}>My Visits</a></li>
			</ul>
		</nav>
	</header>

	<main>
		{@render children()}
	</main>

	<footer>
		<p>
			Data from the <a href="https://www.library.ca.gov/services/to-libraries/library-directory/"
				>California State Library</a
			>. Your data stays on your device.
		</p>
	</footer>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	header {
		background: var(--color-primary);
		color: white;
		height: var(--header-height);
		display: flex;
		align-items: center;
		padding: 0 var(--space-lg);
		box-shadow: var(--shadow-md);
		position: sticky;
		top: 0;
		z-index: 100;
	}

	nav {
		display: flex;
		align-items: center;
		gap: var(--space-xl);
		max-width: var(--max-width);
		width: 100%;
		margin: 0 auto;
	}

	.logo {
		color: white;
		font-weight: 700;
		font-size: 1.125rem;
		white-space: nowrap;
	}

	.logo:hover {
		text-decoration: none;
		color: var(--color-accent-light);
	}

	ul {
		display: flex;
		list-style: none;
		gap: var(--space-lg);
	}

	ul a {
		color: rgba(255, 255, 255, 0.85);
		font-size: 0.875rem;
		font-weight: 500;
		transition: color 0.15s;
	}

	ul a:hover {
		color: white;
		text-decoration: none;
	}

	main {
		flex: 1;
		max-width: var(--max-width);
		width: 100%;
		margin: 0 auto;
		padding: var(--space-xl) var(--space-lg);
	}

	footer {
		border-top: 1px solid var(--color-border);
		padding: var(--space-lg);
		text-align: center;
		font-size: 0.8rem;
		color: var(--color-text-muted);
	}

	@media (max-width: 640px) {
		nav {
			gap: var(--space-md);
		}

		ul {
			gap: var(--space-md);
		}
	}
</style>
