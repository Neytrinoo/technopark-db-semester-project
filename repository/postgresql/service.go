package postgresql

import (
	"github.com/jmoiron/sqlx"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	DeleteTablesCommand    = "TRUNCATE TABLE Forums, Posts, Threads, Users, Votes CASCADE"
	GetCountRecordsCommand = "SELECT (SELECT count(*) FROM Users), (SELECT count(*) FROM Forums), (SELECT count(*) FROM Threads), (SELECT count(*) FROM Posts)"
)

type ServicePostgresRepo struct {
	Db *sqlx.DB
}

func NewServicePostgresRepo(db *sqlx.DB) domain.ServiceRepo {
	return &ServicePostgresRepo{Db: db}
}

func (a *ServicePostgresRepo) GetInfo() (*models.Service, error) {
	var result models.Service
	_ = a.Db.QueryRow(GetCountRecordsCommand).Scan(&result.User, &result.Forum, &result.Thread, &result.Post)

	return &result, nil
}

func (a *ServicePostgresRepo) Clear() error {
	_, _ = a.Db.Exec(DeleteTablesCommand)

	return nil
}
