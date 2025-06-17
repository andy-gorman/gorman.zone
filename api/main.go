package main

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	server := &http.Server{
		Addr: ":3333",
	}

	env.updateLivetrackUrl()
	done := env.setupTimer()
	shutdownChan := make(chan bool, 1)

	http.HandleFunc("/garmin-live-track", env.handler())

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error:", "error", err)
		}

		slog.Info("Stopped serving new connections.")
		done <- true
		shutdownChan <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	<-shutdownChan
	slog.Info("Graceful shutdown complete")
}

func (e *Env) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, *e.LivetrackUrl)
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
