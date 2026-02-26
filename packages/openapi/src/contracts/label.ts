import { initContract } from "@ts-rest/core";
import {
    LabelSchema,
    CreateLabelSchema,
    UpdateLabelSchema,
    AddLabelToIssueSchema,
    z,
} from "@vector/zod";

const c = initContract();

export const labelContract = c.router({
    // --- Label CRUD ---

    createLabel: {
        method: "POST",
        path: "/organizations/:orgId/projects/:projectId/labels",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        body: CreateLabelSchema,
        responses: {
            201: LabelSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Create a new label",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    listLabels: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/labels",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        responses: {
            200: z.array(LabelSchema),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "List all labels for a project",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    getLabel: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/labels/:labelId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            labelId: z.string().uuid(),
        }),
        responses: {
            200: LabelSchema,
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Get a label by ID",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateLabel: {
        method: "PATCH",
        path: "/organizations/:orgId/projects/:projectId/labels/:labelId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            labelId: z.string().uuid(),
        }),
        body: UpdateLabelSchema,
        responses: {
            200: LabelSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Update a label",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    deleteLabel: {
        method: "DELETE",
        path: "/organizations/:orgId/projects/:projectId/labels/:labelId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            labelId: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            204: z.object({}),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Delete a label",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },

    // --- Issue-Label Association ---

    getIssueLabels: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/labels",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        responses: {
            200: z.array(LabelSchema),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Get all labels on an issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    addLabelToIssue: {
        method: "POST",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/labels",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        body: AddLabelToIssueSchema,
        responses: {
            201: z.object({}),
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Add a label to an issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    removeLabelFromIssue: {
        method: "DELETE",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/labels/:labelId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
            labelId: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            204: z.object({}),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Remove a label from an issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
