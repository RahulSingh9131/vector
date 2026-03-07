package repository_test

import (
	"context"
	"testing"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	testutil "github.com/RahulSingh9131/vector/internal/testing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMemberTest creates a test org and user for member tests
func setupMemberTest(t *testing.T, testServer *server.Server) (
	*repository.OrganizationMemberRepository,
	*models.Organization,
	*models.User,
) {
	t.Helper()

	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	memberRepo := repository.NewOrganizationMemberRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, uuid.New().String()[:8])

	user, err := userRepo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_" + uuid.New().String()[:8],
		Email:       "member_" + uuid.New().String()[:8] + "@test.com",
	})
	require.NoError(t, err)

	return memberRepo, org, user
}

func TestOrganizationMemberRepository_AddMember(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	memberRepo, org, user := setupMemberTest(t, testServer)
	ctx := context.Background()

	t.Run("adds member with role", func(t *testing.T) {
		member, err := memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
			OrganizationID: org.ID,
			UserID:         user.ID,
			Role:           "org:member",
		})

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, member.ID)
		assert.Equal(t, org.ID, member.OrganizationID)
		assert.Equal(t, user.ID, member.UserID)
		assert.Equal(t, "org:member", member.Role)
		assert.False(t, member.JoinedAt.IsZero())
	})

	t.Run("upserts role on duplicate", func(t *testing.T) {
		// Same user, same org, different role — should upsert
		member, err := memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
			OrganizationID: org.ID,
			UserID:         user.ID,
			Role:           "org:admin",
		})

		require.NoError(t, err)
		assert.Equal(t, "org:admin", member.Role, "role should be updated to admin")
	})
}

func TestOrganizationMemberRepository_GetMember(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	memberRepo, org, user := setupMemberTest(t, testServer)
	ctx := context.Background()

	t.Run("returns nil for non-member", func(t *testing.T) {
		member, err := memberRepo.GetMember(ctx, org.ID, user.ID)
		require.NoError(t, err)
		assert.Nil(t, member)
	})

	t.Run("returns member after adding", func(t *testing.T) {
		_, err := memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
			OrganizationID: org.ID,
			UserID:         user.ID,
			Role:           "org:owner",
		})
		require.NoError(t, err)

		member, err := memberRepo.GetMember(ctx, org.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, member)
		assert.Equal(t, "org:owner", member.Role)
	})
}

func TestOrganizationMemberRepository_RemoveMember(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	memberRepo, org, user := setupMemberTest(t, testServer)
	ctx := context.Background()

	// Add member first
	_, err := memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Role:           "org:member",
	})
	require.NoError(t, err)

	t.Run("removes existing member", func(t *testing.T) {
		err := memberRepo.RemoveMember(ctx, org.ID, user.ID)
		require.NoError(t, err)

		// Verify member is gone
		member, err := memberRepo.GetMember(ctx, org.ID, user.ID)
		require.NoError(t, err)
		assert.Nil(t, member, "member should be removed")
	})

	t.Run("errors on non-existent member", func(t *testing.T) {
		err := memberRepo.RemoveMember(ctx, org.ID, uuid.New())
		assert.Error(t, err, "should error when removing non-existent member")
	})
}

func TestOrganizationMemberRepository_UpdateRole(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	memberRepo, org, user := setupMemberTest(t, testServer)
	ctx := context.Background()

	// Add member first
	_, err := memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Role:           "org:member",
	})
	require.NoError(t, err)

	t.Run("updates role successfully", func(t *testing.T) {
		member, err := memberRepo.UpdateRole(ctx, org.ID, user.ID, "org:admin")
		require.NoError(t, err)
		assert.Equal(t, "org:admin", member.Role)
	})

	t.Run("errors on non-existent member", func(t *testing.T) {
		_, err := memberRepo.UpdateRole(ctx, org.ID, uuid.New(), "org:admin")
		assert.Error(t, err)
	})
}

func TestOrganizationMemberRepository_GetMembersByOrganization(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	memberRepo := repository.NewOrganizationMemberRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "listmembers")

	t.Run("empty org returns empty list", func(t *testing.T) {
		members, err := memberRepo.GetMembersByOrganization(ctx, org.ID)
		require.NoError(t, err)
		assert.Empty(t, members)
	})

	t.Run("returns members with user details", func(t *testing.T) {
		// Create and add 2 users
		for i := 0; i < 2; i++ {
			user, err := userRepo.Create(ctx, models.CreateUserParams{
				ClerkUserID: "user_list_" + uuid.New().String()[:8],
				Email:       "list_" + uuid.New().String()[:8] + "@test.com",
			})
			require.NoError(t, err)

			_, err = memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
				OrganizationID: org.ID,
				UserID:         user.ID,
				Role:           "org:member",
			})
			require.NoError(t, err)
		}

		members, err := memberRepo.GetMembersByOrganization(ctx, org.ID)
		require.NoError(t, err)
		assert.Len(t, members, 2)

		// Verify user details are populated (JOIN worked)
		for _, m := range members {
			assert.NotNil(t, m.User, "user details should be populated")
			assert.NotEmpty(t, m.User.Email, "user email should be populated")
		}
	})
}

func TestOrganizationMemberRepository_GetOrganizationsByUser(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	memberRepo := repository.NewOrganizationMemberRepository(testServer)
	ctx := context.Background()

	user, err := userRepo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_orgsbyuser",
		Email:       "orgsbyuser@test.com",
	})
	require.NoError(t, err)

	// Add user to 2 different orgs
	for i := 0; i < 2; i++ {
		org := createTestOrg(t, orgRepo, "userorg_"+uuid.New().String()[:8])
		_, err := memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
			OrganizationID: org.ID,
			UserID:         user.ID,
			Role:           "org:member",
		})
		require.NoError(t, err)
	}

	memberships, err := memberRepo.GetOrganizationsByUser(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, memberships, 2)

	// Verify organization details are populated (JOIN worked)
	for _, m := range memberships {
		assert.NotNil(t, m.Organization, "organization details should be populated")
		assert.NotEmpty(t, m.Organization.Name, "org name should be populated")
	}
}

func TestOrganizationMemberRepository_GetMembersByRole(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	memberRepo := repository.NewOrganizationMemberRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "rolefilter")

	// Create users with different roles
	roles := []string{"org:owner", "org:admin", "org:member", "org:member"}
	for _, role := range roles {
		user, err := userRepo.Create(ctx, models.CreateUserParams{
			ClerkUserID: "user_role_" + uuid.New().String()[:8],
			Email:       "role_" + uuid.New().String()[:8] + "@test.com",
		})
		require.NoError(t, err)

		_, err = memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
			OrganizationID: org.ID,
			UserID:         user.ID,
			Role:           role,
		})
		require.NoError(t, err)
	}

	t.Run("filters by org:member", func(t *testing.T) {
		members, err := memberRepo.GetMembersByRole(ctx, org.ID, "org:member")
		require.NoError(t, err)
		assert.Len(t, members, 2, "should have 2 members with org:member role")
	})

	t.Run("filters by org:owner", func(t *testing.T) {
		members, err := memberRepo.GetMembersByRole(ctx, org.ID, "org:owner")
		require.NoError(t, err)
		assert.Len(t, members, 1, "should have 1 owner")
	})

	t.Run("returns empty for unused role", func(t *testing.T) {
		members, err := memberRepo.GetMembersByRole(ctx, org.ID, "org:guest")
		require.NoError(t, err)
		assert.Empty(t, members, "should have no guests")
	})
}
