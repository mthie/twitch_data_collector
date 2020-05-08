package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type TwitchSubscriptions struct {
	Total         int64                 `json:"_total"`
	Subscriptions []*TwitchSubscription `json:"subscriptions"`
}

type TwitchSubscription struct {
	ID          string         `json:"_id"`
	CreatedAt   time.Time      `json:"created_at"`
	IsGift      bool           `json:"is_gift"`
	SubPlan     string         `json:"sub_plan"`
	SubPlanName string         `json:"sub_plan_name"`
	User        *TwitchSubUser `json:"user"`
}

type TwitchSubUser struct {
	ID          string    `json:"_id"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Logo        string    `json:"logo"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func getSubs(u *User) {
	result := &TwitchSubscriptions{}

	limit := 100
	offset := 0

	for {
		client := http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/channels/%s/subscriptions?limit=%d&offset=%d", u.TwitchChannel.ID, limit, offset), nil)
		if err != nil {
			log.WithError(err).Error("Unable to create http request to get twitch subs data")
			return
		}

		req.Header.Set("Client-ID", settings.ClientID)
		req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
		req.Header.Set("Authorization", "OAuth "+u.Token.AccessToken)

		resp, err := client.Do(req)
		if err != nil {
			log.WithError(err).Error("Unable to get twitch subs data")
			return
		}

		t := &TwitchSubscriptions{}
		if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
			log.WithError(err).Error("Unable to parse twitch followers data")
		}
		resp.Body.Close()

		if len(t.Subscriptions) == 0 {
			break
		}

		result.Total = t.Total
		result.Subscriptions = append(result.Subscriptions, t.Subscriptions...)

		offset += limit
	}

	if len(result.Subscriptions) < 1 {
		log.Info("No Subs")
		return
	}

	u.TwitchSubscriptions = result
}
func (s *TwitchSubscriptions) SaveFiles() {
	saveContent("subscriptions", "total", strconv.FormatInt(s.Total, 10))
	saveJSON("subscriptions", "complete_list", s)

	sort.Slice(s.Subscriptions, func(i, j int) bool {
		return s.Subscriptions[i].CreatedAt.Before(s.Subscriptions[j].CreatedAt)
	})
	start := len(s.Subscriptions) - 10
	if start < 0 {
		start = 0
	}
	lastSubs := s.Subscriptions[start:]

	var lastSubsSlice []string
	for _, v := range lastSubs {
		lastSubsSlice = append(lastSubsSlice, v.User.DisplayName)
	}

	saveContent("subscriptions", "last_ten", strings.Join(lastSubsSlice, "\n"))
}
