package controller

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/internal/config"
)

// ControllerInterface the default interface for all controllers
type ControllerInterface interface {
	Start() error
	RunTask(repository *config.Repository) func()
	UpdateRepositories(repositories []*config.Repository)
}

// ControllerBase is the default struct for all controllers
type ControllerBase struct {
	logger       *log.Entry
	repositories []*config.Repository
	schedule     *cron.Cron
}
