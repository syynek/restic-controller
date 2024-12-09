package restic

import (
	"github.com/syynek/restic-controller/config"
)

func RunBackup(repository *config.Repository) (bool, error) {
	args := []string{}
	args = append(args, "--exclude")
	args = append(args, repository.Backup.ExcludeFiles...)
	args = append(args, "backup")
	args = append(args, repository.Backup.IncludeFiles...)

	success, err := runRestic(repository, args)

	return success, err
}
