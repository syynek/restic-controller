package controller

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/config"
	"github.com/syynek/restic-controller/restic"
)

type InitializationController struct {
	ControllerBase
}

func NewInitializationController(repositories []*config.Repository) *InitializationController {
	return &InitializationController{
		ControllerBase: ControllerBase{
			logger:       log.WithFields(log.Fields{"component": "controller/init"}),
			repositories: repositories,
			schedule:     cron.New(),
		},
	}
}

func (controller *InitializationController) Start() error {
	for _, repository := range controller.repositories {
		if repository.AutoInitialize {
			controller.RunTask(repository)()
		}
	}

	return nil
}

func (controller *InitializationController) RunTask(repository *config.Repository) func() {
	return func() {
		controller.logger.WithField("repository", repository.Name).Info("Running auto initialization")

		if !restic.IsURLPath(repository.URL) {
			controller.logger.WithField("repository", repository.Name).Warn("Repository is not local")
			return
		}

		if restic.IsFolderRepository(repository.URL) {
			controller.logger.WithField("repository", repository.Name).Debug("Repository already exists")
			return
		}

		success, err := restic.RunInit(repository)

		if err != nil {
			controller.logger.WithField("repository", repository.Name).Error(err)
		}

		if success {
			controller.logger.WithField("repository", repository.Name).Info("Repository initialized")
		}
	}
}

func (controller *InitializationController) UpdateRepositories(repositories []*config.Repository) {
	controller.repositories = repositories
	controller.Start()
}
