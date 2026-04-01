import { initContract } from "@ts-rest/core";
import { z } from "zod";
import {
    NotificationSchema,
    ListNotificationsRequestSchema,
} from "@vector/zod";

const c = initContract();

export const notificationContract = c.router({
    listNotifications: {
        method: "GET",
        path: "/notifications",
        query: ListNotificationsRequestSchema,
        responses: {
            200: z.array(NotificationSchema),
            401: z.object({ message: z.string() }),
        },
        summary: "List notifications for the current user",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    markAsRead: {
        method: "PATCH",
        path: "/notifications/:id/read",
        pathParams: z.object({
            id: z.string().uuid(),
        }),
        body: z.object({}),
        responses: {
            200: z.object({ message: z.string().optional() }),
            401: z.object({ message: z.string() }),
            404: z.object({ message: z.string() }),
        },
        summary: "Mark a notification as read",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    markAllAsRead: {
        method: "POST",
        path: "/notifications/read-all",
        body: z.object({}),
        responses: {
            200: z.object({ message: z.string().optional() }),
            401: z.object({ message: z.string() }),
        },
        summary: "Mark all notifications as read",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
    streamNotifications: {
        method: "GET",
        path: "/notifications/stream",
        responses: {
            200: z.any(), // SSE stream
            401: z.object({ message: z.string() }),
        },
        summary: "Stream real-time notifications via SSE",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
