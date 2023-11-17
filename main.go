package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/caarlos0/env/v10"
	"github.com/iwanhae/qontainer/cloudinit"
	"github.com/iwanhae/qontainer/config"
	"github.com/vishvananda/netlink"
)

func main() {
	fmt.Println(` / _ \  ___ _ __ ___ | |_ __ _(_)_ __   ___ _ __ `)
	fmt.Println(`| | | |/ _ \ '_ ' _ \| __/ _' | | '_ \ / _ \ '__|`)
	fmt.Println(`| |_| |  __/ | | | | | || (_| | | | | |  __/ |   `)
	fmt.Println(` \__\_\\___|_| |_| |_|\__\__,_|_|_| |_|\___|_|   `)

	// Load Config
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	if err := cfg.AutoComplete(); err != nil {
		panic(err)
	}
	cfg.Print()

	if err := run(context.Background(), cfg); err != nil {
		panic(err)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func run(ctx context.Context, cfg config.Config) error {
	fmt.Println("----------Prepare VM----------")
	fmt.Println("creating cloudinit")
	ciImage, err := createCloudInitISO(&cfg)
	if err != nil {
		return fmt.Errorf("fail to create cloudinit file: %w", err)
	}

	if cfg.Network.Type == config.NetworkType_Bridge {
		if _, err := netlink.LinkByName("br0"); err != nil && errors.As(err, &netlink.LinkNotFoundError{}) {
			fmt.Println("create bridge: br0")
			if err := netlink.LinkAdd(&netlink.Bridge{
				LinkAttrs: netlink.LinkAttrs{Name: "br0"},
			}); err != nil {
				return fmt.Errorf("fail to create br0 bridge: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("fail to create bridge: %w", err)
		} else {
			fmt.Println("skip creating bridge: br0")
		}

		if err := netlink.LinkSetUp(
			must(netlink.LinkByName("br0")),
		); err != nil {
			return fmt.Errorf("'ip link set br0 up' failed: %w", err)
		}

		if err := netlink.LinkSetMaster(
			must(netlink.LinkByName("eth0")),
			must(netlink.LinkByName("br0")),
		); err != nil {
			return fmt.Errorf("'ip link set dev eth0 master br0: %w", err)
		}

	}

	fmt.Println("----------START VM----------")
	defer fmt.Println("----------VM Terminated Bye Bye~ :)----------")
	cmd := exec.Command(cfg.QemuExecutable)
	cmd.Args = append(cmd.Args, "-nographic")
	cmd.Args = append(cmd.Args, "-enable-kvm")
	cmd.Args = append(cmd.Args, "-cpu", "host")
	cmd.Args = append(cmd.Args, "-m", cfg.Memory)
	cmd.Args = append(cmd.Args, "-smp", cfg.CPU)
	cmd.Args = append(cmd.Args, "-nic", fmt.Sprintf("%s,model=virtio-net-pci", cfg.Network.Type))
	cmd.Args = append(cmd.Args, "-cdrom", ciImage)
	cmd.Args = append(cmd.Args, "-hda", cfg.Disk)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func createCloudInitISO(cfg *config.Config) (path string, err error) {
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
					Sudo:         cfg.GuestSudo,
					LockPasswd:   false,
				},
			},
		},
	}
	if cfg.Network.Type == config.NetworkType_Bridge {
		ci.NetworkConfig = &cloudinit.NetworkConfig{
			Network: cloudinit.Network{
				Version: 2,
				Ethernets: map[string]cloudinit.Ethernet{
					// virtio-net-pci => ens3
					"ens3": {
						Addresses: []string{cfg.Network.Address},
						Routes: []cloudinit.Routes{
							{To: "default", Via: cfg.Network.DefaultGateway},
						},
						Nameservers: cloudinit.Nameservers{
							Addresses: cfg.Network.Nameservers,
							Search:    cfg.Search,
						},
					},
				},
			},
		}
	}
	return "./cloudinit.iso", ci.SaveTo("./cloudinit.iso")
}
