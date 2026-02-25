import { z } from "./z.js";
import { UserSchema } from "./user.js";

// --- Enums ---

export const ProjectStatusEnum = z.enum(["active", "archived"]);
export const ProjectMemberRoleEnum = z.enum(["admin", "member", "viewer"]);
export const IssueStatusEnum = z.enum(["backlog", "todo", "in_progress", "in_review", "in_dev", "in_prod", "cancelled"]);
export const IssuePriorityEnum = z.enum(["urgent", "high", "medium", "low", "none"]);
export const IssueTypeEnum = z.enum(["task", "bug", "story", "epic"]);

// --- Project Schemas ---

export const ProjectSchema = z.object({
    id: z.string().uuid(),
    organization_id: z.string().uuid(),
    name: z.string(),
    slug: z.string(),
    description: z.string().nullable(),
    status: ProjectStatusEnum,
    identifier: z.string(),
    issue_counter: z.number(),
    created_by: z.string().uuid(),
    created_at: z.string(),
    updated_at: z.string(),
}).openapi("Project");

export const CreateProjectSchema = z.object({
    name: z.string().min(1).max(255),
    slug: z.string().min(1).max(100),
    description: z.string().optional(),
    identifier: z.string().min(2).max(10),
});

export const UpdateProjectSchema = z.object({
    name: z.string().min(1).max(255).optional(),
    slug: z.string().min(1).max(100).optional(),
    description: z.string().optional(),
    status: ProjectStatusEnum.optional(),
});

// --- Project Member Schemas ---

export const ProjectMemberSchema = z.object({
    id: z.string().uuid(),
    project_id: z.string().uuid(),
    user_id: z.string().uuid(),
    role: ProjectMemberRoleEnum,
    joined_at: z.string(),
}).openapi("ProjectMember");

export const ProjectMemberWithDetailsSchema = ProjectMemberSchema.extend({
    user: UserSchema.optional(),
}).openapi("ProjectMemberWithDetails");

export const AddProjectMemberSchema = z.object({
    user_id: z.string().uuid(),
    role: ProjectMemberRoleEnum,
});

export const UpdateProjectMemberRoleSchema = z.object({
    role: ProjectMemberRoleEnum,
});

// --- Issue Schemas ---

export const IssueSchema = z.object({
    id: z.string().uuid(),
    project_id: z.string().uuid(),
    issue_key: z.string(),
    title: z.string(),
    description: z.string().nullable(),
    status: IssueStatusEnum,
    priority: IssuePriorityEnum,
    type: IssueTypeEnum,
    assignee_id: z.string().uuid().nullable(),
    reporter_id: z.string().uuid(),
    sort_order: z.number(),
    parent_issue_id: z.string().uuid().nullable(),
    due_date: z.string().nullable(),
    created_at: z.string(),
    updated_at: z.string(),
}).openapi("Issue");

export const IssueWithDetailsSchema = IssueSchema.extend({
    assignee: UserSchema.nullable().optional(),
    reporter: UserSchema.optional(),
}).openapi("IssueWithDetails");

export const CreateIssueSchema = z.object({
    title: z.string().min(1).max(500),
    description: z.string().optional(),
    priority: IssuePriorityEnum.optional(),
    type: IssueTypeEnum.optional(),
    assignee_id: z.string().uuid().optional(),
    parent_issue_id: z.string().uuid().optional(),
    due_date: z.string().optional(),
});

export const UpdateIssueSchema = z.object({
    title: z.string().min(1).max(500).optional(),
    description: z.string().optional(),
    status: IssueStatusEnum.optional(),
    priority: IssuePriorityEnum.optional(),
    type: IssueTypeEnum.optional(),
    assignee_id: z.string().uuid().optional(),
    sort_order: z.number().min(0).optional(),
    parent_issue_id: z.string().uuid().optional(),
    due_date: z.string().optional(),
});

export const AssignIssueSchema = z.object({
    assignee_id: z.string().uuid().nullable(),
});

export const UpdateIssueStatusSchema = z.object({
    status: IssueStatusEnum,
});

// --- Label Schemas ---

export const LabelSchema = z.object({
    id: z.string().uuid(),
    project_id: z.string().uuid(),
    name: z.string(),
    color: z.string(),
    description: z.string().nullable(),
    created_at: z.string(),
    updated_at: z.string(),
}).openapi("Label");

