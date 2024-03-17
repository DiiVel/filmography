package main

import (
	"filmography/config"
	"filmography/internal/repository"
	"filmography/service"
	"github.com/sirupsen/logrus"
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

	//srv := &Server{}
	//go func() {
	//	if err := srv.Run(handlersEngine); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//		logrus.WithFields(logrus.Fields{
	//			"error": err,
	//		}).Error("server run failed")
	//		return
	//	}
	//}()
}
