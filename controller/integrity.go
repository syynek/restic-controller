package controller

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/config"
	"github.com/syynek/restic-controller/restic"
)

type IntegrityController struct {
	ControllerBase
}

func NewIntegrityController(repositories []*config.Repository) *IntegrityController {
	return &IntegrityController{
		ControllerBase: ControllerBase{
			logger:       log.WithFields(log.Fields{"component": "controller/integrity"}),
			repositories: repositories,
			schedule:     cron.New(),
		},
	}
}

func (controller *IntegrityController) Start() error {
	for _, repository := range controller.repositories {
		if repository.IntegrityCheck.Schedule == "" {
			continue
		}

		if repository.IntegrityCheck.RunOnStartup {
			go controller.RunTask(repository)()
		}

		_, err := controller.schedule.AddFunc(repository.IntegrityCheck.Schedule, controller.RunTask(repository))
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to add cron for repository %s with schedule %s: %s", repository.Name, repository.IntegrityCheck.Schedule, err)
			controller.logger.WithFields(log.Fields{"repository": repository.Name}).Errorf(errorMessage)
			return errors.New(errorMessage)
		}
	}
	controller.schedule.Start()

	return nil
}

func (controller *IntegrityController) RunTask(repository *config.Repository) func() {
	return func() {
		controller.logger.WithField("repository", repository.Name).Info("Running integrity check")
		success, err := restic.RunIntegrityCheck(repository)

		if err != nil {
			controller.logger.WithField("repository", repository.Name).Error(err)
		}

		if success {
			controller.logger.WithField("repository", repository.Name).Info("Integrity check finished")
		}
	}
}

func (controller *IntegrityController) UpdateRepositories(repositories []*config.Repository) {
	controller.ClearSchedule()
	controller.repositories = repositories
	controller.Start()
}

func (controller *IntegrityController) ClearSchedule() {
	for _, entry := range controller.schedule.Entries() {
		controller.schedule.Remove(entry.ID)
	}
}
