<script lang="ts">
	import type { ApiJob } from '$lib/api';
	import { onMount } from 'svelte';
	import PageAction from '$lib/components/PageAction.svelte';
	import JobTable from './JobTable.svelte';
	import PageIntro from '$lib/components/PageIntro.svelte';

	let jobs = $state<ApiJob[]>([]);
	let loading = $state(true);

	void onMount(() => {
		void fetch('/api/jobs')
			.then(async (response) => {
				if (!response.ok) {
					return [];
				}

				return (await response.json()) as ApiJob[];
			})
			.then((result) => {
				jobs = result;
				loading = false;
			})
			.catch(() => {
				loading = false;
			});

		const eventSource = new EventSource('/api/jobs/events');

		eventSource.onmessage = (event) => {
			jobs = JSON.parse(event.data) as ApiJob[];
		};

		return () => eventSource.close();
	});
</script>

<div class="page">
	<PageIntro title="Jobs" description="Manage Slurm batch jobs" />

	<JobTable {jobs} {loading} />

	<PageAction href="/jobs/new" label="Create job" />
</div>

<style>
	.page {
		display: grid;
		gap: var(--space-6);
		width: min(100%, 68rem);
		margin: 0 auto;
	}
</style>
