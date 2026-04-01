import { z } from "./z.js";

export const SearchResultTypeSchema = z.enum(["issue", "comment"]).openapi("SearchResultType");

export const SearchResultSchema = z.object({
    id: z.string().uuid(),
    type: SearchResultTypeSchema,
    title: z.string(),
    description: z.string(),
    projectId: z.string().uuid().optional().nullable(),
    issueId: z.string().uuid().optional().nullable(),
    issueKey: z.string().optional().nullable(),
    rank: z.number(),
}).openapi("SearchResult");

export const SearchRequestSchema = z.object({
    query: z.string().min(1, "Search query is required"),
    type: SearchResultTypeSchema.optional(),
    projectId: z.string().uuid().optional(),
    limit: z.coerce.number().int().min(1).max(100).default(20),
    offset: z.coerce.number().int().min(0).default(0),
}).openapi("SearchRequest");

export const SearchResponseSchema = z.array(SearchResultSchema).openapi("SearchResponse");
