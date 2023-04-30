package main

import (
	"app/db"
	_ "app/docs"
	"app/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title		EventHandler API
// @version	1.0.0
// @host		localhost:3000
// @description    An event management service API in Go using Gin framework.
// @contact.name  Marek Beck
func main() {
	app := gin.New()
	routes.InitApp(app)
	db.Init()
	server := &http.Server{
		Addr:    ":3000",
		Handler: app,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic("failed to start gin server")
	}
}
