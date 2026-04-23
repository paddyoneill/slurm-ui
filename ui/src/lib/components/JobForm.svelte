<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		name = $bindable(),
		partition = $bindable(),
		currentWorkingDirectory = $bindable(),
		cpusPerTask = $bindable(),
		tasksPerNode = $bindable(),
		memoryPerNode = $bindable(),
		timeLimit = $bindable(),
		envVars = $bindable(),
		errors,
		onSubmit,
		submitting,
		children
	}: {
		name: string;
		partition: string;
		currentWorkingDirectory: string;
		cpusPerTask: number;
		tasksPerNode: number;
		memoryPerNode: number;
		timeLimit: number;
		envVars: { key: string; value: string }[];
		errors: Record<string, string>;
		onSubmit: () => Promise<void> | void;
		submitting: boolean;
		children?: Snippet;
	} = $props();
</script>

<form
	class="form"
	onsubmit={(event) => {
		event.preventDefault();
		onSubmit();
	}}
>
	<div class="fields">
		<div class="field">
			<label class="field-label" for="name">Name</label>
			<input id="name" class="field-input" type="text" bind:value={name} />
			{#if errors.name}<p class="field-error">{errors.name}</p>{/if}
		</div>

		<div class="field">
			<label class="field-label" for="currentWorkingDirectory">Working Directory</label>
			<input
				id="currentWorkingDirectory"
				class="field-input"
				type="text"
				bind:value={currentWorkingDirectory}
			/>
			{#if errors.currentWorkingDirectory}
				<p class="field-error">{errors.currentWorkingDirectory}</p>
			{/if}
		</div>

		<div class="field">
			<label class="field-label" for="partition">Partition</label>
			<input id="partition" class="field-input" type="text" bind:value={partition} />
			{#if errors.partition}<p class="field-error">{errors.partition}</p>{/if}
		</div>

		<div class="field-grid">
			<div class="field">
				<label class="field-label" for="tasksPerNode">Tasks Per Node</label>
				<input id="tasksPerNode" class="field-input" type="number" bind:value={tasksPerNode} />
				{#if errors.tasksPerNode}<p class="field-error">{errors.tasksPerNode}</p>{/if}
			</div>

			<div class="field">
				<label class="field-label" for="cpusPerTask">CPUs Per Task</label>
				<input id="cpusPerTask" class="field-input" type="number" bind:value={cpusPerTask} />
				{#if errors.cpusPerTask}<p class="field-error">{errors.cpusPerTask}</p>{/if}
			</div>

			<div class="field">
				<label class="field-label" for="memoryPerNode">Memory Per Node</label>
				<input id="memoryPerNode" class="field-input" type="number" bind:value={memoryPerNode} />
				{#if errors.memoryPerNode}<p class="field-error">{errors.memoryPerNode}</p>{/if}
			</div>

			<div class="field">
				<label class="field-label" for="timeLimit">Time Limit</label>
				<input id="timeLimit" class="field-input" type="number" bind:value={timeLimit} />
				{#if errors.timeLimit}<p class="field-error">{errors.timeLimit}</p>{/if}
			</div>
		</div>

		<fieldset class="env-fieldset">
			<legend class="field-label">Environment Variables</legend>
			<div class="env-list">
				{#each envVars as env, i (i)}
					<div class="env-row">
						<input
							class="field-input env-input"
							type="text"
							bind:value={env.key}
							placeholder="KEY"
						/>
						<input
							class="field-input env-input"
							type="text"
							bind:value={env.value}
							placeholder="VALUE"
						/>
						<button
							class="env-remove"
							type="button"
							onclick={() => (envVars = envVars.filter((_, j) => j !== i))}
						>
							Remove
						</button>
					</div>
				{/each}
			</div>

			{#if errors.environment}
				<p class="field-error">{errors.environment}</p>
			{/if}

			<button
				class="env-add"
				type="button"
				onclick={() => (envVars = [...envVars, { key: '', value: '' }])}
			>
				Add variable
			</button>
		</fieldset>

		{@render children?.()}
	</div>

	<div class="actions">
		<button class="submit-button" disabled={submitting} type="submit">
			{submitting ? 'Submitting...' : 'Submit'}
		</button>
	</div>
</form>

<style>
	.form {
		display: grid;
		gap: var(--space-5);
		width: min(100%, 44rem);
		padding: var(--space-5);
	}

	.fields {
		display: grid;
		gap: var(--space-5);
	}

	.field,
	.env-fieldset {
		display: grid;
		gap: var(--space-2);
	}

	.field-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: var(--space-4);
	}

	.field-label {
		font-size: 0.95rem;
		font-weight: 700;
	}

	.field-input {
		width: 100%;
		padding: 0.8rem 0.95rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-surface);
		box-shadow: var(--shadow-sm);
		transition:
			border-color var(--transition-fast),
			box-shadow var(--transition-fast);
	}

	.field-input:focus {
		border-color: var(--color-primary);
		outline: none;
		box-shadow: 0 0 0 3px rgb(15 118 110 / 0.14);
	}

	.field-error {
		margin: 0;
		color: var(--color-danger);
		font-size: 0.92rem;
	}

	.env-fieldset {
		margin: 0;
		padding: var(--space-4);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-surface-muted);
	}

	.env-list {
		display: grid;
		gap: var(--space-3);
	}

	.env-row {
		display: grid;
		grid-template-columns: minmax(0, 1fr) minmax(0, 1fr) auto;
		gap: var(--space-3);
		align-items: center;
	}

	.env-input {
		min-width: 0;
	}

	.env-add,
	.env-remove,
	.submit-button {
		border: 0;
		border-radius: var(--radius-md);
		font-weight: 700;
		transition:
			background-color var(--transition-fast),
			opacity var(--transition-fast);
	}

	.env-add,
	.env-remove {
		padding: 0.72rem 0.95rem;
	}

	.env-add {
		justify-self: start;
		background: transparent;
		color: var(--color-primary);
	}

	.env-add:hover {
		background: rgb(15 118 110 / 0.08);
	}

	.env-remove {
		background: rgb(201 66 66 / 0.1);
		color: var(--color-danger);
	}

	.env-remove:hover {
		background: rgb(201 66 66 / 0.18);
	}

	.actions {
		display: flex;
		justify-content: flex-start;
	}

	.submit-button {
		min-width: 11rem;
		padding: 0.9rem 1.25rem;
		background: var(--color-success);
		color: #fff;
	}

	.submit-button:hover:not(:disabled) {
		background: #166534;
	}

	.submit-button:disabled {
		opacity: 0.65;
	}

</style>
