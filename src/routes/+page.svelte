<script lang="ts">
	// import { submitJobForm } from '$lib/rpc/slurm-job.remote';
	import type { SlurmJobSelect } from '$lib/server/db/schema';
	import { onMount } from 'svelte';
	import JobTable from './JobTable.svelte';
	import JobCreateForm from './JobCreateForm.svelte';

	let showJobCreateForm = $state(false);
	let jobs = $state<SlurmJobSelect[]>([]);
	let loading = $state(true);

	onMount(() => {
		fetch('/api/jobs')
			.then((res) => res.json())
			.then((result) => {
				if (result.ok) {
					jobs = result.value;
				}
				loading = false;
			});

		const eventSource = new EventSource('/api/jobs/events');

		eventSource.onmessage = (event) => {
			jobs = JSON.parse(event.data);
		};

		return () => eventSource.close();
	});
</script>

<div class="flex h-full flex-col items-center">
	<JobTable {jobs} {loading} />
	<div class="flex justify-center">
		<button
			class="m-4 rounded bg-green-600 px-4 py-2 text-white hover:bg-green-700"
			onclick={() => (showJobCreateForm = true)}>New Job</button
		>
	</div>
</div>

<JobCreateForm bind:showForm={showJobCreateForm} />
