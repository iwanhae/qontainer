package cloudinit

// https://netplan.io
type NetworkConfig struct {
	Network Network `json:"network"`
}
type Routes struct {
	To  string `json:"to"`
	Via string `json:"via"`
}
type Nameservers struct {
	Addresses []string `json:"addresses"`
}
type Ethernet struct {
	Addresses   []string    `json:"addresses"`
	Routes      []Routes    `json:"routes"`
	Nameservers Nameservers `json:"nameservers"`
}

type Network struct {
	// MUST BE "2"
	Version   int                 `json:"version"`
	Ethernets map[string]Ethernet `json:"ethernets"`
}
