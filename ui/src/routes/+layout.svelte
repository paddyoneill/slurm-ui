<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';

	let { children } = $props();

	const navItems = [
		{ href: '/', label: 'Overview', description: 'Project status' },
		{ href: '/jobs', label: 'Jobs', description: 'Slurm batch jobs' },
		{ href: '/notebooks', label: 'Notebooks', description: 'Jupyter notebooks' }
	] as const;

	type NavHref = (typeof navItems)[number]['href'];

	function isActive(href: NavHref): boolean {
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
	<aside class="sidebar">
		<a class="brand" href={resolve('/')}>
			<img class="brand-icon" src={favicon} alt="" />
			<div class="brand-copy">
				<span class="brand-title">Slurm UI</span>
				<span class="brand-subtitle">Jobs and notebooks</span>
			</div>
		</a>
		<nav class="nav" aria-label="Primary">
			{#each navItems as item, index (index)}
				<a
					href={resolve(item.href)}
					class:active={isActive(item.href)}
					class="nav-link"
					aria-current={isActive(item.href) ? 'page' : undefined}
				>
					<span class="nav-label">{item.label}</span>
					<span class="nav-description">{item.description}</span>
				</a>
			{/each}
		</nav>
	</aside>
	<div class="content">
		<nav class="topbar" aria-label="utility">
			<div class="topbar-section">
				<span class="topbar-chip">Internal</span>
				<span class="topbar-text">Utility navigation placeholder</span>
			</div>
			<div class="topbar-section topbar-section-right">
				<span class="topbar-link">Account</span>
			</div>
		</nav>
		<main class="main">
			<div class="page">
				{@render children()}
			</div>
		</main>
	</div>
</div>

<style>
	.app-shell {
		min-height: 100vh;
		display: grid;
		grid-template-columns: 18rem minmax(0, 1fr);
		background: radial-gradient(circle at top left, rgb(15 118 110 / 0.08), transparent 24rem);
	}

	.sidebar {
		display: flex;
		flex-direction: column;
		gap: var(--space-8);
		padding: var(--space-6);
		background: rgb(255, 255, 255/0.88);
		backdrop-filter: blur(18px);
		border-right: 1px solid var(--colour-border);
	}

	.brand {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-2);
		border-radius: var(--radius-md);
	}

	.brand-icon {
		width: 2.75rem;
		height: 2.75rem;
		padding: var(--space-2);
		border-radius: var(--radius-md);
		background: linear-gradient(135deg, var(--colour-primary-soft), white);
		box-shadow: inset 0 0 0 1px rgb(15 118 120 / 0.08);
	}

	.brand-copy {
		display: grid;
		gap: 0.1rem;
	}

	.brand-title {
		font-size: 1rem;
		font-weight: 700;
		letter-spacing: -0.02em;
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
		display: grid;
		gap: 0.15rem;
		padding: 0.9rem 1rem;
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		background: transparent;
		transition:
			background-color var(--transition-fast),
			border-color var(--transition-fast),
			transform var(--transition-fast);
	}

	.nav-link:hover {
		background: var(--colour-surface-muted);
		border-color: var(--colour-border);
		transform: translateY(-1px);
	}

	.nav-link.active {
		background: linear-gradient(180deg, rgb(255 255 255 / 0.95), var(--colour-primary-soft));
		border-color: rgb(15 118 110 / 0.22);
		box-shadow: var(--shadow-sm);
	}

	.nav-description {
		font-size: 0.875rem;
		color: var(--colour-text-muted);
	}

	.nav-label {
		font-weight: 700;
		letter-spacing: -0.01em;
	}

	.content {
		min-width: 0;
		display: flex;
		flex-direction: column;
	}

	.topbar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-4);
		padding: var(--space-4) var(--space-8);
		border-bottom: 1px solid var(--colour-border);
		background: rgb(255 255 255 / 0.72);
		backdrop-filter: blur(14px);
		box-shadow: var(--shadow-sm);
	}

	.topbar-section {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		min-width: 0;
	}

	.topbar-section-right {
		justify-content: end;
	}

	.topbar-chip {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0.35rem 0.65rem;
		border-radius: 999px;
		background: var(--colour-primary-soft);
		color: var(--colour-primary-strong);
		font-size: 0.75rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
	}

	.topbar-text {
		color: var(--colour-text-muted);
		font-size: 0.95rem;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.topbar-link {
		display: inline-flex;
		align-items: center;
		min-height: 2.25rem;
		padding: 0 0 0.85rem;
		border: 1px solid var(--colour-border);
		border-radius: 999px;
		background: rgb(255 255 255 /0.86);
		color: var(--color-text);
		font-weight: 600;
	}

	.main {
		min-width: 0;
		padding: var(--space-8);
	}

	.page {
		width: min(100%, 80rem);
	}

	@media (max-width: 920px) {
		.app-shell {
			grid-template-columns: 1fr;
		}

		.sidebar {
			gap: var(--space-4);
			border-right: 0;
			border-bottom: 1px solid var(--colour-border);
		}

		.nav {
			grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));
		}

		.topbar,
		.main {
			padding-left: var(--space-4);
			padding-right: var(--space-4);
		}

		.topbar {
			flex-direction: column;
			align-items: stretch;
		}

		.topbar-section-right {
			justify-content: flex-start;
		}
	}
</style>
