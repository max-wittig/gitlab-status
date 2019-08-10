package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/max-wittig/gitlab-status/version"
)

type appOptions struct {
	configPath  string
	gitlabURL   string
	gitlabToken string
	daemon      bool
}

func addEnvironmentVariables(appOptions *appOptions) error {
	gitlabURL := os.Getenv("GITLAB_URL")
	if gitlabURL == "" {
		gitlabURL = "https://gitlab.com"
	}
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabURL == "" || gitlabToken == "" {
		return errors.New("Missing GitLab environment variables")
	}
	appOptions.gitlabURL = gitlabURL
	appOptions.gitlabToken = gitlabToken
	return nil
}

func parseOptions() (*appOptions, error) {
	versionFlag := flag.Bool("version", false, "Version")
	configFlag := flag.String("config", "", "Path to the config.yml file")
	daemonFlag := flag.Bool("daemon", true, "Run as daemon")
	flag.Parse()
	aConfig := appOptions{}

	if *versionFlag {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
		return nil, nil
	}

	if *configFlag != "" {
		// Check if path exists
		if _, err := os.Stat(*configFlag); !os.IsNotExist(err) {
			aConfig.configPath = *configFlag
		}
	}
	aConfig.daemon = *daemonFlag

	addEnvironmentVariables(&aConfig)
	return &aConfig, nil
}

func main() {
	appOptions, err := parseOptions()
	if err != nil {
		log.Fatalln("Config file doesn't exist")
	}
	statusConfig, err := ReadConfig(appOptions.configPath)
	if err != nil {
		log.Fatalln("Could not read config")
	}
	if appOptions.daemon {
		ScheduleDaemons(appOptions, statusConfig)
	} else {
		UpdateGitlabStatus(appOptions, statusConfig)
	}
}
