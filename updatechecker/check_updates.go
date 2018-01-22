package updatechecker

import (
	"context"
	"fmt"

	"github.com/juju/loggo"

	"github.com/google/go-github/github"
	version "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

var Log loggo.Logger

func init() {
	Log = loggo.GetLogger("com.github.okelet.goutils.updatecheck")
}

type UpdateDetectedListener interface {
	OnNewVersionDetecetd(newVersion string)
}

type CheckUpdatesThread struct {
	Cron              *cron.Cron
	Listeners         []UpdateDetectedListener
	IntervalInSeconds int
	GitHubUser        string
	GitHubProject     string
	CurrentVersion    string
}

func NewCheckUpdatesThread(intervalInSeconds int, gitHubUser string, gitHubProject string, currentVersion string) *CheckUpdatesThread {
	t := CheckUpdatesThread{
		Cron:              nil,
		Listeners:         []UpdateDetectedListener{},
		IntervalInSeconds: intervalInSeconds,
		GitHubUser:        gitHubUser,
		GitHubProject:     gitHubProject,
		CurrentVersion:    currentVersion,
	}
	return &t
}

func (t *CheckUpdatesThread) AddListener(listener UpdateDetectedListener) {
	t.Listeners = append(t.Listeners, listener)
}

func (t *CheckUpdatesThread) Check() {

	context := context.Background()
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context, t.GitHubUser, t.GitHubProject)
	if err != nil {
		Log.Errorf("Error getting repository: %v", err)
		return
	}
	Log.Infof("Version from github: %v; local version: %v.", *release.TagName, t.CurrentVersion)

	notify := false
	if t.CurrentVersion == "master" {
		notify = true
	} else {
		currentVersion, err := version.NewVersion(t.CurrentVersion)
		if err != nil {
			Log.Errorf("Error parsing current version %v: %v", t.CurrentVersion, err)
			return
		}
		newVersion, err := version.NewVersion(*release.TagName)
		if err != nil {
			Log.Errorf("Error parsing release version %v: %v", *release.TagName, err)
			return
		}
		if newVersion.GreaterThan(currentVersion) {
			notify = true
		}
	}

	if notify {
		Log.Infof("New version detected; current: %v; new version: %v.", t.CurrentVersion, *release.TagName)
		for _, l := range t.Listeners {
			l.OnNewVersionDetecetd(*release.TagName)
		}
	}

}

func (t *CheckUpdatesThread) SetInterval(seconds int) error {
	t.IntervalInSeconds = seconds
	if t.Cron != nil {
		t.Stop()
		return t.Start()
	}
	return nil
}

func (t *CheckUpdatesThread) Start() error {
	var err error
	if t.Cron == nil {
		t.Cron = cron.New()
		err = t.Cron.AddFunc(fmt.Sprintf("@every %vs", t.IntervalInSeconds), t.Check)
		if err != nil {
			// TODO: translate/i18n
			return errors.Wrap(err, "Error starting cron")
		}
		t.Cron.Start()
	}
	return nil
}

func (t *CheckUpdatesThread) Stop() {
	if t.Cron != nil {
		t.Cron.Stop()
		t.Cron = nil
	}
}
