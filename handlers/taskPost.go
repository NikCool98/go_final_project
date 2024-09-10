package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NikCool98/go_final_project/config"
	"github.com/NikCool98/go_final_project/stor"
)

func TaskPostHandler(store stor.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var t config.Task
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(res, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		id, err := store.CreateTask(t)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		response := config.Response{ID: id}

		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}