//nolint:testpackage // Unexported parsers are exercised from this package.
package domain

import (
	"strings"
	"testing"

	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/stretchr/testify/require"
)

func TestValidateZoneDirectives(t *testing.T) {
	t.Parallel()
	err := validateZoneDirectives("$INCLUDE /etc/passwd\n")
	require.Error(t, err)
	require.Contains(t, err.Error(), "$INCLUDE")

	err = validateZoneDirectives("$GENERATE 1-10 $.example.com A 10.0.0.$\n")
	require.Error(t, err)
	require.Contains(t, err.Error(), "$GENERATE")

	require.NoError(t, validateZoneDirectives("; comment\nwww 3600 IN A 1.2.3.4\n"))
}

func TestParseImportBind(t *testing.T) {
	t.Parallel()
	zone := `example.com. 3600 IN SOA ns1.example.com. hostmaster.example.com. 1 7200 3600 1209600 3600
$TTL 7200
@  IN NS ns1.example.com.
www  IN A 192.0.2.1
www  IN A 192.0.2.2
mail IN MX 10 smtp.example.com.
txt  IN TXT "hello" "world"
`
	recs, err := parseImportBind(zone, "example.com")
	require.NoError(t, err)
	require.NotEmpty(t, recs)

	names := make([]string, 0, len(recs))
	for _, r := range recs {
		names = append(names, r.Name+":"+string(r.Type))
	}
	require.Contains(t, names, "www:"+string(domain.RecordTypeA))
	require.Contains(t, names, "mail:"+string(domain.RecordTypeMX))
	require.Contains(t, names, "txt:"+string(domain.RecordTypeTXT))

	var txt *domain.Record
	for _, r := range recs {
		if r.Type == domain.RecordTypeTXT && r.Name == "txt" {
			txt = r

			break
		}
	}
	require.NotNil(t, txt)
	require.Equal(t, "helloworld", txt.Data)
}

func TestParseImportBindSkipsSOAAndApexNS(t *testing.T) {
	t.Parallel()
	zone := `example.com. 3600 IN SOA ns1.example.com. hostmaster.example.com. 1 7200 3600 1209600 3600
@ 3600 IN NS ns1.example.com.
sub 3600 IN NS ns2.other.net.
`
	recs, err := parseImportBind(zone, "example.com")
	require.NoError(t, err)
	require.Len(t, recs, 1)
	require.Equal(t, "sub", recs[0].Name)
	require.Equal(t, domain.RecordTypeNS, recs[0].Type)
}

func TestParseImportBindUnsupportedRR(t *testing.T) {
	t.Parallel()
	zone := `example.com. 3600 IN SOA ns1.example.com. hostmaster.example.com. 1 7200 3600 1209600 3600
foo 3600 IN SSHFP 1 1 deadbeef
`
	_, err := parseImportBind(zone, "example.com")
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "unsupported")
}

func TestParseImportJSON(t *testing.T) {
	t.Parallel()
	raw := `{
  "records": [
    {"name": "www", "type": "A", "ttl": 600, "data": "203.0.113.1"},
    {"name": "@", "type": "MX", "ttl": 600, "data": "mail.example.net", "priority": 20}
  ]
}`
	recs, err := parseImportJSON(raw)
	require.NoError(t, err)
	require.Len(t, recs, 2)
	require.Equal(t, "www", recs[0].Name)
	require.Equal(t, uint32(600), recs[0].TTL)
	require.Empty(t, recs[1].Name)
	require.Equal(t, uint32(20), recs[1].Priority)
}

func TestParseImportJSONMXRequiresPriority(t *testing.T) {
	t.Parallel()
	raw := `{"records":[{"name":"x","type":"MX","ttl":60,"data":"mx.example.com"}]}`
	_, err := parseImportJSON(raw)
	require.Error(t, err)
	require.Contains(t, err.Error(), "priority")
}
