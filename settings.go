package main

import (
	"bytes"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
	"gopkg.in/yaml.v2"
)

var (
	settings          *Settings
	twitchOauthConfig *oauth2.Config
)

type Settings struct {
	ClientID          string        `yaml:"client_id"`
	ClientSecret      string        `yaml:"client_secret"`
	MaxFollowers      int           `yaml:"max_followers"`
	MaxSubs           int           `yaml:"max_subs"`
	RedirectURL       string        `yaml:"redirect_url"`
	UpdateInterval    time.Duration `yaml:"update_interval"`
	VerificationToken string        `yaml:"verification_token"`
	WebserverPort     string        `yaml:"webserver_port"`
}

func loadSettings() {
	data, err := ioutil.ReadFile(cfg.SettingsFile)
	if err != nil {
		log.WithError(err).Errorf("Unable to read %s", cfg.SettingsFile)
		return
	}

	s := &Settings{}
	b := bytes.NewBuffer(data)
	if err := yaml.NewDecoder(b).Decode(s); err != nil {
		log.WithError(err).Errorf("Unable to decode %s", cfg.SettingsFile)
		return
	}

	settings = s

	twitchOauthConfig = &oauth2.Config{
		RedirectURL:  settings.RedirectURL,
		ClientID:     settings.ClientID,
		ClientSecret: settings.ClientSecret,
		Scopes:       []string{"channel:read:subscriptions", "user:read:broadcast", "chat:read", "chat:edit", "channel_read", "channel_editor", "channel_subscriptions", "channel:moderate", "bits:read", "channel:read:redemptions"},
		Endpoint:     twitch.Endpoint,
	}

}

func settingsUpdater() {
	loadSettings()

	go func() {
		c := time.Tick(time.Minute)
		for range c {
			loadSettings()
		}
	}()
}
