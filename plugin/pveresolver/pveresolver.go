package pveresolver

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type PVEResolver struct {
	next  plugin.Handler
	cache *cache
}

func (p PVEResolver) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	q := r.Question[0]

	if q.Qtype != dns.TypeA {
		return plugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
	}

	ip, ok := p.cache.get(q.Name)
	if !ok {
		return plugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = false
	m.Answer = []dns.RR{
		&dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			A: ip,
		},
	}
	_ = w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

func (p PVEResolver) Name() string {
	return "pveresolver"
}
