import type { Plugin } from 'vite';
import httpProxy from 'http-proxy';
import Database from 'better-sqlite3';

interface NotebookRow {
	host: string;
	port: number;
}

export function notebookWsProxy(): Plugin {
	const proxy = httpProxy.createProxyServer({ ws: true });
	const db = new Database('./local.db');

	proxy.on('error', (err) => {
		console.error('WS proxy error', err);
	});

	return {
		name: 'notebook-ws-proxy',
		configureServer(server) {
			server.httpServer?.on('upgrade', (req, socket, head) => {
				const match = req.url?.match(/^\/api\/notebooks\/([^/]+)\/proxy\/(.*)/);

				if (match) {
					const notebookId = match[1];
					const nb = db.prepare('SELECT host, port FROM notebook WHERE id = ?').get(notebookId) as
						| NotebookRow
						| undefined;

					if (!nb) {
						socket.destroy();
						return;
					}

					proxy.ws(
						req,
						socket,
						head,
						{ target: `http://${nb.host}:${nb.port}`, changeOrigin: true, ws: true },
						(err) => {
							if (err) {
								console.error('WS proxy error', err);
								socket.destroy();
							}
						}
					);
				}
			});
		}
	};
}
