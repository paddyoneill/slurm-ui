<script lang="ts">
	import { createJobSchema } from '$lib/validation/job';
	import * as v from 'valibot';

	let createJobForm = $state() as HTMLDialogElement;
	let {
		showForm = $bindable()
	}: {
		showForm: boolean;
	} = $props();

	$effect(() => {
		if (showForm) createJobForm.showModal();
	});

	let name = $state('');
	let script = $state('');
	let parition = $state('');
	let currentWorkingDirectory = $state('$HOME');
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
		currentWorkingDirectory = '$HOME';
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
			createJobForm.close();
		}
	}
</script>

<dialog
	class="m-auto w-full max-w-md rounded-lg bg-gray-100 p-6 shadow-lg backdrop:bg-black/60"
	bind:this={createJobForm}
	onclose={() => (showForm = false)}
	onclick={(e) => {
		if (e.target === createJobForm) createJobForm.close();
	}}
>
	<form
		class="relative w-full max-w-lg p-6"
		onsubmit={(e) => {
			e.preventDefault();
			handleSubmit();
		}}
	>
		<button
			type="button"
			onclick={() => createJobForm.close()}
			class="absolute top-2 right-2 p-2 text-xl font-semibold text-gray-400 hover:text-gray-600"
			>X</button
		>
		<h2 class="mb-6 text-lg font-semibold">Create New Job</h2>
		<div class="space-y-4">
			<div>
				<label class="block text-sm font-medium text-gray-700"
					>Name
					<input
						type="text"
						bind:value={name}
						class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
					/>
					{#if errors.name}<span class="text-sm text-red-500">{errors.name}</span>{/if}
				</label>
			</div>
			<div>
				<label class="block text-sm font-medium text-gray-700"
					>Working Directory
					<input
						type="text"
						bind:value={currentWorkingDirectory}
						class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
					/>
					{#if errors.currentWorkingDirectory}<span class="text-sm text-red-500"
							>{errors.currentWorkingDirectory}</span
						>{/if}
				</label>
			</div>
			<div>
				<label class="block text-sm font-medium text-gray-700"
					>Partition
					<input
						type="text"
						bind:value={parition}
						class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
					/>
					{#if errors.parition}<span class="text-sm text-red-500">{errors.parition}</span>{/if}
				</label>
			</div>
			<div class="grid grid-cols-2 gap-4">
				<div>
					<label class="block text-sm font-medium text-gray-700"
						>Tasks Per Node
						<input
							type="number"
							bind:value={tasksPerNode}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
						/>
						{#if errors.cpusPerTask}<span class="text-sm text-red-500">{errors.cpusPerTask}</span
							>{/if}
					</label>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700"
						>CPUs Per Task
						<input
							type="number"
							bind:value={cpusPerTask}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
						/>
						{#if errors.cpusPerTask}<span class="text-sm text-red-500">{errors.cpusPerTask}</span
							>{/if}
					</label>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700"
						>Memory Per Node
						<input
							type="number"
							bind:value={memoryPerNode}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
						/>
						{#if errors.memoryPerNode}<span class="text-sm text-red-500"
								>{errors.memoryPerNode}</span
							>{/if}
					</label>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700"
						>Time Limit
						<input
							type="number"
							bind:value={timeLimit}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
						/>
						{#if errors.timeLimit}<span class="text-sm text-red-500">{errors.timeLimit}</span>{/if}
					</label>
				</div>
			</div>

			<fieldset>
				<legend class="block text-sm font-medium text-gray-700">Environment Variables</legend>
				{#each envVars as env, i (i)}
					<div class="mt-1 flex items-center gap-2">
						<input
							class="block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
							type="text"
							bind:value={env.key}
							placeholder="KEY"
						/>
						<input
							class="block w-full rounded-md border-gray-300 shadow-sm focus:border-teal-500 focus:ring-teal-500"
							type="text"
							bind:value={env.value}
							placeholder="VALUE"
						/>
						<button
							type="button"
							onclick={() => (envVars = envVars.filter((_, j) => j != i))}
							class="self-stretch rounded-md bg-red-500 px-3 text-sm text-white hover:bg-red-700"
							>Remove</button
						>
					</div>
				{/each}
				{#if errors.environment}
					<span class="text-sm text-red-500">{errors.environment}</span>
				{/if}
				<button
					type="button"
					onclick={() => (envVars = [...envVars, { key: '', value: '' }])}
					class="text-sm text-teal-600 hover:text-teal-700">+ Add Variable</button
				>
			</fieldset>

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
		</div>
		<div class="flex gap-3 pt-4">
			<button
				disabled={submitting}
				type="submit"
				class="w-full rounded-md bg-green-600 px-4 py-2 text-white hover:bg-green-700"
				>{submitting ? 'Submitting...' : 'Submit'}
			</button>
		</div>
	</form>
</dialog>
