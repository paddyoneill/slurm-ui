import { client } from './client';
import type { Result } from '../../types';
import type { CreateJobInput } from '$lib/validation/job';

export async function cancelJob(jobId: number): Promise<Result<string>> {
	const { error } = await client.DELETE('/slurm/v0.0.43/job/{job_id}', {
		params: { path: { job_id: jobId.toString() } }
	});

	if (error) {
		return { ok: false, error: `failed to cancel job: ${JSON.stringify(error)}` };
	}

	return { ok: true, value: `Job ${jobId} cancelled` };
}

export async function getJobStatus(jobId: number): Promise<Result<string>> {
	const { data, error } = await client.GET('/slurm/v0.0.43/job/{job_id}', {
		params: { path: { job_id: jobId.toString() } }
	});

	if (error) {
		return { ok: false, error: `failed to get job: ${JSON.stringify(error)}` };
	}

	const state = data?.jobs?.[0]?.job_state?.[0];

	if (state == null) {
		return { ok: false, error: 'missing job state in response' };
	}

	return { ok: true, value: state };
}

export async function submitJob(input: CreateJobInput): Promise<Result<number>> {
	const environment = buildEnvironment(input.environment);

	const { data, error } = await client.POST('/slurm/v0.0.43/job/submit', {
		body: {
			script: input.script,
			job: {
				name: input.name,
				partition: input.parition,
				current_working_directory: input.currentWorkingDirectory,
				environment,
				cpus_per_task: input.cpusPerTask,
				tasks_per_node: input.tasksPerNode,
				memory_per_node: input.memoryPerNode,
				time_limit: input.timeLimit
			}
		}
	});

	if (error) {
		return { ok: false, error: `failed to submit job: ${JSON.stringify(error)}` };
	}

	if (!data?.job_id) {
		return { ok: false, error: 'missing job_id in response' };
	}

	return { ok: true, value: data.job_id };
}

const DEFAULT_PATH = 'PATH=/bin:/usr/bin:/usr/local/bin';

function buildEnvironment(userEnv?: string[]): string[] {
	const env = userEnv ?? [];
	const hasPath = env.some((e) => e.startsWith('PATH='));

	if (!hasPath) {
		return [DEFAULT_PATH, ...env];
	}

	return env;
}
