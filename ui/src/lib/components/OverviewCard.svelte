<script lang="ts">
	import { resolve } from '$app/paths';
	import type { AppHref } from '$lib/routes';

	let {
		href,
		label,
		value,
		copy,
		linkLabel
	}: {
		href: AppHref;
		label: string;
		value: string | number;
		copy: string;
		linkLabel: string;
	} = $props();
</script>

<a href={resolve(href)} class="overview-card">
	<span class="card-label">{label}</span>
	<strong class="card-value">{value}</strong>
	<p class="card-copy">{copy}</p>
	<span class="card-link">{linkLabel}</span>
</a>

<style>
	.overview-card {
		display: grid;
		gap: var(--space-3);
		padding: var(--space-6);
		border: 1px solid var(--colour-border);
		border-radius: var(--radius-lg);
		background: rgb(255 255 255 / 0.95);
		box-shadow: var(--shadow-md);
		transition: border-color var(--transition-fast) box-shadow var(--transition-fast) transform
			var(--transition-fast);
	}

	.overview-card:hover {
		border-color: rgb(15 118 110 / 0.25);
		box-shadow: 0 24px 48px rgb(15 23 42 / 0.12);
		transform: translateX(-2px);
	}

	.card-label,
	.card-link {
		font-size: 0.82rem;
		font-weight: 800;
		letter-spacing: 0.04em;
		text-transform: uppercase;
	}

	.card-label {
		color: var(--colour-text-muted);
	}

	.card-value {
		font-size: clamp(2.2rem, 4vw, 3rem);
		line-height: 1;
		letter-spacing: -0.05em;
	}

	.card-copy {
		margin: 0;
		color: var(--colour-text-muted);
		font-size: 0.95rem;
	}

	.card-link {
		color: var(--colour-primary);
	}
</style>
