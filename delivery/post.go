package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
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
func (a *PostHandler) Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)

	slugOrId := ctx.UserValue("slug_or_id").(string)
	postsCreate := make([]models.PostCreate, 0)

	_ = json.Unmarshal(ctx.PostBody(), &postsCreate)

	posts, err := a.postRepo.Create(uctx, slugOrId, &postsCreate)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		if errors.Is(err, postgresql.ErrorThreadDoesNotExist) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		} else if errors.Is(err, postgresql.ErrorParentPostDoesNotExist) {
			ctx.SetStatusCode(fasthttp.StatusConflict)
		} else if errors.Is(err, postgresql.ErrorAuthorDoesNotExist) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
		return
	}

	body, _ := json.Marshal(posts)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusCreated)

	return
}

// GET post/{id}/details
func (a *PostHandler) Get(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))

	related := string(ctx.QueryArgs().Peek("related"))

	postGet := &models.PostGetRequest{
		Related: strings.Split(related, ","),
	}

	posts, err := a.postRepo.Get(uctx, int64(id), postGet)

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

// POST post/{id}/details
func (a *PostHandler) Update(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	var postUpdate models.PostUpdate

	_ = json.Unmarshal(ctx.PostBody(), &postUpdate)

	post, err := a.postRepo.Update(uctx, int64(id), &postUpdate)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	} else {
		body, _ := json.Marshal(post)
		ctx.SetBody(body)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}
}
