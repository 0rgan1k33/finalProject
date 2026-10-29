package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило не указано")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("некорректная дата")
	}

	ruleSlice := strings.Split(repeat, " ")
	rule := ruleSlice[0]
	switch rule {
	case "y":
		if len(ruleSlice) != 1 {
			return "", errors.New("неверный формат правила Y")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "d":
		if len(ruleSlice) != 2 {
			return "", errors.New("неврный формат правила D")
		}
		interval, err := strconv.Atoi(ruleSlice[1])
		if err != nil {
			return "", err
		}
		if interval < 1 || interval >= 400 {
			return "", errors.New("кол-во дней должно быть от 1 до 400")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	default:
		return "", errors.New("неврный формат правила")
	}
}

func afterNow(date, now time.Time) bool {
	if date.After(now) {
		return true
	}
	return false
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowTime, err := time.Parse(dateFormat, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	next, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
