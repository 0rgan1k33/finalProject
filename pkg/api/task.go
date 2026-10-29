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
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendError(w, "оишбка чтения запроса")
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		sendError(w, "оишбка чтения JSON ")
		return
	}
	if task.Title == "" {
		sendError(w, "Не указан заголовок")
		return
	}
	err = checkDate(&task)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		sendError(w, "ошибка добавления задачи")
		return
	}
	idStr := strconv.Itoa(int(id))
	answer, err := json.Marshal(map[string]string{"id": idStr})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		sendError(w, "не указан id")
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, "задача не найдена")
		return
	}
	answer, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendError(w, "ошибка чтения запроса")
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		sendError(w, "ошибка JSON")
		return
	}
	if task.ID == "" {
		sendError(w, "Не указан ID")
		return
	}
	if task.Title == "" {
		sendError(w, "Не указан заголовок")
		return
	}
	err = checkDate(&task)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	answer, err := json.Marshal(map[string]string{})
	if err != nil {
		sendError(w, "ошибка ответа")
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

func sendError(w http.ResponseWriter, msg string) {
	answer, err := json.Marshal(map[string]string{"error": msg})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	if task.ID == "" {
		sendError(w, "Не указан ID")
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(task)
		if err != nil {
			sendError(w, err.Error())
			return
		}
	} else {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		err = db.UpdateDate(nextDate, task)
		if err != nil {
			sendError(w, err.Error())
			return
		}
	}
	answer, err := json.Marshal(map[string]string{})
	if err != nil {
		sendError(w, "ошибка ответа")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	if task.ID == "" {
		sendError(w, "Не указан ID")
		return
	}
	err = db.DeleteTask(task)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	answer, err := json.Marshal(map[string]string{})
	if err != nil {
		sendError(w, "ошибка ответа")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}
