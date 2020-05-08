package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type TwitchChannel struct {
	ID string `json:"_id"`

	BroadcasterLanguage          string `json:"broadcaster_language"`
	BroadcasterType              string `json:"broadcaster_type"`
	CreatedAt                    string `json:"created_at"`
	DisplayName                  string `json:"display_name"`
	Email                        string `json:"email"`
	Followers                    int64  `json:"followers"`
	Game                         string `json:"game"`
	Language                     string `json:"language"`
	Logo                         string `json:"logo"`
	Mature                       bool   `json:"mature"`
	Name                         string `json:"name"`
	Partner                      bool   `json:"partner"`
	ProfileBanner                string `json:"profile_banner"`
	ProfileBannerBackgroundColor string `json:"profile_banner_background_color"`
	Status                       string `json:"status"`
	StreamKey                    string `json:"stream_key"`
	UpdatedAt                    string `json:"updated_at"`
	URL                          string `json:"url"`
	VideoBanner                  string `json:"video_banner"`
	Views                        int64  `json:"views"`
}

func getChannel(u *User) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/channel", nil)
	if err != nil {
		log.WithError(err).Error("Unable to create http request to get twitch channel data")
		return err
	}

	req.Header.Set("Client-ID", settings.ClientID)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Authorization", "OAuth "+u.Token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("Unable to get twitch channel data")
		return err
	}

	defer resp.Body.Close()
	c := &TwitchChannel{}

	if err := json.NewDecoder(resp.Body).Decode(c); err != nil {
		log.WithError(err).Error("Unable to parse twitch channel data")
		return err
	}

	u.TwitchChannel = c

	return nil
}

func (c *TwitchChannel) SaveFiles() {
	data, err := fieldsToMap(c)
	if err != nil {
		log.WithError(err).Error("Unable to convert channel to map")
		return
	}

	for k, v := range data {
		saveContent("channel", k, v)
	}
}
