package cloudinit

import (
	"bytes"
	"fmt"
	"os"

	"github.com/kdomanski/iso9660"
	"gopkg.in/yaml.v3"
)

type CloudConfig struct {
	UserData      UserData
	NetworkConfig *NetworkConfig
	MetaData      *MetaData
}

func (cloudConfig *CloudConfig) SaveTo(path string) error {
	writer, err := iso9660.NewWriter()
	if err != nil {
		return fmt.Errorf("failed to create ISO writer: %w", err)
	}
	defer writer.Cleanup()

	// meta-data
	if err := writer.AddFile(bytes.NewBuffer([]byte{}), "meta-data"); err != nil {
		return fmt.Errorf("failed to add meta-data to ISO: %w", err)
	}

	// user-data
	b, err := yaml.Marshal(cloudConfig.UserData)
	if err != nil {
		return fmt.Errorf("failed to marchal cloud config to yaml: %w", err)
	}
	if err := writer.AddFile(bytes.NewBuffer(append([]byte("#cloud-config\n"), b...)), "user-data"); err != nil {
		return fmt.Errorf("failed to add user-data to ISO: %w", err)
	}

	// meta-data
	if cloudConfig.MetaData != nil {
		b, err = yaml.Marshal(cloudConfig.MetaData)
		if err != nil {
			return fmt.Errorf("failed to marchal cloud config to yaml: %w", err)
		}
		if err := writer.AddFile(bytes.NewBuffer(b), "meta-data"); err != nil {
			return fmt.Errorf("failed to add meta-data to ISO: %w", err)
		}
	}

	// network-config
	if cloudConfig.NetworkConfig != nil {
		b, err = yaml.Marshal(cloudConfig.NetworkConfig)
		if err != nil {
			return fmt.Errorf("failed to marchal cloud config to yaml: %w", err)
		}
		if err := writer.AddFile(bytes.NewBuffer(b), "network-config"); err != nil {
			return fmt.Errorf("failed to add network-config to ISO: %w", err)
		}
	}

	// Write to external file
	outputFile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer outputFile.Close()
	if err := writer.WriteTo(outputFile, "CIDATA"); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
