package config

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/schollz/progressbar/v3"
)

const (
	DefaultDiskPath = "/data/disk.img"
)

type VM struct {
	CPU      string `env:"VM_CPU" envDefault:"1"`
	Memory   string `env:"VM_MEMORY" envDefault:"1G"`
	Disk     string `env:"VM_DISK" envDefault:"/data/disk.img"`
	DiskSize string `env:"VM_DISK_SIZE"`
}

func (c *VM) AutoComplete() error {
	if isValidURL(c.Disk) {
		if _, err := os.Stat(DefaultDiskPath); err != nil {
			if err := downloadFile(c.Disk, DefaultDiskPath); err != nil {
				return fmt.Errorf("fail to download file: %w", err)
			}
		} else {
			log.Println("Skip downloads disk image")
		}
		c.Disk = DefaultDiskPath
	}
	if c.DiskSize != "" {
		cmd := exec.Command("qemu-img")
		cmd.Args = append(cmd.Args, "resize", c.Disk, c.DiskSize)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("fail to resize disk image: %q", err)
		}

		cmd = exec.Command("qemu-img")
		cmd.Args = []string{"info", c.Disk}
		cmd.Run()
	}
	return nil
}

func downloadFile(url string, destination string) error {
	fmt.Printf("Downloads file from %q to %q\n", url, destination)
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data from the URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Get the content length for the progress bar
	fileSize := resp.ContentLength

	// Create a progress bar

	bar := progressbar.DefaultBytes(fileSize, "downloading")

	// Create a multi writer to write to both the file and the progress bar
	writer := io.MultiWriter(out, bar)

	// Copy the content to the file and the progress bar
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("\nDownload complete!")

	return nil
}

func isValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}

	u, err := url.Parse(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
