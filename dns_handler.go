package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
	"golang.org/x/net/publicsuffix"
)

type DNSHandler struct {
	ListenAddr string
	db         *Database
}

func getDnsTypeString(t uint16) string {
	switch t {
	case dns.TypeA:
		return "A"
	case dns.TypeAAAA:
		return "AAAA"
	case dns.TypeCNAME:
		return "CNAME"
	case dns.TypeNS:
		return "NS"
	case dns.TypeTXT:
		return "TXT"
	case dns.TypeSOA:
		return "SOA"
	}
	return ""
}

func (d *DNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	if r.Opcode != dns.OpcodeQuery {
		log.Println("bad opcode: ", r)
		w.WriteMsg(m)
		return
	}

	if len(r.Question) == 0 {
		// send null
		log.Println("no questions: ", r)
		w.WriteMsg(m)
		return
	}
	q := r.Question[0]
	query := q.Name
	query = strings.TrimSuffix(query, ".")

	fld, err := publicsuffix.EffectiveTLDPlusOne(query)

	if err != nil {
		log.Println("error: ", err)
		w.WriteMsg(m)
		return
	}

	// fmt.Println("fld is:", fld)

	zone := d.db.GetMemoryZone(fld)
	if zone == nil {
		// send null
		log.Println("no zone found for: ", r)
		w.WriteMsg(m)
		return
	}
	// fmt.Println(zone)

	ind := strings.LastIndex(query, fld)
	var subdomain string
	if ind > 0 {
		subdomain = query[:ind-1]
	} else {
		subdomain = "@"
	}

	// fmt.Println("subdomain is", subdomain)

	dnsType := getDnsTypeString(q.Qtype)

	// fmt.Println("dnsType:", dnsType)
	// fmt.Println(q.Qtype)

	switch q.Qtype {
	case dns.TypeSOA:
		ns := d.db.FindMemoryRecord(zone, "NS", "")
		email := strings.Replace(strings.ToLower(zone.Owner), "@", ".", 1)
		if len(ns) > 0 {
			v := ns[0]
			rr, err := dns.NewRR(fmt.Sprintf("%s 60 IN SOA %s %s 1 7200 3600 1209600 60", q.Name, v.Value, email))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			} else {
				log.Println(err)
			}
		}
	case dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeNS, dns.TypeTXT:
		// fmt.Println(zone, dnsType, subdomain)
		reList := d.db.FindMemoryRecord(zone, dnsType, subdomain)
		for _, v := range reList {
			rr, err := dns.NewRR(fmt.Sprintf("%s %d IN %s %s", q.Name, v.TTL, dnsType, v.Value))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
	w.WriteMsg(m)
}

func NewDNSHandler(addr string, db *Database) *DNSHandler {
	return &DNSHandler{
		ListenAddr: addr,
		db:         db,
	}
}

func (d *DNSHandler) ListenAndServe() error {
	server := &dns.Server{
		Addr:    d.ListenAddr,
		Net:     "udp",
		Handler: d,
	}
	err := server.ListenAndServe()
	defer server.Shutdown()
	return err
}
