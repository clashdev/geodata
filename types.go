package geodata

import "strconv"

type DomainRecordType string

const (
	Domain        DomainRecordType = "DOMAIN"
	DomainSuffix  DomainRecordType = "DOMAIN-SUFFIX"
	DomainKeyword DomainRecordType = "DOMAIN-KEYWORD"
	DomainRegExp  DomainRecordType = "DOMAIN-REGEXP"
)

type DomainRecord struct {
	Type  DomainRecordType
	Value string
}

func (r DomainRecord) String() string {
	return string(r.Type) + "," + strconv.Quote(r.Value)
}
