import { initContract } from "@ts-rest/core";
import { UserSchema, UpdateUserSchema, z } from "@vector/zod";

const c = initContract();

export const userContract = c.router({
    getCurrentUser: {
        method: "GET",
        path: "/me",
        responses: {
            200: UserSchema,
            401: z.object({ message: z.string() }),
        },
        summary: "Get current user profile",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    updateCurrentUser: {
        method: "PUT",
        path: "/me",
        body: UpdateUserSchema,
        responses: {
            200: UserSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
        },
        summary: "Update current user profile",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    listUsers: {
        method: "GET",
        path: "/users",
        responses: {
            200: z.array(UserSchema),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "List all users",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    createUser: {
        method: "POST",
        path: "/users",
        body: z.object({
            clerk_user_id: z.string(),
            email: z.string().email(),
            first_name: z.string().optional(),
            last_name: z.string().optional(),
            avatar_url: z.string().url().optional(),
        }),
        responses: {
            201: UserSchema,
            400: z.object({ message: z.string() }),
            401: z.object({ message: z.string() }),
            403: z.object({ message: z.string() }),
        },
        summary: "Create a new user manually",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    getUser: {
        method: "GET",
        path: "/users/:id",
        pathParams: z.object({
            id: z.string().uuid(),
        }),
        responses: {
            200: UserSchema,
            401: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Get user by ID",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    deleteUser: {
        method: "DELETE",
        path: "/users/:id",
        pathParams: z.object({
            id: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            204: z.object({}),
            401: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Deactivate user",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
