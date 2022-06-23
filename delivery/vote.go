package delivery

import (
	"context"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

type VoteHandler struct {
	voteRepo domain.VoteRepo
}

func MakeVoteHandler(voteRepo domain.VoteRepo) VoteHandler {
	return VoteHandler{voteRepo: voteRepo}
}

// POST thread/{slug_or_id}/vote
func (a *VoteHandler) Create(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)

	slugOrId := ctx.UserValue("slug_or_id").(string)
	var voteCreate models.VoteCreate
	_ = json.Unmarshal(ctx.PostBody(), &voteCreate)

	thread, err := a.voteRepo.Create(uctx, slugOrId, &voteCreate)

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
