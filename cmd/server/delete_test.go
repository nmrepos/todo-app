package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// setup router including delete endpoint
func setupDelete() *gin.Engine {
	// reset state
	todos = make([]Todo, 0)
	nextID = 1

	r := gin.Default()
	r.Use(corsMiddleware)
	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)
	return r
}

func TestDeleteTodo_Success(t *testing.T) {
	r := setupDelete()
	// add a task
	addReq := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(`{"title":"Task1"}`))
	addReq.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, addReq)

	// delete it
	delReq := httptest.NewRequest("DELETE", "/todos/1", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, delReq)

	if w2.Code != http.StatusNoContent {
		t.Fatalf("DELETE /todos/1 status = %d; want %d", w2.Code, http.StatusNoContent)
	}

	// verify it's gone
	getReq := httptest.NewRequest("GET", "/todos", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, getReq)

	if w3.Code != http.StatusOK {
		t.Fatalf("GET /todos status = %d; want %d", w3.Code, http.StatusOK)
	}
	var list []Todo
	if err := json.Unmarshal(w3.Body.Bytes(), &list); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("Expected empty list after delete; got %v", list)
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	r := setupDelete()
	// delete a non-existent task
	req := httptest.NewRequest("DELETE", "/todos/99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("DELETE /todos/99 status = %d; want %d", w.Code, http.StatusNotFound)
	}
}
