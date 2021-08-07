package main

import (
	"encoding/json"
	"math/rand"
	"sort"
	"time"

	"github.com/miekg/dns"
)

type Record struct {
	Name    string
	Rules   map[uint16]string
	Entries map[uint16]*RecordEntries
}

func (r *Record) UnmarshalJSON(data []byte) error {
	var d map[string]interface{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	r.Name = d["name"].(string)

	r.Rules = make(map[uint16]string)

	r.Rules[dns.TypeA] = d["A"].(string)
	r.Rules[dns.TypeAAAA] = d["AAAA"].(string)
	r.Rules[dns.TypeCNAME] = d["CNAME"].(string)
	r.Rules[dns.TypeNS] = d["NS"].(string)
	r.Rules[dns.TypeTXT] = d["TXT"].(string)

	r.Entries = make(map[uint16]*RecordEntries)

	f := func(dnsTypeStr string, dnsTypeUint16 uint16) error {
		e, err := json.Marshal(d["entries"].(map[string]interface{})[dnsTypeStr])
		if err != nil {
			return err
		}
		var re RecordEntries
		re.Init()
		err = json.Unmarshal(e, &re)
		if err != nil {
			return err
		}
		r.Entries[dnsTypeUint16] = &re
		return nil
	}

	err = f("A", dns.TypeA)
	if err != nil {
		return err
	}
	err = f("AAAA", dns.TypeAAAA)
	if err != nil {
		return err
	}
	err = f("CNAME", dns.TypeCNAME)
	if err != nil {
		return err
	}
	err = f("NS", dns.TypeNS)
	if err != nil {
		return err
	}
	err = f("TXT", dns.TypeTXT)
	if err != nil {
		return err
	}
	return nil
}

type RecordEntries []RecordEntry

func (r *RecordEntries) Init() *RecordEntries {
	i := 0
	for _, x := range *r {
		if x.Weight > 0 {
			(*r)[i] = x
			i++
		}
	}
	*r = (*r)[:i]
	for i := range *r {
		(*r)[i].w = (*r)[i].Weight
	}
	sort.Slice(*r, func(i, j int) bool {
		return (*r)[i].Order < (*r)[j].Order
	})
	return r
}

func (r *RecordEntries) GetBalanced() *RecordEntry {
	if len(*r) == 0 {
		return nil
	}
	for i := range *r {
		if (*r)[i].w > 0 {
			(*r)[i].w--
			ret := &(*r)[i]
			return ret
		}
	}
	// all zeros?
	for i := range *r {
		(*r)[i].w = (*r)[i].Weight
	}
	return r.GetBalanced()
}

func (r *RecordEntries) GetRandom() *RecordEntry {
	if len(*r) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	return &(*r)[rand.Intn(len((*r)))]
}

type RecordEntry struct {
	Type   string `json:"type"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
	Weight int    `json:"weight"`
	Order  int    `json:"order"`

	w int
}

type Zone struct {
	ID      int      `json:"id"`
	UserID  int      `json:"user_id"`
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	Records []Record `json:"records"`
}

type Zones []Zone

func (z *Zones) Init() {
	for i := range *z {
		v := &(*z)[i]
		for j := range v.Records {
			r := &v.Records[j]
			r.Entries[dns.TypeA].Init()
			r.Entries[dns.TypeAAAA].Init()
			r.Entries[dns.TypeCNAME].Init()
			r.Entries[dns.TypeNS].Init()
			r.Entries[dns.TypeTXT].Init()
		}
	}
}
