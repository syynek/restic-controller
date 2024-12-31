package controller

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/internal/config"
	"github.com/syynek/restic-controller/pkg/restic"
)

type RetentionController struct {
	ControllerBase
}

func NewRetentionController(repositories []*config.Repository) *RetentionController {
	return &RetentionController{
		ControllerBase: ControllerBase{
			logger:       log.WithFields(log.Fields{"component": "controller/retention"}),
			repositories: repositories,
			schedule:     cron.New(),
		},
	}
}

func (controller *RetentionController) Start() error {
	for _, repository := range controller.repositories {
		if repository.Retention.Schedule == "" {
			continue
		}

		if repository.Retention.RunOnStartup {
			go controller.RunTask(repository)()
		}

		_, err := controller.schedule.AddFunc(repository.Retention.Schedule, controller.RunTask(repository))
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to add cron for repository %s with schedule %s: %s", repository.Name, repository.Retention.Schedule, err)
			controller.logger.WithFields(log.Fields{"repository": repository.Name}).Errorf(errorMessage)
			return errors.New(errorMessage)
		}
	}
	controller.schedule.Start()

	return nil
}

func (controller *RetentionController) RunTask(repository *config.Repository) func() {
	return func() {
		controller.logger.WithField("repository", repository.Name).Info("Running retention")
		success, err := restic.RunForget(repository)

		if err != nil {
			controller.logger.WithField("repository", repository.Name).Error(err)
		}

		if success {
			controller.logger.WithField("repository", repository.Name).Info("Retention finished")
		}
	}
}

func (controller *RetentionController) UpdateRepositories(repositories []*config.Repository) {
	controller.ClearSchedule()
	controller.repositories = repositories
	controller.Start()
}

func (controller *RetentionController) ClearSchedule() {
	for _, entry := range controller.schedule.Entries() {
		controller.schedule.Remove(entry.ID)
	}
}
