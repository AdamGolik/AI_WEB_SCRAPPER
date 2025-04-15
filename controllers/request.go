package controllers

import (
	"AI_WEB_SCRAPPER/initlizers"
	"AI_WEB_SCRAPPER/models"

	"github.com/gin-gonic/gin"
)

func IsValidStatus(s models.Status) bool {
	switch s {
	case models.Done, models.InProgress, models.Todo:
		return true
	default:
		return false
	}
}
func AddTask(c *gin.Context) {
	userId := c.MustGet("userId")
	var body struct {
		Title  string        `json:"title"`
		Body   string        `json:"body"`
		Status models.Status `json:"status"`
	}
	if !IsValidStatus(body.Status) {
		c.JSON(400, gin.H{"error": "Invalid status"})
		return
	}
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	var task models.Todos
	initlizers.DB.Create(&task)
	task.UserID = userId.(uint)
	task.Title = body.Title
	task.Body = body.Body
	task.Status = body.Status
	if err := initlizers.DB.Save(&task).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(200, gin.H{"task": task})
}
func GetTasks(c *gin.Context) {
	userId := c.MustGet("userId")
	var tasks []models.Todos

	if initlizers.DB.Where("user_id = ?", userId).Find(&tasks).Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(200, gin.H{"tasks": tasks})
}
