package api

import (
	"encoding/json"
	"finalProject/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	tasks, err := db.Tasks(50)
	if err != nil {
		sendError(w, "ошибка получения задачи", http.StatusInternalServerError)
		return
	}
	response := TasksResp{Tasks: tasks}
	answer, err := json.Marshal(response)
	if err != nil {
		sendError(w, "ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}
