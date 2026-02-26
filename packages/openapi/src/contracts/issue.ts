import { initContract } from "@ts-rest/core";
import {
    IssueSchema,
    IssueWithDetailsSchema,
    CreateIssueSchema,
    UpdateIssueSchema,
    AssignIssueSchema,
    UpdateIssueStatusSchema,
    schemaWithPagination,
    z,
} from "@vector/zod";

const c = initContract();

export const issueContract = c.router({
    createIssue: {
        method: "POST",
        path: "/organizations/:orgId/projects/:projectId/issues",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        body: CreateIssueSchema,
        responses: {
            201: IssueSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Create a new issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    listIssues: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/issues",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        query: z.object({
            status: z.string().optional(),
            priority: z.string().optional(),
            type: z.string().optional(),
            assignee_id: z.string().uuid().optional(),
            page: z.coerce.number().min(1).optional(),
            limit: z.coerce.number().min(1).max(100).optional(),
        }),
        responses: {
            200: schemaWithPagination(IssueWithDetailsSchema),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "List issues for a project with optional filters",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    getIssue: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        responses: {
            200: IssueWithDetailsSchema,
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Get issue by ID with details",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateIssue: {
        method: "PATCH",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        body: UpdateIssueSchema,
        responses: {
            200: IssueSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Update an issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    deleteIssue: {
        method: "DELETE",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            204: z.object({}),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Delete an issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    assignIssue: {
        method: "PATCH",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/assign",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        body: AssignIssueSchema,
        responses: {
            200: IssueSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Assign or unassign an issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateIssueStatus: {
        method: "PATCH",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/status",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        body: UpdateIssueStatusSchema,
        responses: {
            200: IssueSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Update an issue's status",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
