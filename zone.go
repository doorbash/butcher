package main

import (
	"math/rand"
	"sort"
	"time"
)

type Record struct {
	Name    string `json:"name"`
	A       string `json:"A"`
	AAAA    string `json:"AAAA"`
	CNAME   string `json:"CNAME"`
	NS      string `json:"NS"`
	TXT     string `json:"TXT"`
	Entries struct {
		A     RecordEntries `json:"A"`
		AAAA  RecordEntries `json:"AAAA"`
		CNAME RecordEntries `json:"CNAME"`
		NS    RecordEntries `json:"NS"`
		TXT   RecordEntries `json:"TXT"`
	} `json:"entries"`
	// entriesMutex struct {
	// 	A     sync.Mutex
	// 	AAAA  sync.Mutex
	// 	CNAME sync.Mutex
	// 	NS    sync.Mutex
	// 	TXT   sync.Mutex
	// }
}

func (r *Record) GetRuleByDnsType(dnsType string) string {
	switch dnsType {
	case "A":
		return r.A
	case "AAAA":
		return r.AAAA
	case "CNAME":
		return r.CNAME
	case "NS":
		return r.NS
	case "TXT":
		return r.TXT
	}
	return ""
}

func (r *Record) GetRecordEntriesByDnsType(dnsType string) *RecordEntries {
	switch dnsType {
	case "A":
		return &r.Entries.A
	case "AAAA":
		return &r.Entries.AAAA
	case "CNAME":
		return &r.Entries.CNAME
	case "NS":
		return &r.Entries.NS
	case "TXT":
		return &r.Entries.TXT
	}
	return nil
}

// func (r *Record) GetEntriesMutexByDnsType(dnsType string) *sync.Mutex {
// 	switch dnsType {
// 	case "A":
// 		return &r.entriesMutex.A
// 	case "AAAA":
// 		return &r.entriesMutex.AAAA
// 	case "CNAME":
// 		return &r.entriesMutex.CNAME
// 	case "NS":
// 		return &r.entriesMutex.NS
// 	case "TXT":
// 		return &r.entriesMutex.TXT
// 	}
// 	return nil
// }

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
			// fmt.Printf("GetBalanced(): %p\n", ret)
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
			r.Entries.A.Init()
			r.Entries.AAAA.Init()
			r.Entries.CNAME.Init()
			r.Entries.NS.Init()
			r.Entries.TXT.Init()
		}
	}
}
