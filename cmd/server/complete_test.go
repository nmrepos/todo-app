package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// setup router including update endpoint
func setupComplete() *gin.Engine {
	// reset state
	todos = make([]Todo, 0)
	nextID = 1

	r := gin.Default()
	r.Use(corsMiddleware)
	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	// include update but skip delete
	r.PUT("/todos/:id", updateTodo)
	return r
}

func TestCompleteTodo_Success(t *testing.T) {
	r := setupComplete()
	// add a task
	addReq := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(`{"title":"Task1"}`))
	addReq.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, addReq)

	// complete it
	compReq := httptest.NewRequest("PUT", "/todos/1", bytes.NewBufferString(`{"title":"Task1","done":true}`))
	compReq.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, compReq)

	if w2.Code != http.StatusOK {
		t.Fatalf("PUT /todos/1 status = %d; want %d", w2.Code, http.StatusOK)
	}
	var todo Todo
	if err := json.Unmarshal(w2.Body.Bytes(), &todo); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if todo.ID != 1 || !todo.Done {
		t.Errorf("Expected Todo ID=1 Done=true; got %+v", todo)
	}

	// verify via GET
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
	if len(list) != 1 || !list[0].Done {
		t.Errorf("Expected list[0].Done=true; got %v", list)
	}
}

func TestCompleteTodo_BadRequest(t *testing.T) {
	r := setupComplete()
	// invalid JSON payload
	req := httptest.NewRequest("PUT", "/todos/1", bytes.NewBufferString(`{"done":}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("PUT /todos/1 bad JSON status = %d; want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCompleteTodo_NotFound(t *testing.T) {
	r := setupComplete()
	// complete non-existent
	req := httptest.NewRequest("PUT", "/todos/99", bytes.NewBufferString(`{"title":"X","done":true}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("PUT /todos/99 status = %d; want %d", w.Code, http.StatusNotFound)
	}
}
