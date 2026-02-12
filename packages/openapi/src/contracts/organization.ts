import { initContract } from "@ts-rest/core";
import {
    OrganizationSchema,
    UpdateOrganizationSchema,
    OrganizationMemberWithDetailsSchema,
    z
} from "@vector/zod";

const c = initContract();

export const organizationContract = c.router({
    listOrganizations: {
        method: "GET",
        path: "/organizations",
        responses: {
            200: z.array(OrganizationMemberWithDetailsSchema),
            401: z.object({ message: z.string() }),
        },
        summary: "List current user's organizations",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    getOrganization: {
        method: "GET",
        path: "/organizations/:id",
        pathParams: z.object({
            id: z.string().uuid(),
        }),
        responses: {
            200: OrganizationSchema,
            401: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Get organization by ID",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateOrganization: {
        method: "PUT",
        path: "/organizations/:id",
        pathParams: z.object({
            id: z.string().uuid(),
        }),
        body: UpdateOrganizationSchema,
        responses: {
            200: OrganizationSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Update organization details",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    listMembers: {
        method: "GET",
        path: "/organizations/:id/members",
        pathParams: z.object({
            id: z.string().uuid(),
        }),
        responses: {
            200: z.array(OrganizationMemberWithDetailsSchema),
            401: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "List organization members",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    createOrganization: {
        method: "POST",
        path: "/organizations",
        body: z.object({
            name: z.string().min(1),
            slug: z.string().min(1),
            logo_url: z.string().url().optional(),
        }),
        responses: {
            201: OrganizationSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
        },
        summary: "Create a new organization",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    addMember: {
        method: "POST",
        path: "/organizations/:id/members",
        pathParams: z.object({
            id: z.string().uuid(),
        }),
        body: z.object({
            user_id: z.string().uuid(),
            role: z.string(),
        }),
        responses: {
            201: OrganizationMemberWithDetailsSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Add a member to an organization",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateMemberRole: {
        method: "PATCH",
        path: "/organizations/:id/members/:userId",
        pathParams: z.object({
            id: z.string().uuid(),
            userId: z.string().uuid(),
        }),
        body: z.object({
            role: z.string(),
        }),
        responses: {
            200: OrganizationMemberWithDetailsSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Update a member's role",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    removeMember: {
        method: "DELETE",
        path: "/organizations/:id/members/:userId",
        pathParams: z.object({
            id: z.string().uuid(),
            userId: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            204: z.object({}),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Remove a member from an organization",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