export const CreateLabelSchema = z.object({
    name: z.string().min(1).max(100),
    color: z.string().length(7).regex(/^#[0-9A-Fa-f]{6}$/),
    description: z.string().optional(),
});

export const UpdateLabelSchema = z.object({
    name: z.string().min(1).max(100).optional(),
    color: z.string().length(7).regex(/^#[0-9A-Fa-f]{6}$/).optional(),
    description: z.string().optional(),
});

export const AddLabelToIssueSchema = z.object({
    label_id: z.string().uuid(),
});

// --- Comment Schemas ---

export const CommentSchema = z.object({
    id: z.string().uuid(),
    issue_id: z.string().uuid(),
    author_id: z.string().uuid(),
    body: z.string(),
    parent_comment_id: z.string().uuid().nullable(),
    is_edited: z.boolean(),
    edited_at: z.string().nullable(),
    is_deleted: z.boolean(),
    deleted_at: z.string().nullable(),
    created_at: z.string(),
    updated_at: z.string(),
}).openapi("Comment");

export const CommentWithAuthorSchema = CommentSchema.extend({
    author_first_name: z.string().nullable(),
    author_last_name: z.string().nullable(),
    author_avatar_url: z.string().nullable(),
    author_email: z.string(),
}).openapi("CommentWithAuthor");

export const CommentThreadSchema = z.object({
    comment: CommentWithAuthorSchema,
    replies: z.array(CommentWithAuthorSchema),
}).openapi("CommentThread");

export const CreateCommentSchema = z.object({
    body: z.string().min(1),
    parent_comment_id: z.string().uuid().optional(),
});

export const UpdateCommentSchema = z.object({
    body: z.string().min(1),
});

// --- Activity Schemas ---

export const ActivitySchema = z.object({
    id: z.string().uuid(),
    project_id: z.string().uuid(),
    issue_id: z.string().uuid().nullable(),
    actor_id: z.string().uuid(),
    action: z.string(),
    entity_type: z.string(),
    entity_id: z.string().uuid(),
    old_value: z.any().nullable(),
    new_value: z.any().nullable(),
    metadata: z.any().nullable(),
    created_at: z.string().datetime(),
}).openapi("Activity");

export const ActivityWithActorSchema = ActivitySchema.extend({
    actor_first_name: z.string().nullable(),
    actor_last_name: z.string().nullable(),
    actor_avatar_url: z.string().nullable(),
    actor_email: z.string(),
}).openapi("ActivityWithActor");

export const ActivitySummaryItemSchema = z.object({
    key: z.string(),
    count: z.number(),
}).openapi("ActivitySummaryItem");

export const ActivitySummaryResponseSchema = z.object({
    data: z.array(ActivitySummaryItemSchema),
    total_count: z.number(),
}).openapi("ActivitySummaryResponse");

// --- Type Exports ---

export type Project = z.infer<typeof ProjectSchema>;
export type CreateProject = z.infer<typeof CreateProjectSchema>;
export type UpdateProject = z.infer<typeof UpdateProjectSchema>;
export type ProjectMember = z.infer<typeof ProjectMemberSchema>;
export type ProjectMemberWithDetails = z.infer<typeof ProjectMemberWithDetailsSchema>;
export type AddProjectMember = z.infer<typeof AddProjectMemberSchema>;
export type UpdateProjectMemberRole = z.infer<typeof UpdateProjectMemberRoleSchema>;
export type Issue = z.infer<typeof IssueSchema>;
export type IssueWithDetails = z.infer<typeof IssueWithDetailsSchema>;
export type CreateIssue = z.infer<typeof CreateIssueSchema>;
export type UpdateIssue = z.infer<typeof UpdateIssueSchema>;
export type AssignIssue = z.infer<typeof AssignIssueSchema>;
export type UpdateIssueStatus = z.infer<typeof UpdateIssueStatusSchema>;
export type Label = z.infer<typeof LabelSchema>;
export type CreateLabel = z.infer<typeof CreateLabelSchema>;
export type UpdateLabel = z.infer<typeof UpdateLabelSchema>;
export type AddLabelToIssue = z.infer<typeof AddLabelToIssueSchema>;
export type Comment = z.infer<typeof CommentSchema>;
export type CommentWithAuthor = z.infer<typeof CommentWithAuthorSchema>;
export type CommentThread = z.infer<typeof CommentThreadSchema>;
export type CreateComment = z.infer<typeof CreateCommentSchema>;
export type UpdateComment = z.infer<typeof UpdateCommentSchema>;
export type Activity = z.infer<typeof ActivitySchema>;
export type ActivityWithActor = z.infer<typeof ActivityWithActorSchema>;
export type ActivitySummaryItem = z.infer<typeof ActivitySummaryItemSchema>;
export type ActivitySummaryResponse = z.infer<typeof ActivitySummaryResponseSchema>;
