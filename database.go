package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Database struct {
	zones Zones
}

func (d *Database) GetAllZones() *Zones {
	return &d.zones
}

func (d *Database) GetMemoryZone(fld string) *Zone {
	for _, v := range d.zones {
		fmt.Println(v.Name)
		if v.Name == fld {
			return &v
		}
	}
	return nil
}

func (d *Database) FindMemoryRecord(zone *Zone, dnsType, subdomain string) []RecordEntry {
	if subdomain == "" {
		subdomain = "@"
	}
	for i := range zone.Records {
		record := &zone.Records[i]
		if record.Name == subdomain {
			rule := record.GetRuleByDnsType(dnsType)
			switch rule {
			case "loadbalance":
				// mutex := *record.GetEntriesMutexByDnsType(dnsType)
				// mutex.Lock()
				r := record.GetRecordEntriesByDnsType(dnsType).GetBalanced()
				if r == nil {
					return nil
				}
				ret := []RecordEntry{
					*r,
				}
				// mutex.Unlock()
				return ret
			case "random":
				r := record.GetRecordEntriesByDnsType(dnsType).GetRandom()
				if r == nil {
					return nil
				}
				return []RecordEntry{
					*r,
				}
			case "all":
				return *record.GetRecordEntriesByDnsType(dnsType)
			}
		}
	}
	return nil
}

func NewDatabase(config_path string) (*Database, error) {
	d := Database{}
	file, err := os.Open(config_path)
	defer file.Close()

	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, &d.zones)
	if err != nil {
		return nil, err
	}
	d.zones.Init()

	return &d, nil
}
