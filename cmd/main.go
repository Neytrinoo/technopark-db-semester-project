package main

import (
	_ "github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"log"
	"technopark-db-semester-project/system"
	"time"
)

func addTimeLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		now := time.Now()
		result := next(c)
		duration := time.Now().Sub(now)

		log.Println("Request:", c.Request().RequestURI, "time:", duration)

		return result
	}
}

func main() {
	e := echo.New()

	db := system.InitDb()
	defer db.Close()
	userRepo, forumRepo, threadRepo, postRepo, voteRepo, serviceRepo := system.InitRepos(db)
	userHandler, forumHandler, threadHandler, postHandler, voteHandler, serviceHandler := system.InitHandlers(userRepo, forumRepo, threadRepo, postRepo, voteRepo, serviceRepo)

	// api routes

	e.POST("api/forum/create", forumHandler.Create)
	e.GET("api/forum/:slug/details", forumHandler.Get)
	e.POST("api/forum/:slug/create", threadHandler.Create)
	e.GET("api/forum/:slug/users", forumHandler.GetUsers)
	e.GET("api/forum/:slug/threads", forumHandler.GetThreads)
	e.GET("api/post/:id/details", postHandler.Get)
	e.POST("api/post/:id/details", postHandler.Update)

	e.POST("api/thread/:slug_or_id/create", postHandler.Create)
	e.GET("api/thread/:slug_or_id/details", threadHandler.Get)
	e.POST("api/thread/:slug_or_id/details", threadHandler.Update)
	e.GET("api/thread/:slug_or_id/posts", threadHandler.GetPosts)
	e.POST("api/thread/:slug_or_id/vote", voteHandler.Create)
	e.POST("api/user/:nickname/create", userHandler.Create)
	e.GET("api/user/:nickname/profile", userHandler.Get)
	e.POST("api/user/:nickname/profile", userHandler.Update)

	e.GET("api/service/status", serviceHandler.GetInfo)
	e.POST("api/service/clear", serviceHandler.Clear)

	if err := e.Start("0.0.0.0:5000"); err != nil {
		log.Fatalln("server error:", err)
	}
}
