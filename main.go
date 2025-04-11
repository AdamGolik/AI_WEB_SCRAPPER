package main

import (
	"AI_WEB_SCRAPPER/auth"
	"AI_WEB_SCRAPPER/controllers"
	"AI_WEB_SCRAPPER/initlizers"

	"github.com/gin-gonic/gin"
)

// func that will execute before main like loading files or so
func init() {
	// load env
	initlizers.LoadEnv()
	// connect to the db
	initlizers.ConnectDB()
	// create table
	initlizers.CreateTables()
}

func main() {
	r := gin.Default()
	// tests:w

	// register
	r.POST("/register", controllers.Register)
	// login
	r.POST("/login", controllers.Login)
	// jwt
	// ✅ Grupa z autoryzacją JWT
	web := r.Group("/web")
	web.Use(auth.AuthHandler)
	{
		// 🔍 AI analiza
		web.POST("/analyze", controllers.Analyze)
	}
	// ading task
	r.Run(":8080")
}
