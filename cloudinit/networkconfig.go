package cloudinit

// https://netplan.io
type NetworkConfig struct {
	Network Network `yaml:"network"`
}
type Routes struct {
	To    string `yaml:"to"`
	Via   string `yaml:"via,omitempty"`
	Scope string `yaml:"scope,omitempty"`
}
type Nameservers struct {
	Search    []string `yaml:"search,omitempty"`
	Addresses []string `yaml:"addresses"`
}
type Ethernet struct {
	Addresses   []string    `yaml:"addresses"`
	Routes      []Routes    `yaml:"routes"`
	Nameservers Nameservers `yaml:"nameservers,omitempty"`
}

type Network struct {
	// MUST BE "2"
	Version   int                 `yaml:"version"`
	Ethernets map[string]Ethernet `yaml:"ethernets"`
}
