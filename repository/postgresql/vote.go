package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
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
	Db *pgxpool.Pool
}

func NewVotePostgresRepo(db *pgxpool.Pool) domain.VoteRepo {
	return &VotePostgresRepo{Db: db}
}

func (a *VotePostgresRepo) Create(threadSlugOrId string, vote *models.VoteCreate) (*models.Thread, error) {
	var thread models.Thread
	id, err := strconv.Atoi(threadSlugOrId)
	if err != nil {
		err = a.Db.QueryRow(context.Background(), GetThreadBySlugCommand, threadSlugOrId).Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	} else {
		err = a.Db.QueryRow(context.Background(), GetThreadByIdCommand, id).Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	}

	if err != nil {
		return nil, ErrorThreadDoesNotExist
	}

	var checkVote models.Vote
	err = a.Db.QueryRow(context.Background(), GetVoteByNicknameAndThreadCommand, vote.Nickname, thread.Id).Scan(&checkVote.Nickname, &checkVote.Thread, &checkVote.Voice)
	if err != nil {
		_, err = a.Db.Exec(context.Background(), CreateVoteCommand, vote.Nickname, thread.Id, vote.Voice)
		if err != nil {
			return nil, ErrorUserDoesNotExist
		}
		thread.Votes += vote.Voice
	} else {
		_, _ = a.Db.Exec(context.Background(), UpdateVoteCommand, vote.Voice, vote.Nickname, thread.Id)
		if vote.Voice != checkVote.Voice {
			thread.Votes += 2 * vote.Voice
		}
	}

	return &thread, nil
}
