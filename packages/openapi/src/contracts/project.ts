import { initContract } from "@ts-rest/core";
import {
    ProjectSchema,
    CreateProjectSchema,
    UpdateProjectSchema,
    ProjectMemberSchema,
    ProjectMemberWithDetailsSchema,
    AddProjectMemberSchema,
    UpdateProjectMemberRoleSchema,
    z,
} from "@vector/zod";

const c = initContract();

export const projectContract = c.router({
    createProject: {
        method: "POST",
        path: "/organizations/:orgId/projects",
        pathParams: z.object({
            orgId: z.string().uuid(),
        }),
        body: CreateProjectSchema,
        responses: {
            201: ProjectSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Create a new project",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    listProjects: {
        method: "GET",
        path: "/organizations/:orgId/projects",
        pathParams: z.object({
            orgId: z.string().uuid(),
        }),
        responses: {
            200: z.array(ProjectSchema),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "List projects the current user is a member of",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    getProject: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        responses: {
            200: ProjectSchema,
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Get project by ID",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateProject: {
        method: "PATCH",
        path: "/organizations/:orgId/projects/:projectId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        body: UpdateProjectSchema,
        responses: {
            200: ProjectSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Update a project",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    deleteProject: {
        method: "DELETE",
        path: "/organizations/:orgId/projects/:projectId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            204: z.object({}),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Delete a project",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },

    // --- Project Members ---

    listProjectMembers: {
        method: "GET",
        path: "/organizations/:orgId/projects/:projectId/members",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        responses: {
            200: z.array(ProjectMemberWithDetailsSchema),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "List project members",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    addProjectMember: {
        method: "POST",
        path: "/organizations/:orgId/projects/:projectId/members",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
        }),
        body: AddProjectMemberSchema,
        responses: {
            201: ProjectMemberSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Add a member to a project",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateProjectMemberRole: {
        method: "PATCH",
        path: "/organizations/:orgId/projects/:projectId/members/:userId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            userId: z.string().uuid(),
        }),
        body: UpdateProjectMemberRoleSchema,
        responses: {
            200: ProjectMemberSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Update a project member's role",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    removeProjectMember: {
        method: "DELETE",
        path: "/organizations/:orgId/projects/:projectId/members/:userId",
        pathParams: z.object({
            orgId: z.string().uuid(),
            projectId: z.string().uuid(),
            userId: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            204: z.object({}),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Remove a member from a project",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
