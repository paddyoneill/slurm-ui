<script lang="ts">
	import type { ApiJob } from '$lib/api';
	import { onMount } from 'svelte';
	import JobTable from './JobTable.svelte';
	import { resolve } from '$app/paths';

	let jobs = $state<ApiJob[]>([]);
	let loading = $state(true);

	void onMount(() => {
		fetch('/api/jobs')
			.then((res) => res.json())
			.then((result: ApiJob[]) => {
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

<section class="page-copy">
	<div>
		<h2>Jobs</h2>
		<p>Browse submitted workloads and open the create flow when you need a new batch job</p>
	</div>
	<a href={resolve('/jobs/new')} class="page-action">Create Job</a>
</section>

<div class="table">
	<JobTable {jobs} {loading} />
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
