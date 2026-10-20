package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var tasks = []Task{}
var nextID = 1

func main() {
	router := gin.Default()

	// CREATE
	router.POST("/tasks", createTask)

	// READ
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByID)

	// UPDATE
	router.PUT("/tasks/:id", updateTask)

	// DELETE
	router.DELETE("/tasks/:id", deleteTask)

	router.Run(":8080")
}

func createTask(c *gin.Context) {
	var task Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	task.ID = nextID
	nextID++

	if task.Status == "" {
		task.Status = "pending"
	}

	tasks = append(tasks, task)

	c.JSON(http.StatusCreated, task)
}

func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func getTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			c.JSON(http.StatusOK, task)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func updateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Title = updatedTask.Title
			tasks[i].Status = updatedTask.Status

			c.JSON(http.StatusOK, tasks[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			c.JSON(http.StatusOK, gin.H{
				"message": "Task deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}
