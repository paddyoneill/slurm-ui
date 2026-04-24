<script lang="ts">
	import CreatePage from '$lib/components/CreatePage.svelte';
	import JobForm from '$lib/components/JobForm.svelte';
	import { createJobSchema } from '$lib/validation/job';
	import * as v from 'valibot';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let name = $state('');
	let script = $state('');
	let partition = $state('');
	let currentWorkingDirectory = $state('/home');
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
			script,
			partition: partition || undefined,
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

		await response.json().catch(() => null);
		submitting = false;

		if (response.ok) {
			goto(resolve('/jobs'));
		}
	}
</script>

<CreatePage
	title="Create Job"
	description="Submit a new batch job to Slurm"
	backHref="/jobs"
	backLabel="Back to Jobs"
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
				>Script
				<textarea bind:value={script} rows="10" class="field-text-area"></textarea>
				{#if errors.script}<span class="field-error">{errors.script}</span>{/if}
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

	.field-text-area {
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
