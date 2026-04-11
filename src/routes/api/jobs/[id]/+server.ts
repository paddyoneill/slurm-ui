import { json } from '@sveltejs/kit';
import { type RequestHandler } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { eq } from 'drizzle-orm';
import { slurmJob } from '$lib/server/db/schema';
import { notifyJobsChanged } from '$lib/server/events';
import { cancelJob } from '$lib/server/slurm/jobs';

export const DELETE: RequestHandler = async ({ params }) => {
	if (!params.id) {
		return json({ ok: false, error: 'Missing ID of job' }, { status: 400 });
	}

	const [job] = await db.select().from(slurmJob).where(eq(slurmJob.id, params.id));
	const cancelResult = await cancelJob(job.jobId);
	if (!cancelResult.ok) {
		return json({ ok: false, error: cancelResult.error }, { status: 400 });
	}

	await db.delete(slurmJob).where(eq(slurmJob.id, params.id));
	notifyJobsChanged();

	return json({ ok: true }, { status: 200 });
};
