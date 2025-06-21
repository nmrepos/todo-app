package main

import (
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

    // Only listing & adding:
    r.GET("/todos", getTodos)
    r.POST("/todos", addTodo)

    // Serve UI
    r.Static("/static", "./static")
    r.GET("/", func(c *gin.Context) {
        c.File("./static/index.html")
    })

    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

// CORS for our simple UI
func corsMiddleware(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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
    var t struct {
        Title string `json:"title"`
    }
    if err := c.BindJSON(&t); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    mu.Lock()
    defer mu.Unlock()
    todo := Todo{ID: nextID, Title: t.Title}
    nextID++
    todos = append(todos, todo)
    c.JSON(http.StatusCreated, todo)
}
