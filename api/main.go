package main

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/andy-gorman/gorman.zone/api/fastmail"
	"github.com/andy-gorman/gorman.zone/api/garmin"
)

type Env struct {
	Client       *fastmail.FastmailClient
	LivetrackUrl *string
}

func main() {
	client := fastmail.NewEmailClient(
		os.Getenv("FASTMAIL_AUTH_URL"),
		os.Getenv("FASTMAIL_AUTH_TOKEN"),
		os.Getenv("FASTMAIL_ACCOUNT_ID"),
		os.Getenv("FASTMAIL_LIVETRACK_FOLDER_ID"),
	)

	env := &Env{
		Client: client,
	}

	env.updateLivetrackUrl()
	done := env.setupTimer()

	http.HandleFunc("/garmin-live-track", env.handler())

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server closed\n")

	} else if err != nil {
		slog.Error("error starting server", "error", err.Error())
		os.Exit(1)
	}
	done <- true
}

func (e *Env) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if garmin.IsLivetrackLinkActive(*e.LivetrackUrl) {
			io.WriteString(w, *e.LivetrackUrl)
		} else {
			io.WriteString(w, "")
		}
	}
}

func (e *Env) updateLivetrackUrl() {
	emailBodyText, err := e.Client.FetchLivetrackEmail()
	if err != nil {
		slog.Error("Unable to fetch email", "error", err.Error())
		return
	}
	link, err := garmin.ParseLivetrackLinkFromEmail(emailBodyText)
	if err != nil {
		slog.Error("Unable To Parse Livetrack Link", "error", err.Error())
		return
	}
	e.LivetrackUrl = &link
	slog.Info(*e.LivetrackUrl)
}

func (e *Env) setupTimer() chan bool {
	done := make(chan bool)
	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				e.updateLivetrackUrl()
			}
		}
	}()

	return done
}
