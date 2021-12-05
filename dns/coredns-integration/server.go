package coredns_integration

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	_ "github.com/coredns/coredns/plugin/file"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"
	"github.com/elnoro/foxylock/m/v2/dns"

	_ "github.com/coredns/example"
)

func NewCoreDns(db BlockedHostsList) dns.Dns {
	// not a good idea, but I'm not sure how to pass dependencies into CoreDNS plugins
	globalBlockedHostsList = NewAdapter(db)

	return &CoreDnsPlugin{}
}

type CoreDnsPlugin struct{}

func (c *CoreDnsPlugin) Start() error {
	c.Init()
	coremain.Run()

	return nil
}

func (c *CoreDnsPlugin) Init() {
	dnsserver.Directives = []string{"foxylock", "log", "forward", "startup", "shutdown"}
}
