package cmd

import (
	"github.com/labstack/echo/v4"
	"log"
	"technopark-db-semester-project/system"
)

func main() {
	e := echo.New()

	db := system.InitDb()
	userRepo, forumRepo, threadRepo, postRepo, voteRepo := system.InitRepos(db)
	userHandler, forumHandler, threadHandler, postHandler, voteHandler := system.InitHandlers(userRepo, forumRepo, threadRepo, postRepo, voteRepo)

	// api routes

	e.POST("/forum/create", forumHandler.Create)
	e.GET("/forum/:slug/details", forumHandler.Get)
	e.POST("/forum/:slug/create", threadHandler.Create)
	e.GET("/forum/:slug/users", forumHandler.GetUsers)
	e.GET("/forum/:slug/threads", forumHandler.GetThreads)
	e.GET("/post/:id/details", postHandler.Get)
	e.POST("/post/:id/details", postHandler.Update)
	// TODO: /service/clear + /service/status
	e.POST("/thread/:slug_or_id/create", postHandler.Create)
	e.GET("/thread/:slug_or_id/details", threadHandler.Get)
	e.POST("/thread/:slug_or_id/details", threadHandler.Update)
	e.GET("/thread/:slug_or_id/posts", threadHandler.GetPosts)
	e.POST("/thread/:slug_or_id/vote", voteHandler.Create)
	e.POST("/user/:nickname/create", userHandler.Create)
	e.GET("/user/:nickname/profile", userHandler.Get)
	e.POST("/user/:nickname/profile", userHandler.Update)

	if err := e.Start("0.0.0.0:5000"); err != nil {
		log.Fatalln("server error:", err)
	}
}
