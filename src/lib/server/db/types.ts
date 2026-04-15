import { notebook, slurmJob } from './schema';

export type NotebookInsert = typeof notebook.$inferInsert;
export type NotebookSelect = typeof notebook.$inferSelect;
export type SlurmJobInsert = typeof slurmJob.$inferInsert;
export type SlurmJobSelect = typeof slurmJob.$inferSelect;
