<script lang="ts">
	import { onMount } from 'svelte';
	import type { NotebookSelect } from '$lib/server/db/types';
	import Modal from '$lib/components/Modal.svelte';
	import NotebookTable from './NotebookTable.svelte';
	import JobForm from '$lib/components/JobForm.svelte';
	import { createNotebookScheme } from '$lib/validation/notebook';
	import * as v from 'valibot';

	let showNotebookCreateForm = $state(false);
	let notebooks = $state<NotebookSelect[]>([]);
	let loading = $state(true);
	let name = $state('');
	let parition = $state('');
	let currentWorkingDirectory = $state('/home');
	let baseEnv = $state('/home/hpcadmin/paddy/base-notebook');
	let cpusPerTask = $state(1);
	let tasksPerNode = $state(1);
	let memoryPerNode = $state(1);
	let timeLimit = $state(1);
	let envVars = $state<{ key: string; value: string }[]>([]);
	let errors = $state<Record<string, string>>({});
	let submitting = $state(false);

	function resetForm() {
		name = '';
		parition = '';
		currentWorkingDirectory = '/home';
		baseEnv = '/home/hpcadmin/paddy/base-notebook';
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
			parition: parition || undefined,
			currentWorkingDirectory,
			baseEnv,
			environment: envVars
				.filter((e) => e.key.trim())
				.map((e) => `${e.key.trim()}=${e.value.trim()}`),
			cpusPerTask: cpusPerTask > 1 ? cpusPerTask : undefined,
			tasksPerNode: tasksPerNode > 1 ? tasksPerNode : undefined,
			memoryPerNode: memoryPerNode > 1 ? memoryPerNode : undefined,
			timeLimit: timeLimit > 1 ? timeLimit : undefined
		};

		const result = v.safeParse(createNotebookScheme, input);

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

		const response = await fetch('/api/notebooks', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(result.output)
		});

		const data = await response.json();
		submitting = false;

		if (data.ok) {
			resetForm();
			showNotebookCreateForm = false;
		}
	}

	void onMount(() => {
		fetch('/api/notebooks')
			.then((res) => res.json())
			.then((result) => {
				if (result.ok) {
					notebooks = result.value;
				}
				loading = false;
			});

		const eventSource = new EventSource('/api/notebooks/events');

		eventSource.onmessage = (event) => {
			notebooks = JSON.parse(event.data);
		};

		return () => eventSource.close();
	});
</script>

<div class="flex h-full flex-col items-center">
	<NotebookTable {notebooks} {loading} />
	<div class="flex justify-center">
		<button
			class="m-4 rounded bg-green-600 px-4 py-2 text-white hover:bg-green-700"
			onclick={() => (showNotebookCreateForm = true)}>New Notebook</button
		>
	</div>
</div>

<Modal bind:showModal={showNotebookCreateForm} title="Create New Notebook">
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
				>Base Environment
				<input
					type="text"
					bind:value={baseEnv}
					class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
				/>
				{#if errors.baseEnv}<span class="text-sm text-red-500">{errors.baseEnv}</span>{/if}
			</label>
		</div>
	</JobForm>
</Modal>
