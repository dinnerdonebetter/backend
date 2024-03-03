import { z } from 'zod';

export const inputSlug = z
  .string()
  .trim()
  .min(1, 'slug is required')
  .regex(new RegExp(/^[a-zA-Z0-9\-]{1,}$/gm), 'must match expected URL slug pattern');
