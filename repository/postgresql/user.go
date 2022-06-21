package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	CreateUserCommand               = "INSERT INTO Users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)"
	UpdateUserCommand               = "UPDATE Users SET (fullname, about, email) = ($1, $2, $3) WHERE nickname = $4"
	GetUserByNicknameCommand        = "SELECT nickname, fullname, about, email FROM Users WHERE nickname = $1"
	GetUserByEmailCommand           = "SELECT nickname, fullname, about, email FROM Users WHERE email = $1"
	GetUserByNicknameOrEmailCommand = "SELECT nickname, fullname, about, email FROM Users WHERE nickname = $1 OR email = $2"
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

func (a *UserPostgresRepo) getUserByNicknameOrEmail(nickname string, email string) (*[]models.User, error) {
	rows, err := a.Db.Query(context.Background(), GetUserByNicknameOrEmailCommand, nickname, email)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func (a *UserPostgresRepo) Create(user *models.User) (*[]models.User, error) {
	_, err := a.Db.Exec(context.Background(), CreateUserCommand, user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		checkAlreadyExist, err := a.getUserByNicknameOrEmail(user.Nickname, user.Email)
		if err == nil && len(*checkAlreadyExist) > 0 {
			return checkAlreadyExist, ErrorUserAlreadyExist
		} else {
			return nil, ErrorUserAlreadyExist
		}
	}

	userToReturn := make([]models.User, 0, 1)
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
