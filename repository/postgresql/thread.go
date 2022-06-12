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
	CreateThreadCommand                               = "INSERT INTO Threads (title, author, message, created, slug, forum) VALUES (:title, :author, :message, :created, :slug, :forum) RETURNING id"
	GetThreadByIdCommand                              = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE id = $1"
	GetThreadBySlugCommand                            = "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE slug = $1"
	UpdateThreadByIdCommand                           = "UPDATE Threads SET (title, message) = ($1, $2) WHERE id = $3"
	GetPostsOnThreadFlatCommand                       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND id > $2 ORDER BY created, id LIMIT $3"
	GetPostsOnThreadFlatDescCommand                   = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND id < $2 ORDER BY created DESC, id DESC LIMIT $3"
	GetPostsOnThreadTreeCommand                       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND parent_path > (SELECT parent_path FROM Posts WHERE id = $2) ORDER BY parent_path, id LIMIT $3" // TODO: попробовать потом убрать id в order by, т.к. id и так содержится в parent_path
	GetPostsOnThreadTreeDescCommand                   = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND parent_path < (SELECT parent_path FROM Posts WHERE id = $2) ORDER BY parent_path DESC LIMIT $3"
	GetPostsOnThreadParentTreeCommand                 = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE parent_path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 AND id > COALESCE((SELECT parent_path[1] FROM Posts WHERE id = $2), 0) ORDER BY id LIMIT $3) ORDER BY parent_path, id"              // TODO: также попробовать убрать id в order by
	GetPostsOnThreadParentTreeDescWithSinceCommand    = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE parent_path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 AND id < (SELECT parent_path[1] FROM Posts WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY parent_path[1] DESC, parent_path, id" // TODO: также попробовать убрать id в order by
	GetPostsOnThreadParentTreeDescWithoutSinceCommand = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE parent_path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2) ORDER BY parent_path[1] DESC, parent_path, id"                                                           // TODO: также попробовать убрать id в order by
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
	thread.Forum = forumSlug

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

	_, _ = a.Db.Exec(UpdateThreadByIdCommand, updateData.Title, updateData.Message, thread.Id)

	return thread, nil
}

func (a *ThreadPostgresRepo) GetPosts(slugOrId string, getSettings *models.ThreadPostRequest) (*[]models.Post, error) {
	thread, err := a.Get(slugOrId)
	if err != nil {
		return nil, ErrorThreadDoesNotExist
	}

	posts := make([]models.Post, 0)

	if getSettings.Sort == models.Flat {
		if getSettings.Desc {
			_ = a.Db.Get(&posts, GetPostsOnThreadFlatDescCommand, thread.Id, getSettings.Since, getSettings.Limit)
		} else {
			_ = a.Db.Get(&posts, GetPostsOnThreadFlatCommand, thread.Id, getSettings.Since, getSettings.Limit)
		}
	} else if getSettings.Sort == models.Tree {
		if getSettings.Desc {
			_ = a.Db.Get(&posts, GetPostsOnThreadTreeDescCommand, thread.Id, getSettings.Since, getSettings.Limit)
		} else {
			_ = a.Db.Get(&posts, GetPostsOnThreadTreeCommand, thread.Id, getSettings.Since, getSettings.Limit)
		}
	} else if getSettings.Sort == models.ParentTree {
		if getSettings.Desc {
			if getSettings.Since > 0 {
				_ = a.Db.Get(&posts, GetPostsOnThreadParentTreeDescWithSinceCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				_ = a.Db.Get(&posts, GetPostsOnThreadParentTreeDescWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		} else {
			_ = a.Db.Get(&posts, GetPostsOnThreadParentTreeCommand, thread.Id, getSettings.Since, getSettings.Limit)
		}
	}

	return &posts, nil
}
