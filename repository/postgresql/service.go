package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/domain/models"
)

const (
	DeleteTablesCommand    = "TRUNCATE TABLE Users, Forums, Threads, Posts, ForumUsers, Votes CASCADE;"
	GetCountRecordsCommand = "SELECT (SELECT count(*) FROM Users), (SELECT count(*) FROM Forums), (SELECT count(*) FROM Threads), (SELECT count(*) FROM Posts);"
)

type ServicePostgresRepo struct {
	Db *pgxpool.Pool
}

func NewServicePostgresRepo(db *pgxpool.Pool) domain.ServiceRepo {
	return &ServicePostgresRepo{Db: db}
}

func (a *ServicePostgresRepo) GetInfo(ctx context.Context) (*models.Service, error) {
	var result models.Service
	_ = a.Db.QueryRow(ctx, GetCountRecordsCommand).Scan(&result.User, &result.Forum, &result.Thread, &result.Post)

	return &result, nil
}

func (a *ServicePostgresRepo) Clear(ctx context.Context) error {
	_, _ = a.Db.Exec(ctx, DeleteTablesCommand)

	return nil
}
