import type { RequestHandler } from '@sveltejs/kit';
import { onJobsChanged } from '$lib/server/events';
import { db } from '$lib/server/db';
import { slurmJob } from '$lib/server/db/schema';

export const GET: RequestHandler = () => {
	const encoder = new TextEncoder();
	let closed = false;

	const stream = new ReadableStream({
		start(controller) {
			const send = async () => {
				if (closed) return;

				try {
					const jobs = await db.select().from(slurmJob).orderBy(slurmJob.createdAt);
					controller.enqueue(encoder.encode(`data: ${JSON.stringify(jobs)}\n\n`));
				} catch {
					closed = true;
					unsubscribe();
				}
			};

			const unsubscribe = onJobsChanged(send);

			return () => {
				closed = true;
				unsubscribe();
			};
		},
		cancel() {
			closed = true;
		}
	});

	return new Response(stream, {
		headers: {
			'Content-Type': 'text/event-stream',
			'Cache-Control': 'no-cache',
			Connection: 'keep-alive'
		}
	});
};
