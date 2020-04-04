package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	//now := time.Now()
	//currentYear, currentMonth, _ := now.Date()
	//currentLocation := now.Location()
	//fmt.Println(currentLocation)
	timeStr := "2020-03"
	index := strings.Index(timeStr, "-")
	year:= timeStr[:index]
	month := timeStr[index+1:]
	yearInt,_ := strconv.Atoi(year)
	monthInt,_ := strconv.Atoi(month)
	fmt.Println(yearInt,monthInt)

	firstOfMonth := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	fmt.Printf("firstOfMonth:%v\n",firstOfMonth)
	fmt.Printf("lastOfMonth:%v\n",lastOfMonth)
}