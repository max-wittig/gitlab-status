package main

import (
	"log"

	"time"

	cron "github.com/robfig/cron/v3"
	"github.com/xanzy/go-gitlab"
)

// ScheduleDaemons schedules the status crons
func ScheduleDaemons(appOptions *appOptions, statusConfig *Config) {
	location, err := time.LoadLocation(statusConfig.Timezone)
	if err != nil {
		log.Fatalln("Invalid timezone")
	}
	cronScheduler := cron.New(cron.WithLocation(location))
	cronScheduler.Location()
	client := gitlab.NewClient(nil, appOptions.gitlabToken)
	client.SetBaseURL(appOptions.gitlabURL)
	for k, v := range statusConfig.Crons {
		cronString := k
		status := v
		_, err := cronScheduler.AddFunc(cronString, func() { SendUpdateRequest(client, &status) })
		if err != nil {
			log.Fatalln(err)
		}
	}
	for _, entry := range cronScheduler.Entries() {
		log.Println(entry.Schedule.Next(time.Now()))
	}
	cronScheduler.Run()
}
