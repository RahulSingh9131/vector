package handler

import (
	"net/http"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Handler
	services *service.Services
}

func NewUserHandler(s *server.Server, services *service.Services) *UserHandler {
	return &UserHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

// GetCurrentUser returns the profile of the currently authenticated user
func (h *UserHandler) GetCurrentUser(c echo.Context, req *EmptyRequest) (*models.User, error) {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return user, nil
}

// ListUsers returns a list of all users
func (h *UserHandler) ListUsers(c echo.Context, req *EmptyRequest) ([]models.User, error) {
	return h.services.User.ListUsers(c.Request().Context())
}

// CreateUser creates a new user manually
func (h *UserHandler) CreateUser(c echo.Context, req *CreateUserRequest) (*models.User, error) {
	return h.services.User.CreateUser(c.Request().Context(), models.CreateUserParams{
		ClerkUserID: req.ClerkUserID,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		AvatarURL:   req.AvatarURL,
	})
}

// GetUser returns details of a specific user
func (h *UserHandler) GetUser(c echo.Context, req *GetUserRequest) (*models.User, error) {
	id, _ := uuid.Parse(req.ID)
	return h.services.User.GetByID(c.Request().Context(), id)
}

// DeleteUser deactivates a user account
func (h *UserHandler) DeleteUser(c echo.Context, req *DeleteUserRequest) (*EmptyResponse, error) {
	id, _ := uuid.Parse(req.ID)
	if err := h.services.User.DeactivateUser(c.Request().Context(), id); err != nil {
		return nil, err
	}
	return &EmptyResponse{}, nil
}

// UpdateCurrentUser updates the profile of the currently authenticated user
func (h *UserHandler) UpdateCurrentUser(c echo.Context, req *UpdateCurrentUserRequest) (*models.User, error) {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	updatedUser, err := h.services.User.UpdateProfile(c.Request().Context(), user.ID, models.UpdateUserParams{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		return nil, err // base.go handles logging + New Relic; GlobalErrorHandler maps sqlerr
	}

	return updatedUser, nil
}

// RegisterRoutes registers all user routes
func (h *UserHandler) RegisterRoutes(g *echo.Group) {
	// Profile routes
	g.GET("/me", Handle(h.Handler, h.GetCurrentUser, http.StatusOK, &EmptyRequest{}))
	g.PUT("/me", Handle(h.Handler, h.UpdateCurrentUser, http.StatusOK, &UpdateCurrentUserRequest{}))

	// Management routes
	g.GET("/users", Handle(h.Handler, h.ListUsers, http.StatusOK, &EmptyRequest{}))
	g.POST("/users", Handle(h.Handler, h.CreateUser, http.StatusCreated, &CreateUserRequest{}))
	g.GET("/users/:id", Handle(h.Handler, h.GetUser, http.StatusOK, &GetUserRequest{}))
	g.DELETE("/users/:id", Handle(h.Handler, h.DeleteUser, http.StatusNoContent, &DeleteUserRequest{}))
}
