import type { RequestHandler } from '@sveltejs/kit';
import { onNotebooksChanged } from '$lib/server/events';
import { db } from '$lib/server/db';
import { notebook } from '$lib/server/db/schema';

export const GET: RequestHandler = () => {
	const encoder = new TextEncoder();
	let closed = false;

	const stream = new ReadableStream({
		start(controller) {
			const send = async () => {
				if (closed) return;

				try {
					const notebooks = await db.select().from(notebook).orderBy(notebook.createdAt);
					controller.enqueue(encoder.encode(`data: ${JSON.stringify(notebooks)}\n\n`));
				} catch {
					closed = true;
					unsubscribe();
				}
			};

			const unsubscribe = onNotebooksChanged(send);

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
