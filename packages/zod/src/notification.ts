import { z } from "./z.js";

// --- Notification Schemas ---

export const NotificationSchema = z.object({
    id: z.string().uuid(),
    user_id: z.string().uuid(),
    actor_id: z.string().uuid().nullable(),
    project_id: z.string().uuid().nullable(),
    issue_id: z.string().uuid().nullable(),
    type: z.string(),
    title: z.string(),
    message: z.string(),
    payload: z.any(), // json.RawMessage in Go
    is_read: z.boolean(),
    created_at: z.string(),
}).openapi("Notification");

export const ListNotificationsRequestSchema = z.object({
    is_read: z.boolean().optional(),
    limit: z.coerce.number().min(1).max(100).optional().default(20),
    offset: z.coerce.number().min(0).optional().default(0),
});

export const MarkAsReadSchema = z.object({
    id: z.string().uuid(),
});

// --- Type Exports ---

export type Notification = z.infer<typeof NotificationSchema>;
export type ListNotificationsRequest = z.infer<typeof ListNotificationsRequestSchema>;
export type MarkAsRead = z.infer<typeof MarkAsReadSchema>;
