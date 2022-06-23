package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
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
func (a *UserHandler) Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)

	nickname := ctx.UserValue("nickname").(string)
	var user models.User
	_ = json.Unmarshal(ctx.PostBody(), &user)
	user.Nickname = nickname

	userAfterCreate, err := a.userRepo.Create(uctx, &user)
	if err != nil {
		body, _ := json.Marshal(userAfterCreate)
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return
	}

	body, _ := json.Marshal((*userAfterCreate)[0])
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusCreated)

	return
}

// GET /user/{nickname}/profile
func (a *UserHandler) Get(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)

	nickname := ctx.UserValue("nickname").(string)
	user, err := a.userRepo.Get(uctx, nickname)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(user)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}

// POST /user/{nickname}/profile
func (a *UserHandler) Update(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)

	nickname := ctx.UserValue("nickname").(string)
	var updateData models.UserUpdate
	_ = json.Unmarshal(ctx.PostBody(), &updateData)

	user, err := a.userRepo.Update(uctx, nickname, &updateData)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		if errors.Is(err, postgresql.ErrorUserDoesNotExist) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		} else if errors.Is(err, postgresql.ErrorConflictUpdateUser) {
			ctx.SetStatusCode(fasthttp.StatusConflict)
		}
		return
	}

	body, _ := json.Marshal(user)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}
