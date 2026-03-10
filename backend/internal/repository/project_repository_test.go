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

func createTestProject(t *testing.T, repo *repository.ProjectRepository, orgID uuid.UUID, userID uuid.UUID, suffix string) *models.Project {
	t.Helper()

	project, err := repo.Create(context.Background(), models.CreateProjectParams{
		OrganizationID: orgID,
		Name:           "Test Project " + suffix,
		Slug:           "test-project-" + suffix,
		Description:    testutil.Ptr("Description " + suffix),
		Identifier:     "TP" + suffix[:2],
		CreatedBy:      userID,
	})
	require.NoError(t, err)
	return project
}

func TestProjectRepository_Create(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "projcreate")
	user, err := userRepo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_proj_create",
		Email:       "projcreate@test.com",
	})
	require.NoError(t, err)

	params := models.CreateProjectParams{
		OrganizationID: org.ID,
		Name:           "New Project",
		Slug:           "new-project",
		Description:    testutil.Ptr("Project description"),
		Identifier:     "NP",
		CreatedBy:      user.ID,
	}

	project, err := repo.Create(ctx, params)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, project.ID)
	assert.Equal(t, params.OrganizationID, project.OrganizationID)
	assert.Equal(t, params.Name, project.Name)
	assert.Equal(t, params.Slug, project.Slug)
	assert.Equal(t, params.Description, project.Description)
	assert.Equal(t, "active", project.Status)
	assert.Equal(t, "NP", project.Identifier)
	assert.Equal(t, 0, project.IssueCounter)
	assert.Equal(t, user.ID, project.CreatedBy)
	assert.False(t, project.CreatedAt.IsZero())
	assert.False(t, project.UpdatedAt.IsZero())
}

func TestProjectRepository_GetByID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "projget")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u1", Email: "u1@t.com"})
	created := createTestProject(t, repo, org.ID, user.ID, "getbyid")

	t.Run("existing project", func(t *testing.T) {
		project, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, project)
		assert.Equal(t, created.ID, project.ID)
	})

	t.Run("non-existent project returns nil", func(t *testing.T) {
		project, err := repo.GetByID(ctx, uuid.New())
		require.NoError(t, err)
		assert.Nil(t, project)
	})
}

func TestProjectRepository_GetBySlug(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "projgetslug")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u2", Email: "u2@t.com"})
	created := createTestProject(t, repo, org.ID, user.ID, "slug")

	t.Run("existing slug", func(t *testing.T) {
		project, err := repo.GetBySlug(ctx, org.ID, created.Slug)
		require.NoError(t, err)
		require.NotNil(t, project)
		assert.Equal(t, created.ID, project.ID)
	})

	t.Run("wrong org returns nil", func(t *testing.T) {
		project, err := repo.GetBySlug(ctx, uuid.New(), created.Slug)
		require.NoError(t, err)
		assert.Nil(t, project)
	})
}

func TestProjectRepository_ListByOrganization(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "projlist")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u3", Email: "u3@t.com"})

	// Create 3 projects
	for i := 1; i <= 3; i++ {
		createTestProject(t, repo, org.ID, user.ID, uuid.New().String()[:8])
	}

	projects, err := repo.ListByOrganization(ctx, org.ID)
	require.NoError(t, err)
	assert.Len(t, projects, 3)
}

func TestProjectRepository_Update(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "projupdate")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u4", Email: "u4@t.com"})
	created := createTestProject(t, repo, org.ID, user.ID, "update")

	newName := "Updated Project Name"
	newStatus := "archived"
	project, err := repo.Update(ctx, created.ID, models.UpdateProjectParams{
		Name:   &newName,
		Status: &newStatus,
	})

	require.NoError(t, err)
	assert.Equal(t, newName, project.Name)
	assert.Equal(t, newStatus, project.Status)
	assert.Equal(t, created.Slug, project.Slug, "slug should be unchanged")
}

func TestProjectRepository_Delete(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "projdelete")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u5", Email: "u5@t.com"})
	created := createTestProject(t, repo, org.ID, user.ID, "delete")

	err := repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "deleted", updated.Status)
}

func TestProjectRepository_IncrementIssueCounter(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "projinc")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u6", Email: "u6@t.com"})
	created := createTestProject(t, repo, org.ID, user.ID, "inc")

	assert.Equal(t, 0, created.IssueCounter)

	c1, err := repo.IncrementIssueCounter(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, c1)

	c2, err := repo.IncrementIssueCounter(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, c2)

	updated, _ := repo.GetByID(ctx, created.ID)
	assert.Equal(t, 2, updated.IssueCounter)
}

func TestProjectRepository_GetProjectCount(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "projcount")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u7", Email: "u7@t.com"})

	// Initially 0
	count, _ := repo.GetProjectCount(ctx, org.ID)
	assert.Equal(t, 0, count)

	// Create 2 projects
	createTestProject(t, repo, org.ID, user.ID, "p1")
	createTestProject(t, repo, org.ID, user.ID, "p2")

	count, _ = repo.GetProjectCount(ctx, org.ID)
	assert.Equal(t, 2, count)

	// One deleted project
	p3 := createTestProject(t, repo, org.ID, user.ID, "p3")
	repo.Delete(ctx, p3.ID)

	count, _ = repo.GetProjectCount(ctx, org.ID)
	assert.Equal(t, 2, count, "deleted projects should not be counted")
}
