package taskrepeater

import (
	"fmt"
	"github.com/NikCool98/go_final_project/config"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	//проверка на пустой repeat
	if repeat == "" {
		return "", fmt.Errorf("не указан repeat")
	}

	startDate, err := time.Parse(config.DateFormat, date)
	//проверка на неверный формат даты
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}
	//разделяем правило на тип правила и его аргумент
	ruleSplited := strings.Split(repeat, " ")
	//тип правила
	ruleType := ruleSplited[0]

	switch ruleType {
	case "d":
		if len(ruleSplited) < 2 {
			return "", fmt.Errorf("не указано количество дней")
		}
		//количество дней для переноса задачи
		daysToMove, err := strconv.Atoi(ruleSplited[1])

		if err != nil {
			return "", err
		}
		if daysToMove > 400 {

			return "", fmt.Errorf("количество дней не должно превышать 400")
		}
		newDate := startDate.AddDate(0, 0, daysToMove)
		//проверяем, что возвращаемая дата не меньше текущей, если меньше - сдвигаем на указанное количество дней
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysToMove)
		}
		return newDate.Format(config.DateFormat), nil

	case "y":
		newDate := startDate.AddDate(1, 0, 0)
		//проверяем, что возвращаемая дата не меньше текущей, если меньше - сдвигаем еще на год
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(config.DateFormat), nil

	default:
		return "", fmt.Errorf("некорректный тип правила")
	}
}
