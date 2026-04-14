package domain

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/miekg/dns"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

const dnsImportDefaultTTL = uint32(3600)

// jsonImportFile is the expected shape for format=json imports.
type jsonImportFile struct {
	Records []jsonImportRecord `json:"records"`
}

type jsonImportRecord struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	TTL      uint32  `json:"ttl"`
	Data     string  `json:"data"`
	Priority *uint32 `json:"priority,omitempty"`
}

func validateZoneDirectives(content string) error {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	const maxScan = 32 * 1024 * 1024 // 32 MiB lines are unexpected; avoid unbounded buffer
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxScan)

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		upper := strings.ToUpper(line)
		if strings.HasPrefix(upper, "$INCLUDE") {
			return fmt.Errorf(
				"line %d: $INCLUDE is not supported; expand includes before importing",
				lineNum,
			)
		}
		if strings.HasPrefix(upper, "$GENERATE") {
			return fmt.Errorf(
				"line %d: $GENERATE is not supported; expand generated records before importing",
				lineNum,
			)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read zone file: %w", err)
	}

	return nil
}

func parseImportJSON(content string) ([]*domain.Record, error) {
	var doc jsonImportFile
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}
	if len(doc.Records) == 0 {
		return nil, errors.New(`no records found in JSON (expected a top-level "records" array)`)
	}
	records := make([]*domain.Record, 0, len(doc.Records))
	for i, r := range doc.Records {
		rec, err := jsonRecordToDomain(r, i)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func jsonRecordToDomain(r jsonImportRecord, index int) (*domain.Record, error) {
	typ := strings.TrimSpace(strings.ToUpper(r.Type))
	if typ == "" {
		return nil, fmt.Errorf("records[%d]: missing type", index)
	}
	if !isAllowedImportRecordType(typ) {
		return nil, fmt.Errorf("records[%d]: unsupported type %q", index, r.Type)
	}
	data := strings.TrimSpace(r.Data)
	if data == "" {
		return nil, fmt.Errorf("records[%d]: missing data", index)
	}
	ttl := r.TTL
	if ttl == 0 {
		ttl = dnsImportDefaultTTL
	}
	name := strings.TrimSpace(r.Name)
	if name == "@" {
		name = ""
	}
	if err := validateRecordOwnerName(name); err != nil {
		return nil, fmt.Errorf("records[%d]: %w", index, err)
	}
	rt := domain.RecordType(typ)
	rec := &domain.Record{
		Data:     data,
		Name:     name,
		TTL:      ttl,
		Type:     rt,
		Priority: 0,
	}
	if r.Priority != nil {
		rec.Priority = *r.Priority
	}
	if typ == "MX" && r.Priority == nil {
		return nil, fmt.Errorf("records[%d]: MX records require priority", index)
	}

	return rec, nil
}

func validateRecordOwnerName(name string) error {
	if name == "" || name == "@" {
		return nil
	}
	if strings.Contains(name, "..") || strings.ContainsAny(name, " \t") {
		return fmt.Errorf("invalid owner name %q", name)
	}
	// Reject absolute FQDNs: we expect short names relative to the zone (like the rest of scw dns record).
	if dns.IsFqdn(name) {
		return fmt.Errorf(
			"owner name %q looks like an FQDN; use a relative name (e.g. www) or @ for apex",
			name,
		)
	}

	return nil
}

func isAllowedImportRecordType(typ string) bool {
	return slices.Contains(domainTypes, typ)
}

func parseImportBind(content, dnsZone string) ([]*domain.Record, error) {
	if err := validateZoneDirectives(content); err != nil {
		return nil, err
	}
	origin := dns.Fqdn(dnsZone)
	zp := dns.NewZoneParser(strings.NewReader(content), origin, "")
	var records []*domain.Record
	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		recs, err := dnsRRToRecords(rr, dnsZone)
		if err != nil {
			return nil, err
		}
		records = append(records, recs...)
	}
	if err := zp.Err(); err != nil {
		return nil, fmt.Errorf("parse BIND zone: %w", err)
	}

	return records, nil
}

