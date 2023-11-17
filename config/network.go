package config

import (
	"fmt"
	"log"
	"net"

	"github.com/iwanhae/qontainer/config/dns"
	"github.com/vishvananda/netlink"
)

type Network struct {
	Type NetworkType `env:"NETWORK_TYPE" envDefault:"user"`

	Interface      string   `env:"NETWORK_INTERFACE" envDefault:"eth0"`
	Address        string   `env:"NETWORK_ADDRESS"`
	DefaultGateway string   `env:"NETWORK_DEFAULT_GATEWAY"`
	Nameservers    []string `env:"NETWORK_NAMESERVERS"`
	Search         []string `env:"NETWORK_SEARCH"`
}

type NetworkType string

const (
	NetworkType_Bridge NetworkType = "bridge"
	NetworkType_User   NetworkType = "user"
)

func (c *Network) AutoComplete() error {
	if c.Type == NetworkType_User {
		// QEMU will handle network settings
		return nil
	}

	if c.Type == NetworkType_Bridge {
		// Will use containers eth0 network settings
		nic, err := net.InterfaceByName(c.Interface)
		if err != nil {
			return fmt.Errorf("fail to get network interface %q's info: %w", c.Interface, err)
		}

		// ip addr
		if c.Address == "" {
			addr, err := nic.Addrs()
			if err != nil {
				return fmt.Errorf("fail to get address of %q: %w", c.Interface, err)
			} else if len(addr) == 0 {
				return fmt.Errorf("no ip address in interface %q: %w", c.Interface, err)
			} else if len(addr) != 1 {
				log.Printf("WARN: more than one ip address in %q, will use %q", c.Interface, addr[0].String())
			}
			c.Address = addr[0].String()
		}

		// gateway
		if c.DefaultGateway == "" {
			routes, err := netlink.RouteGet(net.ParseIP("1.1.1.1"))
			if err != nil {
				return fmt.Errorf("fail to infer default gateway: %w", err)
			} else if len(routes) == 0 {
				return fmt.Errorf("fail to infer default gateway: no routes to 1.1.1.1 found")
			}
			c.DefaultGateway = routes[0].Gw.String()
		}

		// nameservers
		if len(c.Nameservers) == 0 {
			resolvConf := dns.DnsReadConfig("/etc/resolv.conf")
			c.Nameservers = resolvConf.Servers
			c.Search = resolvConf.Search
		}
		return nil
	}

	return fmt.Errorf("unknown network type %q", c.Type)
}
