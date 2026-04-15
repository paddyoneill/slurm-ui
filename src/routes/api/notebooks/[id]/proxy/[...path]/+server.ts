import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { notebook } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import type { RequestEvent } from './$types';

async function proxyRequstHandler({ params, request, url }: RequestEvent) {
	if (!params.id) {
		return json({ ok: false, error: 'Missing ID of job' }, { status: 400 });
	}
	const [nb] = await db.select().from(notebook).where(eq(notebook.id, params.id));
	if (!nb) {
		return new Response('Notebook not found', { status: 404 });
	}

	if (!params.path) {
		return new Response(null, {
			status: 302,
			headers: { Location: `/api/notebooks/${nb.id}/proxy/tree${url.search}` }
		});
	}

	const target = `http://${nb.host}:${nb.port}${url.pathname}${url.search}`;

	try {
		const headers = new Headers(request.headers);
		headers.set('host', `${nb.host}:${nb.port}`);

		const fetchOptions: RequestInit = {
			method: request.method,
			headers
		};

		if (request.method !== 'GET' && request.method !== 'HEAD') {
			fetchOptions.body = await request.arrayBuffer();
		}

		const res = await fetch(target, fetchOptions);
		const responseHeaders = new Headers(res.headers);
		responseHeaders.delete('content-encoding');
		responseHeaders.delete('content-length');
		responseHeaders.delete('transfer-encoding');

		if (res.status === 304) {
			return new Response(null, {
				status: 304,
				headers: responseHeaders
			});
		}

		const body = await res.arrayBuffer();

		return new Response(body, {
			status: res.status,
			headers: responseHeaders
		});
	} catch (err) {
		console.error('Proxy error', request.method, url.pathname, err);
		return new Response('Failed to proxy notebook', {
			status: 502
		});
	}
}

export const DELETE = proxyRequstHandler;
export const GET = proxyRequstHandler;
export const PATCH = proxyRequstHandler;
export const POST = proxyRequstHandler;
export const PUT = proxyRequstHandler;
