package main

import (
	"errors"
	"filmography/config"
	"filmography/internal/handlers"
	"filmography/internal/repository"
	"filmography/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("config new failed")
	}

	repo, err := repository.New(cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("repository new failed")
	}
	defer func() {
		err := repo.Close()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("repo close failed")
		}
	}()

	svc := service.New(repo)
	handlersEngine, err := handlers.SetRequestHandlers(svc, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("set request handlers failed")
	}

	srv := &Server{}
	go func() {
		if err := srv.Run(handlersEngine); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("server run failed")
			return
		}
	}()
	if err := srv.WaitForShutDown(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("server shut down failed")
	}
}
