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

func TestLabelRepository_Create(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewLabelRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "lblcreate")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u_lbl", Email: "lbl@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "lblcreate")

	params := models.CreateLabelParams{
		ProjectID:   proj.ID,
		Name:        "Bug",
		Color:       "#FF0000",
		Description: testutil.Ptr("Critical bugs"),
	}

	label, err := repo.Create(ctx, params)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, label.ID)
	assert.Equal(t, params.ProjectID, label.ProjectID)
	assert.Equal(t, params.Name, label.Name)
	assert.Equal(t, params.Color, label.Color)
	assert.Equal(t, params.Description, label.Description)
	assert.False(t, label.CreatedAt.IsZero())
}

func TestLabelRepository_GetByID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewLabelRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "lblget")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u2", Email: "u2@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "lblget")

	created, _ := repo.Create(ctx, models.CreateLabelParams{
		ProjectID: proj.ID,
		Name:      "Feature",
		Color:     "#00FF00",
	})

	t.Run("existing label", func(t *testing.T) {
		label, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, label.ID)
	})

	t.Run("non-existent label", func(t *testing.T) {
		label, err := repo.GetByID(ctx, uuid.New())
		require.NoError(t, err)
		assert.Nil(t, label)
	})
}

func TestLabelRepository_ListByProject(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewLabelRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "lbllist")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u3", Email: "u3@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "lbllist")

	names := []string{"Backend", "Frontend", "DevOps"}
	for _, name := range names {
		repo.Create(ctx, models.CreateLabelParams{
			ProjectID: proj.ID,
			Name:      name,
			Color:     "#FFFFFF",
		})
	}

	labels, err := repo.ListByProject(ctx, proj.ID)
	require.NoError(t, err)
	assert.Len(t, labels, 3)
}

func TestLabelRepository_Update(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewLabelRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "lblupd")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u4", Email: "u4@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "lblupd")

	created, _ := repo.Create(ctx, models.CreateLabelParams{
		ProjectID: proj.ID,
		Name:      "Old Name",
		Color:     "#000000",
	})

	newName := "New Name"
	newColor := "#FFFFFF"
	updated, err := repo.Update(ctx, created.ID, models.UpdateLabelParams{
		Name:  &newName,
		Color: &newColor,
	})

	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, newColor, updated.Color)
}

func TestLabelRepository_Delete(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewLabelRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "lbldel")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u5", Email: "u5@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "lbldel")

	created, _ := repo.Create(ctx, models.CreateLabelParams{
		ProjectID: proj.ID,
		Name:      "To Delete",
		Color:     "#000000",
	})

	err := repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	label, _ := repo.GetByID(ctx, created.ID)
	assert.Nil(t, label)
}

func TestLabelRepository_IssueAssociations(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewLabelRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "lbliss")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u6", Email: "u6@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "lbliss")

	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{
		ProjectID:  proj.ID,
		Title:      "Test Issue",
		ReporterID: user.ID,
		Priority:   "high",
		Type:       "bug",
	}, "ISS-1")

	l1, _ := repo.Create(ctx, models.CreateLabelParams{ProjectID: proj.ID, Name: "L1", Color: "#111111"})
	l2, _ := repo.Create(ctx, models.CreateLabelParams{ProjectID: proj.ID, Name: "L2", Color: "#222222"})

	t.Run("add labels", func(t *testing.T) {
		err := repo.AddLabelToIssue(ctx, issue.ID, l1.ID)
		require.NoError(t, err)
		err = repo.AddLabelToIssue(ctx, issue.ID, l2.ID)
		require.NoError(t, err)

		labels, err := repo.GetLabelsByIssue(ctx, issue.ID)
		require.NoError(t, err)
		assert.Len(t, labels, 2)
	})

	t.Run("remove label", func(t *testing.T) {
		err := repo.RemoveLabelFromIssue(ctx, issue.ID, l1.ID)
		require.NoError(t, err)

		labels, err := repo.GetLabelsByIssue(ctx, issue.ID)
		require.NoError(t, err)
		assert.Len(t, labels, 1)
		assert.Equal(t, l2.ID, labels[0].ID)
	})
}
