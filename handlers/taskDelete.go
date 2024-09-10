package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NikCool98/go_final_project/config"
	"github.com/NikCool98/go_final_project/stor"
)

func TaskDelHandler(store stor.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		err := store.DeleteTask(id)
		if err != nil {
			err := errors.New("Задача с таким id не найдена")
			config.ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(config.ErrorResponse)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{}); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
