<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { AppHref } from '$lib/routes';
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
	<PageIntro
		{title}
		{description}
		actionHref={backHref}
		actionLabel={backLabel}
		actionVariant="secondary"
	/>

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
		background: rgb(255 255 255 / 0.95);
		box-shadow: var(--shadow-md);
	}
</style>
