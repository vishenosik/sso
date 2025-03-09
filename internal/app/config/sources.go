package config

import (
	configFS "github.com/blacksmith-vish/sso/internal/store/config/fs"
)

type ConfigSources struct {
	FS *configFS.ConfigFS
}

func NewConfigSources() ConfigSources {
	return ConfigSources{
		FS: configFS.MustLoad(),
	}
}
