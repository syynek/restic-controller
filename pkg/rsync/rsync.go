package rsync

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/syynek/restic-controller/internal/config"
)

// Making possible to mock exec.CommandContext
var execCommandContext = exec.CommandContext

func RunRsync(repository *config.Repository) (bool, error) {
	ctx := context.TODO()

	args := []string{ // rsync -e 'ssh -p[port]' --recursive <local folder> <username>@<host>:<target folder>
		"-e", "'ssh -p" + strconv.Itoa(repository.Rsync.Port) + "'",
		"--recursive", repository.URL,
		repository.Rsync.User + "@" + repository.Rsync.Host + ":" + repository.Rsync.TargetFolder,
	}

	cmd := execCommandContext(ctx, "rsync", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)

	log.WithFields(log.Fields{"component": "rsync", "cmd": strings.Join(cmd.Args, " ")}).Debug("Running rsync command")

	_, err := cmd.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return false, fmt.Errorf("rsync command returned with code %d: %s", exiterr.ExitCode(), exiterr.Stderr)
		}
		return false, err
	}

	return true, nil
}
