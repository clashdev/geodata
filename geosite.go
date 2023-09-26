//go:generate protoc --proto_path=internal --go_out=internal --go_opt=paths=source_relative routercommon.proto
package geodata

import (
	"github.com/clashdev/geodata/internal"
	"google.golang.org/protobuf/proto"
	"io"
	"slices"
	"strings"
	"sync/atomic"
)

type GeoSiteSeeker struct {
	entries atomic.Pointer[map[string][]*internal.Domain]
}

func (s *GeoSiteSeeker) ReadFrom(r io.Reader) (n int64, err error) {
	var data []byte
	if data, err = io.ReadAll(r); err == nil {
		err = s.UnmarshalBinary(data)
	}
	return
}

func (s *GeoSiteSeeker) UnmarshalBinary(data []byte) (err error) {
	sites := new(internal.GeoSiteList)
	entries := make(map[string][]*internal.Domain)
	if err = proto.Unmarshal(data, sites); err == nil {
		for _, site := range sites.Entry {
			entries[strings.ToLower(site.CountryCode)] = site.Domain
		}
	}
	s.entries.Store(&entries)
	return
}

func (s *GeoSiteSeeker) IsZero() bool {
	if s == nil {
		return true
	}
	entries := s.entries.Load()
	return entries == nil || *entries == nil
}

func (s *GeoSiteSeeker) Query(query string) []*DomainRecord {
	if s.IsZero() {
		return nil
	}
	entries := *s.entries.Load()
	query = strings.TrimPrefix(query, "geosite:")
	parts := strings.Split(query, "@")
	if domainList, ok := entries[parts[0]]; ok {
		return toDomainRecordSet(domainList, parts[1:])
	}
	return nil
}

func toDomainRecordSet(domains []*internal.Domain, attributes []string) []*DomainRecord {
	var records []*DomainRecord
	for _, domain := range domains {
		if !isGeoSiteMatch(domain, attributes) {
			continue
		}
		switch domain.Type {
		case internal.Domain_Plain:
			records = append(records, &DomainRecord{Type: DomainKeyword, Value: domain.Value})
		case internal.Domain_RootDomain:
			records = append(records, &DomainRecord{Type: DomainSuffix, Value: domain.Value})
		case internal.Domain_Regex:
			records = append(records, &DomainRecord{Type: DomainRegExp, Value: domain.Value})
		case internal.Domain_Full:
			records = append(records, &DomainRecord{Type: Domain, Value: domain.Value})
		}
	}
	return records
}

func isGeoSiteMatch(domain *internal.Domain, attributes []string) bool {
	if len(attributes) == 0 {
		return true
	}
	for _, attribute := range domain.Attribute {
		if slices.Contains(attributes, attribute.Key) {
			return true
		}
	}
	return false
}
