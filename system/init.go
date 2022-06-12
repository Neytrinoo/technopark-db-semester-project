package system

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
	"technopark-db-semester-project/delivery"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/repository/postgresql"
)

func InitDb() *sqlx.DB {
	dsn := "user=postgres dbname=postgres password=admin host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)

	return sqlx.NewDb(db, "pgx")
}

func InitRepos(db *sqlx.DB) (domain.UserRepo, domain.ForumRepo, domain.ThreadRepo, domain.PostRepo, domain.VoteRepo) {
	userRepo := postgresql.NewUserPostgresRepo(db)
	forumRepo := postgresql.NewForumPostgresRepo(db)
	threadRepo := postgresql.NewThreadPostgresRepo(db)
	postRepo := postgresql.NewPostPostgresRepo(db)
	voteRepo := postgresql.NewVotePostgresRepo(db)

	return userRepo, forumRepo, threadRepo, postRepo, voteRepo
}

func InitHandlers(userRepo domain.UserRepo, forumRepo domain.ForumRepo, threadRepo domain.ThreadRepo, postRepo domain.PostRepo, voteRepo domain.VoteRepo) (delivery.UserHandler, delivery.ForumHandler, delivery.ThreadHandler, delivery.PostHandler, delivery.VoteHandler) {
	userHandler := delivery.MakeUserHandler(userRepo)
	forumHandler := delivery.MakeForumHandler(forumRepo)
	threadHandler := delivery.MakeThreadHandler(threadRepo)
	postHandler := delivery.MakePostHandler(postRepo)
	voteHandler := delivery.MakeVoteHandler(voteRepo)

	return userHandler, forumHandler, threadHandler, postHandler, voteHandler
}
