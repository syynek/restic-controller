package restic

import (
	"github.com/syynek/restic-controller/config"
)

func RunIntegrityCheck(repository *config.Repository) (bool, error) {
	args := []string{}
	args = append(args, "check", "-q", "--no-lock")

	success, err := runRestic(repository, args)

	return success, err
}
