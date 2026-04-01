import { initContract } from "@ts-rest/core";
import { z } from "zod";
import {
    SearchResponseSchema,
    SearchRequestSchema,
} from "@vector/zod";

const c = initContract();

export const searchContract = c.router({
    search: {
        method: "GET",
        path: "/search",
        query: SearchRequestSchema,
        responses: {
            200: SearchResponseSchema,
            401: z.object({ message: z.string() }),
        },
        summary: "Unified full-text search across issues and comments",
        metadata: {
            openApiSecurity: [{ bearerAuth: [] }],
        },
    },
});
