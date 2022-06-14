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
	GetUsers(getSettings *models.GetForumUsers) (*[]models.User, error)                    // получение пользователей форума
	GetThreads(slug string, getSettings *models.GetForumThreads) (*[]models.Thread, error) // получение веток обсуждения форума
}

type PostRepo interface {
	Get(id int64, getSettings *models.PostGetRequest) (*models.PostGetResult, error)
	Update(id int64, updateDate *models.PostUpdate) (*models.Post, error)
	Create(threadSlugOrId string, posts *[]models.PostCreate) (*[]models.Post, error) // создание постов для ветки. created у post'ов должен быть одинаковый
}

type ThreadRepo interface {
	Create(forumSlug string, thread *models.ThreadCreate) (*models.Thread, error) // создаст ветку в нужном форуме
	Get(threadSlugOrId string) (*models.Thread, error)
	Update(threadSlugOrId string, updateData *models.ThreadUpdate) (*models.Thread, error)
	GetPosts(slugOrId string, getSettings *models.ThreadPostRequest) (*[]models.Post, error) // все сообщения данной ветки
}

type VoteRepo interface {
	Create(threadSlugOrId string, vote *models.VoteCreate) (*models.Thread, error) // пользователь должен учитываться только один раз
}

type ServiceRepo interface {
	GetInfo() (*models.Service, error)
	Clear() error
}
