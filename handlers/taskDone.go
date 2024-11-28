package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NikCool98/go_final_project/stor"
)

func TaskDoneHandler(store stor.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		err := store.TaskDone(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{}); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
