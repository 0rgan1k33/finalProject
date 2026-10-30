package api

import (
	"bytes"
	"encoding/json"
	"finalProject/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendError(w, "оишбка чтения запроса", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		sendError(w, "оишбка чтения JSON", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		sendError(w, "Не указан заголовок", http.StatusBadRequest)
		return
	}
	err = checkDate(&task)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		sendError(w, "ошибка добавления задачи", http.StatusInternalServerError)
		return
	}
	idStr := strconv.Itoa(int(id))
	answer, err := json.Marshal(map[string]string{"id": idStr})
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		sendError(w, "не указан id", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, "задача не найдена", http.StatusBadRequest)
		return
	}
	answer, err := json.Marshal(task)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendError(w, "ошибка чтения запроса", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		sendError(w, "ошибка JSON", http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		sendError(w, "Не указан ID", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		sendError(w, "Не указан заголовок", http.StatusBadRequest)
		return
	}
	err = checkDate(&task)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	answer, err := json.Marshal(map[string]string{})
	if err != nil {
		sendError(w, "ошибка ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}
	taskDate, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	var nextDate string
	if task.Repeat != "" {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if taskDate.Format(dateFormat) < now.Format(dateFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = nextDate
		}
	}
	return nil
}

func sendError(w http.ResponseWriter, msg string, code int) {
	answer, err := json.Marshal(map[string]string{"error": msg})
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if task.ID == "" {
		sendError(w, "Не указан ID", http.StatusBadRequest)
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(task)
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.UpdateDate(nextDate, task)
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	answer, err := json.Marshal(map[string]string{})
	if err != nil {
		sendError(w, "ошибка ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if task.ID == "" {
		sendError(w, "Не указан ID", http.StatusBadRequest)
		return
	}
	err = db.DeleteTask(task)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	answer, err := json.Marshal(map[string]string{})
	if err != nil {
		sendError(w, "ошибка ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}
