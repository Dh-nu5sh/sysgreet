package collectors

import (
	"context"

	"github.com/veteranbv/hostinfo/internal/network"
)

// DefaultNetworkCollector gathers address information via gopsutil.
type DefaultNetworkCollector struct {
	MaxInterfaces int
}

// NewNetworkCollector returns the default network collector.
func NewNetworkCollector(maxInterfaces int) NetworkCollector {
	return DefaultNetworkCollector{MaxInterfaces: maxInterfaces}
}

// CollectNetwork implements NetworkCollector.
func (c DefaultNetworkCollector) CollectNetwork(ctx context.Context) (NetworkInfo, error) {
	primary, additional, err := network.CollectAddresses(ctx, c.MaxInterfaces)
	if err != nil {
		recordError("network", err)
		return NetworkInfo{}, nil
	}
	info := NetworkInfo{}
	if primary != nil {
		info.Primary = &Address{Interface: primary.Interface, IP: primary.IP}
	}
	for _, addr := range additional {
		info.Additional = append(info.Additional, Address{Interface: addr.Interface, IP: addr.IP})
	}
	return info, nil
}
