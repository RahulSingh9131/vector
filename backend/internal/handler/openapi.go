package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/labstack/echo/v4"
)

type OpenAPIHandler struct {
	Handler
}

func NewOpenAPIHandler(s *server.Server) *OpenAPIHandler {
	return &OpenAPIHandler{
		Handler: NewHandler(s),
	}
}

func (h *OpenAPIHandler) ServeOpenAPIUI(c echo.Context) error {
	templateBytes, err := os.ReadFile("static/openapi.html")
	// do not cache the file so that we can see the new changes with refresh
	c.Response().Header().Set("Cache-Control", "no-cache")
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI UI template: %w", err)
	}

	templateString := string(templateBytes)

	err = c.HTML(http.StatusOK, templateString)
	if err != nil {
		return fmt.Errorf("failed to serve OpenAPI UI: %w", err)
	}

	return nil
}
