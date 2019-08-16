package main

import (
	"errors"
	"time"
)

func getNextCronTime(cron string, now time.Time) (time.Time, error) {
	cronTime, err := cronParser.Parse(cron)
	if err != nil {
		return time.Time{}, errors.New("Invalid cron syntax")
	}
	return cronTime.Next(now), nil
}

func getNow(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	//set timezone,
	return time.Now().In(loc), nil
}
