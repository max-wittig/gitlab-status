package main

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/xanzy/go-gitlab"
)

var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func getNextCronTime(cron string, now time.Time) (time.Time, error) {
	cronTime, err := cronParser.Parse(cron)
	if err != nil {
		return time.Time{}, errors.New("Invalid cron syntax")
	}
	return cronTime.Next(now), nil
}

// sortedStatus sorts status after nearest occurence
func sortedStatus(now time.Time, statusCrons *Config) ([]string, error) {
	cronTimeMap := make(map[time.Time]string)
	var keys []time.Time
	for k := range statusCrons.Crons {
		nextCronTime, err := getNextCronTime(k, now)
		if err != nil {
			return nil, err
		}
		cronTimeMap[nextCronTime] = k
		keys = append(keys, nextCronTime)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })
	var sortedKeys []string
	for _, k := range keys {
		sortedKeys = append(sortedKeys, cronTimeMap[k])
	}
	return sortedKeys, nil
}

// findCurrentShouldStatus returns the status, as if it would have been running before
func findCurrentShouldStatus(statusConfig *Config, now time.Time) (*Status, error) {
	sortedCrons, err := sortedStatus(now, statusConfig)
	if err != nil {
		return nil, err
	}
	var foundConfig Status
	if len(sortedCrons) >= 2 {
		foundConfig = statusConfig.Crons[sortedCrons[len(sortedCrons)-1]]
	} else {
		foundConfig = statusConfig.Crons[sortedCrons[0]]
	}

	return &foundConfig, nil
}

func getNow(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	//set timezone,
	return time.Now().In(loc), nil
}

// UpdateGitlabStatus updates the status on Gitlab, based on the Options
func UpdateGitlabStatus(appOptions *appOptions, statusConfig *Config) error {
	client := gitlab.NewClient(nil, appOptions.gitlabToken)
	client.SetBaseURL(appOptions.gitlabURL)
	now, err := getNow(statusConfig.Timezone)
	if err != nil {
		return err
	}
	currentStatus, err := findCurrentShouldStatus(statusConfig, now)
	if err != nil {
		return err
	}
	return SendUpdateRequest(client, currentStatus)
}

// SendUpdateRequest uses the GitLab API to update your status
func SendUpdateRequest(client *gitlab.Client, statusOption *Status) error {
	userOptions := gitlab.UserStatusOptions{Emoji: &statusOption.Emoji, Message: &statusOption.Message}
	status, _, err := client.Users.SetUserStatus(&userOptions)
	if err == nil {
		fmt.Printf("Updated status to: %s | %s\n", statusOption.Emoji, status.Message)
	}
	return err
}
