<script lang="ts">
	import type { SlurmJobSelect } from '$lib/server/db/types';
	import { onMount } from 'svelte';
	import JobTable from './JobTable.svelte';
	import Modal from '$lib/components/Modal.svelte';
	import JobForm from '$lib/components/JobForm.svelte';
	import { createJobSchema } from '$lib/validation/job';
	import * as v from 'valibot';

	let showJobCreateForm = $state(false);
	let jobs = $state<SlurmJobSelect[]>([]);
	let loading = $state(true);
	let name = $state('');
	let script = $state('');
	let parition = $state('');
	let currentWorkingDirectory = $state('/home');
	let cpusPerTask = $state(1);
	let tasksPerNode = $state(1);
	let memoryPerNode = $state(1);
	let timeLimit = $state(1);
	let envVars = $state<{ key: string; value: string }[]>([]);
	let errors = $state<Record<string, string>>({});
	let submitting = $state(false);

	function resetForm() {
		name = '';
		script = '';
		parition = '';
		currentWorkingDirectory = '/home';
		cpusPerTask = 1;
		tasksPerNode = 1;
		memoryPerNode = 1;
		timeLimit = 1;
		errors = {};
		envVars = [];
	}

	async function handleSubmit() {
		errors = {};

		const input = {
			name,
			script,
			parition: parition || undefined,
			currentWorkingDirectory,
			environment: envVars
				.filter((e) => e.key.trim())
				.map((e) => `${e.key.trim()}=${e.value.trim()}`),
			cpusPerTask: cpusPerTask > 1 ? cpusPerTask : undefined,
			tasksPerNode: tasksPerNode > 1 ? tasksPerNode : undefined,
			memoryPerNode: memoryPerNode > 1 ? memoryPerNode : undefined,
			timeLimit: timeLimit > 1 ? timeLimit : undefined
		};

		const result = v.safeParse(createJobSchema, input);

		if (!result.success) {
			for (const issue of result.issues) {
				const key = issue.path?.[0]?.key as string;
				if (key && !errors[key]) {
					errors[key] = issue.message;
				}
			}
			return;
		}

		submitting = true;

		const response = await fetch('/api/jobs', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(result.output)
		});

		const data = await response.json();
		submitting = false;

		if (data.ok) {
			resetForm();
			showJobCreateForm = false;
		}
	}

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

<Modal bind:showModal={showJobCreateForm} title="Create New Job">
	<JobForm
		bind:name
		bind:parition
		bind:currentWorkingDirectory
		bind:cpusPerTask
		bind:tasksPerNode
		bind:memoryPerNode
		bind:timeLimit
		bind:envVars
		{errors}
		onSubmit={handleSubmit}
		{submitting}
	>
		<div>
			<label class="block text-sm font-medium text-gray-700"
				>Script
				<textarea
					bind:value={script}
					rows="10"
					class="mt-1 block w-full rounded-md border-gray-300 font-mono shadow-sm focus:border-teal-500 focus:ring-teal-500"
				></textarea>
				{#if errors.script}<span class="text-sm text-red-500">{errors.script}</span>{/if}
			</label>
		</div>
	</JobForm>
</Modal>
