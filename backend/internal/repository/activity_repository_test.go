package repository_test

import (
	"context"
	"testing"
	"time"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	testutil "github.com/RahulSingh9131/vector/internal/testing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivityRepository_Create(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewActivityRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "actcreate")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u_act", Email: "act@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "actcreate")

	params := models.CreateActivityParams{
		ProjectID:  proj.ID,
		ActorID:    user.ID,
		Action:     "issue.created",
		EntityType: "issue",
		EntityID:   uuid.New(),
		NewValue:   map[string]string{"title": "New Issue"},
		Metadata:   map[string]string{"ip": "127.0.0.1"},
	}

	activity, err := repo.Create(ctx, params)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, activity.ID)
	assert.Equal(t, proj.ID, activity.ProjectID)
	assert.Equal(t, user.ID, activity.ActorID)
	assert.Equal(t, "issue.created", activity.Action)
	assert.NotNil(t, activity.NewValue)
	assert.NotNil(t, activity.Metadata)
}

func TestActivityRepository_ListByIssue(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewActivityRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "actiss")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u_act", Email: "act@t.com", FirstName: testutil.Ptr("Actor")})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "actiss")
	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "I", ReporterID: user.ID, Priority: "low", Type: "task"}, "K-1")

	// Create 2 activities for this issue
	repo.Create(ctx, models.CreateActivityParams{
		ProjectID:  proj.ID,
		IssueID:    &issue.ID,
		ActorID:    user.ID,
		Action:     "issue.created",
		EntityType: "issue",
		EntityID:   issue.ID,
	})
	repo.Create(ctx, models.CreateActivityParams{
		ProjectID:  proj.ID,
		IssueID:    &issue.ID,
		ActorID:    user.ID,
		Action:     "issue.updated",
		EntityType: "issue",
		EntityID:   issue.ID,
	})

	activities, total, err := repo.ListByIssue(ctx, issue.ID, 1, 10, models.ActivityFilters{})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, activities, 2)
	assert.Equal(t, "Actor", *activities[0].ActorFirstName)
}

func TestActivityRepository_ListByOrganization(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewActivityRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	memberRepo := repository.NewProjectMemberRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "actor")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "actor")

	// User MUST be a member of the project to see its activities in ListByOrganization
	memberRepo.AddMember(ctx, models.CreateProjectMemberParams{
		ProjectID: proj.ID,
		UserID:    user.ID,
		Role:      "admin",
	})

	repo.Create(ctx, models.CreateActivityParams{
		ProjectID:  proj.ID,
		ActorID:    user.ID,
		Action:     "project.updated",
		EntityType: "project",
		EntityID:   proj.ID,
	})

	activities, total, err := repo.ListByOrganization(ctx, org.ID, user.ID, 1, 10, models.ActivityFilters{})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, activities, 1)
}

func TestActivityRepository_Summary(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewActivityRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "actsum")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "actsum")

	actions := []string{"a1", "a1", "a2"}
	for _, a := range actions {
		repo.Create(ctx, models.CreateActivityParams{
			ProjectID:  proj.ID,
			ActorID:    user.ID,
			Action:     a,
			EntityType: "test",
			EntityID:   uuid.New(),
		})
	}

	summary, err := repo.SummaryByProject(ctx, proj.ID, "action", "", models.ActivityFilters{})
	require.NoError(t, err)
	assert.Equal(t, 3, summary.TotalCount)
	assert.Len(t, summary.Data, 2) // a1 and a2
}

func TestActivityRepository_DeleteOlderThan(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewActivityRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "actclean")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "actclean")

	repo.Create(ctx, models.CreateActivityParams{
		ProjectID:  proj.ID,
		ActorID:    user.ID,
		Action:     "old",
		EntityType: "test",
		EntityID:   uuid.New(),
	})

	// Manually backdate the activity for testing cleanup
	// Note: We can't easily backdate via the Repository, so we use direct DB access
	backdateQuery := "UPDATE activities SET created_at = NOW() - INTERVAL '100 days'"
	_, err := testServer.DB.Pool.Exec(ctx, backdateQuery)
	require.NoError(t, err)

	repo.Create(ctx, models.CreateActivityParams{
		ProjectID:  proj.ID,
		ActorID:    user.ID,
		Action:     "new",
		EntityType: "test",
		EntityID:   uuid.New(),
	})

	before := time.Now().AddDate(0, 0, -90)
	deleted, err := repo.DeleteOlderThan(ctx, before, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), deleted)

	_, total, _ := repo.ListByProject(ctx, proj.ID, 1, 10, models.ActivityFilters{})
	assert.Equal(t, 1, total)
}
