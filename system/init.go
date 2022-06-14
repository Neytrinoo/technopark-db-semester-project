package system

import (
	"context"
	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"log"
	"technopark-db-semester-project/delivery"
	"technopark-db-semester-project/domain"
	"technopark-db-semester-project/repository/postgresql"
)

func InitDb() *pgxpool.Pool {
	/*dsn := "user=root dbname=postgres password=rootpassword host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping() // вот тут будет первое подключение к базе*/
	/*db, err := sqlx.Connect("pgx", "postgres://root:rootpassword@localhost:5432/technopark-dbms")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln("ping error:", err)
	}
	db.SetMaxOpenConns(10)*/
	dbPool, err := pgxpool.Connect(context.Background(), "postgres://root:rootpassword@localhost:5432/technopark-dbms")
	if err != nil {
		log.Fatalln("conn error:", err)
	}

	return dbPool
}

func InitRepos(db *sqlx.DB) (domain.UserRepo, domain.ForumRepo, domain.ThreadRepo, domain.PostRepo, domain.VoteRepo, domain.ServiceRepo) {
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
