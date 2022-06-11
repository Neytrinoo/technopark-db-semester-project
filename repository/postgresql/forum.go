package postgresql

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	CreateForumCommand = "INSERT INTO Forums (title, user, slug) VALUES ($1, $2, $3) RETURNING id"
	GetForumCommand    = "SELECT title, user, slug, posts, threads FROM Forums WHERE slug = $1"
	GetUsersOnForum    = "SELECT nickname, fullname, about, email FROM " +
		"(SELECT DISTINCT author AS user_nickname FROM Threads WHERE forum=$1 " +
		"UNION DISTINCT " +
		"SELECT DISTINCT author AS user_nickname FROM Posts WHERE forum=$1) " +
		"LEFT JOIN Users ON Users.nickname = user_nickname"
	GetThreadsOnForum = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE forum=$1 ORDER BY created"
)

var (
	ErrorForumAlreadyExist = errors.New("forum already exist")
	ErrorForumDoesNotExist = errors.New("forum does not exist")
)

type ForumPostgresRepo struct {
	Db *sqlx.DB
}

func NewForumPostgresRepo(db *sqlx.DB) domain.ForumRepo {
	return &ForumPostgresRepo{Db: db}
}

func (a *ForumPostgresRepo) Create(forum *models.ForumCreate) (*models.Forum, error) {
	_, err := a.Db.Exec(GetUserByNicknameCommand, forum.User)
	if err != nil {
		return nil, ErrorUserDoesNotExist
	}

	_, err = a.Db.Exec(CreateForumCommand, forum.Title, forum.User, forum.Slug)
	if err != nil {
		forumAlreadyExist, _ := a.Get(forum.Slug)
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

func (a *ForumPostgresRepo) GetUsers(slug string) (*[]models.User, error) {
	users := make([]models.User, 0)

	err := a.Db.Get(&users, GetUsersOnForum, slug)
	if err != nil {
		return nil, ErrorForumDoesNotExist
	}

	return &users, nil
}

func (a *ForumPostgresRepo) GetThreads(slug string) (*[]models.Thread, error) {
	threads := make([]models.Thread, 0)

	err := a.Db.Get(&threads, GetThreadsOnForum, slug)
	if err != nil {
		return nil, ErrorForumDoesNotExist
	}

	return &threads, nil
}
