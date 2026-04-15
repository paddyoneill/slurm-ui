CREATE TABLE `notebook` (
	`id` text PRIMARY KEY NOT NULL,
	`slurm_job_id` integer NOT NULL,
	`state` text NOT NULL,
	`host` text,
	`port` integer NOT NULL,
	`token` text NOT NULL,
	`created_at` integer NOT NULL
);
--> statement-breakpoint
ALTER TABLE `slurm_job` ADD `state` text NOT NULL;--> statement-breakpoint
ALTER TABLE `slurm_job` DROP COLUMN `host`;