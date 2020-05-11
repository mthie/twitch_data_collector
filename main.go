package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Luzifer/rconfig"
	gorilla "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var cfg = struct {
	SettingsFile string `default:"settings.yml" flag:"settings-file" description:"Path to settings file"`
}{}

func main() {
	rconfig.Parse(&cfg)
	settingsUpdater()

	loadUser()

	log.Infof("User: %+v", user)

	mux := gorilla.NewRouter()
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/login", handleTwitchLogin)
	mux.HandleFunc("/callback", handleTwitchCallback)
	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		if user == nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if err := json.NewEncoder(w).Encode(user.getWebUser()); err != nil {
			log.WithError(err).Error("Unable to encode user data")
			http.Error(w, "Unable to encode user data", http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi!"))
	})

	go func() {
		if user != nil {
			handleSaves()
		}
		c := time.Tick(20 * time.Second)
		for range c {
			if user == nil {
				continue
			}
			handleSaves()
		}
	}()

	log.Info("Starting webserver...")
	log.Fatal(http.ListenAndServe(":"+settings.WebserverPort, mux))
}

func handleSaves() {
	if user == nil {
		return
	}

	getUser(user)

	if user.TwitchUser != nil {
		user.TwitchUser.SaveFiles()
	}

	getChannel(user)
	if user.TwitchChannel != nil {
		user.TwitchChannel.SaveFiles()
	}

	getFollows(user)
	if user.TwitchFollowers != nil {
		user.TwitchFollowers.SaveFiles()
	}

	getSubs(user)
	if user.TwitchSubscriptions != nil {
		user.TwitchSubscriptions.SaveFiles()
	}

	getStreams(user)
	if user.TwitchStream != nil {
		user.TwitchStream.SaveFiles()
	}
}

func handleTwitchLogin(w http.ResponseWriter, r *http.Request) {
	url := twitchOauthConfig.AuthCodeURL(settings.VerificationToken, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleTwitchCallback(w http.ResponseWriter, r *http.Request) {
	twitchAuthToToken(r.FormValue("state"), r.FormValue("code"), w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func twitchAuthToToken(state string, code string, w http.ResponseWriter, r *http.Request) {
	if state != settings.VerificationToken {
		log.Fatal("invalid oauth state")
	}

	timeoutCtx, cancel := context.WithTimeout(oauth2.NoContext, 5*time.Second)
	defer cancel()

	token, err := twitchOauthConfig.Exchange(timeoutCtx, code, oauth2.AccessTypeOffline)
	if err != nil {
		log.WithError(err).Error("code exchange failed")
		return
	}
	u, err := createUser(token)
	if err != nil {
		log.WithError(err).Error("Unable to get user")
		return
	}
	u.Token = token
	user = u
	log.Infof("User: %+v, User: %+v", user, user.TwitchUser)
	saveUser()
}
