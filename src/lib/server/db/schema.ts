import { integer, sqliteTable, text } from 'drizzle-orm/sqlite-core';

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

export type SlurmJobInsert = typeof slurmJob.$inferInsert;
export type SlurmJobSelect = typeof slurmJob.$inferSelect;
