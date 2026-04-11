import type { paths } from './v0.0.43';
import createClient from 'openapi-fetch';
import { env } from '$env/dynamic/private';

if (!env.SLURM_RESTAPI_URL) throw new Error('SLURM_RESTAPI_URL is not set');

export const client = createClient<paths>({ baseUrl: env.SLURM_RESTAPI_URL });
