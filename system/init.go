package system

import (
	"context"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"technopark-db-semester-project/delivery"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/repository/postgresql"
)

func InitDb() *pgxpool.Pool {
	dbPool, err := pgxpool.Connect(context.Background(), "postgres://root:admin@localhost:5432/forum_db")
	if err != nil {
		log.Fatalln("conn error:", err)
	}

	return dbPool
}

func InitRepos(db *pgxpool.Pool) (domain.UserRepo, domain.ForumRepo, domain.ThreadRepo, domain.PostRepo, domain.VoteRepo, domain.ServiceRepo) {
	userRepo := postgresql.NewUserPostgresRepo(db)
	forumRepo := postgresql.NewForumPostgresRepo(db)
	threadRepo := postgresql.NewThreadPostgresRepo(db)
	postRepo := postgresql.NewPostPostgresRepo(db)
	voteRepo := postgresql.NewVotePostgresRepo(db)
	serviceRepo := postgresql.NewServicePostgresRepo(db)

	return userRepo, forumRepo, threadRepo, postRepo, voteRepo, serviceRepo
}

func InitHandlers(userRepo domain.UserRepo, forumRepo domain.ForumRepo, threadRepo domain.ThreadRepo, postRepo domain.PostRepo, voteRepo domain.VoteRepo, serviceRepo domain.ServiceRepo) (delivery.UserHandler, delivery.ForumHandler, delivery.ThreadHandler, delivery.PostHandler, delivery.VoteHandler, delivery.ServiceHandler) {
	userHandler := delivery.MakeUserHandler(userRepo)
	forumHandler := delivery.MakeForumHandler(forumRepo)
	threadHandler := delivery.MakeThreadHandler(threadRepo)
	postHandler := delivery.MakePostHandler(postRepo)
	voteHandler := delivery.MakeVoteHandler(voteRepo)
	serviceHandler := delivery.MakeServiceHandler(serviceRepo)

	return userHandler, forumHandler, threadHandler, postHandler, voteHandler, serviceHandler
}
