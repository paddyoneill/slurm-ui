import * as v from 'valibot';
import { baseJobFields } from './base';

export const createNotebookScheme = v.object({ ...baseJobFields });

export type CreateNotebookInput = v.InferOutput<typeof createNotebookScheme>;
