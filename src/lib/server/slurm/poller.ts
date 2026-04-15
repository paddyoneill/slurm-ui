import { getJob, getJobStatus } from '$lib/server/slurm';
import { db } from '$lib/server/db';
import { notebook, slurmJob } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import { notifyJobsChanged, notifyNotebooksChanged } from '../events';
import type { NotebookInsert } from '../db/types';

async function pollSlurmJobs() {
	let updated = false;

	const jobs = await db.select().from(slurmJob);
	for (const job of jobs) {
		const result = await getJobStatus(job.jobId);

		if (!result.ok) continue;

		if (result.value !== job.state) {
			await db.update(slurmJob).set({ state: result.value }).where(eq(slurmJob.id, job.id));
			updated = true;
		}
	}

	if (updated) notifyJobsChanged();
}

async function pollNotebooks() {
	let updated = false;

	const notebooks = await db.select().from(notebook);
	for (const nb of notebooks) {
		const result = await getJob(nb.slurmJobId);

		if (!result.ok) continue;

		const updates: Partial<NotebookInsert> = {};

		if (result.value.host !== nb.host) {
			updates.host = result.value.host;
		}

		if (result.value.state != nb.state) {
			updates.state = result.value.state;
		}

		if (Object.keys(updates).length > 0) {
			await db.update(notebook).set(updates).where(eq(notebook.id, nb.id));
			updated = true;
		}
	}

	if (updated) notifyNotebooksChanged();
}

export function startSlurmPoller() {
	const jobPoller = setInterval(pollSlurmJobs, 10000);
	const notebookPoller = setInterval(pollNotebooks, 5000);
	return () => {
		clearInterval(jobPoller);
		clearInterval(notebookPoller);
	};
}
