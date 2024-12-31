package restic

import (
	"github.com/syynek/restic-controller/internal/config"
)

// RunInit prepares the restic init command and returns the result
func RunInit(repository *config.Repository) (bool, error) {
	args := []string{}
	args = append(args, "init", "-r", repository.URL)

	success, err := runRestic(repository, args)

	return success, err
}
