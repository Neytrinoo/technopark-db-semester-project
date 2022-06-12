package postgresql

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
	"time"
)

const (
	GetPostCommand       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE id = $1"
	GetPostAuthorCommand = "SELECT nickname, fullname, about, email FROM Users WHERE nickname = $1"
	GetPostForumCommand  = "SELECT title, user, slug, posts, threads FROM Forums WHERE slug = $1"
	GetPostThreadCommand = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE id = $1"
	UpdatePostCommand    = "UPDATE Posts SET (message, isEdited) = ($1, true) WHERE id = $2"
)

var (
	ErrorPostDoesNotExist   = errors.New("post does not exist")
	ErrorAuthorDoesNotExist = errors.New("author does not exist")
)

type PostPostgresRepo struct {
	Db *sqlx.DB
}

func NewPostPostgresRepo(db *sqlx.DB) domain.PostRepo {
	return &PostPostgresRepo{Db: db}
}

func isIn(arr *[]string, find string) bool {
	for _, str := range *arr {
		if str == find {
			return true
		}
	}

	return false
}

func (a *PostPostgresRepo) Get(id int64, getSettings *models.PostGetRequest) (*models.PostGetResult, error) {
	var post models.Post
	err := a.Db.Get(&post, GetPostCommand, id)
	if err != nil {
		return nil, ErrorPostDoesNotExist
	}

	var postResult models.PostGetResult
	postResult.Post = &post

	if isIn(&getSettings.Related, models.RelatedUser) {
		postResult.Author = &models.User{}
		_ = a.Db.Get(postResult.Author, GetPostAuthorCommand, post.Author)
	}
	if isIn(&getSettings.Related, models.RelatedThread) {
		postResult.Thread = &models.Thread{}
		_ = a.Db.Get(postResult.Thread, GetPostThreadCommand, post.Thread)
	}
	if isIn(&getSettings.Related, models.RelatedForum) {
		postResult.Forum = &models.Forum{}
		_ = a.Db.Get(postResult.Forum, GetPostForumCommand, post.Forum)
	}

	return &postResult, nil
}

func (a *PostPostgresRepo) Update(id int64, updateDate *models.PostUpdate) (*models.Post, error) {
	var post models.Post
	err := a.Db.Get(&post, GetPostCommand, id)
	if err != nil {
		return nil, ErrorPostDoesNotExist
	}

	if updateDate.Message == "" {
		return &post, nil
	}

	_, _ = a.Db.Exec(UpdatePostCommand, updateDate.Message, id)
	post.Message = updateDate.Message
	post.IsEdited = true

	return &post, nil
}

func (a *PostPostgresRepo) Create(threadSlugOrId string, posts *[]models.PostCreate) (*[]models.Post, error) {
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

	// TODO: пока что сделаю без проверки наличия каждого родительского поста
	postsResult := make([]models.Post, len(*posts))
	argsForCommand := make([]interface{}, len(*posts))
	dbCommand := "INSERT INTO Posts (parent, author, message, forum, thread, created) VALUES "

	createdTime := time.Now()
	for ind, post := range *posts {
		sixInd := ind * 6
		dbCommand += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", sixInd+1, sixInd+2, sixInd+3, sixInd+4, sixInd+5, sixInd+6)
		argsForCommand = append(argsForCommand, post.Parent, post.Author, post.Message, thread.Forum, thread.Id, createdTime)
		postsResult = append(postsResult, models.Post{Parent: post.Parent, Author: post.Author, Message: post.Message, Forum: thread.Forum, Thread: thread.Id, Created: createdTime})
	}

	dbCommand = dbCommand[:len(dbCommand)-1] + " RETURNING id"

	rows, err := a.Db.Query(dbCommand, argsForCommand...)
	if err != nil {
		return nil, ErrorAuthorDoesNotExist
	}

	for ind := 0; rows.Next(); ind++ {
		err = rows.Scan(&postsResult[ind].Id)
		if err != nil {
			return nil, ErrorAuthorDoesNotExist
		}
	}

	return &postsResult, nil
}
