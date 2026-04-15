<script lang="ts">
	import type { SlurmJobSelect } from '$lib/server/db/types';
	import { formatAge } from '$lib/date';
	let {
		jobs,
		loading
	}: {
		jobs: SlurmJobSelect[];
		loading: boolean;
	} = $props();

	let deletingIds = $state<Set<string>>(new Set());

	async function deleteJob(id: string) {
		deletingIds = new Set([...deletingIds, id]);
		const response = await fetch(`/api/jobs/${id}`, { method: 'DELETE' });
		const result = await response.json();

		if (!result.ok) {
			deletingIds = new Set([...deletingIds].filter((i) => i !== id));
		}
	}
</script>

<div
	class="w-full max-w-3xl overflow-hidden overflow-x-auto rounded-lg border border-gray-200 shadow-sm"
>
	<div class="h-full min-w-160 overflow-y-auto">
		<table class="w-full">
			<thead class="sticky top-0 z-50 bg-gray-100 text-left text-sm font-semibold text-gray-700">
				<tr>
					<th class="px-4 py-3">ID</th>
					<th class="px-4 py-3">Job ID</th>
					<th class="px-4 py-3">State</th>
					<th class="px-4 py-3">Age</th>
					<th class="px-4 py-3"><span class="sr-only">Delete</span></th>
				</tr>
			</thead>
			<tbody class="whitespace-nowrap">
				{#if loading}
					<tr><td colspan="5" class="text-center text-2xl font-medium">Loading jobs...</td></tr>
				{:else}
					{#each jobs as { id, state, jobId, createdAt }, i (id)}
						<tr
							class="border-t border-gray-100 {i % 2 === 0 ? 'bg-white' : 'bg-gray-50'}
								{deletingIds.has(id) ? 'pointer-events-none opacity-50' : ''}"
						>
							<td class="px-4 py-3">{id}</td>
							<td class="px-4 py-3">{jobId}</td>
							<td class="px-4 py-3">{state}</td>
							<td class="px-4 py-3">
								<span class="group relative cursor-default">
									{formatAge(createdAt)}
									<span
										class="absolute bottom-full left-0 z-50 hidden rounded-md bg-gray-900 px-3 py-2 text-xs text-white shadow-lg group-hover:block"
									>
										{createdAt.toLocaleString()}
									</span>
								</span>
							</td>
							<td class="w-[20%] px-4 py-3"
								><button
									class="h-10 w-20 rounded bg-red-600 px-4 py-2 font-semibold text-white hover:bg-red-700"
									disabled={deletingIds.has(id)}
									onclick={() => {
										deleteJob(id);
									}}
									>{#if deletingIds.has(id)}
										<svg class="mx-auto h-4 w-4 animate-spin" viewBox="0 0 24 24">
											<circle
												class="opacity-25"
												cx="12"
												cy="12"
												r="10"
												stroke="currentColor"
												stroke-width="4"
												fill="none"
											/>
											<path
												class="opacity-75"
												fill="currentColor"
												d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
											/>
										</svg>
									{:else}
										Delete
									{/if}</button
								></td
							>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>
