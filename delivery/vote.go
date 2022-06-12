package delivery

import (
	"github.com/labstack/echo/v4"
	"net/http"
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
func (a *VoteHandler) Create(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")
	var voteCreate models.VoteCreate
	_ = c.Bind(voteCreate)

	thread, err := a.voteRepo.Create(slugOrId, &voteCreate)
	
	if err != nil {
		return c.JSON(http.StatusNotFound, GetErrorMessage(err))
	}

	return c.JSON(http.StatusOK, thread)
}
