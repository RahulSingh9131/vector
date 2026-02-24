import { initContract } from "@ts-rest/core";
import { z } from "zod";
import { ActivityWithActorSchema } from "@vector/zod";

const c = initContract();

const activityFilterQuery = {
    page: z.coerce.number().min(1).optional().default(1),
    limit: z.coerce.number().min(1).max(100).optional().default(20),
    action: z.string().optional(),
    entity_type: z.string().optional(),
    actor_id: z.string().uuid().optional(),
    from: z.string().datetime().optional(),
    to: z.string().datetime().optional(),
};

export const activityContract = c.router({
    listIssueActivity: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/activity",
        query: z.object(activityFilterQuery),
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
        query: z.object(activityFilterQuery),
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
