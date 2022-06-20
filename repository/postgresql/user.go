package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	CreateUserCommand        = "INSERT INTO Users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)"
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
	Db *pgxpool.Pool
}

func NewUserPostgresRepo(db *pgxpool.Pool) domain.UserRepo {
	return &UserPostgresRepo{Db: db}
}

func (a *UserPostgresRepo) Create(user *models.User) (*[]models.User, error) {
	userToReturn := make([]models.User, 0, 2)
	_, err := a.Db.Exec(context.Background(), CreateUserCommand, user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		var firstNickname string

		userAlreadyExist, err := a.Get(user.Nickname)

		if err == nil {
			firstNickname = userAlreadyExist.Nickname
			userToReturn = append(userToReturn, *userAlreadyExist)
		}
		userAlreadyExist, err = a.Get(user.Email)
		if err == nil && userAlreadyExist.Nickname != firstNickname {
			userToReturn = append(userToReturn, *userAlreadyExist)
		}

		return &userToReturn, ErrorUserAlreadyExist
	}

	userToReturn = append(userToReturn, *user)
	return &userToReturn, nil
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

	_, err = a.Db.Exec(context.Background(), UpdateUserCommand, updateData.Fullname, updateData.About, updateData.Email, nickname)

	if err != nil {
		return nil, ErrorConflictUpdateUser
	}

	return user, nil
}

func (a *UserPostgresRepo) Get(nicknameOrEmail string) (*models.User, error) {
	var user models.User
	err := a.Db.QueryRow(context.Background(), GetUserByNicknameCommand, nicknameOrEmail).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		err = a.Db.QueryRow(context.Background(), GetUserByEmailCommand, nicknameOrEmail).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	}

	if err != nil {
		return nil, ErrorUserDoesNotExist
	}

	return &user, nil
}
