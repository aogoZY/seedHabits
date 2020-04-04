package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetStartDayAndEndDayByMonth(date string) (StartDate string, EndDate string, err error) {
	index := strings.Index(date, "-")
	year := date[:index]
	month := date[index+1:]
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)

	firstOfMonth := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	lastofMonthDay := lastOfMonth.Format("2006-01-02")
	fmt.Printf("lastofMonthDay:%s", lastofMonthDay)
	EndDate = lastofMonthDay + " 23:59:59"
	StartDate = date + " 00:00:00"
	fmt.Printf("StartDate:%s\n", StartDate)
	fmt.Printf("EndDate:%s\n", EndDate)
	return StartDate, EndDate, nil

}
