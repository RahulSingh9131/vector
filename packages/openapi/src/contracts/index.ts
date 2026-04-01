import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { userContract } from "./user.js";
import { organizationContract } from "./organization.js";
import { projectContract } from "./project.js";
import { issueContract } from "./issue.js";
import { labelContract } from "./label.js";
import { commentContract } from "./comment.js";
import { activityContract } from "./activity.js";
import { notificationContract } from "./notification.js";

const c = initContract();

const apiV1 = c.router(
    {
        User: userContract,
        Organization: organizationContract,
        Project: projectContract,
        Issue: issueContract,
        Label: labelContract,
        Comment: commentContract,
        Activity: activityContract,
        Notification: notificationContract,
    },
    {
        pathPrefix: "/api/v1",
    }
);

export const apiContract = c.router({
    Health: healthContract,
    User: apiV1.User,
    Organization: apiV1.Organization,
    Project: apiV1.Project,
    Issue: apiV1.Issue,
    Label: apiV1.Label,
    Comment: apiV1.Comment,
    Activity: apiV1.Activity,
    Notification: apiV1.Notification,
});