<script lang="ts">
	import type { ApiNotebook } from '$lib/api';

	let {
		notebooks,
		loading
	}: {
		notebooks: ApiNotebook[];
		loading: boolean;
	} = $props();

	let deletingIds = $state<Set<string>>(new Set());

	async function deleteNotebook(id: string) {
		deletingIds = new Set([...deletingIds, id]);
		const response = await fetch(`/api/notebooks/${id}`, { method: 'DELETE' });

		if (!response.ok && response.status !== 204) {
			deletingIds = new Set([...deletingIds].filter((item) => item !== id));
		}
	}
</script>

<section class="table-shell">
	<div class="table-scroll">
		<table>
			<thead>
				<tr>
					<th>ID</th>
					<th>State</th>
					<th>Target</th>
					<th><span class="sr-only">Connect</span></th>
					<th><span class="sr-only">Delete</span></th>
				</tr>
			</thead>
			<tbody>
				{#if loading}
					<tr>
						<td class="empty" colspan="5">Loading notebooks...</td>
					</tr>
				{:else if notebooks.length === 0}
					<tr>
						<td class="empty" colspan="5">No notebooks launched yet.</td>
					</tr>
				{:else}
					{#each notebooks as { id, state, host, port, token }, i (id)}
						<tr class:alt={i % 2 === 1} class:muted={deletingIds.has(id)}>
							<td class="mono">{id}</td>
							<td>{state}</td>
							<td>{host && port ? `${host}:${port}` : 'Waiting for registration'}</td>
							<td class="action-cell">
								<button
									class="primary-button"
									disabled={!host || !port}
									onclick={() => {
										window.open(
											`/api/notebooks/${id}/proxy/tree?token=${encodeURIComponent(token)}`,
											'_blank',
											'noopener,noreferrer'
										);
									}}
								>
									{host && port ? 'Connect' : 'Starting...'}
								</button>
							</td>
							<td class="action-cell">
								<button
									class="danger-button"
									disabled={deletingIds.has(id)}
									onclick={() => {
										void deleteNotebook(id);
									}}
								>
									{deletingIds.has(id) ? 'Deleting...' : 'Delete'}
								</button>
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</section>

<style>
	.table-shell {
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		background: rgb(255 255 255 / 0.94);
		box-shadow: var(--shadow-md);
		overflow: hidden;
	}

	.table-scroll {
		overflow-x: auto;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		min-width: 56rem;
	}

	th,
	td {
		padding: 1rem 1.1rem;
		border-top: 1px solid var(--color-border);
		text-align: left;
	}

	thead th {
		border-top: 0;
		background: var(--color-surface-muted);
		color: var(--color-text-muted);
		font-size: 0.82rem;
		font-weight: 800;
		letter-spacing: 0.04em;
		text-transform: uppercase;
	}

	tbody tr.alt {
		background: rgb(248 250 252 / 0.7);
	}

	tbody tr.muted {
		opacity: 0.6;
	}

	.empty {
		padding: 2.5rem 1rem;
		color: var(--color-text-muted);
		text-align: center;
	}

	.mono {
		font-family: var(--font-mono);
		font-size: 0.9rem;
	}

	.action-cell {
		width: 1%;
		white-space: nowrap;
	}

	.primary-button,
	.danger-button {
		padding: 0.72rem 0.95rem;
		border: 0;
		border-radius: var(--radius-md);
		font-weight: 700;
	}

	.primary-button {
		background: rgb(15 118 110 / 0.12);
		color: var(--color-primary);
	}

	.primary-button:hover:not(:disabled) {
		background: rgb(15 118 110 / 0.2);
	}

	.danger-button {
		background: rgb(201 66 66 / 0.12);
		color: var(--color-danger);
	}

	.danger-button:hover:not(:disabled) {
		background: rgb(201 66 66 / 0.2);
	}

	button:disabled {
		opacity: 0.6;
	}

	.sr-only {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}
</style>
