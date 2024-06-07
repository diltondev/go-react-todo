package router

import (
	"go-react-todo/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", middleware.GetAllTasks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", middleware.CreateTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/task/{id}", middleware.CompleteTask).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/undoTask/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/task/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/deleteAllTasks", middleware.DeleteAllTasks).Methods("DELETE", "OPTIONS")
	return router
}
