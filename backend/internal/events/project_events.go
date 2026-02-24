package events

import "github.com/google/uuid"

// --- Project Events ---

// ProjectCreatedPayload is published when a new project is created.
type ProjectCreatedPayload struct {
	ProjectID  uuid.UUID `json:"project_id"`
	Name       string    `json:"name"`
	Identifier string    `json:"identifier"`
	OrgID      uuid.UUID `json:"org_id"`
}

// ProjectUpdatedPayload is published when project fields are modified.
type ProjectUpdatedPayload struct {
	ProjectID uuid.UUID              `json:"project_id"`
	OldValue  map[string]interface{} `json:"old_value"`
	NewValue  map[string]interface{} `json:"new_value"`
}

// ProjectDeletedPayload is published when a project is deleted.
type ProjectDeletedPayload struct {
	ProjectID  uuid.UUID `json:"project_id"`
	Name       string    `json:"name"`
	Identifier string    `json:"identifier"`
}

// --- Member Events ---

// MemberAddedPayload is published when a user is added to a project.
type MemberAddedPayload struct {
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
}

// MemberRemovedPayload is published when a user is removed from a project.
type MemberRemovedPayload struct {
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// MemberRoleChangedPayload is published when a member's role changes.
type MemberRoleChangedPayload struct {
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `json:"user_id"`
	OldRole   string    `json:"old_role"`
	NewRole   string    `json:"new_role"`
}
