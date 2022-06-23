package domain

import (
	"context"
	"technopark-db-semester-project/domain/models"
)

type UserRepo interface {
	Create(ctx context.Context, user *models.User) (*[]models.User, error)
	Update(ctx context.Context, nickname string, updateData *models.UserUpdate) (*models.User, error)
	Get(ctx context.Context, nicknameOrEmail string) (*models.User, error)
}

type ForumRepo interface {
	Create(ctx context.Context, forum *models.ForumCreate) (*models.Forum, error)
	Get(ctx context.Context, slug string) (*models.Forum, error)
	GetUsers(ctx context.Context, getSettings *models.GetForumUsers) (*[]models.User, error)                    // получение пользователей форума
	GetThreads(ctx context.Context, slug string, getSettings *models.GetForumThreads) (*[]models.Thread, error) // получение веток обсуждения форума
}

type PostRepo interface {
	Get(ctx context.Context, id int64, getSettings *models.PostGetRequest) (*models.PostGetResult, error)
	Update(ctx context.Context, id int64, updateDate *models.PostUpdate) (*models.Post, error)
	Create(ctx context.Context, threadSlugOrId string, posts *[]models.PostCreate) (*[]models.Post, error) // создание постов для ветки. created у post'ов должен быть одинаковый
}

type ThreadRepo interface {
	Create(ctx context.Context, forumSlug string, thread *models.ThreadCreate) (*models.Thread, error) // создаст ветку в нужном форуме
	Get(ctx context.Context, threadSlugOrId string) (*models.Thread, error)
	Update(ctx context.Context, threadSlugOrId string, updateData *models.ThreadUpdate) (*models.Thread, error)
	GetPosts(ctx context.Context, slugOrId string, getSettings *models.ThreadPostRequest) (*[]models.Post, error) // все сообщения данной ветки
}

type VoteRepo interface {
	Create(ctx context.Context, threadSlugOrId string, vote *models.VoteCreate) (*models.Thread, error) // пользователь должен учитываться только один раз
}

type ServiceRepo interface {
	GetInfo(ctx context.Context) (*models.Service, error)
	Clear(ctx context.Context) error
}
