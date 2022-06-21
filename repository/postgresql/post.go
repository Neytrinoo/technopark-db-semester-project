package postgresql

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
	"time"
)

const (
	GetPostCommand       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE id = $1;"
	GetPostAuthorCommand = "SELECT nickname, fullname, about, email FROM Users WHERE nickname = $1;"
	GetPostForumCommand  = "SELECT title, \"user\", slug, posts, threads FROM Forums WHERE slug = $1;"
	GetPostThreadCommand = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE id = $1;"
	UpdatePostCommand    = "UPDATE Posts SET (message, isEdited) = ($1, true) WHERE id = $2;"
)

var (
	ErrorPostDoesNotExist       = errors.New("post does not exist")
	ErrorAuthorDoesNotExist     = errors.New("author does not exist")
	ErrorParentPostDoesNotExist = errors.New("parent post does not exist")
)

type PostPostgresRepo struct {
	Db *pgxpool.Pool
}

func NewPostPostgresRepo(db *pgxpool.Pool) domain.PostRepo {
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
	err := a.Db.QueryRow(context.Background(), GetPostCommand, id).Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, ErrorPostDoesNotExist
	}

	var postResult models.PostGetResult
	postResult.Post = &post

	if isIn(&getSettings.Related, models.RelatedUser) {
		author := &models.User{}
		_ = a.Db.QueryRow(context.Background(), GetPostAuthorCommand, post.Author).Scan(&author.Nickname, &author.Fullname, &author.About, &author.Email)
		postResult.Author = author
	}
	if isIn(&getSettings.Related, models.RelatedThread) {
		thread := &models.Thread{}
		_ = a.Db.QueryRow(context.Background(), GetPostThreadCommand, post.Thread).Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		postResult.Thread = thread
	}
	if isIn(&getSettings.Related, models.RelatedForum) {
		forum := &models.Forum{}
		_ = a.Db.QueryRow(context.Background(), GetPostForumCommand, post.Forum).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
		postResult.Forum = forum
	}

	return &postResult, nil
}

func (a *PostPostgresRepo) Update(id int64, updateDate *models.PostUpdate) (*models.Post, error) {
	var post models.Post
	err := a.Db.QueryRow(context.Background(), GetPostCommand, id).Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, ErrorPostDoesNotExist
	}

	if updateDate.Message == "" || updateDate.Message == post.Message {
		return &post, nil
	}

	_ = a.Db.QueryRow(context.Background(), UpdatePostCommand, updateDate.Message, id)
	post.Message = updateDate.Message
	post.IsEdited = true

	return &post, nil
}

func (a *PostPostgresRepo) CheckParentAndAuthor(post *models.PostCreate) error {
	if post.Parent != 0 {
		_, err := a.Db.Exec(context.Background(), GetPostCommand, post.Parent)
		if err != nil {
			return err
		}
	}
	_, err := a.Db.Exec(context.Background(), GetUserByNicknameCommand, post.Author)

	return err
}

func (a *PostPostgresRepo) Create(ctx context.Context, threadSlugOrId string, posts *[]models.PostCreate) (*[]models.Post, error) {
	var thread models.Thread
	id, err := strconv.Atoi(threadSlugOrId)
	if err != nil {
		err = a.Db.QueryRow(ctx, GetThreadBySlugCommand, threadSlugOrId).Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	} else {
		err = a.Db.QueryRow(ctx, GetThreadByIdCommand, id).Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	}

	if err != nil {
		return nil, ErrorThreadDoesNotExist
	}

	if len(*posts) == 0 {
		postsToRet := make([]models.Post, 0)
		return &postsToRet, nil
	}

	if len(*posts) == 0 {
		emptyPosts := make([]models.Post, 0)
		return &emptyPosts, nil
	}

	if a.CheckParentAndAuthor(&(*posts)[0]) != nil {
		return nil, ErrorParentPostDoesNotExist
	}

	command := strings.Builder{}
	command.WriteString("INSERT INTO Posts (parent, author, message, forum, thread, created) VALUES ")

	argsForCommand := make([]interface{}, 0, len(*posts))
	postsToReturn := make([]models.Post, 0, len(*posts))
	createdTime := time.Unix(0, time.Now().UnixNano()/1e6*1e6)
	for ind, post := range *posts {
		if post.Parent != 0 {
			var parentPost models.Post
			err = a.Db.QueryRow(context.Background(), GetPostCommand, post.Parent).Scan(&parentPost.Id, &parentPost.Parent, &parentPost.Author, &parentPost.Message, &parentPost.IsEdited, &parentPost.Forum, &parentPost.Thread, &parentPost.Created)
			if err != nil || parentPost.Thread != thread.Id {
				return nil, ErrorParentPostDoesNotExist
			}
		}
		var author models.User
		err = a.Db.QueryRow(context.Background(), GetUserByNicknameCommand, post.Author).Scan(&author.Nickname, &author.Fullname, &author.About, &author.Email)
		if err != nil {
			return nil, ErrorAuthorDoesNotExist
		}
		sixInd := ind * 6

		postsToReturn = append(postsToReturn, models.Post{Parent: post.Parent, Author: post.Author, Message: post.Message, Forum: thread.Forum, Thread: thread.Id, Created: createdTime})
		fmt.Fprintf(&command, "($%d, $%d, $%d, $%d, $%d, $%d),", sixInd+1, sixInd+2, sixInd+3, sixInd+4, sixInd+5, sixInd+6)
		argsForCommand = append(argsForCommand, post.Parent, post.Author, post.Message, thread.Forum, thread.Id, createdTime)
	}

	qs := command.String()
	qs = qs[:len(qs)-1] + " RETURNING id"

	rows, err := a.Db.Query(ctx, qs, argsForCommand...)
	if err != nil {
		return nil, ErrorAuthorDoesNotExist
	}
	defer rows.Close()

	for ind := 0; rows.Next(); ind++ {
		err = rows.Scan(&postsToReturn[ind].Id)

		if err != nil {
			return nil, ErrorAuthorDoesNotExist
		}
	}

	return &postsToReturn, nil
}
