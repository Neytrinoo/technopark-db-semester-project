package main

import (
	"context"
	"fmt"
	"github.com/fasthttp/router"
	_ "github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"log"
	"technopark-db-semester-project/system"
	"time"
)

func addTimeLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("Start. Request:", c.Request().RequestURI)
		now := time.Now()
		result := next(c)
		duration := time.Now().Sub(now)

		log.Println("End. Request:", c.Request().RequestURI, "time:", duration, "Status:", c.Response().Status)

		return result
	}
}

func main() {
	//router := echo.New()

	db := system.InitDb()
	defer db.Close()
	userRepo, forumRepo, threadRepo, postRepo, voteRepo, serviceRepo := system.InitRepos(db)
	userHandler, forumHandler, threadHandler, postHandler, voteHandler, serviceHandler := system.InitHandlers(userRepo, forumRepo, threadRepo, postRepo, voteRepo, serviceRepo)
	fasthttpRouter := router.New()

	// api routes

	fasthttpRouter.POST("/api/forum/create", forumHandler.Create)
	fasthttpRouter.GET("/api/forum/{slug}/details", forumHandler.Get)
	fasthttpRouter.POST("/api/forum/{slug}/create", threadHandler.Create)
	fasthttpRouter.GET("/api/forum/{slug}/users", forumHandler.GetUsers)
	fasthttpRouter.GET("/api/forum/{slug}/threads", forumHandler.GetThreads)
	fasthttpRouter.GET("/api/post/{id}/details", postHandler.Get)
	fasthttpRouter.POST("/api/post/{id}/details", postHandler.Update)

	fasthttpRouter.POST("/api/thread/{slug_or_id}/create", postHandler.Create)
	fasthttpRouter.GET("/api/thread/{slug_or_id}/details", threadHandler.Get)
	fasthttpRouter.POST("/api/thread/{slug_or_id}/details", threadHandler.Update)
	fasthttpRouter.GET("/api/thread/{slug_or_id}/posts", threadHandler.GetPosts)
	fasthttpRouter.POST("/api/thread/{slug_or_id}/vote", voteHandler.Create)
	fasthttpRouter.POST("/api/user/{nickname}/create", userHandler.Create)
	fasthttpRouter.GET("/api/user/{nickname}/profile", userHandler.Get)
	fasthttpRouter.POST("/api/user/{nickname}/profile", userHandler.Update)

	fasthttpRouter.GET("/api/service/status", serviceHandler.GetInfo)
	fasthttpRouter.POST("/api/service/clear", serviceHandler.Clear)

	ctx := context.Background()

	err := fasthttp.ListenAndServe(
		"0.0.0.0:5000",
		func(fasthttpCtx *fasthttp.RequestCtx) {
			fasthttpCtx.SetUserValue("ctx", ctx)
			fasthttpRouter.Handler(fasthttpCtx)
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	/*
		if err := router.Start("0.0.0.0:5000"); err != nil {
			log.Fatalln("server error:", err)
		}*/
}
