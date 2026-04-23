<script lang="ts">
	import { onMount } from 'svelte';
	import type { ApiNotebook } from '$lib/api';
	import NotebookTable from './NotebookTable.svelte';
	import PageIntro from '$lib/components/PageIntro.svelte';

	let notebooks = $state<ApiNotebook[]>([]);
	let loading = $state(true);

	void onMount(() => {
		fetch('/api/notebooks')
			.then(async (response) => {
				if (!response.ok) {
					return [];
				}

				return (await response.json()) as ApiNotebook[];
			})
			.then((result) => {
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

<div class="page">
	<PageIntro
		title="Notebooks"
		description="Browse active sessions and start a new notebook server from the create page"
		actionHref="/notebooks/new"
		actionLabel="Create notebook"
	/>

	<NotebookTable {notebooks} {loading} />
</div>

<style>
	.page {
		display: grid;
		gap: var(--space-6);
		width: min(100%, 68rem);
		margin: 0 auto;
	}
</style>
