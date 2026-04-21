<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		name = $bindable(),
		parition = $bindable(),
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
		parition: string;
		currentWorkingDirectory: string;
		cpusPerTask: number;
		tasksPerNode: number;
		memoryPerNode: number;
		timeLimit: number;
		envVars: { key: string; value: string }[];
		errors: Record<string, string>;
		onSubmit: () => void;
		submitting: boolean;
		children?: Snippet;
	} = $props();
</script>

<form
	class="relative w-full max-w-lg p-6"
	onsubmit={(e) => {
		e.preventDefault();
		onSubmit();
	}}
>
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
					{#if errors.memoryPerNode}<span class="text-sm text-red-500">{errors.memoryPerNode}</span
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
		{@render children?.()}
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
