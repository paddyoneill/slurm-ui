import * as v from 'valibot';
import { baseJobFields } from './base';

export const createNotebookScheme = v.object({
	...baseJobFields,
	baseEnv: v.pipe(v.string(), v.trim(), v.minLength(1, 'Base Environment is required'))
});

export type CreateNotebookInput = v.InferOutput<typeof createNotebookScheme>;
