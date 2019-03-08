package config

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	msgEmptyValue = "empty config value '%s'"

	defaultTimeout = 10
	defaultLang    = "en-US"
)

var mandatoryParams = []string{
	"telegram.token",
	"telegram.whitelist",
	"yandex.token",
	"yandex.folder_id",
	"commands.start",
	"messages.hello",
	"messages.in_progress",
	"errors.download",
	"errors.failed",
	"errors.whitelist",
	"errors.unknown",
}

// New creates new viper config instance
func New(name string) (*viper.Viper, error) {
	if name == "" {
		return nil, errors.New("empty config name")
	}

	cfg := viper.New()

	cfg.SetConfigName(name)
	cfg.SetConfigType("toml")
	cfg.AddConfigPath("$HOME/.config")
	cfg.AddConfigPath("/etc")
	cfg.AddConfigPath(".")

	if err := cfg.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "unable to read config")
	}
	cfg.WatchConfig()

	cfg.SetDefault("telegram.timeout", time.Duration(defaultTimeout)*time.Second)
	cfg.SetDefault("yandex.lang", defaultLang)

	if err := validateConfig(cfg); err != nil {
		return nil, errors.Wrap(err, "unable to validate config")
	}

	return cfg, nil
}

func validateConfig(cfg *viper.Viper) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	for _, p := range mandatoryParams {
		if cfg.Get(p) == nil {
			return errors.Errorf(msgEmptyValue, p)
		}
	}

	if len(cfg.GetStringSlice("telegram.whitelist")) < 1 {
		return errors.Errorf("'telegram.whitelist' should contain at least one chatID")
	}

	return nil
}
