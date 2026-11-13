package domain_test

import (
	"strings"
	"testing"

	domain "github.com/scaleway/scaleway-cli/v2/internal/namespaces/domain/v2beta1"
	domainsdk "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseImportBind(t *testing.T) {
	t.Parallel()

	const zone = "example.com"
	content := strings.TrimSpace(`
$TTL 3600
@       IN SOA ns1.example.com. hostmaster.example.com. (1 7200 3600 1209600 3600)
@       IN NS  ns1.example.com.
www     IN A   1.2.3.4
mail    IN MX  10 mx.example.com.
txt     IN TXT "hello"
sub     IN NS  ns.sub.example.com.
_sip._tcp IN SRV 10 20 443 sip.example.com.
`) + "\n"

	records, err := domain.ParseImportBind(content, zone)
	require.NoError(t, err)
	require.Len(t, records, 5)

	byNameType := map[string]*domainsdk.Record{}
	for _, r := range records {
		byNameType[r.Name+"/"+string(r.Type)] = r
	}

	assert.Equal(t, "1.2.3.4", byNameType["www/A"].Data)
	assert.Equal(t, uint32(3600), byNameType["www/A"].TTL)
	assert.Equal(t, "mx.example.com", byNameType["mail/MX"].Data)
	assert.Equal(t, uint32(10), byNameType["mail/MX"].Priority)
	assert.Equal(t, "hello", byNameType["txt/TXT"].Data)
	assert.Equal(t, "ns.sub.example.com", byNameType["sub/NS"].Data)
	assert.Equal(t, "20 443 sip.example.com", byNameType["_sip._tcp/SRV"].Data)
	assert.Equal(t, uint32(10), byNameType["_sip._tcp/SRV"].Priority)
}

func TestParseImportBindCaseInsensitiveZone(t *testing.T) {
	t.Parallel()

	content := "WWW IN A 1.2.3.4\n"
	records, err := domain.ParseImportBind(content, "Example.COM")
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "WWW", records[0].Name)
	assert.Equal(t, "1.2.3.4", records[0].Data)
}

func TestParseImportBindRejectsInclude(t *testing.T) {
	t.Parallel()

	_, err := domain.ParseImportBind("$INCLUDE other.zone\n", "example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "$INCLUDE")
}

func TestParseImportJSON(t *testing.T) {
	t.Parallel()

	content := `{
  "records": [
    {"name": "www", "type": "A", "ttl": 300, "data": "1.2.3.4"},
    {"name": "@", "type": "MX", "ttl": 3600, "data": "mx.example.com", "priority": 10},
    {"name": "", "type": "TXT", "data": "v=spf1 -all"}
  ]
}`

	records, err := domain.ParseImportJSON(content)
	require.NoError(t, err)
	require.Len(t, records, 3)

	assert.Equal(t, "www", records[0].Name)
	assert.Equal(t, domainsdk.RecordTypeA, records[0].Type)
	assert.Equal(t, uint32(300), records[0].TTL)

	assert.Empty(t, records[1].Name)
	assert.Equal(t, domainsdk.RecordTypeMX, records[1].Type)
	assert.Equal(t, uint32(10), records[1].Priority)

	assert.Empty(t, records[2].Name)
	assert.Equal(t, domainsdk.RecordTypeTXT, records[2].Type)
	assert.Equal(t, domain.DNSImportDefaultTTL, records[2].TTL)
}

func TestParseImportJSONRequiresMXPriority(t *testing.T) {
	t.Parallel()

	_, err := domain.ParseImportJSON(
		`{"records":[{"name":"@","type":"MX","data":"mx.example.com"}]}`,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "priority")
}

func TestParseImportJSONRejectsFQDNName(t *testing.T) {
	t.Parallel()

	_, err := domain.ParseImportJSON(
		`{"records":[{"name":"www.example.com.","type":"A","data":"1.2.3.4"}]}`,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "FQDN")
}

func TestRelativeOwnerName(t *testing.T) {
	t.Parallel()

	rel, err := domain.RelativeOwnerName("www.example.com.", "example.com")
	require.NoError(t, err)
	assert.Equal(t, "www", rel)

	rel, err = domain.RelativeOwnerName("WWW.Example.COM.", "example.com")
	require.NoError(t, err)
	assert.Equal(t, "WWW", rel)

	rel, err = domain.RelativeOwnerName("example.com.", "example.com")
	require.NoError(t, err)
	assert.Empty(t, rel)

	_, err = domain.RelativeOwnerName("www.other.com.", "example.com")
	require.Error(t, err)
}
