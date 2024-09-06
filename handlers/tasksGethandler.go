package handlers

import (
	"encoding/json"
	"errors"
	"github.com/NikCool98/go_final_project/config"
	"github.com/NikCool98/go_final_project/stor"
	"net/http"
)

func TasksGetHandler(store stor.Stor) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		searchParams := req.URL.Query().Get("search")
		tasks, err := store.GetTasks(searchParams)
		if err != nil {
			if err != nil {
				err := errors.New("Ошибка запроса к базе данных")
				config.ErrorResponse.Error = err.Error()
				json.NewEncoder(res).Encode(config.ErrorResponse)
				return
			}
		}
		response := map[string][]config.Task{
			"tasks": tasks,
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
