package network

import (
	"context"
	"net"
	"sort"
	"strings"

	gnet "github.com/shirou/gopsutil/v3/net"
)

// Address represents an IP address bound to an interface.
type Address struct {
	Interface string
	IP        string
}

var virtualPrefixes = []string{"docker", "veth", "vbox", "vmnet", "br-", "utun", "tap", "wg", "zt"}

// CollectAddresses gathers interface data from the OS and applies filtering rules.
func CollectAddresses(ctx context.Context, maxAdditional int) (primary *Address, additional []Address, err error) {
	stats, err := gnet.InterfacesWithContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	candidates := filterStats(stats)
	if len(candidates) == 0 {
		return nil, nil, nil
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return rankInterface(candidates[i].Interface) < rankInterface(candidates[j].Interface)
	})
	primary = &candidates[0]
	if len(candidates) > 1 {
		limit := maxAdditional
		if limit <= 0 || limit > len(candidates)-1 {
			limit = len(candidates) - 1
		}
		additional = append([]Address{}, candidates[1:limit+1]...)
	}
	return primary, additional, nil
}

func filterStats(stats []gnet.InterfaceStat) []Address {
	var out []Address
	for _, stat := range stats {
		if !isUp(stat) || isLoopback(stat) || isVirtual(stat) {
			continue
		}
		for _, addr := range stat.Addrs {
			ip := parseIP(addr.Addr)
			if ip == "" {
				continue
			}
			out = append(out, Address{Interface: stat.Name, IP: ip})
		}
	}
	return out
}

func parseIP(cidr string) string {
	if cidr == "" {
		return ""
	}
	host, _, err := net.ParseCIDR(cidr)
	if err != nil {
		host = net.ParseIP(cidr)
	}
	if host == nil {
		return ""
	}
	if host.IsLoopback() {
		return ""
	}
	if host.IsLinkLocalUnicast() {
		return ""
	}
	if host.To4() == nil {
		return ""
	}
	return host.String()
}

func isUp(stat gnet.InterfaceStat) bool {
	for _, flag := range stat.Flags {
		if flag == "up" || flag == "UP" {
			return true
		}
	}
	return false
}

func isLoopback(stat gnet.InterfaceStat) bool {
	for _, flag := range stat.Flags {
		if strings.EqualFold(flag, "loopback") {
			return true
		}
	}
	return strings.HasPrefix(strings.ToLower(stat.Name), "lo")
}

func isVirtual(stat gnet.InterfaceStat) bool {
	lower := strings.ToLower(stat.Name)
	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

func rankInterface(name string) int {
	lower := strings.ToLower(name)
	switch {
	case strings.HasPrefix(lower, "eth"), strings.HasPrefix(lower, "en"), strings.HasPrefix(lower, "eno"), strings.HasPrefix(lower, "ens"), strings.HasPrefix(lower, "wlan"), strings.HasPrefix(lower, "wifi"):
		return 0
	case strings.HasPrefix(lower, "em"), strings.HasPrefix(lower, "p"):
		return 1
	default:
		return 5
	}
}
