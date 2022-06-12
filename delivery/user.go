package delivery

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
	"technopark-db-semester-project/repository/postgresql"
)

type UserHandler struct {
	userRepo domain.UserRepo
}

func MakeUserHandler(userRepo domain.UserRepo) UserHandler {
	return UserHandler{userRepo: userRepo}
}

// POST user/{nickname}/create
func (a *UserHandler) Create(c echo.Context) error {
	nickname := c.Param("nickname")
	var user models.User
	_ = c.Bind(&user)
	user.Nickname = nickname

	userAfterCreate, err := a.userRepo.Create(&user)
	if err != nil {
		return c.JSON(http.StatusConflict, userAfterCreate)
	}

	return c.JSON(http.StatusCreated, userAfterCreate)
}

// GET /user/{nickname}/profile
func (a *UserHandler) Get(c echo.Context) error {
	nickname := c.Param("nickname")
	user, err := a.userRepo.Get(nickname)
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, user)
}

// POST /user/{nickname}/profile
func (a *UserHandler) Update(c echo.Context) error {
	nickname := c.Param("nickname")
	var updateData models.UserUpdate
	_ = c.Bind(&updateData)

	user, err := a.userRepo.Update(nickname, &updateData)
	if err != nil {
		if errors.Is(err, postgresql.ErrorUserDoesNotExist) {
			return c.JSON(http.StatusNotFound, GetErrorMessage(err))
		} else if errors.Is(err, postgresql.ErrorConflictUpdateUser) {
			return c.JSON(http.StatusConflict, GetErrorMessage(err))
		}
	}

	return c.JSON(http.StatusOK, user)
}
