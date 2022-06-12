package delivery

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
	"technopark-db-semester-project/repository/postgresql"
)

type ThreadHandler struct {
	threadRepo domain.ThreadRepo
}

func MakeThreadHandler(threadRepo domain.ThreadRepo) ThreadHandler {
	return ThreadHandler{threadRepo: threadRepo}
}

// POST forum/{slug}/create
func (a *ThreadHandler) Create(c echo.Context) error {
	slug := c.Param("slug")
	var threadCreate models.ThreadCreate
	_ = c.Bind(&threadCreate)

	thread, err := a.threadRepo.Create(slug, &threadCreate)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNoAuthorOrForum) {
			return c.JSON(http.StatusNotFound, GetErrorMessage(err))
		} else if errors.Is(err, postgresql.ErrorThreadAlreadyExist) {
			return c.JSON(http.StatusConflict, thread)
		}
	}

	return c.JSON(http.StatusCreated, thread)
}

// GET thread/{slug_or_id}/details
func (a *ThreadHandler) Get(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")

	thread, err := a.threadRepo.Get(slugOrId)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, thread)
}

// POST thread/{slug_or_id}/details
func (a *ThreadHandler) Update(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")

	var threadUpdate models.ThreadUpdate
	_ = c.Bind(&threadUpdate)

	thread, err := a.threadRepo.Update(slugOrId, &threadUpdate)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, thread)
}

// GET thread/{slug_or_id}/posts
func (a *ThreadHandler) GetPosts(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	var threadGetPosts models.ThreadPostRequest
	_ = c.Bind(&threadGetPosts)

	posts, err := a.threadRepo.GetPosts(slugOrId, &threadGetPosts)

	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, posts)
}
