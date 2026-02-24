package handler

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// --- User Request Structs ---

// EmptyRequest is a no-op request for endpoints that don't accept input
type EmptyRequest struct{}

func (r *EmptyRequest) Validate() error {
	return nil
}

// EmptyResponse is a no-op response for endpoints that don't return data
type EmptyResponse struct{}

// UpdateCurrentUserRequest represents the request body for updating the current user's profile
type UpdateCurrentUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	FirstName *string `json:"firstName" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"lastName" validate:"omitempty,min=1,max=100"`
	AvatarURL *string `json:"avatarUrl" validate:"omitempty,url"`
}

func (r *UpdateCurrentUserRequest) Validate() error {
	return validate.Struct(r)
}

// CreateUserRequest represents the request body for manually creating a user
type CreateUserRequest struct {
	ClerkUserID string  `json:"clerk_user_id" validate:"required"`
	Email       string  `json:"email" validate:"required,email"`
	FirstName   *string `json:"firstName" validate:"omitempty,min=1,max=100"`
	LastName    *string `json:"lastName" validate:"omitempty,min=1,max=100"`
	AvatarURL   *string `json:"avatarUrl" validate:"omitempty,url"`
}

func (r *CreateUserRequest) Validate() error {
	return validate.Struct(r)
}

// GetUserRequest represents path params for getting a user by ID
type GetUserRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetUserRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteUserRequest represents path params for deleting a user
type DeleteUserRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteUserRequest) Validate() error {
	return validate.Struct(r)
}

// --- Organization Request Structs ---

// GetOrganizationRequest represents path params for getting an organization by ID
type GetOrganizationRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetOrganizationRequest) Validate() error {
	return validate.Struct(r)
}

// ListMembersRequest represents path params for listing members of an organization
type ListMembersRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *ListMembersRequest) Validate() error {
	return validate.Struct(r)
}

// CreateOrganizationRequest represents the request body for creating an organization
type CreateOrganizationRequest struct {
	Name    string  `json:"name" validate:"required,min=1,max=100"`
	Slug    string  `json:"slug" validate:"required,min=1,max=100"`
	LogoURL *string `json:"logoUrl" validate:"omitempty,url"`
}

func (r *CreateOrganizationRequest) Validate() error {
	return validate.Struct(r)
}

// AddMemberRequest represents the request body for adding a member to an organization
type AddMemberRequest struct {
	ID     string `param:"id" validate:"required,uuid"`
	UserID string `json:"user_id" validate:"required,uuid"`
	Role   string `json:"role" validate:"required,oneof=admin member guest"`
}

func (r *AddMemberRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateMemberRoleRequest represents the request body for updating a member's role
type UpdateMemberRoleRequest struct {
	ID     string `param:"id" validate:"required,uuid"`
	UserID string `param:"userId" validate:"required,uuid"`
	Role   string `json:"role" validate:"required,oneof=admin member guest"`
}

func (r *UpdateMemberRoleRequest) Validate() error {
	return validate.Struct(r)
}

// RemoveMemberRequest represents path params for removing a member
type RemoveMemberRequest struct {
	ID     string `param:"id" validate:"required,uuid"`
	UserID string `param:"userId" validate:"required,uuid"`
}

func (r *RemoveMemberRequest) Validate() error {
	return validate.Struct(r)
}

// --- Project Request Structs ---

// CreateProjectRequest represents the request body for creating a project
type CreateProjectRequest struct {
	OrgID       string  `param:"orgId" validate:"required,uuid"`
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Slug        string  `json:"slug" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty"`
	Identifier  string  `json:"identifier" validate:"required,min=2,max=10,uppercase"`
}

func (r *CreateProjectRequest) Validate() error {
	return validate.Struct(r)
}

// ListProjectsRequest represents path params for listing projects
type ListProjectsRequest struct {
	OrgID string `param:"orgId" validate:"required,uuid"`
}

func (r *ListProjectsRequest) Validate() error {
	return validate.Struct(r)
}

// GetProjectRequest represents path params for getting a project
type GetProjectRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
}

func (r *GetProjectRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateProjectRequest represents the request body for updating a project
type UpdateProjectRequest struct {
	OrgID       string  `param:"orgId" validate:"required,uuid"`
	ProjectID   string  `param:"projectId" validate:"required,uuid"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Slug        *string `json:"slug" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty"`
	Status      *string `json:"status" validate:"omitempty,oneof=active archived"`
}

func (r *UpdateProjectRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteProjectRequest represents path params for deleting a project
type DeleteProjectRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
}

func (r *DeleteProjectRequest) Validate() error {
	return validate.Struct(r)
}

// --- Project Member Request Structs ---

// ListProjectMembersRequest represents path params for listing project members
type ListProjectMembersRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
}

func (r *ListProjectMembersRequest) Validate() error {
	return validate.Struct(r)
}

