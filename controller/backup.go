package controller

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/config"
	"github.com/syynek/restic-controller/restic"
)

type BackupController struct {
	ControllerBase
}

func NewBackupController(repositories []*config.Repository) *BackupController {
	return &BackupController{
		ControllerBase: ControllerBase{
			logger:       log.WithFields(log.Fields{"component": "controller/backup"}),
			repositories: repositories,
			schedule:     cron.New(),
		},
	}
}

func (controller *BackupController) StartSchedule() error {
	for _, repository := range controller.repositories {
		if repository.Backup.Schedule == "" {
			continue
		}

		if repository.Backup.RunOnStartup {
			go controller.RunTask(repository)()
		}

		_, err := controller.schedule.AddFunc(repository.Backup.Schedule, controller.RunTask(repository))
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to add cron for repository %s with schedule %s: %s", repository.Name, repository.Backup.Schedule, err)
			controller.logger.WithFields(log.Fields{"repository": repository.Name}).Errorf(errorMessage)
			return errors.New(errorMessage)
		}
	}
	controller.schedule.Start()

	return nil
}

func (controller *BackupController) RunTask(repository *config.Repository) func() {
	return func() {
		controller.logger.WithField("repository", repository.Name).Info("Running backup")
		success, err := restic.RunBackup(repository)
		if success {
			controller.logger.WithField("repository", repository.Name).Info("Backup finished")
		}
	}
}

func (controller *BackupController) UpdateRepositories(repositories []*config.Repository) {
	controller.ClearSchedule()
	controller.repositories = repositories
	controller.StartSchedule()
}

func (controller *BackupController) ClearSchedule() {
	for _, entry := range controller.schedule.Entries() {
		controller.schedule.Remove(entry.ID)
	}
}
