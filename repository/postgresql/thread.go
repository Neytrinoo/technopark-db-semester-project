package postgresql

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

type ThreadPostgresRepo struct {
	Db *sqlx.DB
}

const (
	CreateThreadCommand    = "INSERT INTO Threads (title, author, message, created, slug, forum) VALUES (:title, :author, :message, :created, :slug, :forum) RETURNING id"
	GetThreadByIdCommand   = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE id = $1"
	GetThreadBySlugCommand = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE slug = $1"
	UpdateThreadById       = "UPDATE Threads SET (title, message) = ($1, $2) WHERE id = $3"
)

var (
	ErrorNoAuthorOrForum    = errors.New("author or forum does not exist")
	ErrorThreadAlreadyExist = errors.New("thread already exist")
	ErrorThreadDoesNotExist = errors.New("thread does not exist")
)

func NewThreadPostgresRepo(db *sqlx.DB) domain.ThreadRepo {
	return &ThreadPostgresRepo{Db: db}
}

func (a *ThreadPostgresRepo) Create(forumSlug string, thread *models.ThreadCreate) (*models.Thread, error) {
	_, err := a.Db.Exec(GetForumCommand, forumSlug)
	if err != nil {
		return nil, ErrorNoAuthorOrForum
	}
	_, err = a.Db.Exec(GetUserByNicknameCommand, thread.Author)
	if err != nil {
		return nil, ErrorNoAuthorOrForum
	}

	stmt, _ := a.Db.PrepareNamed(CreateThreadCommand)
	var id int
	err = stmt.Get(&id, thread) // выполнит запрос и вернет id

	if err != nil {
		threadAlreadyExist, _ := a.Get(thread.Slug)
		return threadAlreadyExist, ErrorThreadAlreadyExist
	}

	threadToReturn, _ := a.Get(strconv.Itoa(id))

	return threadToReturn, nil
}

func (a *ThreadPostgresRepo) Get(threadSlugOrId string) (*models.Thread, error) {
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

	return &thread, nil
}

func (a *ThreadPostgresRepo) Update(threadSlugOrId string, updateData *models.ThreadUpdate) (*models.Thread, error) {
	thread, err := a.Get(threadSlugOrId)
	if err != nil {
		return nil, ErrorThreadDoesNotExist
	}

	if updateData.Message == "" {
		updateData.Message = thread.Message
	} else {
		thread.Message = updateData.Message
	}
	if updateData.Title == "" {
		updateData.Title = thread.Title
	} else {
		thread.Title = updateData.Title
	}

	_, _ = a.Db.Exec(UpdateThreadById, updateData.Title, updateData.Message, thread.Id)

	return thread, nil
}

func (a *ThreadPostgresRepo) GetPosts(getSettings *models.ThreadPostRequest) (*[]models.Post, error) {
	thread, err := a.Get(getSettings.SlugOrId)
	if err != nil {
		return nil, ErrorThreadDoesNotExist
	}
	// TODO
	return nil, nil
}
