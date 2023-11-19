package cloudinit

// https://cloudinit.readthedocs.io/en/latest/reference/modules.html
type UserData struct {
	Hostname         string              `yaml:"hostname"`
	DisableRoot      bool                `yaml:"disable_root"`
	PreserveHostname bool                `yaml:"preserve_hostname"`
	Users            []UserCoinfig       `yaml:"users"`
	GrowPartition    GrowPartitionConfig `yaml:"growpart"`
	RunCMD           []string            `yaml:"runcmd,omitempty"`
	BootCMD          []string            `yaml:"bootcmd,omitempty"`
}

type UserCoinfig struct {
	Name              string   `yaml:"name"`
	Shell             string   `yaml:"shell"`
	Sudo              string   `yaml:"sudo"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	HashedPasswd      string   `yaml:"hashed_passwd,omitempty"`
	LockPasswd        bool     `yaml:"lock_passwd"`
}

type GrowPartitionConfig struct {
	Mode    GrowPartitionMode `yaml:"mode"`
	Devices []string          `yaml:"devices"`
}

type GrowPartitionMode string

const (
	GrowPartitionMode_Auto GrowPartitionMode = "auto"
)
