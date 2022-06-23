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

type ThreadHandler struct {
	threadRepo domain.ThreadRepo
}

func MakeThreadHandler(threadRepo domain.ThreadRepo) ThreadHandler {
	return ThreadHandler{threadRepo: threadRepo}
}

// POST forum/{slug}/create
func (a *ThreadHandler) Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	slug := ctx.UserValue("slug").(string)

	var threadCreate models.ThreadCreate
	_ = json.Unmarshal(ctx.PostBody(), &threadCreate)

	thread, err := a.threadRepo.Create(uctx, slug, &threadCreate)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNoAuthorOrForum) {
			body, _ := json.Marshal(GetErrorMessage(err))
			ctx.SetBody(body)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		} else if errors.Is(err, postgresql.ErrorThreadAlreadyExist) {
			body, _ := json.Marshal(thread)
			ctx.SetBody(body)
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return
		}
	}

	body, _ := json.Marshal(thread)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusCreated)

	return
}

// GET thread/{slug_or_id}/details
func (a *ThreadHandler) Get(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	slugOrId := ctx.UserValue("slug_or_id").(string)

	thread, err := a.threadRepo.Get(uctx, slugOrId)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(thread)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}

// POST thread/{slug_or_id}/details
func (a *ThreadHandler) Update(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	slugOrId := ctx.UserValue("slug_or_id").(string)

	var threadUpdate models.ThreadUpdate
	_ = json.Unmarshal(ctx.PostBody(), &threadUpdate)

	thread, err := a.threadRepo.Update(uctx, slugOrId, &threadUpdate)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(thread)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}

// GET thread/{slug_or_id}/posts
func (a *ThreadHandler) GetPosts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	slugOrId := ctx.UserValue("slug_or_id").(string)

	limit, err := strconv.Atoi(string(ctx.QueryArgs().Peek("limit")))
	if err != nil {
		limit = 100
	}

	since, err := strconv.Atoi(string(ctx.QueryArgs().Peek("since")))
	if err != nil {
		since = -1
	}

	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort == "" {
		sort = models.Flat
	}

	desc, err := strconv.ParseBool(string(ctx.QueryArgs().Peek("desc")))
	if err != nil {
		desc = false
	}

	threadGetPosts := &models.ThreadPostRequest{
		Limit: int32(limit),
		Since: int64(since),
		Sort:  sort,
		Desc:  desc,
	}

	posts, err := a.threadRepo.GetPosts(uctx, slugOrId, threadGetPosts)

	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(posts)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}
