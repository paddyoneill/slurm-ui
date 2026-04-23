<script lang="ts">
	import { resolve } from '$app/paths';
	import type { AppHref } from '$lib/routes';

	type ActionVariant = 'primary' | 'secondary';

	let {
		title,
		description,
		actionHref,
		actionLabel,
		actionVariant = 'primary'
	}: {
		title: string;
		description: string;
		actionHref?: AppHref;
		actionLabel?: string;
		actionVariant?: ActionVariant;
	} = $props();
</script>

<div class="intro">
	<div class="copy">
		<h1>{title}</h1>
		<p>{description}</p>
	</div>

	{#if actionHref && actionLabel}
		<a class={`action ${actionVariant}`} href={resolve(actionHref)}>{actionLabel}</a>
	{/if}
</div>

<style>
	.intro {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-4);
		padding: var(--space-6);
		border: 1px solid var(--colour-border);
		border-radius: var(--radius-lg);
		background: linear-gradient(135deg, rgb(255 255 255 / 0.96), rgb(248 250 252 / 0.96));
		box-shadow: var(--shadow-md);
	}

	.copy {
		display: grid;
		gap: var(--space-2);
	}

	h1,
	p {
		margin: 0;
	}

	h1 {
		font-size: clamp(1.75rem, 3vw, 2.35rem);
		line-height: 1.1;
	}

	p {
		max-width: 42rem;
		color: var(--colour-text-muted);
	}

	.action {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0.8rem 1.15rem;
		border-radius: var(--radius-md);
		border: 1px solid transparent;
		font-weight: 700;
		white-space: nowrap;
		transition:
			background-color var(--transition-fast),
			color var(--transition-fast),
			border-color var(--transition-fast);
	}

	.action.primary {
		background: var(--colour-primary);
		color: #fff;
	}

	.action.primary:hover {
		background: var(--colour-primary-strong);
	}

	.action.secondary {
		border-color: var(--colour-border);
		background: var(--colour-surface);
		color: var(--colour-text);
	}

	.action.secondary:hover {
		background: var(--colour-surface-muted);
	}
</style>
