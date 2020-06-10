package recon

import (
	"context"
	n "net"
)

func resolveHost(net, addr string) ([]n.IPAddr, error) {
	ctx := context.Background()
	var host string
	var err error
	ips := []n.IPAddr{}

	if net == "" {
		net = "tcp"
	}

	switch net {
	case "tcp", "tcp4", "tcp6":
		if addr != "" {
			if host, _, err = n.SplitHostPort(addr); err != nil {
				return nil, err
			}
		}
	default:
		return nil, n.UnknownNetworkError(net)
	}

	if host == "" {
		return ips, nil
	}

	ip := n.ParseIP(host)
	var rs []n.IPAddr
	if ip != nil {
		rs = append(rs, n.IPAddr{IP: ip})
	} else {
		rs, err = n.DefaultResolver.LookupIPAddr(ctx, host)
	}
	nettype := net[len(net)-1]
	if len(rs) > 0 {
		for _, ipaddr := range rs {
			switch nettype {
			case '4':
				ip4 := ipaddr.IP.To4()
				if ip4 != nil && len(ip4) == n.IPv4len {
					ips = append(ips, ipaddr)
				}
			case '6':
				ip4 := ipaddr.IP.To4()
				if ip4 == nil && len(ipaddr.IP) == n.IPv6len {
					ips = append(ips, ipaddr)
				}
			default:
				ips = append(ips, ipaddr)
			}
		}
	}

	return ips, err
}
