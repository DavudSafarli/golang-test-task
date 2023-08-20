package main

import (
	"net/http"
	"twitch_chat_analysis/adapters/storage"
	"twitch_chat_analysis/api/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	redisConnStr := "redis:6379"
	redis := storage.NewRedis(redisConnStr)
	reportingController := controller.NewReportingController(redis)

	r.GET("/message/list", reportingController.Get)

	srv := http.Server{
		Addr:    ":3000",
		Handler: r,
	}
	srv.ListenAndServe()
}
