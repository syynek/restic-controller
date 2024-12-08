package controller

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/config"
)

type ControllerInterface interface {
	StartSchedule() error
	RunTask(repository *config.Repository) func()
	UpdateRepositories(repositories []*config.Repository)
}

type ControllerBase struct {
	logger       *log.Entry
	repositories []*config.Repository
	schedule     *cron.Cron
}
