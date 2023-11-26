package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"reflect"
)

const QemuExecutable = "qemu-system-x86_64"

type Config struct {
	VM
	Network

	GuestShell             string   `env:"GUEST_SHELL" envDefault:"/bin/bash"`
	GuestHostname          string   `env:"GUEST_HOSTNAME"`
	GuestUsername          string   `env:"GUEST_USERNAME" envDefault:"deploy"`
	GuestPassword          string   `env:"GUEST_PASSWORD" envDefault:"$6$rounds=4096$KUjo2cumnYaz0fmk$EsoVV1xP/FXIkv5mm4V26CR3qJrDZhs3Rga8OfBKNBUSsmCM7OHouHMHHz8lApGsD835DqpFvAgqJv1Hq5J.k0"`
	GuestSSHAuthorizedKeys []string `env:"GUEST_SSH_AUTHORIZED_KEYS"`
	GuestSudo              string   `env:"GUEST_SUDO" envDefault:"ALL=(ALL) NOPASSWD:ALL"`
	GuestUserScript        string   `env:"GUEST_USERSCRIPT_BASE64"`

	QemuExecutable string `env:"QEMU_EXECUTABLE" envDefault:"qemu-system-x86_64"`
}

func (c *Config) AutoComplete() error {
	if err := c.VM.AutoComplete(); err != nil {
		return err
	}
	if err := c.Network.AutoComplete(); err != nil {
		return err
	}
	if hostname, err := os.Hostname(); err != nil {
		return fmt.Errorf("fail to get hostname: %w", err)
	} else {
		c.GuestHostname = hostname
	}
	if path, err := exec.LookPath(QemuExecutable); err != nil {
		return fmt.Errorf("fail to get path of %q: %w", QemuExecutable, err)
	} else {
		c.QemuExecutable = path
	}
	if c.GuestUserScript != "" {
		if b, err := base64.StdEncoding.DecodeString(c.GuestUserScript); err != nil {
			return fmt.Errorf("fail to decode base64 encoded user script %q: %w", c.GuestUserScript, err)
		} else {
			c.GuestUserScript = string(b)
		}
	}

	return nil
}

func (cfg Config) Print() {
	fmt.Println("-----------CONFIG-----------")
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)
	for _, f := range reflect.VisibleFields(t) {
		if f.Tag.Get("env") != "" {
			fmt.Printf("%s=%q\n", f.Tag.Get("env"), v.FieldByName(f.Name))
		}
	}
}
