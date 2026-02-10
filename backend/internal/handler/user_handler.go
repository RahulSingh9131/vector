package handler

import (
	"net/http"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
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
	g.GET("/me", Handle(h.Handler, h.GetCurrentUser, http.StatusOK, &EmptyRequest{}))
	g.PUT("/me", Handle(h.Handler, h.UpdateCurrentUser, http.StatusOK, &UpdateCurrentUserRequest{}))
}
