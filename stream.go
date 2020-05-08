package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type TwitchStreams struct {
	Data       []*TwitchStream   `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

type TwitchStream struct {
	GameID       string    `json:"game_id"`
	ID           string    `json:"id"`
	Language     string    `json:"language"`
	StartedAt    time.Time `json:"started_at"`
	TagIDs       []string  `json:"tag_ids"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	ViewerCount  int       `json:"viewer_count"`
}

func getStreams(u *User) {
	result := &TwitchStreams{}

	after := ""
	for {
		client := twitchOauthConfig.Client(context.Background(), u.Token)
		req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_id="+u.ID+after, nil)
		if err != nil {
			log.WithError(err).Error("Unable to create http request to get twitch streams data")
			return
		}
		req.Header.Set("Client-ID", settings.ClientID)

		resp, err := client.Do(req)
		if err != nil {
			log.WithError(err).Error("Unable to get twitch stream data")
			return
		}

		t := &TwitchStreams{}
		if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
			log.WithError(err).Error("Unable to parse twitch streams data")
		}
		resp.Body.Close()

		if len(t.Data) == 0 {
			break
		}

		result.Data = append(result.Data, t.Data...)
		if t.Pagination == nil || t.Pagination.Cursor == "" {
			break
		}
		after = "&after=" + t.Pagination.Cursor
	}

	if len(result.Data) < 1 {
		log.Info("No streams")
		return
	}

	if len(result.Data) > 0 {
		u.TwitchStream = result.Data[0]
	}
}

func (s *TwitchStream) SaveFiles() {
	data, err := fieldsToMap(s)
	if err != nil {
		log.WithError(err).Error("Unable to convert stream to map")
		return
	}

	for k, v := range data {
		saveContent("stream", k, v)
	}
}
