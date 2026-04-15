import * as v from 'valibot';

export const baseJobFields = {
	name: v.pipe(v.string(), v.trim(), v.minLength(1, 'Job name is required')),
	parition: v.optional(v.pipe(v.string(), v.minLength(1, 'Parition cannot be empty'))),
	currentWorkingDirectory: v.pipe(
		v.string(),
		v.trim(),
		v.minLength(1, 'Working directory is required'),
		v.startsWith('/', 'Must be an absolute path')
	),
	environment: v.optional(
		v.array(v.pipe(v.string(), v.regex(/^\S+=.*$/, 'Key cannot contain spaces')))
	),
	cpusPerTask: v.optional(
		v.pipe(
			v.number(),
			v.integer('Must be an integer'),
			v.minValue(1, 'Must request at least 1 CPU')
		)
	),
	tasksPerNode: v.optional(
		v.pipe(
			v.number(),
			v.integer('Must be an integer'),
			v.minValue(1, 'Must request at least 1 task')
		)
	),
	memoryPerNode: v.optional(
		v.pipe(
			v.number(),
			v.integer('Must be an integer'),
			v.minValue(1, 'Memory must be at least 1 MB')
		)
	),
	timeLimit: v.optional(
		v.pipe(
			v.number(),
			v.integer('Must be an integer'),
			v.minValue(1, 'Time limit must be at least 1 minute')
		)
	)
};
