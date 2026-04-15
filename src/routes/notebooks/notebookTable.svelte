<script lang="ts">
	import type { NotebookSelect } from '$lib/server/db/types';

	let {
		notebooks,
		loading
	}: {
		notebooks: NotebookSelect[];
		loading: boolean;
	} = $props();

	let deletingIds = $state<Set<string>>(new Set());

	async function deleteNotebook(id: string) {
		deletingIds = new Set([...deletingIds, id]);
		const response = await fetch(`/api/notebooks/${id}`, { method: 'DELETE' });
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
					<th class="px-4 py-3">Job State</th>
					<th class="px-4 py-3"></th>
					<th class="px-4 py-3"></th>
				</tr>
			</thead>
			<tbody class="whitespace-nowrap">
				{#if loading}
					<tr><td colspan="5" class="text-center text-2xl font-medium">Loading notebooks...</td></tr
					>
				{:else}
					{#each notebooks as { id, state, host, token }, i (id)}
						<tr class="border-t border-gray-100 {i % 2 === 0 ? 'bg-white' : 'bg-gray-50'}">
							<td class="px-4 py-3">{id}</td>
							<td class="px-4 py-3">{state}</td>
							<td class="px-4 py-3"
								><button
									class="h-10 w-24 rounded bg-teal-600 px-4 py-2 font-semibold text-white hover:bg-teal-700 disabled:cursor-not-allowed disabled:opacity-50"
									disabled={!host}
									onclick={() =>
										window.open(`/api/notebooks/${id}/proxy/?token=${token}`, '_blank')}
									>{#if !host}
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
										Connect
									{/if}</button
								></td
							>

							<td class="px-4 py-3"
								><button
									class="h-10 w-20 rounded bg-red-600 px-4 py-2 font-semibold text-white hover:bg-red-700"
									disabled={deletingIds.has(id)}
									onclick={() => {
										deleteNotebook(id);
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
