import { z } from "./z.js";

export const OrganizationSchema = z.object({
    id: z.string().uuid(),
    clerk_org_id: z.string(),
    name: z.string(),
    slug: z.string(),
    logo_url: z.string().url().nullable(),
    subscription_tier: z.enum(["free", "pro", "enterprise"]),
    max_members: z.number(),
    max_projects: z.number(),
    is_active: z.boolean(),
    created_at: z.string(),
    updated_at: z.string(),
}).openapi("Organization");

export const UpdateOrganizationSchema = z.object({
    name: z.string().optional(),
    slug: z.string().optional(),
    logo_url: z.string().url().optional(),
    subscription_tier: z.enum(["free", "pro", "enterprise"]).optional(),
    max_members: z.number().optional(),
    max_projects: z.number().optional(),
    is_active: z.boolean().optional(),
});

export const CreateOrganizationSchema = z.object({
    name: z.string().min(1),
    slug: z.string().min(1),
    logo_url: z.string().url().optional(),
});

export const AddMemberSchema = z.object({
    user_id: z.string().uuid(),
    role: z.string(),
});

export const UpdateMemberRoleSchema = z.object({
    role: z.string(),
});

export type Organization = z.infer<typeof OrganizationSchema>;
export type UpdateOrganization = z.infer<typeof UpdateOrganizationSchema>;
export type CreateOrganization = z.infer<typeof CreateOrganizationSchema>;

export const OrganizationMemberSchema = z.object({
    id: z.string().uuid(),
    organization_id: z.string().uuid(),
    user_id: z.string().uuid(),
    role: z.string(),
    joined_at: z.string(),
}).openapi("OrganizationMember");

export const OrganizationMemberWithDetailsSchema = OrganizationMemberSchema.extend({
    user_email: z.string().email().optional(),
    user_name: z.string().optional(),
    organization_name: z.string().optional(),
}).openapi("OrganizationMemberWithDetails");

export type OrganizationMember = z.infer<typeof OrganizationMemberSchema>;
export type OrganizationMemberWithDetails = z.infer<typeof OrganizationMemberWithDetailsSchema>;
export type AddMember = z.infer<typeof AddMemberSchema>;
export type UpdateMemberRole = z.infer<typeof UpdateMemberRoleSchema>;
