package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func DaysInMonth(t time.Time) int {
	firstOfNextMonth := t.AddDate(0, 1, -t.Day()+1)
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfMonth.Day()
}

func NextDate(now time.Time, dstart string, repeat string) (result string, err error) {
	// некорректное правило
	if repeat != "" && repeat != "y" && !strings.HasPrefix(repeat, "d") && !strings.HasPrefix(repeat, "w") &&
		!strings.HasPrefix(repeat, "m") {
		return "", errors.New("invalid repeat")
	}

	//нет переноса
	if repeat == "" {
		return "", nil
	}

	curDeadLine, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid start date")
	}

	//перенос на год
	if repeat == "y" {
		for {
			curDeadLine = curDeadLine.AddDate(1, 0, 0)
			if curDeadLine.After(now) {
				break
			}
		}
		result = curDeadLine.Format(dateFormat)
	}

	repeats := strings.Split(repeat, " ")

	//перенос на заданное кол-во дней
	if repeats[0] == "d" {
		if len(repeats) != 2 {
			return "", errors.New("invalid repeat days format")
		}
		days, err := strconv.Atoi(repeats[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid days count")
		}
		for {
			curDeadLine = curDeadLine.AddDate(0, 0, days)
			if curDeadLine.After(now) {
				break
			}
		}
		result = curDeadLine.Format(dateFormat)
	}

	//перенос на заданные дни недели
	if repeats[0] == "w" {
		if len(repeats) != 2 {
			return "", errors.New("invalid repeat weekdays format")
		}
		var (
			weekdaysString []string = strings.Split(repeats[1], ",")
			weekdays       []int
		)
		//проверяем валидность
		if len(weekdaysString) > 7 || len(weekdaysString) < 1 {
			return "", errors.New("invalid weekdays count")
		}
		for i := 0; i < len(weekdaysString); i++ {
			weekDay, err := strconv.Atoi(weekdaysString[i])
			if err != nil || weekDay < 1 || weekDay > 7 {
				return "", errors.New("invalid weekday")
			}
			if weekDay == 7 {
				weekDay = 0 // формат Weekday пакета time
			}
			weekdays = append(weekdays, weekDay)
		}
		//если все нормально рассчитываем новую дату
		curWeekDay := int(curDeadLine.Weekday())
		candidate := curDeadLine.AddDate(0, 0, 1)
		count := 0
	labelWeekDays:
		for {
			curWeekDay = (curWeekDay + 1) % 7
			count++
			for _, v := range weekdays {
				candidate = curDeadLine.AddDate(0, 0, count)
				if curWeekDay == v && candidate.After(now) {
					break labelWeekDays
				}
			}
		}
		result = candidate.Format(dateFormat)
	}

	// перенос на заданный день месяца
	if repeats[0] == "m" {
		if len(repeats) != 3 && len(repeats) != 2 {
			return "", errors.New("invalid repeat monthdays format")
		}

		var (
			days   []int
			months []int

			//	count  int
			//	curRes time.Time
		)

		// парсим дни, они есть всегда
		daysStrings := strings.Split(repeats[1], ",")

		if len(daysStrings) < 1 || len(daysStrings) > 31 {
			return "", errors.New("invalid days count")
		}
		for _, v := range daysStrings {
			day, err := strconv.Atoi(v)
			if err != nil || day < -2 || day > 31 || day == 0 {
				return "", errors.New("invalid days format")
			}
			days = append(days, day)
		}

		// парсим месяцы, если они есть
		if len(repeats) == 3 {
			monthsStrings := strings.Split(repeats[2], ",")

			if len(monthsStrings) > 12 {
				return "", errors.New("invalid months or days count")
			}
			for _, v := range monthsStrings {
				month, err := strconv.Atoi(v)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("invalid months format")
				}
				months = append(months, month)
			}
		}

		candidate := curDeadLine.AddDate(0, 0, 1)
		if candidate.Before(now) {
			candidate = now.AddDate(0, 0, 1)
		}

	labelMonthDays:
		for {
			candidateDay := candidate.Day()
			candidateMonth := int(candidate.Month())
			daysInCandidateMonth := DaysInMonth(candidate)

			monthMatches := len(months) == 0 //если месяцы не даны, подходит любой
			if len(months) > 0 {
				for _, month := range months {
					if month == candidateMonth {
						monthMatches = true
						break
					}
				}
			}

			if monthMatches {
				for _, day := range days {
					var actualDay int
					if day > 0 {
						actualDay = day
					} else {
						actualDay = daysInCandidateMonth + day + 1 // для отрицательных дней
					}
					if actualDay <= daysInCandidateMonth && actualDay == candidateDay {
						result = candidate.Format(dateFormat)
						break labelMonthDays
					}
				}
			}
			candidate = candidate.AddDate(0, 0, 1)
		}
	}
	return result, nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(dateFormat, r.FormValue("now"))
	if err != nil {
		now = time.Now()
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
