<script lang="ts">
	import { onMount } from 'svelte';
	import type { ApiJob, ApiNotebook } from '$lib/api';
	import PageIntro from '$lib/components/PageIntro.svelte';
	import OverviewCard from '$lib/components/OverviewCard.svelte';

	let jobs = $state<ApiJob[]>([]);
	let notebooks = $state<ApiNotebook[]>([]);
	let loadingJobs = $state(true);
	let loadingNotebooks = $state(true);

	function countByState(items: { state: string }[], states: string[]): number {
		const targets = new Set(states.map((state) => state.toUpperCase()));
		return items.filter((item) => targets.has(item.state.toUpperCase())).length;
	}

	const runningJobs = $derived(countByState(jobs, ['RUNNING']));
	const queuedJobs = $derived(countByState(jobs, ['CONFIGURING', 'PENDING']));
	const readyNotebooks = $derived(
		notebooks.filter((notebook) => Boolean(notebook.host && notebook.port)).length
	);
	const startingNotebooks = $derived(notebooks.length - readyNotebooks);

	void onMount(() => {
		void Promise.all([
			fetch('/api/jobs')
				.then(async (response) => {
					if (!response.ok) {
						return [];
					}

					return (await response.json()) as ApiJob[];
				})
				.then((result) => {
					jobs = result;
					loadingJobs = false;
				})
				.catch(() => {
					loadingJobs = false;
				}),
			fetch('/api/notebooks')
				.then(async (response) => {
					if (!response.ok) {
						return [];
					}

					return (await response.json()) as ApiNotebook[];
				})
				.then((result) => {
					notebooks = result;
					loadingNotebooks = false;
				})
				.catch(() => {
					loadingNotebooks = false;
				})
		]);

		const jobsEvents = new EventSource('/api/jobs/events');
		const notebooksEvents = new EventSource('/api/notebooks/events');

		jobsEvents.onmessage = (event) => {
			jobs = JSON.parse(event.data) as ApiJob[];
			loadingJobs = false;
		};

		notebooksEvents.onmessage = (event) => {
			notebooks = JSON.parse(event.data) as ApiNotebook[];
			loadingNotebooks = false;
		};

		return () => {
			jobsEvents.close();
			notebooksEvents.close();
		};
	});
</script>

<div class="page">
	<PageIntro
		title="Overview"
		description="High level overview of current Slurm jobs and Jupyter notebooks"
	/>
	<section class="summary-grid" aria-label="Overviews summaries">
		<OverviewCard
			href="/jobs"
			label="Jobs"
			value={loadingJobs ? '...' : jobs.length}
			copy={loadingJobs ? 'Loading job activity' : `${runningJobs} running, ${queuedJobs} queued`}
			linkLabel="Open Jobs"
		/>
		<OverviewCard
			href="/notebooks"
			label="Notebooks"
			value={loadingNotebooks ? '...' : notebooks.length}
			copy={loadingNotebooks
				? 'Loading notebook activity'
				: `${readyNotebooks} ready, ${startingNotebooks} starting`}
			linkLabel="Open Notebooks"
		/>
	</section>
</div>

<style>
	.page {
		display: grid;
		gap: var(--space-6);
		width: min(100%, 68rem);
		margin: 0 auto;
	}

	.summary-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: var(--space-6);
	}
</style>
