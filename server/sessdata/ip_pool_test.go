package sessdata

import (
	"fmt"
	"net"
	"os"
	"path"
	"testing"

	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/dbdata"
	"github.com/stretchr/testify/assert"
)

func preData(tmpDir string) {
	base.Test()
	tmpDb := path.Join(tmpDir, "test.db")
	base.Cfg.DbType = "sqlite3"
	base.Cfg.DbSource = tmpDb
	base.Cfg.Ipv4CIDR = "192.168.3.0/24"
	base.Cfg.Ipv4Start = "192.168.3.1"
	base.Cfg.Ipv4End = "192.168.3.199"
	base.Cfg.MaxClient = 100
	base.Cfg.MaxUserClient = 3

	dbdata.Start()
	group := dbdata.Group{
		Name:      "group1",
		Bandwidth: 1000,
	}
	_ = dbdata.Add(&group)
	initIpPool()
}

func cleardata(tmpDir string) {
	_ = dbdata.Stop()
	tmpDb := path.Join(tmpDir, "test.db")
	os.Remove(tmpDb)
}

func TestIpPool(t *testing.T) {
	assert := assert.New(t)
	tmp := t.TempDir()
	preData(tmp)
	defer cleardata(tmp)

	var ip net.IP

	for i := 1; i <= 100; i++ {
		_ = AcquireIp("user", fmt.Sprintf("mac-%d", i), true)
	}
	ip = AcquireIp("user", "mac-new", true)
	assert.True(net.IPv4(192, 168, 3, 101).Equal(ip))
	for i := 102; i <= 199; i++ {
		ip = AcquireIp("user", fmt.Sprintf("mac-%d", i), true)
	}
	assert.True(net.IPv4(192, 168, 3, 199).Equal(ip))
	ip = AcquireIp("user", "mac-nil", true)
	assert.Nil(ip)

	ReleaseIp(net.IPv4(192, 168, 3, 88), "mac-88")
	ReleaseIp(net.IPv4(192, 168, 3, 188), "mac-188")
	// 从头循环获取可用ip
	ip = AcquireIp("user", "mac-188", true)
	t.Log("mac-188", ip)
	assert.True(net.IPv4(192, 168, 3, 188).Equal(ip))
}
