import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { userContract } from "./user.js";
import { organizationContract } from "./organization.js";

const c = initContract();

const apiV1 = c.router(
    {
        User: userContract,
        Organization: organizationContract,
    },
    {
        pathPrefix: "/api/v1",
    }
);

export const apiContract = c.router({
    Health: healthContract,
    User: apiV1.User,
    Organization: apiV1.Organization,
});