package geodata

import (
	"github.com/clashdev/geodata/internal"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"strings"
	"sync/atomic"
)

type GeoIPSeeker struct {
	entries atomic.Pointer[map[string][]*net.IPNet]
}

func (s *GeoIPSeeker) ReadFrom(r io.Reader) (n int64, err error) {
	var data []byte
	if data, err = io.ReadAll(r); err == nil {
		err = s.UnmarshalBinary(data)
	}
	return
}

func (s *GeoIPSeeker) UnmarshalBinary(data []byte) (err error) {
	ipset := new(internal.GeoIPList)
	entries := make(map[string][]*net.IPNet)
	if err = proto.Unmarshal(data, ipset); err == nil {
		for _, entry := range ipset.Entry {
			networks := make([]*net.IPNet, len(entry.Cidr))
			for index, cidr := range entry.Cidr {
				networks[index] = &net.IPNet{
					IP:   cidr.Ip,
					Mask: net.CIDRMask(int(cidr.Prefix), len(cidr.Ip)*8),
				}
			}
			entries[strings.ToLower(entry.CountryCode)] = networks
		}
	}
	s.entries.Store(&entries)
	return
}

func (s *GeoIPSeeker) IsZero() bool {
	if s == nil {
		return true
	}
	entries := s.entries.Load()
	return entries == nil || *entries == nil
}

func (s *GeoIPSeeker) Query(name string) []*net.IPNet {
	if s.IsZero() {
		return nil
	}
	name = strings.TrimPrefix(name, "geoip:")
	entries := *s.entries.Load()
	networks, _ := entries[name]
	return networks
}
