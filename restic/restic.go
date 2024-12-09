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

func buildCmdEnv(repositoryPassword string, env *map[string]string) []string {
	var cmdEnv []string
	cmdEnv = append(cmdEnv, "RESTIC_PASSWORD="+repositoryPassword)
	for k, v := range *env {
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, v))
	}

	return cmdEnv
}

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
			log.WithFields(log.Fields{"err": err, "repository": repository.Name}).Warning("Failed to read password file")
			return "", err
		}
		return string(password), nil
	}

	return "", nil
}
