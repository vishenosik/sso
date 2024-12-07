package config

import (
	"time"

	"github.com/blacksmith-vish/sso/internal/lib/config"
)

type ServerOptions struct {
	Port    uint16        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type GRPCConfig struct {
	ServerOptions `yaml:"server_options"`
}

type RESTConfig struct {
	ServerOptions `yaml:"server_options"`
}

type Servers struct {
	GRPC GRPCConfig `yaml:"grpc"`
	REST RESTConfig `yaml:"rest"`
}

func (srv ServerOptions) options() config.ServerOptions {
	return config.ServerOptions{
		Port:    srv.Port,
		Timeout: srv.Timeout,
	}
}

func (srvs Servers) getRestConfig() config.RESTConfig {
	return config.RESTConfig{
		ServerOptions: srvs.REST.options(),
	}
}

func (srvs Servers) getGrpcConfig() config.GRPCConfig {
	return config.GRPCConfig{
		ServerOptions: srvs.GRPC.options(),
	}
}
