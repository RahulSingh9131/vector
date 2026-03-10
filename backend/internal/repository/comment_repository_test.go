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

func TestCommentRepository_Create(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewCommentRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "comcreate")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u_com", Email: "com@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "comcreate")
	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{
		ProjectID:  proj.ID,
		Title:      "Issue with comments",
		ReporterID: user.ID,
		Priority:   "low",
		Type:       "task",
	}, "K-1")

	params := models.CreateCommentParams{
		IssueID: issue.ID,
		AuthorID: user.ID,
		Body: "My first comment",
	}

	comment, err := repo.Create(ctx, params)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, comment.ID)
	assert.Equal(t, params.IssueID, comment.IssueID)
	assert.Equal(t, params.AuthorID, comment.AuthorID)
	assert.Equal(t, params.Body, comment.Body)
	assert.Nil(t, comment.ParentCommentID)
	assert.False(t, comment.IsEdited)
	assert.False(t, comment.IsDeleted)
	assert.False(t, comment.CreatedAt.IsZero())
}

func TestCommentRepository_GetByID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewCommentRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "comget")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u2", Email: "u2@t.com", FirstName: testutil.Ptr("Jane")})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "comget")
	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "I", ReporterID: user.ID, Priority: "low", Type: "task"}, "K-1")

	created, _ := repo.Create(ctx, models.CreateCommentParams{
		IssueID: issue.ID,
		AuthorID: user.ID,
		Body: "Hello",
	})

	t.Run("existing comment with author", func(t *testing.T) {
		comment, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, comment)
		assert.Equal(t, created.ID, comment.ID)
		assert.Equal(t, "Jane", *comment.AuthorFirstName)
	})

	t.Run("non-existent comment returns nil", func(t *testing.T) {
		comment, err := repo.GetByID(ctx, uuid.New())
		require.NoError(t, err)
		assert.Nil(t, comment)
	})
}

func TestCommentRepository_ListByIssue(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewCommentRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "comlist")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "comlist")
	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "I", ReporterID: user.ID, Priority: "low", Type: "task"}, "K-1")

	// Create a thread
	c1, _ := repo.Create(ctx, models.CreateCommentParams{IssueID: issue.ID, AuthorID: user.ID, Body: "Top 1"})
	c2, _ := repo.Create(ctx, models.CreateCommentParams{IssueID: issue.ID, AuthorID: user.ID, Body: "Top 2"})
	repo.Create(ctx, models.CreateCommentParams{IssueID: issue.ID, AuthorID: user.ID, Body: "Reply to 1", ParentCommentID: &c1.ID})

	threads, total, err := repo.ListByIssue(ctx, issue.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total, "only top-level comments should be counted in total")
	assert.Len(t, threads, 2)

	// Verify threading
	var t1 models.CommentThread
	for _, t := range threads {
		if t.Comment.ID == c1.ID {
			t1 = t
		}
	}
	assert.Len(t, t1.Replies, 1)
	assert.Equal(t, "Reply to 1", t1.Replies[0].Body)

	var t2 models.CommentThread
	for _, t := range threads {
		if t.Comment.ID == c2.ID {
			t2 = t
		}
	}
	assert.Len(t, t2.Replies, 0)
}

func TestCommentRepository_Update(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewCommentRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "comupd")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "comupd")
	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "I", ReporterID: user.ID, Priority: "low", Type: "task"}, "K-1")

	created, _ := repo.Create(ctx, models.CreateCommentParams{IssueID: issue.ID, AuthorID: user.ID, Body: "Old"})

	newBody := "New"
	updated, err := repo.Update(ctx, created.ID, models.UpdateCommentParams{Body: newBody})

	require.NoError(t, err)
	assert.Equal(t, newBody, updated.Body)
	assert.True(t, updated.IsEdited)
	assert.NotNil(t, updated.EditedAt)
}

func TestCommentRepository_SoftDelete(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewCommentRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "comdel")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "comdel")
	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "I", ReporterID: user.ID, Priority: "low", Type: "task"}, "K-1")

	created, _ := repo.Create(ctx, models.CreateCommentParams{IssueID: issue.ID, AuthorID: user.ID, Body: "Bye"})

	err := repo.SoftDelete(ctx, created.ID)
	require.NoError(t, err)

	comment, _ := repo.GetByID(ctx, created.ID)
	assert.True(t, comment.IsDeleted)
	assert.Equal(t, "[deleted]", comment.Body)
}

func TestCommentRepository_CountByIssue(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewCommentRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "comcount")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "comcount")
	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "I", ReporterID: user.ID, Priority: "low", Type: "task"}, "K-1")

	repo.Create(ctx, models.CreateCommentParams{IssueID: issue.ID, AuthorID: user.ID, Body: "C1"})
	c2, _ := repo.Create(ctx, models.CreateCommentParams{IssueID: issue.ID, AuthorID: user.ID, Body: "C2"})
	repo.SoftDelete(ctx, c2.ID)

	count, _ := repo.CountByIssue(ctx, issue.ID)
	assert.Equal(t, 1, count, "soft deleted comments should be excluded from count")
}

func TestCommentRepository_GetAuthorID(t *testing.T) {
	_, testServer, cleanup := testutil.SetupTest(t)
	defer cleanup()

	repo := repository.NewCommentRepository(testServer)
	orgRepo := repository.NewOrganizationRepository(testServer)
	userRepo := repository.NewUserRepository(testServer)
	projRepo := repository.NewProjectRepository(testServer)
	issueRepo := repository.NewIssueRepository(testServer)
	ctx := context.Background()

	org := createTestOrg(t, orgRepo, "comauth")
	user, _ := userRepo.Create(ctx, models.CreateUserParams{ClerkUserID: "u", Email: "u@t.com"})
	proj := createTestProject(t, projRepo, org.ID, user.ID, "comauth")
	issue, _ := issueRepo.Create(ctx, models.CreateIssueParams{ProjectID: proj.ID, Title: "I", ReporterID: user.ID, Priority: "low", Type: "task"}, "K-1")

	created, _ := repo.Create(ctx, models.CreateCommentParams{IssueID: issue.ID, AuthorID: user.ID, Body: "C1"})

	authorID, err := repo.GetAuthorID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, *authorID)
}
