import { client } from './client';
import type { Result } from '../../types';
import type { CreateJobInput } from '$lib/validation/job';
import { env } from '$env/dynamic/private';

type SlurmJob = {
	jobId: number;
	state: string;
	host?: string;
	name?: string;
};

export async function cancelJob(jobId: number): Promise<Result<string>> {
	const { error } = await client.DELETE('/slurm/v0.0.43/job/{job_id}', {
		headers: {
			'X-SLURM-USER-TOKEN': env.SLURM_JWT
		},
		params: { path: { job_id: jobId.toString() } }
	});

	if (error) {
		return { ok: false, error: `failed to cancel job: ${JSON.stringify(error)}` };
	}

	return { ok: true, value: `Job ${jobId} cancelled` };
}

export async function getJob(jobId: number): Promise<Result<SlurmJob>> {
	const { data, error } = await client.GET('/slurm/v0.0.43/job/{job_id}', {
		headers: {
			'X-SLURM-USER-TOKEN': env.SLURM_JWT
		},
		params: { path: { job_id: jobId.toString() } }
	});

	if (error) {
		return { ok: false, error: `failed to get job: ${JSON.stringify(error)}` };
	}

	const job = data?.jobs?.[0];

	if (!job) return { ok: false, error: `failed to get data for job ${jobId}` };

	return {
		ok: true,
		value: {
			jobId: job.job_id ?? jobId,
			state: job.job_state?.[0] ?? 'UNKNOWN',
			host: job.nodes ?? undefined,
			name: job.name ?? undefined
		}
	};
}

export async function getJobStatus(jobId: number): Promise<Result<string>> {
	const result = await getJob(jobId);

	if (!result.ok) return result;

	return { ok: true, value: result.value.state };
}

export async function submitJob(input: CreateJobInput): Promise<Result<number>> {
	const environment = buildEnvironment(input.environment);

	const { data, error } = await client.POST('/slurm/v0.0.43/job/submit', {
		headers: {
			'X-SLURM-USER-TOKEN': env.SLURM_JWT
		},
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

function buildEnvironment(userEnv?: string[]): string[] {
	const DEFAULT_PATH = 'PATH=/bin:/usr/bin:/usr/local/bin';
	const env = userEnv ?? [];
	const hasPath = env.some((e) => e.startsWith('PATH='));

	if (!hasPath) {
		return [DEFAULT_PATH, ...env];
	}

	return env;
}
