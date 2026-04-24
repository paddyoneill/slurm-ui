<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { AppHref } from '$lib/routes';
	import PageAction from './PageAction.svelte';
	import PageIntro from './PageIntro.svelte';

	let {
		title,
		description,
		backHref,
		backLabel,
		children
	}: {
		title: string;
		description: string;
		backHref: AppHref;
		backLabel: string;
		children: Snippet;
	} = $props();
</script>

<div class="create-page">
	<PageIntro {title} {description} />

	<PageAction href={backHref} label={backLabel} variant="secondary" />

	<section class="form-panel">
		{@render children()}
	</section>
</div>

<style>
	.create-page {
		display: grid;
		gap: var(--space-6);
		width: min(100%, 68rem);
		margin: 0 auto;
	}

	.form-panel {
		padding: var(--space-4);
		border: 1px solid var(--colour-border);
		border-radius: var(--radius-lg);
		background: var(--colour-surface);
		box-shadow: var(--shadow-md);
	}
</style>
