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

func TestUserRepository_Create(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	params := models.CreateUserParams{
		ClerkUserID: "user_clerk_1",
		Email:       "user1@test.com",
		FirstName:   testutil.Ptr("John"),
		LastName:    testutil.Ptr("Doe"),
		AvatarURL:   testutil.Ptr("https://example.com/avatar.png"),
	}

	user, err := repo.Create(ctx, params)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, params.ClerkUserID, user.ClerkUserID)
	assert.Equal(t, params.Email, user.Email)
	assert.Equal(t, params.FirstName, user.FirstName)
	assert.Equal(t, params.LastName, user.LastName)
	assert.Equal(t, params.AvatarURL, user.AvatarURL)
	assert.True(t, user.IsActive)
	assert.Nil(t, user.LastLoginAt)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestUserRepository_GetByID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	// Create user
	created, err := repo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_getbyid",
		Email:       "getbyid@test.com",
	})
	require.NoError(t, err)

	t.Run("existing user", func(t *testing.T) {
		user, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, "user_getbyid", user.ClerkUserID)
	})

	t.Run("non-existent user returns nil", func(t *testing.T) {
		user, err := repo.GetByID(ctx, uuid.New())
		require.NoError(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetByClerkID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	created, err := repo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_clerk_unique",
		Email:       "clerk@test.com",
	})
	require.NoError(t, err)

	t.Run("existing clerk ID", func(t *testing.T) {
		user, err := repo.GetByClerkID(ctx, created.ClerkUserID)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, created.ID, user.ID)
	})

	t.Run("non-existent clerk ID returns nil", func(t *testing.T) {
		user, err := repo.GetByClerkID(ctx, "clerk_nonexistent")
		require.NoError(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	email := "unique_email@test.com"
	created, err := repo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_email",
		Email:       email,
	})
	require.NoError(t, err)

	t.Run("existing email", func(t *testing.T) {
		user, err := repo.GetByEmail(ctx, email)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, created.ID, user.ID)
	})

	t.Run("non-existent email returns nil", func(t *testing.T) {
		user, err := repo.GetByEmail(ctx, "missing@test.com")
		require.NoError(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_Update(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	created, err := repo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_update",
		Email:       "before@test.com",
		FirstName:   testutil.Ptr("Before"),
	})
	require.NoError(t, err)

	t.Run("partial update", func(t *testing.T) {
		newName := "After"
		newEmail := "after@test.com"
		user, err := repo.Update(ctx, created.ID, models.UpdateUserParams{
			FirstName: &newName,
			Email:     &newEmail,
		})

		require.NoError(t, err)
		assert.Equal(t, newName, *user.FirstName)
		assert.Equal(t, newEmail, user.Email)
		assert.Equal(t, created.ClerkUserID, user.ClerkUserID, "clerk_user_id should not change")
	})
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	created, err := repo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_login",
		Email:       "login@test.com",
	})
	require.NoError(t, err)

	assert.Nil(t, created.LastLoginAt)

	err = repo.UpdateLastLogin(ctx, created.ID)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, updated.LastLoginAt)
	assert.WithinDuration(t, time.Now().UTC(), *updated.LastLoginAt, 5*time.Second)
}

func TestUserRepository_Delete(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	created, err := repo.Create(ctx, models.CreateUserParams{
		ClerkUserID: "user_delete",
		Email:       "delete@test.com",
	})
	require.NoError(t, err)
	assert.True(t, created.IsActive)

	err = repo.Delete(ctx, created.ID)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.False(t, updated.IsActive, "user should be inactive after soft delete")
}

func TestUserRepository_GetOrCreate(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	params := models.CreateUserParams{
		ClerkUserID: "unique_clerk_id",
		Email:       "getorcreate@test.com",
		FirstName:   testutil.Ptr("GetOrCreate"),
	}

	t.Run("creates new user", func(t *testing.T) {
		user, err := repo.GetOrCreate(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, params.ClerkUserID, user.ClerkUserID)
		assert.Equal(t, params.Email, user.Email)
	})

	t.Run("returns existing on duplicate clerk_id", func(t *testing.T) {
		// Try to "create" again with same clerk_id but different email
		params2 := params
		params2.Email = "different@test.com"

		user, err := repo.GetOrCreate(ctx, params2)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, params.Email, user.Email, "should return original user even if params differ")
	})
}

func TestUserRepository_List(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewUserRepository(testServer)
	ctx := context.Background()

	// Create 3 users
	for i := 1; i <= 3; i++ {
		_, err := repo.Create(ctx, models.CreateUserParams{
			ClerkUserID: uuid.New().String(),
			Email:       uuid.New().String() + "@test.com",
		})
		require.NoError(t, err)
	}

	users, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 3)
}
