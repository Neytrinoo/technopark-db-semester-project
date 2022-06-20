package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	DeleteTablesCommand    = "TRUNCATE TABLE Forums, Posts, Threads, Users, Votes CASCADE"
	GetCountRecordsCommand = "SELECT (SELECT count(*) FROM Users), (SELECT count(*) FROM Forums), (SELECT count(*) FROM Threads), (SELECT count(*) FROM Posts)"
)

type ServicePostgresRepo struct {
	Db *pgxpool.Pool
}

func NewServicePostgresRepo(db *pgxpool.Pool) domain.ServiceRepo {
	return &ServicePostgresRepo{Db: db}
}

func (a *ServicePostgresRepo) GetInfo() (*models.Service, error) {
	var result models.Service
	_ = a.Db.QueryRow(context.Background(), GetCountRecordsCommand).Scan(&result.User, &result.Forum, &result.Thread, &result.Post)

	return &result, nil
}

func (a *ServicePostgresRepo) Clear() error {
	_ = a.Db.QueryRow(context.Background(), DeleteTablesCommand)

	return nil
}
