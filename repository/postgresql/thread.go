package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

type ThreadPostgresRepo struct {
	Db *pgxpool.Pool
}

const (
	CreateThreadCommand     = "INSERT INTO Threads (title, author, message, created, slug, forum) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	GetThreadByIdCommand    = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE id = $1"
	GetThreadBySlugCommand  = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE slug = $1"
	UpdateThreadByIdCommand = "UPDATE Threads SET (title, message) = ($1, $2) WHERE id = $3"

	GetPostsOnThreadFlatCommand                    = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND id > $2 ORDER BY created, id LIMIT $3"
	GetPostsOnThreadFlatDescCommand                = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND id < $2 ORDER BY created DESC, id DESC LIMIT $3"
	GetPostsOnThreadTreeCommand                    = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND parent_path > (SELECT parent_path FROM Posts WHERE id = $2) ORDER BY parent_path, id LIMIT $3" // TODO: попробовать потом убрать id в order by, т.к. id и так содержится в parent_path
	GetPostsOnThreadTreeDescCommand                = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND parent_path < (SELECT parent_path FROM Posts WHERE id = $2) ORDER BY parent_path DESC LIMIT $3"
	GetPostsOnThreadParentTreeCommand              = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE parent_path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 AND id > (SELECT parent_path[1] FROM Posts WHERE id = $2) ORDER BY id LIMIT $3) ORDER BY parent_path, id"                           // TODO: также попробовать убрать id в order by
	GetPostsOnThreadParentTreeDescWithSinceCommand = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE parent_path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 AND id < (SELECT parent_path[1] FROM Posts WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY parent_path[1] DESC, parent_path, id" // TODO: также попробовать убрать id в order by

	GetPostsOnThreadFlatWithoutSinceCommand           = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 ORDER BY created, id LIMIT $2"
	GetPostsOnThreadFlatDescWithoutSinceCommand       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 ORDER BY created DESC, id DESC LIMIT $2"
	GetPostsOnThreadTreeWithoutSinceCommand           = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 ORDER BY parent_path, id LIMIT $2" // TODO: попробовать потом убрать id в order by, т.к. id и так содержится в parent_path
	GetPostsOnThreadTreeDescWithoutSinceCommand       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 ORDER BY parent_path DESC LIMIT $2"
	GetPostsOnThreadParentTreeWithoutSinceCommand     = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE parent_path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 ORDER BY id LIMIT $2) ORDER BY parent_path, id"                           // TODO: также попробовать убрать id в order by
	GetPostsOnThreadParentTreeDescWithoutSinceCommand = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE parent_path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2) ORDER BY parent_path[1] DESC, parent_path, id" // TODO: также попробовать убрать id в order by
)

var (
	ErrorNoAuthorOrForum    = errors.New("author or forum does not exist")
	ErrorThreadAlreadyExist = errors.New("thread already exist")
	ErrorThreadDoesNotExist = errors.New("thread does not exist")
)

func NewThreadPostgresRepo(db *pgxpool.Pool) domain.ThreadRepo {
	return &ThreadPostgresRepo{Db: db}
}

func (a *ThreadPostgresRepo) Create(forumSlug string, thread *models.ThreadCreate) (*models.Thread, error) {
	var forum models.Forum
	err := a.Db.QueryRow(context.Background(), GetForumCommand, forumSlug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, ErrorNoAuthorOrForum
	}

	var user models.User
	err = a.Db.QueryRow(context.Background(), GetUserByNicknameCommand, thread.Author).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return nil, ErrorNoAuthorOrForum
	}
	thread.Forum = forum.Slug

	if thread.Slug != "" {
		var threadAlreadyExist models.Thread
		err = a.Db.QueryRow(context.Background(), GetThreadBySlugCommand, thread.Slug).Scan(&threadAlreadyExist.Id, &threadAlreadyExist.Title, &threadAlreadyExist.Author, &threadAlreadyExist.Forum, &threadAlreadyExist.Message, &threadAlreadyExist.Votes, &threadAlreadyExist.Slug, &threadAlreadyExist.Created)
		if err == nil {
			return &threadAlreadyExist, ErrorThreadAlreadyExist
		}
	}

	var id int32
	err = a.Db.QueryRow(context.Background(), CreateThreadCommand, thread.Title, thread.Author, thread.Message, thread.Created, thread.Slug, thread.Forum).Scan(&id)
	if err != nil {
		threadAlreadyExist, _ := a.Get(thread.Slug)
		return threadAlreadyExist, ErrorThreadAlreadyExist
	}

	threadToReturn := &models.Thread{
		Id:      id,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Message: thread.Message,
		Votes:   0,
		Slug:    thread.Slug,
		Created: thread.Created,
	}

	return threadToReturn, nil
}

func (a *ThreadPostgresRepo) Get(threadSlugOrId string) (*models.Thread, error) {
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

	_, _ = a.Db.Exec(context.Background(), UpdateThreadByIdCommand, updateData.Title, updateData.Message, thread.Id)

	return thread, nil
}

func (a *ThreadPostgresRepo) GetPosts(slugOrId string, getSettings *models.ThreadPostRequest) (*[]models.Post, error) {
	thread, err := a.Get(slugOrId)
	if err != nil {
		return nil, ErrorThreadDoesNotExist
	}

	var rows pgx.Rows

	if getSettings.Sort == models.Flat {
		if getSettings.Desc {
			if getSettings.Since != -1 {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadFlatDescCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadFlatDescWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		} else {
			if getSettings.Since != -1 {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadFlatCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadFlatWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		}
	} else if getSettings.Sort == models.Tree {
		if getSettings.Desc {
			if getSettings.Since != -1 {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadTreeDescCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadTreeDescWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		} else {
			if getSettings.Since != -1 {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadTreeCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadTreeWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		}
	} else if getSettings.Sort == models.ParentTree {
		if getSettings.Desc {
			if getSettings.Since > 0 {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadParentTreeDescWithSinceCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadParentTreeDescWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		} else {
			if getSettings.Since != -1 {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadParentTreeCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = a.Db.Query(context.Background(), GetPostsOnThreadParentTreeWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		}
	}
	defer rows.Close()

	posts := make([]models.Post, 0, rows.CommandTag().RowsAffected())

	for rows.Next() {
		post := models.Post{}
		_ = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		posts = append(posts, post)
	}

	return &posts, nil
}
