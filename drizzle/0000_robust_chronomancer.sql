CREATE TABLE `slurm_job` (
	`id` text PRIMARY KEY NOT NULL,
	`job_id` integer NOT NULL,
	`host` text NOT NULL,
	`created_at` integer NOT NULL
);
