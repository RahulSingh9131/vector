import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { userContract } from "./user.js";
import { organizationContract } from "./organization.js";

const c = initContract();

export const apiContract = c.router({
    Health: healthContract,
    User: userContract,
    Organization: organizationContract,
});