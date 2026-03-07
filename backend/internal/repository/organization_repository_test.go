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

// createTestOrg is a helper to create an organization for testing
func createTestOrg(t *testing.T, repo *repository.OrganizationRepository, suffix string) *models.Organization {
	t.Helper()

	org, err := repo.Create(context.Background(), models.CreateOrganizationParams{
		ClerkOrgID:       "org_test_" + suffix,
		Name:             "Test Org " + suffix,
		Slug:             "test-org-" + suffix,
		SubscriptionTier: "free",
		MaxMembers:       10,
		MaxProjects:      5,
	})
	require.NoError(t, err)
	return org
}

func TestOrganizationRepository_Create(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewOrganizationRepository(testServer)
	ctx := context.Background()

	logoURL := "https://example.com/logo.png"
	org, err := repo.Create(ctx, models.CreateOrganizationParams{
		ClerkOrgID:       "org_clerk_123",
		Name:             "My Organization",
		Slug:             "my-org",
		LogoURL:          &logoURL,
		SubscriptionTier: "free",
		MaxMembers:       10,
		MaxProjects:      5,
	})

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, org.ID, "should auto-generate UUID")
	assert.Equal(t, "org_clerk_123", org.ClerkOrgID)
	assert.Equal(t, "My Organization", org.Name)
	assert.Equal(t, "my-org", org.Slug)
	assert.Equal(t, &logoURL, org.LogoURL)
	assert.Equal(t, "free", org.SubscriptionTier)
	assert.Equal(t, 10, org.MaxMembers)
	assert.Equal(t, 5, org.MaxProjects)
	assert.True(t, org.IsActive, "should default to active")
	assert.False(t, org.CreatedAt.IsZero(), "should auto-set created_at")
	assert.False(t, org.UpdatedAt.IsZero(), "should auto-set updated_at")
}

func TestOrganizationRepository_GetByID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewOrganizationRepository(testServer)
	ctx := context.Background()

	// Create an org first
	created := createTestOrg(t, repo, "getbyid")

	t.Run("existing org", func(t *testing.T) {
		org, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, org)
		assert.Equal(t, created.ID, org.ID)
		assert.Equal(t, created.Name, org.Name)
	})

	t.Run("non-existent org returns nil", func(t *testing.T) {
		org, err := repo.GetByID(ctx, uuid.New())
		require.NoError(t, err)
		assert.Nil(t, org)
	})
}

func TestOrganizationRepository_GetByClerkID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewOrganizationRepository(testServer)
	ctx := context.Background()

	created := createTestOrg(t, repo, "getbyclerk")

	t.Run("existing clerk ID", func(t *testing.T) {
		org, err := repo.GetByClerkID(ctx, created.ClerkOrgID)
		require.NoError(t, err)
		require.NotNil(t, org)
		assert.Equal(t, created.ID, org.ID)
	})

	t.Run("non-existent clerk ID returns nil", func(t *testing.T) {
		org, err := repo.GetByClerkID(ctx, "org_nonexistent")
		require.NoError(t, err)
		assert.Nil(t, org)
	})
}

func TestOrganizationRepository_GetBySlug(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewOrganizationRepository(testServer)
	ctx := context.Background()

	created := createTestOrg(t, repo, "getbyslug")

	t.Run("existing slug", func(t *testing.T) {
		org, err := repo.GetBySlug(ctx, created.Slug)
		require.NoError(t, err)
		require.NotNil(t, org)
		assert.Equal(t, created.ID, org.ID)
	})

	t.Run("non-existent slug returns nil", func(t *testing.T) {
		org, err := repo.GetBySlug(ctx, "nonexistent-slug")
		require.NoError(t, err)
		assert.Nil(t, org)
	})
}

func TestOrganizationRepository_Update(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewOrganizationRepository(testServer)
	ctx := context.Background()

	created := createTestOrg(t, repo, "update")

	t.Run("partial update - only name", func(t *testing.T) {
		newName := "Updated Name"
		org, err := repo.Update(ctx, created.ID, models.UpdateOrganizationParams{
			Name: &newName,
		})

		require.NoError(t, err)
		assert.Equal(t, "Updated Name", org.Name)
		assert.Equal(t, created.Slug, org.Slug, "slug should remain unchanged")
		assert.Equal(t, created.MaxMembers, org.MaxMembers, "max_members should remain unchanged")
	})

	t.Run("partial update - name and slug", func(t *testing.T) {
		newName := "Final Name"
		newSlug := "final-slug"
		org, err := repo.Update(ctx, created.ID, models.UpdateOrganizationParams{
			Name: &newName,
			Slug: &newSlug,
		})

		require.NoError(t, err)
		assert.Equal(t, "Final Name", org.Name)
		assert.Equal(t, "final-slug", org.Slug)
	})
}

func TestOrganizationRepository_Delete(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewOrganizationRepository(testServer)
	ctx := context.Background()

	created := createTestOrg(t, repo, "delete")

	// Delete (soft delete)
	err := repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify it's soft deleted
	org, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, org, "soft deleted org should still be retrievable")
	assert.False(t, org.IsActive, "is_active should be false after delete")
}

func TestOrganizationRepository_GetOrCreate(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewOrganizationRepository(testServer)
	ctx := context.Background()

	params := models.CreateOrganizationParams{
		ClerkOrgID:       "org_getorcreate",
		Name:             "GetOrCreate Org",
		Slug:             "getorcreate-org",
		SubscriptionTier: "free",
		MaxMembers:       10,
		MaxProjects:      5,
	}

	t.Run("creates new org", func(t *testing.T) {
		org, err := repo.GetOrCreate(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, org)
		assert.Equal(t, "GetOrCreate Org", org.Name)
		assert.Equal(t, "org_getorcreate", org.ClerkOrgID)
	})

	t.Run("returns existing on duplicate clerk_org_id", func(t *testing.T) {
		// Call again with same ClerkOrgID but different name
		dupParams := params
		dupParams.Name = "Different Name"

		org, err := repo.GetOrCreate(ctx, dupParams)
		require.NoError(t, err)
		require.NotNil(t, org)
		assert.Equal(t, "GetOrCreate Org", org.Name, "should return the original, not the new name")
	})
}

func TestOrganizationRepository_GetMemberCount(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	orgRepo := repository.NewOrganizationRepository(testServer)
	memberRepo := repository.NewOrganizationMemberRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "membercount")

	t.Run("empty org returns 0", func(t *testing.T) {
		count, err := orgRepo.GetMemberCount(ctx, org.ID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("returns correct count after adding members", func(t *testing.T) {
		// Create test users and add them to the org
		for i := 0; i < 3; i++ {
			user, err := userRepo.Create(ctx, models.CreateUserParams{
				ClerkUserID: "user_count_" + uuid.New().String()[:8],
				Email:       "count" + uuid.New().String()[:8] + "@test.com",
			})
			require.NoError(t, err)

			_, err = memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
				OrganizationID: org.ID,
				UserID:         user.ID,
				Role:           "org:member",
			})
			require.NoError(t, err)
		}

		count, err := orgRepo.GetMemberCount(ctx, org.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})
}
