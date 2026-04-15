import { db } from '$lib/server/db';
import { notebook } from '$lib/server/db/schema';
import { notifyNotebooksChanged } from '$lib/server/events';
import { getJobStatus, submitJob } from '$lib/server/slurm';
import { createNotebookScheme } from '$lib/validation/notebook';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import crypto from 'crypto';
import * as v from 'valibot';

export const GET: RequestHandler = async () => {
	const notebooks = await db.select().from(notebook).orderBy(notebook.id);

	return json({ ok: true, value: notebooks }, { status: 200 });
};

export const POST: RequestHandler = async ({ request }) => {
	const body = await request.json();
	const parsed = v.safeParse(createNotebookScheme, body);
	const token = crypto.randomUUID();
	const notebookId = crypto.randomUUID();

	if (!parsed.success) {
		const errors = parsed.issues.map((issue) => ({
			field: issue.path?.[0]?.key,
			message: issue.message
		}));
		return json({ ok: false, error: 'Request failed validation', errors }, { status: 400 });
	}

	const input = parsed.output;

	const script = `#!/usr/bin/env bash
    mkdir -p ${input.currentWorkingDirectory}/notebooks/${notebookId}
    cp -r /shared/base-notebook-env ${input.currentWorkingDirectory}/notebooks/${notebookId}/venv
    source ${input.currentWorkingDirectory}/notebooks/${notebookId}/venv/bin/activate
    trap "rm -rf ${input.currentWorkingDirectory}/notebooks/${notebookId}" EXIT
    jupyter notebook --ip=0.0.0.0 --port=8888 --no-browser --ServerApp.token='${token}' --ServerApp.base_url=/api/notebooks/${notebookId}/proxy/ --ServerApp.allow_origin='*'`;

	const result = await submitJob({
		...input,
		script
	});

	if (!result.ok) {
		return json(result, { status: 400 });
	}

	const getStatusResult = await getJobStatus(result.value);
	if (!getStatusResult.ok) {
		return json(
			{ ok: false, error: `Failed to get job status: ${getStatusResult.error}` },
			{ status: 400 }
		);
	}

	await db.insert(notebook).values({
		id: notebookId,
		slurmJobId: result.value,
		state: getStatusResult.value,
		port: 8888,
		token
	});

	notifyNotebooksChanged();

	return json({ ok: true, value: { notebookId, token } });
};
