package domain

import "technopark-db-semester-project/domain/models"

type UserRepo interface {
	Create(user *models.User) (*models.User, error)
	Update(nickname string, updateData *models.UserUpdate) (*models.User, error)
	Get(nicknameOrEmail string) (*models.User, error)
}

type ForumRepo interface {
	Create(forum *models.ForumCreate) (*models.Forum, error)
	Get(slug string) (*models.Forum, error)
	GetUsers(slug string) (*[]models.User, error)     // получение пользователей форума
	GetThreads(slug string) (*[]models.Thread, error) // получение веток обсуждения форума
}

type PostRepo interface {
	Get(id int) (*models.Post, error)
	Update(id int, updateDate *models.Post) error
	Create(threadSlugOrId string, posts *[]models.Post) error // создание постов для ветки. created у post'ов должен быть одинаковый
}

type ThreadRepo interface {
	Create(forumSlug string, thread *models.ThreadCreate) (*models.Thread, error) // создаст ветку в нужном форуме
	Get(threadSlugOrId string) (*models.Thread, error)
	Update(threadSlugOrId string, updateData *models.ThreadUpdate) (*models.Thread, error)
	GetPosts(getSettings *models.ThreadPostRequest) (*[]models.Post, error) // все сообщения данной ветки
}

type VoteRepo interface {
	Create(threadSlugOrId string, vote *models.Vote) error // пользователь должен учитываться только один раз
}
