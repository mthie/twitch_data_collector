package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-irc/irc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var (
	user *User
)

type User struct {
	ID                  string
	Name                string
	DisplayName         string
	Token               *oauth2.Token
	IRCClient           *irc.Client          `json:"-"`
	TwitchUser          *TwitchUser          `json:"-"`
	TwitchChannel       *TwitchChannel       `json:"-"`
	TwitchFollowers     *TwitchFollowers     `json:"-"`
	TwitchSubscriptions *TwitchSubscriptions `json:"-"`
	TwitchStream        *TwitchStream        `json:"-"`
}

type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImage    string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int64  `json:"view_count"`
}

func createUser(token *oauth2.Token) (*User, error) {
	u := &User{
		Token: token,
	}

	if err := getUser(u); err != nil {
		log.WithError(err).Error("Unable to get user")
		return nil, err
	}

	u.ID = u.TwitchUser.ID
	u.Name = u.TwitchUser.Login
	u.DisplayName = u.TwitchUser.DisplayName

	return u, nil
}

func getUser(u *User) error {
	client := twitchOauthConfig.Client(context.Background(), u.Token)
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		log.WithError(err).Error("Unable to create http request to get twitch user data")
		return err
	}

	req.Header.Set("Client-ID", settings.ClientID)
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("Unable to get twitch user data")
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Unable to get response body for user data")
		return err
	}

	users := map[string][]*TwitchUser{}

	if err := json.Unmarshal(b, &users); err != nil {
		log.WithError(err).Error("Unable to parse twitch user data")
		return err
	}

	u.TwitchUser = users["data"][0]

	return nil
}

func saveUser() {
	log.Info("Save user")
	f, err := os.OpenFile(".user.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.WithError(err).Fatal("Unable to open file")
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(user); err != nil {
		log.WithError(err).Fatal("Unable to encode users data")
	}
}

func loadUser() {
	data, err := ioutil.ReadFile(".user.json")
	if err != nil {
		log.WithError(err).Error("Unable to read .users.json")
		return
	}

	var u *User

	b := bytes.NewBuffer(data)
	if err := json.NewDecoder(b).Decode(&u); err != nil {
		log.WithError(err).Error("Unable to decode .users.json")
		return
	}

	if u == nil {
		log.Warning("User not existent")
		return
	}

	if err := getUser(u); err != nil {
		log.WithError(err).Error("Unable to get user information")
		return
	}

	if err := getChannel(u); err != nil {
		log.WithError(err).Error("Unable to get channel information")
		return
	}

	user = u
}

func (t *TwitchUser) SaveFiles() {
	data, err := fieldsToMap(t)
	if err != nil {
		log.WithError(err).Error("Unable to convert user to map")
		return
	}

	for k, v := range data {
		saveContent("user", k, v)
	}
}
