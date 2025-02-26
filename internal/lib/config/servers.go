package config

import (
	"iter"
	"slices"
	"time"

	"github.com/blacksmith-vish/sso/pkg/collections"
	"github.com/pkg/errors"
)

var (
	ErrServerPortMustBeUnique = errors.New("port numbers must be unique")
)

type ServerOptions struct {
	Port    uint16 `validate:"gte=2999,lte=65535"`
	Timeout time.Duration
}

type GRPCConfig struct {
	ServerOptions
}

type RESTConfig struct {
	ServerOptions
}

type Servers []ServerOptions

func comparePorts(servers ...ServerOptions) error {
	if collections.HasDuplicates(Servers(servers).Ports()) {
		return ErrServerPortMustBeUnique
	}
	return nil
}

func (srvs Servers) Ports() []uint16 {
	return slices.Collect(srvs.PortsIter())
}

func (srvs Servers) PortsIter() iter.Seq[uint16] {
	return func(yield func(uint16) bool) {
		for _, srv := range srvs {
			if !yield(srv.Port) {
				return
			}
		}
	}
}
