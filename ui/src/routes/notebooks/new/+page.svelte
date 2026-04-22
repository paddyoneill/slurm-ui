<script lang="ts">
	import JobForm from '$lib/components/JobForm.svelte';
	import { createNotebookScheme } from '$lib/validation/notebook';
	import * as v from 'valibot';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

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

		await response.json().catch(() => null);
		submitting = false;

		if (response.ok) {
			goto(resolve('/notebooks'));
		}
	}
</script>

<section class="page-copy">
	<div>
		<h2>Create Notebook</h2>
		<p>Submit a new Jupyter notebook server</p>
	</div>
	<a href={resolve('/notebooks')} class="page-link">Back to Notebooks</a>
</section>

<div class="form-shell">
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
			<label class="field-label"
				>Base Environment
				<input type="text" bind:value={baseEnv} class="field-input" />
				{#if errors.baseEnv}<span class="field-error">{errors.baseEnv}</span>{/if}
			</label>
		</div>
	</JobForm>
</div>

<style>
	.page-copy {
		display: flex;
		align-items: flex-start;
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

	.page-link {
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

	.form-shell {
		display: flex;
		align-items: flex-start;
	}

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
		color: #dc2626;
		font-width: 0.875rem;
	}

	@media (max-width: 760px) {
		.page-copy {
			flex-direction: column;
		}
	}
</style>
