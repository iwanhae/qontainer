package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"reflect"

	"github.com/caarlos0/env/v10"
	"github.com/iwanhae/qemtainer/cloudinit"
	"github.com/vishvananda/netlink"
)

type Config struct {
	CPU    string `env:"VM_CPU" envDefault:"1"`
	Memory string `env:"VM_MEMORY" envDefault:"1G"`
	Disk   string `env:"VM_DISK" envDefault:"/data/disk.img"`

	NetworkInterface      string   `env:"NETWORK_INTERFACE" envDefault:"eth0"`
	NetworkAddress        string   `env:"NETWORK_ADDRESS"`
	NetworkDefaultGateway string   `env:"NETWORK_DEFAULT_GATEWAY"`
	NetworkNameservers    []string `env:"NETWORK_NAMESERVERS" envDefault:"8.8.8.8,8.8.4.4"`

	GuestHostname string `env:"GUEST_HOSTNAME"`
	GuestUsername string `env:"GUEST_USERNAME" envDefault:"deploy"`
	GuestPassword string `env:"GUEST_PASSWORD" envDefault:"$6$rounds=4096$KUjo2cumnYaz0fmk$EsoVV1xP/FXIkv5mm4V26CR3qJrDZhs3Rga8OfBKNBUSsmCM7OHouHMHHz8lApGsD835DqpFvAgqJv1Hq5J.k0"`
	GuestSudo     string `env:"GUEST_SUDO" envDefault:"ALL=(ALL) NOPASSWD:ALL"`

	QemuExecutable string `env:"QEMU_EXECUTABLE" envDefault:"qemu-system-x86_64"`
}

func main() {
	fmt.Println(` / _ \  ___ _ __ ___ | |_ __ _(_)_ __   ___ _ __ `)
	fmt.Println(`| | | |/ _ \ '_ ' _ \| __/ _' | | '_ \ / _ \ '__|`)
	fmt.Println(`| |_| |  __/ | | | | | || (_| | | | | |  __/ |   `)
	fmt.Println(` \__\_\\___|_| |_| |_|\__\__,_|_|_| |_|\___|_|   `)

	// Load Config
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// Networks
	{
		// Get IP Address
		nic, err := net.InterfaceByName(cfg.NetworkInterface)
		if err != nil {
			panic(
				fmt.Errorf("fail to get network interface named %q: %w", cfg.NetworkInterface, err),
			)
		}
		addr, err := nic.Addrs()
		if err != nil {
			panic(
				fmt.Errorf("fail to get ip address of %q: %w", cfg.NetworkInterface, err),
			)
		}
		cfg.NetworkAddress = addr[0].String()
	}

	{
		routes, err := netlink.RouteGet(net.ParseIP(cfg.NetworkNameservers[0]))
		if err != nil {
			panic(
				fmt.Errorf("fail to infer default gateway to use: %w", err),
			)
		}
		cfg.NetworkDefaultGateway = routes[0].Gw.String()
	}

	// Guset
	{
		hostname, err := os.Hostname()
		if err != nil {
			panic(
				fmt.Errorf("fail fetch hostname: %w", err),
			)
		}
		cfg.GuestHostname = hostname
	}

	{
		fmt.Println("----------CONFIG----------")
		t := reflect.TypeOf(cfg)
		v := reflect.ValueOf(cfg)
		for i, f := range reflect.VisibleFields(t) {
			fmt.Printf("%s=%q\n", f.Tag.Get("env"), v.Field(i))
		}
	}

	if err := run(context.Background(), cfg); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, cfg Config) error {
	if err := createCloudInitISO(&cfg); err != nil {
		return fmt.Errorf("fail to create cloudinit file: %w", err)
	}
	fmt.Println("----------START VM----------")
	defer fmt.Println("----------VM Terminated Bye Bye~ :)----------")
	cmd := exec.Command(cfg.QemuExecutable)
	cmd.Args = append(cmd.Args, "-nographic")
	cmd.Args = append(cmd.Args, "-enable-kvm")
	cmd.Args = append(cmd.Args, "-cpu", "host")
	cmd.Args = append(cmd.Args, "-m", cfg.Memory)
	cmd.Args = append(cmd.Args, "-smp", cfg.CPU)
	cmd.Args = append(cmd.Args, "-nic", "user,model=virtio-net-pci")
	cmd.Args = append(cmd.Args, "-cdrom", "./cloudinit.iso")
	cmd.Args = append(cmd.Args, "-hda", cfg.Disk)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func createCloudInitISO(cfg *Config) error {
	ci := cloudinit.CloudConfig{
		UserData: cloudinit.UserData{
			Hostname:         cfg.GuestHostname,
			DisableRoot:      true,
			PreserveHostname: false,
			GrowPartition: cloudinit.GrowPartitionConfig{
				Mode:    cloudinit.GrowPartitionMode_Auto,
				Devices: []string{"/"},
			},
			Users: []cloudinit.UserCoinfig{
				{
					Name:         cfg.GuestUsername,
					HashedPasswd: cfg.GuestPassword,
					Sudo:         "ALL=(ALL) NOPASSWD:ALL",
					LockPasswd:   false,
				},
			},
		},
		NetworkConfig: cloudinit.NetworkConfig{
			Network: cloudinit.Network{
				Version: 2,
				Ethernets: map[string]cloudinit.Ethernet{
					// virtio-net-pci => ens3
					"ens3": {
						Addresses: []string{cfg.NetworkAddress},
						Routes: []cloudinit.Routes{
							{To: "default", Via: cfg.NetworkDefaultGateway},
						},
						Nameservers: cloudinit.Nameservers{
							Addresses: cfg.NetworkNameservers,
						},
					},
				},
			},
		},
	}
	return ci.SaveTo("./cloudinit.iso")
}
