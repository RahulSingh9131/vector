import { initContract } from "@ts-rest/core";
import {
    CommentWithAuthorSchema,
    CommentThreadSchema,
    CreateCommentSchema,
    UpdateCommentSchema,
    schemaWithPagination,
    z,
} from "@vector/zod";

const c = initContract();

export const commentContract = c.router({
    createComment: {
        method: "POST",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/comments",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        body: CreateCommentSchema,
        responses: {
            201: CommentWithAuthorSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Create a comment on an issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    listComments: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/comments",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
        }),
        query: z.object({
            page: z.coerce.number().min(1).optional(),
            limit: z.coerce.number().min(1).max(100).optional(),
        }),
        responses: {
            200: schemaWithPagination(CommentThreadSchema),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "List threaded comments for an issue",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    getComment: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/comments/:commentId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
            commentId: z.string().uuid(),
        }),
        responses: {
            200: CommentWithAuthorSchema,
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Get a single comment",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateComment: {
        method: "PATCH",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/comments/:commentId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
            commentId: z.string().uuid(),
        }),
        body: UpdateCommentSchema,
        responses: {
            200: CommentWithAuthorSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Edit a comment (author only)",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    deleteComment: {
        method: "DELETE",
        path: "/organizations/:orgId/projects/:projectId/issues/:issueId/comments/:commentId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            issueId: z.string().uuid(),
            commentId: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            204: z.object({}),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Soft-delete a comment (author or project admin)",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
