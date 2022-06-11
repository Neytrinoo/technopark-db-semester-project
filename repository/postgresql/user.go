package postgresql

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	CreateUserCommand        = "INSERT INTO Users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4) RETURNING id"
	UpdateUserCommand        = "UPDATE Users SET (fullname, about, email) = ($1, $2, $3) WHERE nickname = $4"
	GetUserByNicknameCommand = "SELECT nickname, fullname, about, email FROM Users WHERE nickname = $1"
	GetUserByEmailCommand    = "SELECT nickname, fullname, about, email FROM Users WHERE email = $1"
)

var (
	ErrorUserAlreadyExist   = errors.New("user already exist")
	ErrorUserDoesNotExist   = errors.New("user does not exist")
	ErrorConflictUpdateUser = errors.New("data conflicts with existing users")
)

type UserPostgresRepo struct {
	Db *sqlx.DB
}

func NewUserPostgresRepo(db *sqlx.DB) domain.UserRepo {
	return &UserPostgresRepo{Db: db}
}

func (a *UserPostgresRepo) Create(user *models.User) (*models.User, error) {
	_, err := a.Db.Exec(CreateUserCommand, user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		userAlreadyExist, err := a.Get(user.Nickname)
		if err != nil {
			userAlreadyExist, _ = a.Get(user.Email)
		}
		return userAlreadyExist, ErrorUserAlreadyExist
	}

	return user, nil
}

func (a *UserPostgresRepo) Update(nickname string, updateData *models.UserUpdate) (*models.User, error) {
	user, err := a.Get(nickname)
	if err != nil {
		return nil, ErrorUserDoesNotExist
	}

	if updateData.Fullname == "" {
		updateData.Fullname = user.Fullname
	} else {
		user.Fullname = updateData.Fullname
	}
	if updateData.About == "" {
		updateData.About = user.About
	} else {
		user.About = updateData.About
	}
	if updateData.Email == "" {
		updateData.Email = user.Email
	} else {
		user.Email = updateData.Email
	}

	_, err = a.Db.Exec(UpdateUserCommand, updateData.Fullname, updateData.About, updateData.Email, nickname)

	if err != nil {
		return nil, ErrorConflictUpdateUser
	}

	return user, nil
}

func (a *UserPostgresRepo) Get(nicknameOrEmail string) (*models.User, error) {
	var user models.User
	err := a.Db.Get(&user, GetUserByNicknameCommand, nicknameOrEmail)
	if err != nil {
		err = a.Db.Get(&user, GetUserByEmailCommand, nicknameOrEmail)
	}

	if err != nil {
		return nil, ErrorUserDoesNotExist
	}

	return &user, nil
}
