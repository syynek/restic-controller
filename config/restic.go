package config

// ForgetPolicy specifies to restic retention policy rules
type ForgetPolicy struct {
	KeepLast    int      `mapstructure:"keep_last"`
	KeepDaily   int      `mapstructure:"keep_daily"`
	KeepHourly  int      `mapstructure:"keep_hourly"`
	KeepWeekly  int      `mapstructure:"keep_weekly"`
	KeepMonthly int      `mapstructure:"keep_monthly"`
	KeepYearly  int      `mapstructure:"keep_yearly"`
	KeepTags    []string `mapstructure:"keep_tags"`
	KeepWithin  string   `mapstructure:"keep_within"`
}

// Repository contains the configuration for a restic repository
type Repository struct {
	Name           string            `mapstructure:"name" validate:"required"`
	URL            string            `mapstructure:"url" validate:"required"`
	Password       string            `mapstructure:"password" validate:"required_without=PasswordFile"`
	PasswordFile   string            `mapstructure:"password_file" validate:"required_without=Password"`
	EnvFromFile    map[string]string `mapstructure:"env_from_file"`
	Env            map[string]string `mapstructure:"env"`
	AutoInitialize bool              `mapstructure:"auto_initialize"`
	Backup         struct {
		Schedule     string   `mapstructure:"schedule" validate:"required"`
		RunOnStartup bool     `mapstructure:"run_on_startup"`
		IncludeFiles []string `mapstructure:"include_files" validate:"required"`
		ExcludeFiles []string `mapstructure:"exclude_files"`
	} `mapstructure:"backup" validate:"required_without_all=IntegrityCheck Retention"`
	IntegrityCheck struct {
		Schedule     string `mapstructure:"schedule" validate:"required"`
		RunOnStartup bool   `mapstructure:"run_on_startup"`
	} `mapstructure:"integrity_check" validate:"required_without_all=Backup Retention"`
	Retention struct {
		Schedule     string        `mapstructure:"schedule" validate:"required"`
		RunOnStartup bool          `mapstructure:"run_on_startup"`
		Policy       *ForgetPolicy `mapstructure:"policy" validate:"required"`
	} `mapstructure:"retention" validate:"required_without_all=Backup IntegrityCheck"`
}
