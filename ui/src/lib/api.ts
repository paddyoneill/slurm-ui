export type ApiJob = {
	id: string;
	slurmJobId: number;
	name: string;
	state: string;
	createdAt: string;
	script?: string;
	currentWorkingDirectory?: string;
	partition?: string;
	environment?: string[];
	cpusPerTask?: number;
	tasksPerNode?: number;
	memoryPerNode?: number;
	timeLimit?: number;
};

export type ApiNotebook = {
	id: string;
	slurmJobId: number;
	name: string;
	state: string;
	host?: string | null;
	port: number;
	token: string;
	baseEnv: string;
	createdAt: string;
	script?: string;
	currentWorkingDirectory?: string;
	partition?: string;
	environment?: string[];
	cpusPerTask?: number;
	tasksPerNode?: number;
	memoryPerNode?: number;
	timeLimit?: number;
};
