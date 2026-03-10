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

func TestProjectMemberRepository_AddMember(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectMemberRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "pmadd")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "pm1", Email: "pm1@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "pmadd")

	params := models.CreateProjectMemberParams{
		ProjectID: proj.ID,
		UserID:    user.ID,
		Role:      "admin",
	}

	member, err := repo.AddMember(ctx, params)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, member.ID)
	assert.Equal(t, proj.ID, member.ProjectID)
	assert.Equal(t, user.ID, member.UserID)
	assert.Equal(t, "admin", member.Role)
	assert.False(t, member.JoinedAt.IsZero())
}

func TestProjectMemberRepository_GetMember(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectMemberRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "pmget")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "pm2", Email: "pm2@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "pmget")
	repo.AddMember(ctx, models.CreateProjectMemberParams{ProjectID: proj.ID, UserID: user.ID, Role: "member"})

	t.Run("existing member", func(t *testing.T) {
		member, err := repo.GetMember(ctx, proj.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, member)
		assert.Equal(t, "member", member.Role)
	})

	t.Run("non-existent member returns nil", func(t *testing.T) {
		member, err := repo.GetMember(ctx, proj.ID, uuid.New())
		require.NoError(t, err)
		assert.Nil(t, member)
	})
}

func TestProjectMemberRepository_UpdateRole(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectMemberRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "pmupd")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "pm3", Email: "pm3@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "pmupd")
	repo.AddMember(ctx, models.CreateProjectMemberParams{ProjectID: proj.ID, UserID: user.ID, Role: "member"})

	_, err := repo.UpdateRole(ctx, proj.ID, user.ID, "admin")
	require.NoError(t, err)

	member, _ := repo.GetMember(ctx, proj.ID, user.ID)
	assert.Equal(t, "admin", member.Role)
}

func TestProjectMemberRepository_RemoveMember(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectMemberRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "pmrem")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "pm4", Email: "pm4@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "pmrem")
	repo.AddMember(ctx, models.CreateProjectMemberParams{ProjectID: proj.ID, UserID: user.ID, Role: "member"})

	err := repo.RemoveMember(ctx, proj.ID, user.ID)
	require.NoError(t, err)

	member, _ := repo.GetMember(ctx, proj.ID, user.ID)
	assert.Nil(t, member)
}

func TestProjectMemberRepository_GetMembersByProject(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectMemberRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "pmpro")
	owner, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "owner_pro", Email: "owner_pro@t.com"})
	proj := createTestProject(t, projRepo, org.ID, owner.ID, "pmpro")

	// Add 3 members
	for i := 1; i <= 3; i++ {
		u, _ := userRepo.Create(ctx, models.CreateUserParams{
			ClerkUserID: uuid.New().String(),
			Email:       uuid.New().String() + "@t.com",
			FirstName:   testutil.Ptr("Member"),
			LastName:    testutil.Ptr(uuid.New().String()[:4]),
		})
		repo.AddMember(ctx, models.CreateProjectMemberParams{ProjectID: proj.ID, UserID: u.ID, Role: "member"})
	}

	members, err := repo.GetMembersByProject(ctx, proj.ID)
	require.NoError(t, err)
	assert.Len(t, members, 3)
	assert.NotEmpty(t, members[0].User.FirstName)
}

func TestProjectMemberRepository_GetProjectsByUser(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectMemberRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "pmusr")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "pmu", Email: "pmu@t.com"})

	// Add to 2 projects
	for i := 1; i <= 2; i++ {
		p := createTestProject(t, projRepo, org.ID, user.ID, uuid.New().String()[:8])
		repo.AddMember(ctx, models.CreateProjectMemberParams{ProjectID: p.ID, UserID: user.ID, Role: "member"})
	}

	projects, err := repo.GetProjectsByUser(ctx, user.ID, org.ID)
	require.NoError(t, err)
	assert.Len(t, projects, 2)
	assert.NotEmpty(t, projects[0].Project.Name)
}

func TestProjectMemberRepository_Aggregates(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewProjectMemberRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "pmagg")
	owner, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "owner_agg", Email: "owner_agg@t.com"})
	proj := createTestProject(t, projRepo, org.ID, owner.ID, "pmagg")

	u1, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "a1", Email: "a1@t.com"})
	u2, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "a2", Email: "a2@t.com"})

	repo.AddMember(ctx, models.CreateProjectMemberParams{ProjectID: proj.ID, UserID: u1.ID, Role: "admin"})
	repo.AddMember(ctx, models.CreateProjectMemberParams{ProjectID: proj.ID, UserID: u2.ID, Role: "member"})

	count, _ := repo.GetMemberCount(ctx, proj.ID)
	assert.Equal(t, 2, count)

	admins, _ := repo.GetAdminCount(ctx, proj.ID)
	assert.Equal(t, 1, admins)
}
