<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { AppHref } from '$lib/routes';
	import favicon from '$lib/assets/favicon.svg';
	import './layout.css';

	let { children } = $props();

	const navItems: { href: AppHref; label: string }[] = [
		{ href: '/', label: 'Overview' },
		{ href: '/jobs', label: 'Jobs' },
		{ href: '/notebooks', label: 'Notebooks' }
	];

	function isActive(href: AppHref): boolean {
		if (href === '/') {
			return page.url.pathname === '/';
		}

		return page.url.pathname === href || page.url.pathname.startsWith(`${href}/`);
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Slurm UI</title>
</svelte:head>

<div class="app-shell">
	<div class="utility-bar">
		<a class="brand" href={resolve('/')}>
			<div class="brand-mark">S</div>
			<div class="brand-copy">
				<span class="brand-title">Slurm UI</span>
				<span class="brand-subtitle">Jobs and notebooks</span>
			</div>
		</a>

		<div class="utility-actions">
			<div class="utility-chip">Internal</div>
			<div class="utility-meta">Account</div>
		</div>
	</div>

	<aside class="app-sidebar">
		<nav class="nav" aria-label="Primary">
			{#each navItems as item, index (index)}
				<a href={resolve(item.href)} class:active={isActive(item.href)} class="nav-link">
					{item.label}
				</a>
			{/each}
		</nav>
	</aside>

	<div class="content-shell">
		<main class="content">
			{@render children()}
		</main>
	</div>
</div>

<style>
	.app-shell {
		display: grid;
		grid-template-columns: 17rem 1fr;
		grid-template-rows: auto 1fr;
		min-height: 100vh;
		background:
			linear-gradient(180deg, rgb(255 255 255 / 0.55), transparent 12rem),
			radial-gradient(circle at top right, rgb(var(--colour-primary-rgb) / 0.1), transparent 24%),
			var(--colour-bg);
	}

	.utility-bar {
		grid-column: 1 / -1;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-3);
		padding: 0.75rem var(--space-8);
		background: var(--colour-surface);
		border-bottom: 1px solid var(--colour-border);
		box-shadow: var(--shadow-sm);
	}

	.app-sidebar {
		grid-row: 2;
		display: flex;
		flex-direction: column;
		padding: var(--space-6);
		background: linear-gradient(180deg, var(--colour-surface), var(--colour-surface-muted));
		border-right: 1px solid var(--colour-border);
	}

	.brand {
		display: flex;
		align-items: center;
		gap: var(--space-3);
	}

	.brand-mark {
		display: grid;
		place-items: center;
		width: 2.75rem;
		height: 2.75rem;
		border-radius: var(--radius-md);
		background: linear-gradient(135deg, var(--colour-brand-start), var(--colour-brand-end));
		color: #fff;
		font-size: 1.1rem;
		font-weight: 700;
		box-shadow:
			inset 0 0 0 1px rgb(255 255 255 / 0.28),
			0 10px 24px rgb(var(--colour-primary-rgb) / 0.22);
	}

	.brand-copy {
		display: grid;
		gap: var(--space-1);
	}

	.brand-title,
	.brand-subtitle {
		margin: 0;
	}

	.brand-title {
		font-size: 1rem;
		font-weight: 700;
	}

	.brand-subtitle {
		font-size: 0.875rem;
		color: var(--colour-text-muted);
	}

	.utility-actions {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: var(--space-3);
	}

	.nav {
		display: grid;
		gap: 0.35rem;
	}

	.nav-link {
		position: relative;
		display: flex;
		align-items: center;
		padding: 0.8rem 1rem 0.8rem 1.25rem;
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		background: transparent;
		color: var(--colour-text-muted);
		font-weight: 600;
		transition:
			background-color var(--transition-fast),
			border-color var(--transition-fast),
			color var(--transition-fast),
			box-shadow var(--transition-fast);
	}

	.nav-link:hover {
		border-color: var(--colour-border);
		background: var(--colour-surface);
		color: var(--colour-text);
		box-shadow: var(--shadow-sm);
	}

	.nav-link.active {
		border-color: rgb(var(--colour-primary-rgb) / 0.22);
		background: var(--colour-primary-soft);
		color: var(--colour-primary-strong);
		box-shadow: var(--shadow-sm);
	}

	.nav-link.active::before {
		content: '';
		position: absolute;
		left: 0.5rem;
		width: 0.3rem;
		height: 1.2rem;
		border-radius: 999px;
		background: var(--colour-primary);
	}

	.content-shell {
		grid-row: 2;
		display: flex;
		flex-direction: column;
		min-width: 0;
		padding-top: var(--space-5);
	}

	.utility-chip,
	.utility-meta {
		display: inline-flex;
		align-items: center;
		padding: 0.45rem 0.9rem;
		border: 1px solid var(--colour-border);
		border-radius: 999px;
		background: var(--colour-surface);
		box-shadow: var(--shadow-sm);
	}

	.utility-chip {
		background-color: rgb(var(--colour-primary-rgb) / 0.18);
		background: var(--colour-primary-soft);
		color: var(--colour-primary-strong);
		font-size: 0.85rem;
		font-weight: 700;
	}

	.utility-meta {
		color: var(--colour-text-muted);
		font-size: 0.9rem;
		font-weight: 600;
	}

	.content {
		flex: 1;
		padding: 0 var(--space-8) var(--space-8);
	}
</style>
