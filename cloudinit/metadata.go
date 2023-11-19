package cloudinit

import (
	"fmt"
	"net"
	"strings"
)

type MetaData struct {
	// https://cdn.amazonlinux.com/os-images/2.0.20231101.0/README.cloud-init
	LocalHostname     string `yaml:"local-hostname"`
	NetworkInterfaces string `yaml:"network-interfaces,omitempty"`
}

func NetworkInterfaces(addressWcidr string, gateway string, nameservers []string, search []string) string {
	ip, ipNet, err := net.ParseCIDR(addressWcidr)
	if err != nil {
		panic(err)
	}

	network := ip.Mask(ipNet.Mask)
	broadcast := net.IP(make([]byte, 4))
	for i := range broadcast {
		broadcast[i] = network[i] | ^ipNet.Mask[i]
	}

	result := fmt.Sprintf(`
auto eth0
iface eth0 inet static
address %s
network %s
netmask %s
broadcast %s
gateway %s
	`,
		ip.String(),
		network.String(),
		net.IPv4(ipNet.Mask[0], ipNet.Mask[1], ipNet.Mask[2], ipNet.Mask[3]).String(),
		broadcast.String(),
		gateway)

	result = strings.TrimSpace(result)
	result += "\n"

	result += fmt.Sprintf("dns-nameservers %s\n",
		strings.Join(nameservers, " "),
	)
	if len(search) != 0 {
		result += fmt.Sprintf("dns-search %s\n",
			strings.Join(search, " "),
		)
	}
	return result
}