// AddProjectMemberRequest represents the request body for adding a project member
type AddProjectMemberRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	UserID    string `json:"user_id" validate:"required,uuid"`
	Role      string `json:"role" validate:"required,oneof=admin member viewer"`
}

func (r *AddProjectMemberRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateProjectMemberRoleRequest represents the request body for updating a project member's role
type UpdateProjectMemberRoleRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	UserID    string `param:"userId" validate:"required,uuid"`
	Role      string `json:"role" validate:"required,oneof=admin member viewer"`
}

func (r *UpdateProjectMemberRoleRequest) Validate() error {
	return validate.Struct(r)
}

// RemoveProjectMemberRequest represents path params for removing a project member
type RemoveProjectMemberRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	UserID    string `param:"userId" validate:"required,uuid"`
}

func (r *RemoveProjectMemberRequest) Validate() error {
	return validate.Struct(r)
}

// --- Issue Request Structs ---

// CreateIssueRequest represents the request body for creating an issue
type CreateIssueRequest struct {
	OrgID         string     `param:"orgId" validate:"required,uuid"`
	ProjectID     string     `param:"projectId" validate:"required,uuid"`
	Title         string     `json:"title" validate:"required,min=1,max=500"`
	Description   *string    `json:"description" validate:"omitempty"`
	Priority      string     `json:"priority" validate:"omitempty,oneof=urgent high medium low none"`
	Type          string     `json:"type" validate:"omitempty,oneof=task bug story epic"`
	AssigneeID    *string    `json:"assignee_id" validate:"omitempty,uuid"`
	ParentIssueID *string    `json:"parent_issue_id" validate:"omitempty,uuid"`
	DueDate       *string    `json:"due_date" validate:"omitempty"`
}

func (r *CreateIssueRequest) Validate() error {
	return validate.Struct(r)
}

// ListIssuesRequest represents query params for listing issues
type ListIssuesRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	ProjectID  string  `param:"projectId" validate:"required,uuid"`
	Status     *string `query:"status" validate:"omitempty,oneof=backlog todo in_progress in_review in_dev in_prod cancelled"`
	Priority   *string `query:"priority" validate:"omitempty,oneof=urgent high medium low none"`
	Type       *string `query:"type" validate:"omitempty,oneof=task bug story epic"`
	AssigneeID *string `query:"assignee_id" validate:"omitempty,uuid"`
	Page       int     `query:"page" validate:"omitempty,min=1"`
	Limit      int     `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (r *ListIssuesRequest) Validate() error {
	return validate.Struct(r)
}

// GetIssueRequest represents path params for getting an issue
type GetIssueRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
}

func (r *GetIssueRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateIssueRequest represents the request body for updating an issue
type UpdateIssueRequest struct {
	OrgID         string  `param:"orgId" validate:"required,uuid"`
	ProjectID     string  `param:"projectId" validate:"required,uuid"`
	IssueID       string  `param:"issueId" validate:"required,uuid"`
	Title         *string `json:"title" validate:"omitempty,min=1,max=500"`
	Description   *string `json:"description" validate:"omitempty"`
	Status        *string `json:"status" validate:"omitempty,oneof=backlog todo in_progress in_review in_dev in_prod cancelled"`
	Priority      *string `json:"priority" validate:"omitempty,oneof=urgent high medium low none"`
	Type          *string `json:"type" validate:"omitempty,oneof=task bug story epic"`
	AssigneeID    *string `json:"assignee_id" validate:"omitempty,uuid"`
	SortOrder     *int    `json:"sort_order" validate:"omitempty,min=0"`
	ParentIssueID *string `json:"parent_issue_id" validate:"omitempty,uuid"`
	DueDate       *string `json:"due_date" validate:"omitempty"`
}

func (r *UpdateIssueRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteIssueRequest represents path params for deleting an issue
type DeleteIssueRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
}

func (r *DeleteIssueRequest) Validate() error {
	return validate.Struct(r)
}

// AssignIssueRequest represents the request body for assigning an issue
type AssignIssueRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	ProjectID  string  `param:"projectId" validate:"required,uuid"`
	IssueID    string  `param:"issueId" validate:"required,uuid"`
	AssigneeID *string `json:"assignee_id" validate:"omitempty,uuid"`
}

func (r *AssignIssueRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateIssueStatusRequest represents the request body for updating an issue's status
type UpdateIssueStatusRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	Status    string `json:"status" validate:"required,oneof=backlog todo in_progress in_review in_dev in_prod cancelled"`
}

func (r *UpdateIssueStatusRequest) Validate() error {
	return validate.Struct(r)
}

// --- Label Request Structs ---

// CreateLabelRequest represents the request body for creating a label
type CreateLabelRequest struct {
	OrgID       string  `param:"orgId" validate:"required,uuid"`
	ProjectID   string  `param:"projectId" validate:"required,uuid"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Color       string  `json:"color" validate:"required,min=7,max=7"`
	Description *string `json:"description" validate:"omitempty"`
}

