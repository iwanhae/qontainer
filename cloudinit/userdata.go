package cloudinit

// https://cloudinit.readthedocs.io/en/latest/reference/modules.html
type UserData struct {
	Hostname         string              `yaml:"hostname" json:"hostname"`
	DisableRoot      bool                `yaml:"disable_root" json:"disable_root"`
	PreserveHostname bool                `yaml:"preserve_hostname" json:"preserve_hostname"`
	Users            []UserCoinfig       `yaml:"users" json:"users"`
	GrowPartition    GrowPartitionConfig `yaml:"growpart" json:"growpart"`
}

type UserCoinfig struct {
	Name              string   `yaml:"name" json:"name"`
	Sudo              string   `yaml:"sudo" json:"sudo"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty" json:"ssh_authorized_keys,omitempty"`
	HashedPasswd      string   `yaml:"hashed_passwd" json:"hashed_passwd"`
	LockPasswd        bool     `yaml:"lock_passwd" json:"lock_passwd"`
}

type GrowPartitionConfig struct {
	Mode    GrowPartitionMode `yaml:"mode" json:"mode"`
	Devices []string          `yaml:"devices" json:"devices"`
}

type GrowPartitionMode string

const (
	GrowPartitionMode_Auto GrowPartitionMode = "auto"
)
