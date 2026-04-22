<script lang="ts">
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import type { ApiNotebook } from '$lib/api';
	import NotebookTable from './NotebookTable.svelte';

	let notebooks = $state<ApiNotebook[]>([]);
	let loading = $state(true);

	void onMount(() => {
		fetch('/api/notebooks')
			.then((res) => res.json())
			.then((result: ApiNotebook[]) => {
				notebooks = result;
				loading = false;
			})
			.catch(() => {
				loading = false;
			});

		const eventSource = new EventSource('/api/notebooks/events');

		eventSource.onmessage = (event) => {
			notebooks = JSON.parse(event.data) as ApiNotebook[];
		};

		return () => eventSource.close();
	});
</script>

<section class="page-copy">
	<div>
		<h2>Notebooks</h2>
		<p>Browse active sessions and start a new notebook server from the create page</p>
	</div>
	<a href={resolve('/notebooks/new')} class="page-action">Create Notebook</a>
</section>

<div class="table">
	<NotebookTable {notebooks} {loading} />
</div>

<style>
	.page-copy {
		display: flex;
		align-content: flex-start;
		justify-content: space-between;
		gap: var(--space-4);
		margin-bottom: var(--space-6);
	}

	h2 {
		margin: 0;
		font-size: 1.75rem;
		letter-spacing: -0.03em;
	}

	p {
		margin: var(--space-2) 0 0;
		color: var(--colour-text-muted);
	}

	.page-action {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-height: 2.75rem;
		padding: 0 1rem;
		border: 1px solid var(--colour-border);
		border-radius: var(--border-md);
		background: var(--colour-surface);
		font-weight: 600;
		box-shadow: var(--shadow-sm);
	}

	.table {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
	}

	@media (max-width: 760px) {
		.page-copy {
			flex-direction: column;
		}
	}
</style>
