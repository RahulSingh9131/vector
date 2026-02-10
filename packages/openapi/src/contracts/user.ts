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
});
