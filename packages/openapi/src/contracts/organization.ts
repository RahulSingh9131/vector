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
});
