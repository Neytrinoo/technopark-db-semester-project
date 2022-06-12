package delivery

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
	"technopark-db-semester-project/repository/postgresql"
)

type PostHandler struct {
	postRepo domain.PostRepo
}

func MakePostHandler(postRepo domain.PostRepo) PostHandler {
	return PostHandler{postRepo: postRepo}
}

// POST thread/{slug_or_id}/create
func (a *PostHandler) Create(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	postsCreate := make([]models.PostCreate, 0)
	_ = c.Bind(&postsCreate)

	posts, err := a.postRepo.Create(slugOrId, &postsCreate)
	if err != nil {
		if errors.Is(err, postgresql.ErrorThreadDoesNotExist) {
			return c.JSON(http.StatusNotFound, GetErrorMessage(err))
		} else if errors.Is(err, postgresql.ErrorAuthorDoesNotExist) {
			return c.JSON(http.StatusConflict, GetErrorMessage(err))
		}
	}

	return c.JSON(http.StatusCreated, posts)
}

// GET post/{id}/details
func (a *PostHandler) Get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var postGet models.PostGetRequest
	_ = c.Bind(&postGet)

	posts, err := a.postRepo.Get(int64(id), &postGet)

	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, posts)
}

// POST post/{id}/details
func (a *PostHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var postUpdate models.PostUpdate
	_ = c.Bind(&postUpdate)

	post, err := a.postRepo.Update(int64(id), &postUpdate)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	} else {
		return c.JSON(http.StatusOK, post)
	}
}
