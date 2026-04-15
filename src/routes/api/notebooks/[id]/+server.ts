import { json } from '@sveltejs/kit';
import { type RequestHandler } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { eq } from 'drizzle-orm';
import { notebook } from '$lib/server/db/schema';
import { notifyNotebooksChanged } from '$lib/server/events';
import { cancelJob } from '$lib/server/slurm/jobs';

export const DELETE: RequestHandler = async ({ params }) => {
	if (!params.id) {
		return json({ ok: false, error: 'Missing ID of notebook' }, { status: 400 });
	}

	const [nb] = await db.select().from(notebook).where(eq(notebook.id, params.id));
	const cancelResult = await cancelJob(nb.slurmJobId);
	if (!cancelResult.ok) {
		return json({ ok: false, error: cancelResult.error }, { status: 400 });
	}

	await db.delete(notebook).where(eq(notebook.id, params.id));
	notifyNotebooksChanged();

	return json({ ok: true }, { status: 200 });
};
