package main

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/xanzy/go-gitlab"
)

func getNextCronTime(cronParser *cron.Parser, cron string, now time.Time) (time.Time, error) {
	cronTime, err := cronParser.Parse(cron)
	if err != nil {
		return time.Time{}, errors.New("Invalid cron syntax")
	}
	return cronTime.Next(now), nil
}

func sortedStatus(cronParser *cron.Parser, now time.Time, statusCrons *Config) ([]string, error) {
	cronTimeMap := make(map[time.Time]string)
	var keys []time.Time
	for k := range statusCrons.Crons {
		nextCronTime, err := getNextCronTime(cronParser, k, now)
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

func findCurrentShouldStatus(statusConfig *Config) (*StatusOptions, error) {
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	now := time.Now()
	sortedCrons, err := sortedStatus(&cronParser, now, statusConfig)
	if err != nil {
		return nil, err
	}
	for _, cronKey := range sortedCrons {
		cronNext, err := getNextCronTime(&cronParser, cronKey, now)
		if err != nil {
			return nil, err
		}
		tolerance := time.Duration(30 * time.Minute) // 30 minutes
		if cronNext.Add(-tolerance).Sub(now) < 0 {
			// it's in the future
			currentCron := statusConfig.Crons[cronKey]
			return &currentCron, nil
		}
	}
	// nothing found
	return nil, errors.New("No cron match found")
}

// UpdateGitlabStatus updates the status on Gitlab, based on the Options
func UpdateGitlabStatus(appOptions *appOptions, statusConfig *Config) error {
	client := gitlab.NewClient(nil, appOptions.gitlabToken)
	client.SetBaseURL(appOptions.gitlabURL)
	currentStatus, err := findCurrentShouldStatus(statusConfig)
	if err != nil {
		return err
	}
	if currentStatus == nil {
		currentStatus = &statusConfig.Default
	}
	err = SendUpdateRequest(client, currentStatus)
	return err
}

// SendUpdateRequest uses the GitLab API to update your status
func SendUpdateRequest(client *gitlab.Client, statusOptions *StatusOptions) error {
	userOptions := gitlab.UserStatusOptions{Emoji: &statusOptions.Emoji, Message: &statusOptions.Message}
	status, _, err := client.Users.SetUserStatus(&userOptions)
	fmt.Printf("Updated status to: %s | %s\n", statusOptions.Emoji, status.Message)
	return err
}
