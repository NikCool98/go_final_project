package main

import (
	"github.com/NikCool98/go_final_project/config"
	"github.com/NikCool98/go_final_project/handlers"
	"github.com/NikCool98/go_final_project/stor"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
)

func main() {
	// Проверка наличия БД
	dataBase := stor.OpenDb()
	defer dataBase.Close()
	store := stor.NewStore(dataBase)

	fileSrv := http.FileServer(http.Dir(config.WebDir))
	http.Handle("/", fileSrv)
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("POST /api/task", handlers.TaskPostHandler(store))
	http.HandleFunc("/api/tasks", handlers.TasksGetHandler(store))
	http.HandleFunc("GET /api/task", handlers.TaskGetHandler(store))
	http.HandleFunc("PUT /api/task", handlers.TaskPutHandler(store))
	http.HandleFunc("DELETE /api/task", handlers.TaskDelHandler(store))
	http.HandleFunc("/api/task/done", handlers.TaskDoneHandler(store))

	log.Printf("Starting server on port: %s", config.DefaultPort)
	err := http.ListenAndServe(":"+config.DefaultPort, nil)
	if err != nil {
		log.Fatalf("Server run error: %v", err)
	}
}
