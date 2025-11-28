package demo

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// User struct for data mapping
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserHandler struct {
	Repo IUserRepository
}

// Simple API handler example
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := h.Repo.GetByID(context.Background(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// ------- User Repository
type IUserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
}
