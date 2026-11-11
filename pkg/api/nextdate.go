package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat       = "20060102"
	maxRepeatDays    = 400
	minRepeatDays    = 1
	maxWeekdays      = 7
	minWeekdays      = 1
	maxMonthsAmount  = 12
	minMonthAmount   = 1
	searchLimitYears = 10
)

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func searchLimit(date, now time.Time) time.Time {
	base := now
	if afterNow(date, now) {
		base = date
	}
	return base.AddDate(searchLimitYears, 0, 0)
}

func nextByDays(repeat string, date, now time.Time) (time.Time, error) {
	days, err := strconv.Atoi(repeat)
	if err != nil {
		return time.Time{}, err
	}

	if days < minRepeatDays || days > maxRepeatDays {
		return time.Time{}, errors.New("Days count out of range")
	}

	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}

	return date, nil
}

func nextByWeekdays(repeat []string, date, now time.Time) (time.Time, error) {
	var dayMatches [7]bool

	for _, repeatString := range repeat {
		weekday, err := strconv.Atoi(repeatString)
		if err != nil {
			return time.Time{}, err
		}

		if weekday < minWeekdays || weekday > maxWeekdays {
			return time.Time{}, errors.New("weekdays argument is out of range")
		}

		dayMatches[weekday%7] = true
	}

	for {
		date = date.AddDate(0, 0, 1)

		if afterNow(date, now) && dayMatches[int(date.Weekday())] {
			break
		}
	}

	return date, nil
}

func nextByMonthDays(repeatDays, repeatMonths []string, date, now time.Time) (time.Time, error) {
	var allowedMonthDays [32]bool
	var monthsMatches [13]bool
	var lastDay bool
	var secondToLast bool

	if len(repeatMonths) == 0 {
		for i := range monthsMatches {
			monthsMatches[i] = true
		}
	} else {
		for _, v := range repeatMonths {
			month, err := strconv.Atoi(v)
			if err != nil {
				return time.Time{}, err
			}

			if month < minMonthAmount || month > maxMonthsAmount {
				return time.Time{}, errors.New("amount of months out of range")
			}

			monthsMatches[month] = true
		}
	}

	for _, repeatStr := range repeatDays {
		day, err := strconv.Atoi(repeatStr)
		if err != nil {
			return time.Time{}, err
		}

		if day == 0 || day < -2 || day > 31 {
			return time.Time{}, errors.New("month days out of range")
		}

		switch day {
		case -2:
			secondToLast = true
		case -1:
			lastDay = true
		default:
			allowedMonthDays[day] = true
		}
	}

	limit := searchLimit(date, now)

	for {
		date = date.AddDate(0, 0, 1)
		if date.After(limit) {
			return time.Time{}, errors.New("no matching date found")
		}
		dayMatches := allowedMonthDays[date.Day()]

		if lastDay && date.AddDate(0, 0, 1).Day() == 1 {
			dayMatches = true
		}
		if secondToLast && date.AddDate(0, 0, 2).Day() == 1 {
			dayMatches = true
		}

		if afterNow(date, now) && monthsMatches[int(date.Month())] && dayMatches {
			break
		}
	}

	return date, nil
}

func nextYear(date, now time.Time) time.Time {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date
}

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat argument.")
	}
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("time parse error: %w", err)
	}
	repeatString := strings.Split(repeat, " ")

	switch repeatString[0] {
	case "d":

		if len(repeatString) != 2 {
			return "", errors.New("repeat argument has wrong count of value")
		}

		date, err = nextByDays(repeatString[1], date, now)
		if err != nil {
			return "", err
		}

	case "w":
		if len(repeatString) != 2 {
			return "", errors.New("repeat argument has wrong count of value")
		}

		weekString := strings.Split(repeatString[1], ",")

		date, err = nextByWeekdays(weekString, date, now)

		if err != nil {
			return "", err
		}

	case "m":
		if len(repeatString) < 2 || len(repeatString) > 3 {
			return "", errors.New("repeat argument has wrong count of value")
		}
		var months []string
		var days []string

		days = strings.Split(repeatString[1], ",")
		if len(repeatString) == 3 {
			months = strings.Split(repeatString[2], ",")
		}

		date, err = nextByMonthDays(days, months, date, now)
		if err != nil {
			return "", err
		}

	case "y":
		if len(repeatString) > 1 {
			return "", errors.New("too many arguments")
		}
		date = nextYear(date, now)

	default:
		return "", errors.New("wrong type of option")
	}

	return date.Format(dateFormat), nil
}
