import { initContract } from "@ts-rest/core";
import { z } from "zod";
import { ActivityWithActorSchema, ActivitySummaryItemSchema } from "@vector/zod";

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
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "List issue activity timeline",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
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
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "List project activity feed",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    listMyActivity: {
        method: "GET",
        path: "/me/activity",
        query: z.object({
            page: z.coerce.number().min(1).optional().default(1),
            limit: z.coerce.number().min(1).max(100).optional().default(20),
            action: z.string().optional(),
            entity_type: z.string().optional(),
            from: z.string().datetime().optional(),
            to: z.string().datetime().optional(),
        }),
        responses: {
            200: z.object({
                data: z.array(ActivityWithActorSchema),
                total: z.number(),
                page: z.number(),
                limit: z.number(),
                total_pages: z.number(),
            }),
            401: z.object({ message: z.string() }),
        },
        summary: "List authenticated user's activity across all projects",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    listOrgActivity: {
        method: "GET",
        path: "/organizations/:orgId/activity",
        query: z.object(activityFilterQuery),
        responses: {
            200: z.object({
                data: z.array(ActivityWithActorSchema),
                total: z.number(),
                page: z.number(),
                limit: z.number(),
                total_pages: z.number(),
            }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "List organization activity feed (scoped to user's projects)",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    projectActivitySummary: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/activity/summary",
        query: z.object({
            group_by: z
                .enum(["action", "entity_type", "actor_id", "date"])
                .optional()
                .default("action"),
            interval: z
                .enum(["day", "week", "month"])
                .optional()
                .default("day"),
            ...activityFilterQuery,
        }),
        responses: {
            200: z.object({
                data: z.array(ActivitySummaryItemSchema),
                total_count: z.number(),
            }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "Get aggregated activity counts for a project",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    orgActivitySummary: {
        method: "GET",
        path: "/organizations/:orgId/activity/summary",
        query: z.object({
            group_by: z
                .enum(["action", "entity_type", "actor_id", "date"])
                .optional()
                .default("action"),
            interval: z
                .enum(["day", "week", "month"])
                .optional()
                .default("day"),
            ...activityFilterQuery,
        }),
        responses: {
            200: z.object({
                data: z.array(ActivitySummaryItemSchema),
                total_count: z.number(),
            }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "Get aggregated activity counts for an organization",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
