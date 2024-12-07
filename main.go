package main

import (
	"flag"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/config"
)

func main() {
	configFile := flag.String("config", "config.yml", "Specify a config file to load")
	flag.Parse()

	appConfig, err := getConfig(configFile)
	if err != nil {
		log.WithField("err", err).Fatal("Failed to load configuration")
	}

	reloadLogConfig(appConfig)
	addFileWatcher(configFile)
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

func addFileWatcher(configFile *string) {
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
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Debugf("Event: %s", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debugf("File modified: %s", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Debugf("File removed: %s", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Debugf("File renamed: %s", event.Name)
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Debugf("File permissions changed: %s", event.Name)
				}

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
