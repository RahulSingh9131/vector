package handler

import (
	"net/http"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	server   *server.Server
	services *service.Services
}

func NewUserHandler(s *server.Server, services *service.Services) *UserHandler {
	return &UserHandler{
		server:   s,
		services: services,
	}
}

// GetCurrentUser returns the profile of the currently authenticated user
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "User not found in context")
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser updates the profile of the currently authenticated user
func (h *UserHandler) UpdateCurrentUser(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "User not found in context")
	}

	var params models.UpdateUserParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	updatedUser, err := h.services.User.UpdateProfile(c.Request().Context(), user.ID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user profile")
	}

	return c.JSON(http.StatusOK, updatedUser)
}
