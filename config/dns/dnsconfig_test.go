package dns_test

import (
	"fmt"
	"testing"

	"github.com/iwanhae/qemtainer/config/dns"
	"github.com/stretchr/testify/assert"
)

func TestLoadResolveConf(t *testing.T) {
	conf := dns.DnsReadConfig("/etc/resolv.conf")
	fmt.Println(conf.Servers)
	fmt.Println(conf.Search)
	assert.NotZero(t, conf.Servers)
}
