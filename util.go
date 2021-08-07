package main

import "github.com/miekg/dns"

func GetDnsTypeString(t uint16) string {
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
