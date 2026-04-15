import { integer, sqliteTable, text } from 'drizzle-orm/sqlite-core';

export const notebook = sqliteTable('notebook', {
	id: text('id').primaryKey(),
	slurmJobId: integer('slurm_job_id').notNull(),
	state: text('state').notNull(),
	host: text('host'),
	port: integer('port').notNull(),
	token: text('token').notNull(),
	createdAt: integer('created_at', { mode: 'timestamp' })
		.notNull()
		.$defaultFn(() => new Date())
});

export const slurmJob = sqliteTable('slurm_job', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => crypto.randomUUID()),
	jobId: integer('job_id').notNull(),
	state: text('state').notNull(),
	createdAt: integer('created_at', { mode: 'timestamp' })
		.notNull()
		.$defaultFn(() => new Date())
});
