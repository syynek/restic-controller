package restic

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/config"
)

// Making possible to mock exec.CommandContext
var execCommandContext = exec.CommandContext

func runRestic(repository *config.Repository, args []string) (bool, error) {
	ctx := context.TODO()

	args = append(args, "-r", repository.URL)

	cmd := execCommandContext(ctx, "restic", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)

	password, err := getRepositoryPassword(repository)
	if err != nil {

	}
	cmd.Env = append(cmd.Env, "RESTIC_PASSWORD="+password)

	log.WithFields(log.Fields{"component": "restic", "cmd": strings.Join(cmd.Args, " ")}).Debug("Running restic command")

	_, err = cmd.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return false, fmt.Errorf("restic command returned with code %d: %s", exiterr.ExitCode(), exiterr.Stderr)
		}
		return false, err
	}

	return true, nil
}

func getRepositoryPassword(repository *config.Repository) (string, error) {
	if repository.Password != "" {
		return repository.Password, nil
	}
	if repository.PasswordFile != "" {
		password, err := os.ReadFile(repository.PasswordFile)
		if err != nil {
			log.WithFields(log.Fields{"err": err, "repository": repository.Name}).Error("Failed to read password file")
			return "", err
		}
		return string(password), nil
	}

	return "", nil
}

func IsFolderRepository(path string) bool {
	repositoryStructure := map[string]string{
		"data":      "folder",
		"index":     "folder",
		"keys":      "folder",
		"locks":     "folder",
		"snapshots": "folder",

		"config": "file",
	}

	for objectName, objectType := range repositoryStructure {
		info, err := os.Stat(path + "/" + objectName)
		if os.IsNotExist(err) {
			return false
		}
		if objectType == "folder" && !info.IsDir() {
			return false
		}
	}

	return true
}
