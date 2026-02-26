package handler

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

// ListMyActivityRequest represents query params for listing the authenticated user's activity
type ListMyActivityRequest struct {
	Page       int     `query:"page" validate:"omitempty,min=1"`
	Limit      int     `query:"limit" validate:"omitempty,min=1,max=100"`
	Action     *string `query:"action" validate:"omitempty"`
	EntityType *string `query:"entity_type" validate:"omitempty"`
	From       *string `query:"from" validate:"omitempty"`
	To         *string `query:"to" validate:"omitempty"`
}

func (r *ListMyActivityRequest) Validate() error {
	return validate.Struct(r)
}

// ListOrgActivityRequest represents query params for listing organization-level activity
type ListOrgActivityRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	Page       int     `query:"page" validate:"omitempty,min=1"`
	Limit      int     `query:"limit" validate:"omitempty,min=1,max=100"`
	Action     *string `query:"action" validate:"omitempty"`
	EntityType *string `query:"entity_type" validate:"omitempty"`
	ActorID    *string `query:"actor_id" validate:"omitempty,uuid"`
	From       *string `query:"from" validate:"omitempty"`
	To         *string `query:"to" validate:"omitempty"`
}

func (r *ListOrgActivityRequest) Validate() error {
	return validate.Struct(r)
}

// ProjectActivitySummaryRequest represents query params for project activity summary
type ProjectActivitySummaryRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	ProjectID  string  `param:"projectId" validate:"required,uuid"`
	GroupBy    *string `query:"group_by" validate:"omitempty,oneof=action entity_type actor_id date"`
	Interval   *string `query:"interval" validate:"omitempty,oneof=day week month"`
	Action     *string `query:"action" validate:"omitempty"`
	EntityType *string `query:"entity_type" validate:"omitempty"`
	ActorID    *string `query:"actor_id" validate:"omitempty,uuid"`
	From       *string `query:"from" validate:"omitempty"`
	To         *string `query:"to" validate:"omitempty"`
}

func (r *ProjectActivitySummaryRequest) Validate() error {
	return validate.Struct(r)
}

// OrgActivitySummaryRequest represents query params for organization activity summary
type OrgActivitySummaryRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	GroupBy    *string `query:"group_by" validate:"omitempty,oneof=action entity_type actor_id date"`
	Interval   *string `query:"interval" validate:"omitempty,oneof=day week month"`
	Action     *string `query:"action" validate:"omitempty"`
	EntityType *string `query:"entity_type" validate:"omitempty"`
	ActorID    *string `query:"actor_id" validate:"omitempty,uuid"`
	From       *string `query:"from" validate:"omitempty"`
	To         *string `query:"to" validate:"omitempty"`
}

func (r *OrgActivitySummaryRequest) Validate() error {
	return validate.Struct(r)
}
