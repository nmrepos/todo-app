package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	todos  = make([]Todo, 0)
	nextID = 1
	mu     sync.Mutex
)

func main() {
	r := gin.Default()
	r.Use(corsMiddleware)

	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)

	// serve all your JS/CSS/etc. under /static/*
    r.Static("/static", "./static")

    // when someone hits “/”, return your SPA’s index.html
    r.GET("/", func(c *gin.Context) {
        c.File("./static/index.html")
    })

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func getTodos(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, todos)
}

func addTodo(c *gin.Context) {
	var t Todo
	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	t.ID = nextID
	nextID++
	todos = append(todos, t)
	mu.Unlock()
	c.JSON(http.StatusCreated, t)
}

func updateTodo(c *gin.Context) {
	idParam := c.Param("id")
	var t Todo
	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i := range todos {
		if fmt.Sprintf("%d", todos[i].ID) == idParam {
			todos[i].Title = t.Title
			todos[i].Done = t.Done
			c.JSON(http.StatusOK, todos[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func deleteTodo(c *gin.Context) {
	idParam := c.Param("id")
	mu.Lock()
	defer mu.Unlock()
	for i := range todos {
		if fmt.Sprintf("%d", todos[i].ID) == idParam {
			todos = append(todos[:i], todos[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}
