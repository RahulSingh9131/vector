import { z } from "./z.js";

export const UserSchema = z.object({
    id: z.string().uuid(),
    clerk_user_id: z.string(),
    email: z.string().email(),
    first_name: z.string().nullable(),
    last_name: z.string().nullable(),
    avatar_url: z.string().url().nullable(),
    is_active: z.boolean(),
    last_login_at: z.string().nullable(),
    created_at: z.string(),
    updated_at: z.string(),
}).openapi("User");

export const UpdateUserSchema = z.object({
    first_name: z.string().optional(),
    last_name: z.string().optional(),
    avatar_url: z.string().url().optional(),
    is_active: z.boolean().optional(),
});

export type User = z.infer<typeof UserSchema>;
export type UpdateUser = z.infer<typeof UpdateUserSchema>;
