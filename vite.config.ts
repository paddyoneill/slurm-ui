import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { notebookWsProxy } from './vite-ws-proxy';

export default defineConfig({ plugins: [tailwindcss(), notebookWsProxy(), sveltekit()] });
