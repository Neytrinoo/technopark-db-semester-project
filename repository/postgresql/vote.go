package postgresql

import (
	"github.com/jmoiron/sqlx"
	"strconv"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	GetVoteByNicknameAndThreadCommand = "SELECT nickname, thread, voice FROM Votes WHERE nickname = $1 AND thread = $2"
	CreateVoteCommand                 = "INSERT INTO Votes (nickname, thread, voice) VALUES ($1, $2, $3)"
	UpdateVoteCommand                 = "UPDATE Votes SET voice = $1 WHERE nickname = $2 AND thread = $3 AND voice != $1"
)

type VotePostgresRepo struct {
	Db *sqlx.DB
}

func NewVotePostgresRepo(db *sqlx.DB) domain.VoteRepo {
	return &VotePostgresRepo{Db: db}
}

func (a *VotePostgresRepo) Create(threadSlugOrId string, vote *models.VoteCreate) (*models.Thread, error) {
	var thread models.Thread
	id, err := strconv.Atoi(threadSlugOrId)
	if err != nil {
		err = a.Db.Get(&thread, GetThreadBySlugCommand, threadSlugOrId)
	} else {
		err = a.Db.Get(&thread, GetThreadByIdCommand, id)
	}

	if err != nil {
		return nil, ErrorThreadDoesNotExist
	}

	var checkVote models.Vote
	err = a.Db.Get(&checkVote, GetVoteByNicknameAndThreadCommand, vote.Nickname, thread.Id)
	if err != nil {
		_, _ = a.Db.Exec(CreateVoteCommand, vote.Nickname, thread.Id, vote.Voice)
		thread.Votes += vote.Voice
	} else {
		_, _ = a.Db.Exec(UpdateVoteCommand, vote.Voice, vote.Nickname, thread.Id)
		if vote.Voice != checkVote.Voice {
			thread.Votes += 2 * vote.Voice
		}
	}

	return &thread, nil
}
