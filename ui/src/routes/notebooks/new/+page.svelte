<script lang="ts">
	import CreatePage from '$lib/components/CreatePage.svelte';
	import JobForm from '$lib/components/JobForm.svelte';
	import { createNotebookScheme } from '$lib/validation/notebook';
	import * as v from 'valibot';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let name = $state('');
	let partition = $state('');
	let currentWorkingDirectory = $state('/home');
	let baseEnv = $state('/home/hpcadmin/paddy/base-notebook');
	let cpusPerTask = $state(1);
	let tasksPerNode = $state(1);
	let memoryPerNode = $state(1);
	let timeLimit = $state(1);
	let envVars = $state<{ key: string; value: string }[]>([]);
	let errors = $state<Record<string, string>>({});
	let submitting = $state(false);

	async function handleSubmit() {
		errors = {};

		const input = {
			name,
			partition: partition || undefined,
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

		await response.json().catch(() => null);
		submitting = false;

		if (response.ok) {
			goto(resolve('/notebooks'));
		}
	}
</script>

<CreatePage
	title="Create Notebook"
	description="Launch a new Juputer notebook server"
	backHref="/notebooks"
	backLabel="Back to Notebooks"
>
	<JobForm
		bind:name
		bind:partition
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
			<label class="field-label"
				>Base Environment
				<input type="text" bind:value={baseEnv} class="field-input" />
				{#if errors.baseEnv}<span class="field-error">{errors.baseEnv}</span>{/if}
			</label>
		</div>
	</JobForm>
</CreatePage>

<style>
	.field-label {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.field-input {
		display: block;
		width: 100%;
		margin-top: var(--space-2);
		padding: 0.75rem;
		border: 1px solid var(--colour-border);
		border-radius: var(--radius-md);
		font-family: var(--font-mono);
		background: var(--colour-surface);
	}

	.field-error {
		display: block;
		margin-top: var(--space-1);
		color: var(--colour-danger);
		font-size: 0.875rem;
	}
</style>
