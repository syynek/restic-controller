package restic

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/internal/config"
)

// Making possible to mock exec.CommandContext
var execCommandContext = exec.CommandContext

// runRestic runs restic commands and adds the repository
// password to the environment variables of a newly created context
func runRestic(repository *config.Repository, args []string) (bool, error) {
	ctx := context.TODO()

	args = append(args, "-r", repository.URL)

	cmd := execCommandContext(ctx, "restic", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	for k, v := range repository.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	password, err := getRepositoryPassword(repository)
	if err != nil {
		log.WithFields(log.Fields{"component": "restic", "repository": repository.Name}).Error("Invalid repository password")
		return false, err
	}
	cmd.Env = append(cmd.Env, "RESTIC_PASSWORD="+password)

	log.WithFields(log.Fields{"component": "restic", "repository": repository.Name, "cmd": strings.Join(cmd.Args, " ")}).Debug("Running restic command")

	_, err = cmd.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return false, fmt.Errorf("restic command returned with code %d: %s", exiterr.ExitCode(), exiterr.Stderr)
		}
		return false, err
	}

	return true, nil
}

// getRepositoryPassword returns the password of a repository
// by checking the password field in the config or reading
// the contents of a specified password file
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

// IsFolderRepository returns a boolean indicating
// if the folder has the structure of a restic repository
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

// IsURLPath is a modified version of the location.isPath function
// to check if a string a path from the restic repository,
// origin: https://github.com/restic/restic/blob/master/internal/backend/location/location.go
func IsURLPath(s string) bool {
	if strings.HasPrefix(s, "../") || strings.HasPrefix(s, `..\`) {
		return true
	}

	if strings.HasPrefix(s, "./") || strings.HasPrefix(s, `.\`) {
		return true
	}

	if strings.HasPrefix(s, "/") || strings.HasPrefix(s, `\`) {
		return true
	}

	if len(s) < 3 {
		return false
	}

	// check for drive paths
	drive := s[0]
	if !(drive >= 'a' && drive <= 'z') && !(drive >= 'A' && drive <= 'Z') {
		return false
	}

	if s[1] != ':' {
		return false
	}

	if s[2] != '\\' && s[2] != '/' {
		return false
	}

	return true
}
