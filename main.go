package main

import (
	"flag"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/config"
	"github.com/syynek/restic-controller/controller"
)

func main() {
	configFile := flag.String("config", "config.yml", "Specify a config file to load")
	flag.Parse()

	appConfig, err := getConfig(configFile)
	if err != nil {
		log.WithField("err", err).Fatal("Failed to load configuration")
	}

	reloadLogConfig(appConfig)

	initializationController := controller.NewInitializationController(appConfig.Repositories)
	initializationController.Start()

	backupController := controller.NewBackupController(appConfig.Repositories)
	backupController.Start()

	integrityController := controller.NewIntegrityController(appConfig.Repositories)
	integrityController.Start()

	retentionController := controller.NewRetentionController(appConfig.Repositories)
	retentionController.Start()

	controllers := []controller.ControllerInterface{initializationController, backupController, integrityController, retentionController}

	addFileWatcher(configFile, controllers)
}

func getConfig(configFile *string) (*config.AppConfig, error) {
	appConfig, err := config.ReloadConfig(*configFile)
	if err != nil {
		return nil, err
	}

	return appConfig, nil
}

func reloadLogConfig(appConfig *config.AppConfig) {
	err := config.ConfigureLogging(&appConfig.Log)
	if err != nil {
		log.WithField("err", err).Fatal("Failed to configure logging")
	}
}

func addFileWatcher(configFile *string, controllers []controller.ControllerInterface) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.WithField("err", err).Fatal("Failed to add file watcher")
	}
	defer watcher.Close()

	err = watcher.Add(*configFile)
	if err != nil {
		log.WithField("err", err).Fatal("Failed to add config file to file watcher")
	}

	log.Debugf("Watching for changes in %s", *configFile)

	done := make(chan bool)
	go func() {
		var lastEventTime time.Time
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write != fsnotify.Write {
					continue
				}

				now := time.Now()
				if now.Sub(lastEventTime) < 100*time.Millisecond {
					continue
				}
				lastEventTime = now

				log.Debugf("File modified: %s", event.Name)

				time.Sleep(100 * time.Millisecond)

				appConfig, err := getConfig(configFile)
				if err != nil {
					log.WithField("err", err).Fatal("Failed to load configuration")
				}

				updateControllers(controllers, appConfig.Repositories)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Errorf("Error: %s", err)
			}
		}
	}()

	<-done
}

func updateControllers(controllers []controller.ControllerInterface, repositories []*config.Repository) {
	for _, controllerInstance := range controllers {
		controllerInstance.UpdateRepositories(repositories)
	}
}
