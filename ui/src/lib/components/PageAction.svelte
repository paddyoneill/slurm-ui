<script lang="ts">
	import { resolve } from '$app/paths';
	import type { AppHref } from '$lib/routes';

	type ActionVariant = 'primary' | 'secondary';

	let {
		href,
		label,
		variant = 'primary'
	}: {
		href: AppHref;
		label: string;
		variant?: ActionVariant;
	} = $props();
</script>

<div class="action-shell">
	<a class={`page-action ${variant}`} href={resolve(href)}>{label}</a>
</div>

<style>
	.action-shell {
		display: flex;
		justify-content: center;
	}

	.page-action {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0.8rem 1.15rem;
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		font-weight: 700;
		box-shadow: var(--shadow-sm);
		transition:
			border-color var(--transition-fast),
			background-color var(--transition-fast),
			box-shadow var(--transition-fast);
	}

	.page-action.primary {
		background: var(--colour-primary);
		color: #fff;
	}

	.page-action.primary:hover {
		background: var(--colour-primary-strong);
		box-shadow: 0 10px 18px rgb(var(--colour-primary-rgb) / 0.18);
	}

	.page-action.secondary {
		background-color: var(--colour-border);
		background: var(--colour-surface);
		color: var(--colour-text);
	}

	.page-action.secondary:hover {
		background: var(--colour-surface-muted);
	}
</style>
