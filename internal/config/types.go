package config

// Config is the root configuration structure for dcctl.
type Config struct {
	SchemaVersion string `yaml:"schema_version" mapstructure:"schema_version"`
	Common       Common `yaml:"common" mapstructure:"common"`
	Dcctl        Dcctl  `yaml:"dcctl" mapstructure:"dcctl"`
}

// Common holds shared settings for compose and optional dcctl options.
type Common struct {
	Dcctl   CommonDcctl   `yaml:"dcctl" mapstructure:"dcctl"`
	Compose CommonCompose `yaml:"compose" mapstructure:"compose"`
}

// CommonDcctl holds common dcctl options (e.g. verbose level).
type CommonDcctl struct {
	VerboseLevel int `yaml:"verbose_level" mapstructure:"verbose_level"`
}

// CommonCompose holds Docker Compose project settings.
type CommonCompose struct {
	IgnoreOrphans bool   `yaml:"ignore_orphans" mapstructure:"ignore_orphans"`
	ProjectName   string `yaml:"project_name" mapstructure:"project_name"`
}

// Dcctl holds environment selection and environment definitions.
type Dcctl struct {
	DefaultEnvironment string                 `yaml:"default_environment" mapstructure:"default_environment"`
	Environments       map[string]Environment `yaml:"environments" mapstructure:"environments"`
}

// Environment defines the list of services (manifest names) for an environment.
type Environment struct {
	Services []string `yaml:"services" mapstructure:"services"`
}
