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
	r.POST("/register", controllers.RegisterUser)
	// login
	r.POST("/login", controllers.Login)
	// jwt
	// ✅ Grupa z autoryzacją JWT
	a := r.Group("/account")
	a.Use(auth.AuthHandler) // Middleware do autoryzacji JWT
	{
		a.GET("/GetAccount", controllers.GetAccount)
		a.PUT("/UpdateAccount", controllers.UpdateAccount)
		a.DELETE("/DeleteAccount", controllers.DeleteAccount)
	}
	t := r.Group("/tasks")
	t.Use(auth.AuthHandler) // Middleware do autoryzacji JWT
	{
		t.POST("/add", controllers.AddTask)
		t.GET("/get", controllers.GetTasks)
		t.PUT("/update/:id", controllers.UpdateTask)
	}
	// ading task
	r.Run(":8080")
}
