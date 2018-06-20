package config

type ledisConfig struct {
	Path                   string `mapstructure:"path"`
	Database               int    `mapstructure:"database"`
	KeyTimeExpireInSeconds int64  `mapstructure:"key_expire_time"`
}
