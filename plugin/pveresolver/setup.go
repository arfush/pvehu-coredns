package pveresolver

import (
	"context"
	"crypto/tls"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/luthermonson/go-proxmox"
	"net"
	"net/http"
)

var logger = log.NewWithPlugin("pveresolver")

func init() {
	plugin.Register("pveresolver", func(c *caddy.Controller) error {
		c.Next()
		c.Next()
		_, network, err := net.ParseCIDR(c.Val())
		if err != nil {
			panic(err)
		}

		var (
			endpoint string
			token    string
			secret   string
		)
		for c.Next() {
			switch c.Val() {
			case "endpoint":
				c.Next()
				endpoint = c.Val()
			case "token":
				c.Next()
				token = c.Val()
			case "secret":
				c.Next()
				secret = c.Val()
			}
		}

		pve := proxmox.NewClient(endpoint, proxmox.WithAPIToken(token, secret), proxmox.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}))
		_, err = pve.Version(context.Background())
		if err != nil {
			panic(err)
		}

		ch := newCache(pve, network)
		ch.goCaching()

		dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
			return &PVEResolver{
				next:  next,
				cache: ch,
			}
		})

		return nil
	})
}
