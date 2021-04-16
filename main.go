package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// Task struct is task's info
type Task struct {
	ID      string `json:"id" form:"id" binding:"-"`
	Content string `json:"content" form:"content" binding:"required"`
	Done    bool   `json:"done" form:"done"`
}

var tasks = make(map[string]Task)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/tasks", func(c *gin.Context) {
		list := make([]Task, 0, len(tasks))
		for _, task := range tasks {
			list = append(list, task)
		}
		c.JSON(http.StatusOK, list)
	})

	r.GET("/task/:id", func(c *gin.Context) {
		id := c.Param("id")
		task, ok := tasks[id]
		if ok {
			c.JSON(http.StatusOK, task)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	r.POST("/task", func(c *gin.Context) {
		var task Task
		if err := c.ShouldBind(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			guid := xid.New()
			task.ID = guid.String()
			tasks[task.ID] = task
			c.JSON(http.StatusOK, task)
		}
	})

	r.PUT("/task/:id", func(c *gin.Context) {
		var task Task
		if err := c.ShouldBind(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			id := c.Param("id")
			_, ok := tasks[id]
			if ok {
				task.ID = id
				tasks[id] = task
				c.JSON(http.StatusOK, task)
			} else {
				c.Status(http.StatusNotFound)
			}
		}
	})

	r.DELETE("/task/:id", func(c *gin.Context) {
		id := c.Param("id")
		task, ok := tasks[id]
		if ok {
			delete(tasks, id)
			c.JSON(http.StatusOK, task)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run()
}
