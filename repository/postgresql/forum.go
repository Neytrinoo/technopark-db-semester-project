package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	CreateForumCommand     = "INSERT INTO Forums (title, \"user\", slug) VALUES ($1, $2, $3) RETURNING id"
	GetForumCommand        = "SELECT title, \"user\", slug, posts, threads FROM Forums WHERE slug = $1"
	GetUsersOnForumCommand = "SELECT nickname, fullname, about, email FROM " +
		"(SELECT DISTINCT author AS user_nickname FROM Threads WHERE forum = $1 AND author > $2" +
		"UNION DISTINCT " +
		"SELECT DISTINCT author AS user_nickname FROM Posts WHERE forum = $1 AND author > $2) " +
		"LEFT JOIN Users ON Users.nickname = user_nickname " +
		"ORDER BY nickname LIMIT $3"
	GetUsersOnForumDescCommand = "SELECT nickname, fullname, about, email FROM " +
		"(SELECT DISTINCT author AS user_nickname FROM Threads WHERE forum = $1 AND author < $2" +
		"UNION DISTINCT " +
		"SELECT DISTINCT author AS user_nickname FROM Posts WHERE forum = $1 AND author < $2) " +
		"LEFT JOIN Users ON Users.nickname = user_nickname " +
		"ORDER BY nickname DESC LIMIT $3"
	GetUsersOnForumWithoutSinceCommand = "SELECT nickname, fullname, about, email FROM " +
		"(SELECT DISTINCT author AS user_nickname FROM Threads WHERE forum = $1" +
		"UNION DISTINCT " +
		"SELECT DISTINCT author AS user_nickname FROM Posts WHERE forum = $1) " +
		"LEFT JOIN Users ON Users.nickname = user_nickname " +
		"ORDER BY nickname LIMIT $2"
	GetUsersOnForumWithoutSinceDescCommand = "SELECT nickname, fullname, about, email FROM " +
		"(SELECT DISTINCT author AS user_nickname FROM Threads WHERE forum = $1" +
		"UNION DISTINCT " +
		"SELECT DISTINCT author AS user_nickname FROM Posts WHERE forum = $1) " +
		"LEFT JOIN Users ON Users.nickname = user_nickname " +
		"ORDER BY nickname DESC LIMIT $2"

	GetThreadsOnForumCommand                 = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE forum = $1 AND created >= $2 ORDER BY created LIMIT $3"
	GetThreadsOnForumDescCommand             = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE forum = $1 AND created <= $2 ORDER BY created DESC LIMIT $3"
	GetThreadsOnForumWithoutSinceCommand     = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE forum = $1 ORDER BY created LIMIT $2"
	GetThreadsOnForumWithoutSinceDescCommand = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE forum = $1 ORDER BY created DESC LIMIT $2"
)

var (
	ErrorForumAlreadyExist = errors.New("forum already exist")
	ErrorForumDoesNotExist = errors.New("forum does not exist")
)

type ForumPostgresRepo struct {
	Db *pgxpool.Pool
}

func NewForumPostgresRepo(db *pgxpool.Pool) domain.ForumRepo {
	return &ForumPostgresRepo{Db: db}
}

func (a *ForumPostgresRepo) Create(forum *models.ForumCreate) (*models.Forum, error) {
	_, err := a.Db.Query(context.Background(), GetUserByNicknameCommand, forum.User)
	if err != nil {
		return nil, ErrorUserDoesNotExist
	}

	_, err = a.Db.Query(context.Background(), CreateForumCommand, forum.Title, forum.User, forum.Slug)
	if err != nil {
		forumAlreadyExist, _ := a.Get(forum.Slug)
		log.Println("error in create:", err)
		return forumAlreadyExist, ErrorForumAlreadyExist
	}

	forumToReturn, _ := a.Get(forum.Slug)

	return forumToReturn, nil
}

func (a *ForumPostgresRepo) Get(slug string) (*models.Forum, error) {
	var forum models.Forum

	err := a.Db.Get(&forum, GetForumCommand, slug)
	if err != nil {
		return nil, ErrorForumDoesNotExist
	}

	return &forum, nil
}

func (a *ForumPostgresRepo) GetUsers(getSettings *models.GetForumUsers) (*[]models.User, error) {
	var err error
	users := make([]models.User, 0)

	if getSettings.Desc {
		if getSettings.Since != "" {
			err = a.Db.Get(&users, GetUsersOnForumDescCommand, getSettings.Slug, getSettings.Since, getSettings.Limit)
		} else {
			err = a.Db.Get(&users, GetUsersOnForumWithoutSinceDescCommand, getSettings.Slug, getSettings.Limit)
		}
	} else {
		if getSettings.Since != "" {
			err = a.Db.Get(&users, GetUsersOnForumCommand, getSettings.Slug, getSettings.Since, getSettings.Limit)
		} else {
			err = a.Db.Get(&users, GetUsersOnForumWithoutSinceCommand, getSettings.Slug, getSettings.Limit)
		}
	}

	if err != nil {
		return nil, ErrorForumDoesNotExist
	}

	return &users, nil
}

func (a *ForumPostgresRepo) GetThreads(slug string, getSettings *models.GetForumThreads) (*[]models.Thread, error) {
	threads := make([]models.Thread, 0)
	var err error

	if getSettings.Desc {
		if getSettings.Since.IsZero() {
			err = a.Db.Get(&threads, GetThreadsOnForumWithoutSinceDescCommand, slug, getSettings.Limit)
		} else {
			err = a.Db.Get(&threads, GetThreadsOnForumDescCommand, slug, getSettings.Since, getSettings.Limit)
		}
	} else {
		if getSettings.Since.IsZero() {
			err = a.Db.Get(&threads, GetThreadsOnForumWithoutSinceCommand, slug, getSettings.Limit)
		} else {
			err = a.Db.Get(&threads, GetThreadsOnForumCommand, slug, getSettings.Since, getSettings.Limit)
		}
	}

	if err != nil {
		return nil, ErrorForumDoesNotExist
	}

	return &threads, nil
}
