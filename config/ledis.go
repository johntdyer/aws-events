package config

type ledisConfig struct {
	Path     string `mapstructure:"path"`
	Database int    `mapstructure:"database"`
}
