package repository_test

import (
	"context"
	"testing"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	testutil "github.com/RahulSingh9131/vector/internal/testing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueRepository_Create(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewIssueRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "isscreate")
	reporter, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "rep1", Email: "rep1@t.com"})
	assignee, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "ass1", Email: "ass1@t.com"})
	proj := createTestProject(t, projRepo, org.ID, reporter.ID, "isscreate")

	params := models.CreateIssueParams{
		ProjectID:   proj.ID,
		Title:       "Test Issue",
		Description: testutil.Ptr("Detailed report"),
		Priority:    "high",
		Type:        "bug",
		AssigneeID:  &assignee.ID,
		ReporterID:  reporter.ID,
	}

	issue, err := repo.Create(ctx, params, "KEY-1")

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, issue.ID)
	assert.Equal(t, proj.ID, issue.ProjectID)
	assert.Equal(t, "KEY-1", issue.IssueKey)
	assert.Equal(t, params.Title, issue.Title)
	assert.Equal(t, "backlog", issue.Status)
	assert.Equal(t, params.Priority, issue.Priority)
	assert.Equal(t, params.Type, issue.Type)
	assert.Equal(t, &assignee.ID, issue.AssigneeID)
	assert.Equal(t, reporter.ID, issue.ReporterID)
}

func TestIssueRepository_GetByID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewIssueRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "issget")
	reporter, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "rep2", Email: "rep2@t.com", FirstName: testutil.Ptr("John")})
	proj := createTestProject(t, projRepo, org.ID, reporter.ID, "issget")

	created, _ := repo.Create(ctx, models.CreateIssueParams{
		ProjectID:  proj.ID,
		Title:      "Issue 1",
		ReporterID: reporter.ID,
		Priority:   "low",
		Type:       "task",
	}, "K-2")

	t.Run("existing issue with details", func(t *testing.T) {
		issue, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, issue)
		assert.Equal(t, created.ID, issue.ID)
		assert.Equal(t, reporter.ID, issue.Reporter.ID)
		assert.Equal(t, "John", *issue.Reporter.FirstName)
		assert.Nil(t, issue.Assignee)
	})

	t.Run("non-existent issue returns nil", func(t *testing.T) {
		issue, err := repo.GetByID(ctx, uuid.New())
		require.NoError(t, err)
		assert.Nil(t, issue)
	})
}

func TestIssueRepository_ListByProject(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewIssueRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "isslist")
	u, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, u.ID, "isslist")

	// Create diverse issues
	repo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "Bug 1", ReporterID: u.ID, Priority: "high", Type: "bug"}, "K-1")
	repo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "Bug 2", ReporterID: u.ID, Priority: "low", Type: "bug"}, "K-2")
	repo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "Task 1", ReporterID: u.ID, Priority: "medium", Type: "task"}, "K-3")

	t.Run("list all", func(t *testing.T) {
		resp, err := repo.ListByProject(ctx, proj.ID, models.IssueFilters{})
		require.NoError(t, err)
		assert.Len(t, resp.Data, 3)
		assert.Equal(t, 3, resp.Total)
	})

	t.Run("filter by type", func(t *testing.T) {
		bugType := "bug"
		resp, err := repo.ListByProject(ctx, proj.ID, models.IssueFilters{Type: &bugType})
		require.NoError(t, err)
		assert.Len(t, resp.Data, 2)
	})

	t.Run("pagination", func(t *testing.T) {
		resp, err := repo.ListByProject(ctx, proj.ID, models.IssueFilters{Limit: 2, Page: 1})
		require.NoError(t, err)
		assert.Len(t, resp.Data, 2)
		assert.Equal(t, 3, resp.Total)
		assert.Equal(t, 2, resp.TotalPages)
	})
}

func TestIssueRepository_Update(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewIssueRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "issupd")
	u, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, u.ID, "issupd")

	created, _ := repo.Create(ctx, models.CreateIssueParams{
		ProjectID:  proj.ID,
		Title:      "Original Title",
		ReporterID: u.ID,
		Priority:   "low",
		Type:       "task",
	}, "K-1")

	newTitle := "Updated Title"
	newStatus := "in_progress"
	updated, err := repo.Update(ctx, created.ID, models.UpdateIssueParams{
		Title:  &newTitle,
		Status: &newStatus,
	})

	require.NoError(t, err)
	assert.Equal(t, newTitle, updated.Title)
	assert.Equal(t, newStatus, updated.Status)
}

func TestIssueRepository_Delete(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewIssueRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "issdel")
	u, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, u.ID, "issdel")

	created, _ := repo.Create(ctx, models.CreateIssueParams{
		ProjectID:  proj.ID,
		Title:      "To Delete",
		ReporterID: u.ID,
		Priority:   "low",
		Type:       "task",
	}, "K-1")

	err := repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	issue, _ := repo.GetByID(ctx, created.ID)
	assert.Nil(t, issue)
}

func TestIssueRepository_SubIssues(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewIssueRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "isssub")
	u, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, u.ID, "isssub")

	parent, _ := repo.Create(ctx, models.CreateIssueParams{
		ProjectID:  proj.ID,
		Title:      "Parent",
		ReporterID: u.ID,
		Priority:   "high",
		Type:       "story",
	}, "P-1")

	repo.Create(ctx, models.CreateIssueParams{
		ProjectID:     proj.ID,
		Title:         "Sub 1",
		ReporterID:    u.ID,
		Priority:      "medium",
		Type:          "task",
		ParentIssueID: &parent.ID,
	}, "S-1")

	repo.Create(ctx, models.CreateIssueParams{
		ProjectID:     proj.ID,
		Title:         "Sub 2",
		ReporterID:    u.ID,
		Priority:      "low",
		Type:          "task",
		ParentIssueID: &parent.ID,
	}, "S-2")

	subIssues, err := repo.GetSubIssues(ctx, parent.ID)
	require.NoError(t, err)
	assert.Len(t, subIssues, 2)
}
