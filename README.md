# gitlab-status

Automate your GitLab status

## Usage

```
  -config string
        Path to the config.yml file
  -daemon
        Run as daemon (default true)
  -version
        Version
```

### Daemon mode

Daemon mode will keep the tool alive and run the cronjobs at the specified times.

Non-daemon mode can be used for something like GitLab scheduled jobs, so that you don't need your own server.

You also need to set the following environment variables

* GITLAB_URL -> optional, defaults to `https://gitlab.com`
* GITLAB_TOKEN -> GitLab token with `api` scope

### Example usage

GITLAB_TOKEN=your-token gitlab-status -config config.yml

See the [config.example.yml](./config.example.yml) file for example configuration
