package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// afterNow сравнивает даты, не учитывая время
func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return date.After(now)
}

// NextDate рассчитывает ближайшую дату следующего выполнения задачи
// по заданному правилу повторения.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}
	// Условие для повтора события каждый год
	if repeat == "y" {
		for {
			date = date.AddDate(1, 0, 0)

			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	}
	// Условие для повтора события через указанное количество дней
	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Split(repeat, " ")

		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat format")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("invalid day interval")
		}

		for {
			date = date.AddDate(0, 0, interval)

			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	}
	// Повтор задачи в указанные дни недели
	if strings.HasPrefix(repeat, "w ") {
		parts := strings.Split(repeat, " ")

		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat format")
		}

		weekdays := strings.Split(parts[1], ",")

		var days [8]bool

		for _, weekday := range weekdays {
			day, err := strconv.Atoi(weekday)
			if err != nil {
				return "", err
			}

			if day < 1 || day > 7 {
				return "", fmt.Errorf("invalid weekday")
			}

			days[day] = true
		}

		for {
			date = date.AddDate(0, 0, 1)

			weekday := int(date.Weekday())

			if weekday == 0 {
				weekday = 7
			}

			if days[weekday] && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	}
	// Повтор задачи в указанные дни месяца и, если нужно, в месяцы года.
	if strings.HasPrefix(repeat, "m ") {
		parts := strings.Split(repeat, " ")

		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid repeat format")
		}

		monthDays := strings.Split(parts[1], ",")

		var days [32]bool
		var lastDay bool
		var beforeLastDay bool

		for _, monthDay := range monthDays {
			day, err := strconv.Atoi(monthDay)
			if err != nil {
				return "", err
			}

			if day == -1 {
				lastDay = true
				continue
			}

			if day == -2 {
				beforeLastDay = true
				continue
			}

			if day < 1 || day > 31 {
				return "", fmt.Errorf("invalid month day")
			}

			days[day] = true
		}

		var months [13]bool
		useMonths := false

		// Проверяем условие, если указаны конкретные месяцы
		if len(parts) == 3 {
			useMonths = true
			monthList := strings.Split(parts[2], ",")

			for _, monthString := range monthList {
				month, err := strconv.Atoi(monthString)
				if err != nil {
					return "", err
				}

				if month < 1 || month > 12 {
					return "", fmt.Errorf("invalid month")
				}

				months[month] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			month := int(date.Month())

			// Если пользователь задал конкретный месяц, пропускаем остальные
			if useMonths && !months[month] {
				continue
			}

			day := date.Day()

			if days[day] && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}

			// Определим последний день текущего месяца
			nextMonth := time.Date(
				date.Year(),
				date.Month()+1,
				1,
				0, 0, 0, 0,
				date.Location(),
			)

			lastDayOfMonth := nextMonth.AddDate(0, 0, -1).Day()

			if lastDay && day == lastDayOfMonth && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}

			if beforeLastDay && day == lastDayOfMonth-1 && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	}

	return "", fmt.Errorf("unsupported repeat format")
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowString := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	// Если now не указан, используем текущую дату
	if nowString == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
