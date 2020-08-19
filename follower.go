package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type TwitchFollowers struct {
	Total      int64             `json:"total"`
	Data       []*TwitchFollower `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

type TwitchPagination struct {
	Cursor string `json:"cursor"`
}

type TwitchFollower struct {
	FromID     string    `json:"from_id"`
	FromName   string    `json:"from_name"`
	ToID       string    `json:"to_id"`
	ToName     string    `json:"to_name"`
	FollowedAt time.Time `json:"followed_at"`
}

func getFollows(u *User, max int) {
	result := &TwitchFollowers{}

	after := ""
	for {
		client := twitchOauthConfig.Client(context.Background(), u.Token)
		req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users/follows?to_id="+u.ID+after, nil)
		if err != nil {
			log.WithError(err).Error("Unable to create http request to get twitch follower data")
			return
		}
		req.Header.Set("Client-ID", settings.ClientID)

		resp, err := client.Do(req)
		if err != nil {
			log.WithError(err).Error("Unable to get twitch follower data")
			return
		}

		t := &TwitchFollowers{}
		if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
			log.WithError(err).Error("Unable to parse twitch followers data")
		}
		resp.Body.Close()

		if len(t.Data) == 0 {
			break
		}

		result.Total = t.Total
		result.Data = append(result.Data, t.Data...)
		if t.Pagination == nil || t.Pagination.Cursor == "" {
			break
		}

		if max > -1 && len(result.Data) >= max {
			break
		}

		after = "&after=" + t.Pagination.Cursor
	}

	if len(result.Data) < 1 {
		log.Info("No followers")
		return
	}

	u.TwitchFollowers = result
}

func (f *TwitchFollowers) SaveFiles() {
	saveContent("followers", "total", strconv.FormatInt(f.Total, 10))
	saveJSON("followers", "complete_list", f)
	sort.Slice(f.Data, func(i, j int) bool {
		return f.Data[i].FollowedAt.Before(f.Data[j].FollowedAt)
	})
	start := len(f.Data) - 10
	if start < 0 {
		start = 0
	}
	lastFollowers := f.Data[start:]

	var lastFollowerSlice []string
	for _, v := range lastFollowers {
		lastFollowerSlice = append(lastFollowerSlice, v.FromName)
	}

	saveContent("followers", "last_ten", strings.Join(lastFollowerSlice, "\n"))
}