func (r *CreateLabelRequest) Validate() error {
	return validate.Struct(r)
}

// ListLabelsRequest represents path params for listing labels
type ListLabelsRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
}

func (r *ListLabelsRequest) Validate() error {
	return validate.Struct(r)
}

// GetLabelRequest represents path params for getting a label
type GetLabelRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	LabelID   string `param:"labelId" validate:"required,uuid"`
}

func (r *GetLabelRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateLabelRequest represents the request body for updating a label
type UpdateLabelRequest struct {
	OrgID       string  `param:"orgId" validate:"required,uuid"`
	ProjectID   string  `param:"projectId" validate:"required,uuid"`
	LabelID     string  `param:"labelId" validate:"required,uuid"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Color       *string `json:"color" validate:"omitempty,min=7,max=7"`
	Description *string `json:"description" validate:"omitempty"`
}

func (r *UpdateLabelRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteLabelRequest represents path params for deleting a label
type DeleteLabelRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	LabelID   string `param:"labelId" validate:"required,uuid"`
}

func (r *DeleteLabelRequest) Validate() error {
	return validate.Struct(r)
}

// GetIssueLabelsRequest represents path params for getting labels on an issue
type GetIssueLabelsRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
}

func (r *GetIssueLabelsRequest) Validate() error {
	return validate.Struct(r)
}

// AddLabelToIssueRequest represents the request body for adding a label to an issue
type AddLabelToIssueRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	LabelID   string `json:"label_id" validate:"required,uuid"`
}

func (r *AddLabelToIssueRequest) Validate() error {
	return validate.Struct(r)
}

// RemoveLabelFromIssueRequest represents path params for removing a label from an issue
type RemoveLabelFromIssueRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	LabelID   string `param:"labelId" validate:"required,uuid"`
}

func (r *RemoveLabelFromIssueRequest) Validate() error {
	return validate.Struct(r)
}

// --- Comment Request Structs ---

// CreateCommentRequest represents the request body for creating a comment
type CreateCommentRequest struct {
	OrgID           string  `param:"orgId" validate:"required,uuid"`
	ProjectID       string  `param:"projectId" validate:"required,uuid"`
	IssueID         string  `param:"issueId" validate:"required,uuid"`
	Body            string  `json:"body" validate:"required,min=1"`
	ParentCommentID *string `json:"parent_comment_id" validate:"omitempty,uuid"`
}

func (r *CreateCommentRequest) Validate() error {
	return validate.Struct(r)
}

// ListCommentsRequest represents query params for listing comments
type ListCommentsRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (r *ListCommentsRequest) Validate() error {
	return validate.Struct(r)
}

// GetCommentRequest represents path params for getting a comment
type GetCommentRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	CommentID string `param:"commentId" validate:"required,uuid"`
}

func (r *GetCommentRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateCommentRequest represents the request body for updating a comment
type UpdateCommentRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	CommentID string `param:"commentId" validate:"required,uuid"`
	Body      string `json:"body" validate:"required,min=1"`
}

func (r *UpdateCommentRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteCommentRequest represents path params for deleting a comment
type DeleteCommentRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	CommentID string `param:"commentId" validate:"required,uuid"`
}

func (r *DeleteCommentRequest) Validate() error {
	return validate.Struct(r)
}

// --- Activity Request Structs ---

// ListIssueActivityRequest represents query params for listing issue activity
type ListIssueActivityRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	ProjectID  string  `param:"projectId" validate:"required,uuid"`
	IssueID    string  `param:"issueId" validate:"required,uuid"`
	Page       int     `query:"page" validate:"omitempty,min=1"`
	Limit      int     `query:"limit" validate:"omitempty,min=1,max=100"`
	Action     *string `query:"action" validate:"omitempty"`
	EntityType *string `query:"entity_type" validate:"omitempty"`
	ActorID    *string `query:"actor_id" validate:"omitempty,uuid"`
	From       *string `query:"from" validate:"omitempty"`
	To         *string `query:"to" validate:"omitempty"`
}

func (r *ListIssueActivityRequest) Validate() error {
	return validate.Struct(r)
}

// ListProjectActivityRequest represents query params for listing project activity
type ListProjectActivityRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	ProjectID  string  `param:"projectId" validate:"required,uuid"`
	Page       int     `query:"page" validate:"omitempty,min=1"`
	Limit      int     `query:"limit" validate:"omitempty,min=1,max=100"`
	Action     *string `query:"action" validate:"omitempty"`
	EntityType *string `query:"entity_type" validate:"omitempty"`
	ActorID    *string `query:"actor_id" validate:"omitempty,uuid"`
	From       *string `query:"from" validate:"omitempty"`
	To         *string `query:"to" validate:"omitempty"`
}

func (r *ListProjectActivityRequest) Validate() error {
	return validate.Struct(r)
}