func dnsRRToRecords(rr dns.RR, dnsZone string) ([]*domain.Record, error) {
	switch hdr := rr.Header(); hdr.Rrtype {
	case dns.TypeSOA:
		return nil, nil
	case dns.TypeNS:
		name, err := relativeOwnerName(hdr.Name, dnsZone)
		if err != nil {
			return nil, err
		}
		if name == "" {
			// Apex NS are managed by Scaleway for zones with default name servers.
			return nil, nil
		}
		ns, ok := rr.(*dns.NS)
		if !ok {
			return nil, errors.New("internal error: expected NS record")
		}

		return []*domain.Record{{
			Data:     targetToData(ns.Ns),
			Name:     name,
			TTL:      ttlOrDefault(hdr.Ttl),
			Type:     domain.RecordTypeNS,
			Priority: 0,
		}}, nil
	default:
		rec, err := dnsRRToRecord(rr, dnsZone)
		if err != nil {
			return nil, err
		}
		if rec == nil {
			return nil, nil
		}

		return []*domain.Record{rec}, nil
	}
}

func dnsRRToRecord(rr dns.RR, dnsZone string) (*domain.Record, error) {
	hdr := rr.Header()
	name, err := relativeOwnerName(hdr.Name, dnsZone)
	if err != nil {
		return nil, err
	}
	ttl := ttlOrDefault(hdr.Ttl)

	switch v := rr.(type) {
	case *dns.A:
		return &domain.Record{
			Data: v.A.String(),
			Name: name, TTL: ttl, Type: domain.RecordTypeA,
		}, nil
	case *dns.AAAA:
		return &domain.Record{
			Data: v.AAAA.String(),
			Name: name, TTL: ttl, Type: domain.RecordTypeAAAA,
		}, nil
	case *dns.CNAME:
		return &domain.Record{
			Data: targetToData(v.Target),
			Name: name, TTL: ttl, Type: domain.RecordTypeCNAME,
		}, nil
	case *dns.TXT:
		return &domain.Record{
			Data: strings.Join(v.Txt, ""),
			Name: name, TTL: ttl, Type: domain.RecordTypeTXT,
		}, nil
	case *dns.MX:
		return &domain.Record{
			Data:     targetToData(v.Mx),
			Name:     name,
			TTL:      ttl,
			Type:     domain.RecordTypeMX,
			Priority: uint32(v.Preference),
		}, nil
	case *dns.PTR:
		return &domain.Record{
			Data: targetToData(v.Ptr),
			Name: name, TTL: ttl, Type: domain.RecordTypePTR,
		}, nil
	case *dns.SRV:
		return &domain.Record{
			Data: fmt.Sprintf("%d %d %d %s", v.Priority, v.Weight, v.Port, targetToData(v.Target)),
			Name: name, TTL: ttl, Type: domain.RecordTypeSRV,
		}, nil
	case *dns.CAA:
		return &domain.Record{
			Data: fmt.Sprintf("%d %s %s", v.Flag, v.Tag, v.Value),
			Name: name, TTL: ttl, Type: domain.RecordTypeCAA,
		}, nil
	default:
		typeName := dns.TypeToString[hdr.Rrtype]
		if typeName == "" {
			typeName = fmt.Sprintf("TYPE%d", hdr.Rrtype)
		}

		return nil, fmt.Errorf(
			"unsupported record type %s for %s in BIND format; use format=json with a raw data field instead",
			typeName, hdr.Name,
		)
	}
}

func ttlOrDefault(ttl uint32) uint32 {
	if ttl == 0 {
		return dnsImportDefaultTTL
	}

	return ttl
}

func targetToData(target string) string {
	return strings.TrimSuffix(target, ".")
}

func relativeOwnerName(ownerFQN, zone string) (string, error) {
	owner := strings.TrimSuffix(ownerFQN, ".")
	z := strings.TrimSuffix(dns.Fqdn(zone), ".")
	if owner == z {
		return "", nil
	}
	suf := "." + z
	rel, ok := strings.CutSuffix(owner, suf)
	if !ok {
		return "", fmt.Errorf("owner %q is not under DNS zone %q", ownerFQN, zone)
	}

	return rel, nil
}
