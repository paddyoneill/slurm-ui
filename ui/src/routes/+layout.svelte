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
	<aside class="app-sidebar">
		<a class="brand" href={resolve('/')}>
			<div class="brand-mark">S</div>
			<div class="brand-copy">
				<span class="brand-title">Slurm UI</span>
				<span class="brand-subtitle">Jobs and notebooks</span>
			</div>
		</a>

		<nav class="nav" aria-label="Primary">
			{#each navItems as item, index (index)}
				<a href={resolve(item.href)} class:active={isActive(item.href)} class="nav-link">
					{item.label}
				</a>
			{/each}
		</nav>
	</aside>

	<div class="content-shell">
		<div class="utility-bar">
			<div class="utility-chip">Internal</div>
			<div class="utility-meta">Account</div>
		</div>

		<main class="content">
			{@render children()}
		</main>
	</div>
</div>

<style>
	.app-shell {
		display: grid;
		grid-template-columns: 17rem 1fr;
		min-height: 100vh;
		background:
			radial-gradient(circle at top left, rgb(15 118 110 / 0.08), transparent 28%), var(--color-bg);
	}

	.app-sidebar {
		display: flex;
		flex-direction: column;
		gap: var(--space-8);
		padding: var(--space-6);
		background: rgb(255 255 255 / 0.92);
		border-right: 1px solid var(--colour-border);
		box-shadow: var(--shadow-sm);
		backdrop-filter: blur(18px);
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
		background: linear-gradient(135deg, #0f766e, #14532d);
		color: #fff;
		font-size: 1.1rem;
		font-weight: 700;
		box-shadow: inset 0 0 0 1px rgb(15 118 120 / 0.08);
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

	.nav {
		display: grid;
		gap: var(--space-2);
	}

	.nav-link {
		display: flex;
		align-items: center;
		padding: 0.85rem 1rem;
		border-radius: var(--radius-md);
		color: var(--colour-text-muted);
		font-weight: 600;
		transition:
			background-color var(--transition-fast),
			border-color var(--transition-fast),
			transform var(--transition-fast);
	}

	.nav-link:hover {
		background: var(--colour-surface-muted);
		color: var(--colour-text);
		transform: translateX(2px);
	}

	.nav-link.active {
		background: var(--colour-primary-soft);
		color: var(--colour-primary);
	}

	.content-shell {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.utility-bar {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: var(--space-3);
		padding: var(--space-5) var(--space-8);
	}

	.utility-chip,
	.utility-meta {
		border: 1px solid var(--colour-border);
		border-radius: 999px;
		background: rgb(255 255 255 /0.78);
		box-shadow: var(--shadow-sm);
	}

	.utility-chip {
		padding: 0.45rem 0.85rem;
		color: var(--colour-primary);
		font-size: 0.85rem;
		font-weight: 700;
	}

	.utility-chip {
		padding: 0.45rem 0.9rem;
		color: var(--colour-text-muted);
		font-size: 0.9rem;
		font-weight: 600;
	}

	.content {
		flex: 1;
		padding: 0 var(--space-8) var(--space-8);
	}
</style>
