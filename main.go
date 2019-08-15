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
		} else {
			return nil, errors.New("Config file doesn't exist")
		}
	}
	aConfig.daemon = *daemonFlag

	err := addEnvironmentVariables(&aConfig)
	if err != nil {
		return nil, err
	}
	return &aConfig, nil
}

func main() {
	appOptions, err := parseOptions()
	if err != nil {
		log.Fatalln(err)
	}
	statusConfig, err := ReadConfig(appOptions.configPath)
	if err != nil {
		log.Fatalln("Could not read config")
	}
	if appOptions.daemon {
		err = ScheduleDaemons(appOptions, statusConfig)
	} else {
		err = UpdateGitlabStatus(appOptions, statusConfig)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
