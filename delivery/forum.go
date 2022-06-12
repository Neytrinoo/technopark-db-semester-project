package delivery

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
	"technopark-db-semester-project/repository/postgresql"
)

type ForumHandler struct {
	forumRepo domain.ForumRepo
}

func MakeForumHandler(forumRepo domain.ForumRepo) ForumHandler {
	return ForumHandler{forumRepo: forumRepo}
}

// POST forum/create
func (a *ForumHandler) Create(c echo.Context) error {
	var forumCreate models.ForumCreate
	_ = c.Bind(&forumCreate)

	forum, err := a.forumRepo.Create(&forumCreate)
	if err != nil {
		if errors.Is(err, postgresql.ErrorUserDoesNotExist) {
			return c.JSON(http.StatusNotFound, GetErrorMessage(err))
		} else if errors.Is(err, postgresql.ErrorForumAlreadyExist) {
			return c.JSON(http.StatusConflict, forum)
		}
	}

	return c.JSON(http.StatusCreated, forum)
}

// GET forum/{slug}/details
func (a *ForumHandler) Get(c echo.Context) error {
	slug := c.Param("slug")

	forum, err := a.forumRepo.Get(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, forum)
}

// GET forum/{slug}/users
func (a *ForumHandler) GetUsers(c echo.Context) error {
	slug := c.Param("slug")
	var forumUsers models.GetForumUsers
	_ = c.Bind(&forumUsers)
	forumUsers.Slug = slug
	if forumUsers.Limit == 0 {
		forumUsers.Limit = 100
	}

	users, err := a.forumRepo.GetUsers(&forumUsers)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, users)
}

// GET forum/{slug}/threads
func (a *ForumHandler) GetThreads(c echo.Context) error {
	slug := c.Param("slug")
	var forumThreads models.GetForumThreads
	_ = c.Bind(&forumThreads)
	if forumThreads.Limit == 0 {
		forumThreads.Limit = 100
	}

	threads, err := a.forumRepo.GetThreads(slug, &forumThreads)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, threads)
}
