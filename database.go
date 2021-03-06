package main

import (
	"encoding/json"
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
	for i := range d.zones {
		v := &d.zones[i]
		if v.Name == fld {
			return v
		}
	}
	return nil
}

func (d *Database) FindMemoryRecord(zone *Zone, dnsType uint16, subdomain string) []RecordEntry {
	if subdomain == "" {
		subdomain = "@"
	}
	for i := range zone.Records {
		record := &zone.Records[i]
		if record.Name == subdomain {
			rule := record.Rules[dnsType]
			switch rule {
			case "loadbalance":
				r := record.Entries[dnsType].GetBalanced()
				if r == nil {
					return nil
				}
				return []RecordEntry{
					*r,
				}
			case "random":
				r := record.Entries[dnsType].GetRandom()
				if r == nil {
					return nil
				}
				return []RecordEntry{
					*r,
				}
			case "all":
				return *record.Entries[dnsType]
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
