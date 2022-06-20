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

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	since := c.QueryParam("since")

	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	forumUsers := &models.GetForumUsers{
		Slug:  slug,
		Limit: int32(limit),
		Since: since,
		Desc:  desc,
	}

	users, err := a.forumRepo.GetUsers(forumUsers)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, users)
}

// GET forum/{slug}/threads
func (a *ForumHandler) GetThreads(c echo.Context) error {
	slug := c.Param("slug")

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}
	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	forumThreads := &models.GetForumThreads{
		Limit: int32(limit),
		Since: c.QueryParam("since"),
		Desc:  desc,
	}

	threads, err := a.forumRepo.GetThreads(slug, forumThreads)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, threads)
}
