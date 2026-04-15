import * as v from 'valibot';
import { baseJobFields } from './base';

export const createJobSchema = v.object({
	...baseJobFields,
	script: v.pipe(v.string(), v.trim(), v.minLength(1, 'Script is required'))
});

export type CreateJobInput = v.InferOutput<typeof createJobSchema>;
