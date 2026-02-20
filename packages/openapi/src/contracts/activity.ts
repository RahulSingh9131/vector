import { initContract } from "@ts-rest/core";
import { z } from "zod";
import { ActivityWithActorSchema } from "@vector/zod";

const c = initContract();

export const activityContract = c.router({
    listIssueActivity: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/activity",
        query: z.object({
            page: z.coerce.number().min(1).optional().default(1),
            limit: z.coerce.number().min(1).max(100).optional().default(20),
        }),
        responses: {
            200: z.object({
                data: z.array(ActivityWithActorSchema),
                total: z.number(),
                page: z.number(),
                limit: z.number(),
                total_pages: z.number(),
            }),
        },
        summary: "List issue activity timeline",
    },
    listProjectActivity: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/activity",
        query: z.object({
            page: z.coerce.number().min(1).optional().default(1),
            limit: z.coerce.number().min(1).max(100).optional().default(20),
        }),
        responses: {
            200: z.object({
                data: z.array(ActivityWithActorSchema),
                total: z.number(),
                page: z.number(),
                limit: z.number(),
                total_pages: z.number(),
            }),
        },
        summary: "List project activity feed",
    },
});
