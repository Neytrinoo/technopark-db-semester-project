package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
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
func (a *ForumHandler) Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)

	var forumCreate models.ForumCreate
	_ = json.Unmarshal(ctx.PostBody(), &forumCreate)

	forum, err := a.forumRepo.Create(uctx, &forumCreate)

	if err != nil {
		if errors.Is(err, postgresql.ErrorUserDoesNotExist) {
			body, _ := json.Marshal(GetErrorMessage(err))
			ctx.SetBody(body)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		} else if errors.Is(err, postgresql.ErrorForumAlreadyExist) {
			body, _ := json.Marshal(forum)
			ctx.SetBody(body)
			ctx.SetStatusCode(fasthttp.StatusConflict)
		}

		return
	}

	body, _ := json.Marshal(forum)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusCreated)

	return
}

// GET forum/{slug}/details
func (a *ForumHandler) Get(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	slug := ctx.UserValue("slug").(string)

	forum, err := a.forumRepo.Get(uctx, slug)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(forum)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}

// GET forum/{slug}/users
func (a *ForumHandler) GetUsers(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	slug := ctx.UserValue("slug").(string)

	limit, err := strconv.Atoi(string(ctx.QueryArgs().Peek("limit")))
	if err != nil {
		limit = 100
	}

	since := string(ctx.QueryArgs().Peek("since"))

	desc, err := strconv.ParseBool(string(ctx.QueryArgs().Peek("desc")))
	if err != nil {
		desc = false
	}

	forumUsers := &models.GetForumUsers{
		Slug:  slug,
		Limit: int32(limit),
		Since: since,
		Desc:  desc,
	}

	users, err := a.forumRepo.GetUsers(uctx, forumUsers)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(users)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)

	return
}

// GET forum/{slug}/threads
func (a *ForumHandler) GetThreads(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	slug := ctx.UserValue("slug").(string)

	limit, err := strconv.Atoi(string(ctx.QueryArgs().Peek("limit")))
	if err != nil {
		limit = 100
	}

	desc, err := strconv.ParseBool(string(ctx.QueryArgs().Peek("desc")))
	if err != nil {
		desc = false
	}

	forumThreads := &models.GetForumThreads{
		Limit: int32(limit),
		Since: string(ctx.QueryArgs().Peek("since")),
		Desc:  desc,
	}

	threads, err := a.forumRepo.GetThreads(uctx, slug, forumThreads)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(threads)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}
