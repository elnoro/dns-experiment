package coredns_integration

import (
	"context"
	"fmt"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/request"
	"net"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
)

const name = "foxylock"

func init() {
	plugin.Register(name, setup)
}

func setup(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		// If there was another token, return an error, because we don't have any configuration.
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error(name, c.ArgErr())
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return FoxyLock{Next: next}
	})

	// All OK, return a nil error.
	return nil
}

var log = clog.NewWithPlugin(name)

type FoxyLock struct {
	Next plugin.Handler
}

func (e FoxyLock) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	blocked, err := globalBlockedHostsList.Get(state.QName())
	if err != nil {
		// falling through to other plugins
		log.Errorf("cannot access blocked hosts list, %v", err)
	} else if blocked {
		a := new(dns.Msg)
		a.SetReply(r)
		a.Authoritative = true
		rr := new(dns.A)
		rr.Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rr.A = net.ParseIP("127.0.0.1").To4()
		a.Answer = append(a.Answer, rr)
		log.Info(rr.String())

		err := w.WriteMsg(a)

		return 0, fmt.Errorf("failed to write response for blocked host: %w", err)
	}

	return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
}

func (e FoxyLock) Ready() bool { return true }

func (e FoxyLock) Name() string { return name }
