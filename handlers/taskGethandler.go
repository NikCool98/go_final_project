package handlers

import (
	"encoding/json"
	"errors"
	"github.com/NikCool98/go_final_project/config"
	"github.com/NikCool98/go_final_project/stor"
	"net/http"
)

func TaskGetHandler(store stor.Stor) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		//var t configs.Task
		id := req.URL.Query().Get("id")
		task, err := store.GetTask(id)
		if err != nil {
			err := errors.New("Задача с таким id не найдена")
			config.ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(config.ErrorResponse)
			return
		}
		// Возвращаем ответ
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(task); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
