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
	var nextDate time.Time
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}
	if repeat == "y" {
		for {
			nextDate = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return nextDate.Format(dateFormat), nil
	} else {
		ruleSlice := strings.Split(repeat, " ")
		if len(ruleSlice) != 2 {
			return "", errors.New("Неврное условие повторения")
		}
		rule, interval := ruleSlice[0], ruleSlice[1]
		if rule != "d" {
			return "", errors.New("Неверное условие повторения")
		}
		intervalInt, err := strconv.Atoi(interval)
		if err != nil {
			return "", err
		}
		if 0 >= intervalInt && intervalInt > 400 {
			return "", errors.New("Количество дней от 1 до 400")
		}
		for {
			nextDate = date.AddDate(0, 0, intervalInt)
			if afterNow(date, now) {
				break
			}
		}
		return nextDate.Format(dateFormat), nil
	}
}

func afterNow(date, now time.Time) bool {
	if date.After(now) {
		return true
	}
	return false
}

func nextDayHandler(w http.ResponseWriter, r http.Request) {

}
