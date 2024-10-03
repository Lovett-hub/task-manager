package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` 
}

var tasks []Task

// Handle CRUD operations
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if len(r.URL.Path) == len("/tasks") { 
			getTasks(w, r)
		} else { 
			getTaskByID(w, r)
		}
	case http.MethodPost:
		createTask(w, r)
	case http.MethodPut:
		updateTask(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}


func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = len(tasks) + 1
	task.Status = "pending"
	tasks = append(tasks, task)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}


func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}


func getTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):] 
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.NotFound(w, r) 
}


func updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):] 
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = updatedTask.Title
			tasks[i].Description = updatedTask.Description
			tasks[i].Status = updatedTask.Status
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	http.NotFound(w, r) 
}


func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):] 
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...) 
			w.WriteHeader(http.StatusNoContent)       
			return
		}
	}

	http.NotFound(w, r) 
}

func main() {
	http.HandleFunc("/tasks", taskHandler)  
	http.HandleFunc("/tasks/", taskHandler) 

	log.Println("Starting server on :8000...")
	log.Fatal(http.ListenAndServe(":8000", nil)) 
}
