package controller

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/internal/config"
	"github.com/syynek/restic-controller/pkg/restic"
	"github.com/syynek/restic-controller/pkg/rsync"
)

type RsyncController struct {
	ControllerBase
}

func NewRsyncController(repositories []*config.Repository) *RsyncController {
	return &RsyncController{
		ControllerBase: ControllerBase{
			logger:       log.WithFields(log.Fields{"component": "controller/rsync"}),
			repositories: repositories,
			schedule:     cron.New(),
		},
	}
}

func (controller *RsyncController) Start() error {
	for _, repository := range controller.repositories {
		if repository.Rsync.Schedule == "" {
			continue
		}

		_, err := controller.schedule.AddFunc(repository.Rsync.Schedule, controller.RunTask(repository))
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to add cron for repository %s with schedule %s: %s", repository.Name, repository.Rsync.Schedule, err)
			controller.logger.WithFields(log.Fields{"repository": repository.Name}).Errorf(errorMessage)
			return errors.New(errorMessage)
		}
	}
	controller.schedule.Start()

	return nil
}

func (controller *RsyncController) RunTask(repository *config.Repository) func() {
	return func() {
		controller.logger.WithField("repository", repository.Name).Info("Running Rsync")

		if !restic.IsURLPath(repository.URL) {
			controller.logger.WithField("repository", repository.Name).Error("Repository is not local")
			return
		}

		if !restic.IsFolderRepository(repository.URL) {
			controller.logger.WithField("repository", repository.Name).Error("Folder is not a repository or possibly broken")
			return
		}

		success, err := rsync.RunRsync(repository)

		if err != nil {
			controller.logger.WithField("repository", repository.Name).Error(err)
		}

		if success {
			controller.logger.WithField("repository", repository.Name).Info("Rsync finished")
		}
	}
}

func (controller *RsyncController) UpdateRepositories(repositories []*config.Repository) {
	controller.ClearSchedule()
	controller.repositories = repositories
	controller.Start()
}

func (controller *RsyncController) ClearSchedule() {
	for _, entry := range controller.schedule.Entries() {
		controller.schedule.Remove(entry.ID)
	}
}
