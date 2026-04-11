import { getJobStatus } from '$lib/server/slurm';
import { db } from '$lib/server/db';
import { slurmJob } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import { notifyJobsChanged } from '../events';

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

export function startSlurmPoller() {
	const interval = setInterval(pollSlurmJobs, 10000);
	return () => clearInterval(interval);
}
