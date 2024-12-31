package restic

import (
	"strconv"

	"github.com/syynek/restic-controller/internal/config"
)

// RunForget prepares the restic forget command and returns the result
func RunForget(repository *config.Repository) (bool, error) {
	args := []string{}
	args = append(args, "forget", "--prune", "--json", "-q")
	args = append(args, getForgetPolicyArgs(repository.Retention.Policy)...)

	success, err := runRestic(repository, args)

	return success, err
}

// getForgetPolicyArgs returns an array with the forget policy for the restic forget command
func getForgetPolicyArgs(policy *config.ForgetPolicy) []string {
	var args []string

	if policy.KeepLast != 0 {
		args = append(args, "--keep-last="+strconv.Itoa(policy.KeepLast))
	}

	if policy.KeepDaily != 0 {
		args = append(args, "--keep-daily="+strconv.Itoa(policy.KeepDaily))
	}

	if policy.KeepHourly != 0 {
		args = append(args, "--keep-hourly="+strconv.Itoa(policy.KeepHourly))
	}

	if policy.KeepWeekly != 0 {
		args = append(args, "--keep-weekly="+strconv.Itoa(policy.KeepWeekly))
	}

	if policy.KeepMonthly != 0 {
		args = append(args, "--keep-monthly="+strconv.Itoa(policy.KeepMonthly))
	}

	if policy.KeepYearly != 0 {
		args = append(args, "--keep-yearly="+strconv.Itoa(policy.KeepYearly))
	}

	if len(policy.KeepTags) > 0 {
		for _, v := range policy.KeepTags {
			args = append(args, "--keep-tags="+v)
		}
	}

	if len(policy.KeepWithin) > 0 {
		args = append(args, "--keep-within="+policy.KeepWithin)
	}

	return args
}
