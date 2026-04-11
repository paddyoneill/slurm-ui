import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { db } from '$lib/server/db';
import { slurmJob } from '$lib/server/db/schema';
import { type SlurmJobInsert } from '$lib/server/db/schema';
import * as v from 'valibot';
import { createJobSchema } from '$lib/validation/job';
import { getJobStatus, submitJob } from '$lib/server/slurm';
import { notifyJobsChanged } from '$lib/server/events';

export const GET: RequestHandler = async () => {
	const jobs = await db.select().from(slurmJob).orderBy(slurmJob.createdAt);

	return json({ ok: true, value: jobs }, { status: 200 });
};

export const POST: RequestHandler = async ({ request }) => {
	const body = await request.json();
	const parsed = v.safeParse(createJobSchema, body);

	if (!parsed.success) {
		const errors = parsed.issues.map((issue) => ({
			field: issue.path?.[0]?.key,
			message: issue.message
		}));
		return json({ ok: false, error: 'Request failed validation', errors }, { status: 400 });
	}

	const input = parsed.output;

	const submitResult = await submitJob(input);
	if (!submitResult.ok) {
		return json(
			{ ok: false, error: `Failed to submit job: ${submitResult.error}` },
			{ status: 400 }
		);
	}

	const jobId = submitResult.value;
	const getStatusResult = await getJobStatus(jobId);
	if (!getStatusResult.ok) {
		return json(
			{ ok: false, error: `Failed to get job status: ${getStatusResult.error}` },
			{ status: 400 }
		);
	}

	const job: SlurmJobInsert = {
		jobId,
		state: getStatusResult.value
	};

	await db.insert(slurmJob).values(job);
	notifyJobsChanged();

	return json({ ok: true }, { status: 200 });
};
