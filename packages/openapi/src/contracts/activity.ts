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
        },
        summary: "List authenticated user's activity across all projects",
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
        },
        summary: "List organization activity feed (scoped to user's projects)",
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
        },
        summary: "Get aggregated activity counts for a project",
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
        },
        summary: "Get aggregated activity counts for an organization",
    },
});
