package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	CreateForumCommand = "INSERT INTO Forums (title, \"user\", slug) VALUES ($1, $2, $3)"
	GetForumCommand    = "SELECT title, \"user\", slug, posts, threads FROM Forums WHERE slug = $1"

	GetUsersOnForumCommand                 = "SELECT nickname, fullname, about, email FROM ForumUsers WHERE forum = $1 AND nickname > $2 ORDER BY nickname LIMIT $3"
	GetUsersOnForumDescCommand             = "SELECT nickname, fullname, about, email FROM ForumUsers WHERE forum = $1 AND nickname < $2 ORDER BY nickname DESC LIMIT $3"
	GetUsersOnForumWithoutSinceCommand     = "SELECT nickname, fullname, about, email FROM ForumUsers WHERE forum = $1 ORDER BY nickname LIMIT $2"
	GetUsersOnForumWithoutSinceDescCommand = "SELECT nickname, fullname, about, email FROM ForumUsers WHERE forum = $1 ORDER BY nickname DESC LIMIT $2"

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
	var user models.User
	err := a.Db.QueryRow(context.Background(), GetUserByNicknameCommand, forum.User).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return nil, ErrorUserDoesNotExist
	}

	_, err = a.Db.Exec(context.Background(), CreateForumCommand, forum.Title, user.Nickname, forum.Slug)
	if err != nil {
		forumAlreadyExist, _ := a.Get(forum.Slug)
		return forumAlreadyExist, ErrorForumAlreadyExist
	}

	forumToReturn := &models.Forum{
		Title:   forum.Title,
		User:    user.Nickname,
		Slug:    forum.Slug,
		Posts:   0,
		Threads: 0,
	}

	return forumToReturn, nil
}

func (a *ForumPostgresRepo) Get(slug string) (*models.Forum, error) {
	var forum models.Forum

	err := a.Db.QueryRow(context.Background(), GetForumCommand, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, ErrorForumDoesNotExist
	}

	return &forum, nil
}

func (a *ForumPostgresRepo) GetUsers(getSettings *models.GetForumUsers) (*[]models.User, error) {
	var err error
	var rows pgx.Rows

	_, err = a.Get(getSettings.Slug)
	if err != nil {
		return nil, ErrorForumDoesNotExist
	}

	if getSettings.Desc {
		if getSettings.Since == "" {
			rows, err = a.Db.Query(context.Background(), GetUsersOnForumWithoutSinceDescCommand, getSettings.Slug, getSettings.Limit)
		} else {
			rows, err = a.Db.Query(context.Background(), GetUsersOnForumDescCommand, getSettings.Slug, getSettings.Since, getSettings.Limit)
		}
	} else {
		if getSettings.Since == "" {
			rows, err = a.Db.Query(context.Background(), GetUsersOnForumWithoutSinceCommand, getSettings.Slug, getSettings.Limit)
		} else {
			rows, err = a.Db.Query(context.Background(), GetUsersOnForumCommand, getSettings.Slug, getSettings.Since, getSettings.Limit)
		}
	}

	if err != nil {
		return nil, ErrorForumDoesNotExist
	}
	defer rows.Close()

	users := make([]models.User, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, ErrorForumDoesNotExist
		}
		users = append(users, user)
	}

	return &users, nil
}

func (a *ForumPostgresRepo) GetThreads(slug string, getSettings *models.GetForumThreads) (*[]models.Thread, error) {
	var rows pgx.Rows

	_, err := a.Get(slug)
	if err != nil {
		return nil, ErrorForumDoesNotExist
	}

	if getSettings.Desc {
		if getSettings.Since == "" {
			rows, err = a.Db.Query(context.Background(), GetThreadsOnForumWithoutSinceDescCommand, slug, getSettings.Limit)
		} else {
			rows, err = a.Db.Query(context.Background(), GetThreadsOnForumDescCommand, slug, getSettings.Since, getSettings.Limit)
		}
	} else {
		if getSettings.Since == "" {
			rows, err = a.Db.Query(context.Background(), GetThreadsOnForumWithoutSinceCommand, slug, getSettings.Limit)
		} else {
			rows, err = a.Db.Query(context.Background(), GetThreadsOnForumCommand, slug, getSettings.Since, getSettings.Limit)
		}
	}

	if err != nil {
		return nil, ErrorForumDoesNotExist
	}
	defer rows.Close()

	threads := make([]models.Thread, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		thread := models.Thread{}
		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, ErrorForumDoesNotExist
		}

		threads = append(threads, thread)
	}

	return &threads, nil
}
