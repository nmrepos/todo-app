package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// helper to reset state and router
func setup() *gin.Engine {
	// reset global todos
	todos = make([]Todo, 0)
	nextID = 1

	// create router with only our feature handlers
	r := gin.Default()
	r.Use(corsMiddleware)
	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	return r
}

func TestGetTodos_Empty(t *testing.T) {
	r := setup()
	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /todos status = %d; want %d", w.Code, http.StatusOK)
	}
	var got []Todo
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Expected empty list; got %v", got)
	}
}

func TestAddTodo_Success(t *testing.T) {
	r := setup()
	payload := map[string]string{"title": "Buy milk"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /todos status = %d; want %d", w.Code, http.StatusCreated)
	}
	var todo Todo
	if err := json.Unmarshal(w.Body.Bytes(), &todo); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if todo.ID != 1 || todo.Title != "Buy milk" {
		t.Errorf("Unexpected todo = %+v; want ID=1 Title=Buy milk", todo)
	}
}

func TestGetTodos_AfterAdd(t *testing.T) {
	r := setup()
	// add one todo
	addReq := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(`{"title":"Task1"}`))
	addReq.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, addReq)

	// now get
	req := httptest.NewRequest("GET", "/todos", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)

	if w2.Code != http.StatusOK {
		t.Fatalf("GET /todos status = %d; want %d", w2.Code, http.StatusOK)
	}
	var list []Todo
	if err := json.Unmarshal(w2.Body.Bytes(), &list); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if len(list) != 1 || list[0].Title != "Task1" {
		t.Errorf("Unexpected list = %v; want one Task1", list)
	}
}

func TestAddTodo_BadRequest(t *testing.T) {
	r := setup()
	// invalid JSON
	req := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(`{"title":}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("POST /todos bad JSON status = %d; want %d", w.Code, http.StatusBadRequest)
	}
}
