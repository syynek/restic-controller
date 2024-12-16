package restic

import (
	"github.com/syynek/restic-controller/config"
)

func RunInit(repository *config.Repository) (bool, error) {
	args := []string{}
	args = append(args, "init", "-r", repository.URL)

	success, err := runRestic(repository, args)

	return success, err
}
