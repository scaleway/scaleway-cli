package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
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
`) + "\n"

	records, err := parseImportBind(content, zone)
	require.NoError(t, err)
	require.Len(t, records, 4)

	byNameType := map[string]*domain.Record{}
	for _, r := range records {
		byNameType[r.Name+"/"+string(r.Type)] = r
	}

	assert.Equal(t, "1.2.3.4", byNameType["www/A"].Data)
	assert.Equal(t, uint32(3600), byNameType["www/A"].TTL)
	assert.Equal(t, "mx.example.com", byNameType["mail/MX"].Data)
	assert.Equal(t, uint32(10), byNameType["mail/MX"].Priority)
	assert.Equal(t, "hello", byNameType["txt/TXT"].Data)
	assert.Equal(t, "ns.sub.example.com", byNameType["sub/NS"].Data)
}

func TestParseImportBindCaseInsensitiveZone(t *testing.T) {
	t.Parallel()

	content := "WWW IN A 1.2.3.4\n"
	records, err := parseImportBind(content, "Example.COM")
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "WWW", records[0].Name)
	assert.Equal(t, "1.2.3.4", records[0].Data)
}

func TestParseImportBindRejectsInclude(t *testing.T) {
	t.Parallel()

	_, err := parseImportBind("$INCLUDE other.zone\n", "example.com")
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

	records, err := parseImportJSON(content)
	require.NoError(t, err)
	require.Len(t, records, 3)

	assert.Equal(t, "www", records[0].Name)
	assert.Equal(t, domain.RecordTypeA, records[0].Type)
	assert.Equal(t, uint32(300), records[0].TTL)

	assert.Equal(t, "", records[1].Name)
	assert.Equal(t, domain.RecordTypeMX, records[1].Type)
	assert.Equal(t, uint32(10), records[1].Priority)

	assert.Equal(t, "", records[2].Name)
	assert.Equal(t, domain.RecordTypeTXT, records[2].Type)
	assert.Equal(t, dnsImportDefaultTTL, records[2].TTL)
}

func TestParseImportJSONRequiresMXPriority(t *testing.T) {
	t.Parallel()

	_, err := parseImportJSON(`{"records":[{"name":"@","type":"MX","data":"mx.example.com"}]}`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "priority")
}

func TestParseImportJSONRejectsFQDNName(t *testing.T) {
	t.Parallel()

	_, err := parseImportJSON(`{"records":[{"name":"www.example.com.","type":"A","data":"1.2.3.4"}]}`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "FQDN")
}

func TestRelativeOwnerName(t *testing.T) {
	t.Parallel()

	rel, err := relativeOwnerName("www.example.com.", "example.com")
	require.NoError(t, err)
	assert.Equal(t, "www", rel)

	rel, err = relativeOwnerName("WWW.Example.COM.", "example.com")
	require.NoError(t, err)
	assert.Equal(t, "WWW", rel)

	rel, err = relativeOwnerName("example.com.", "example.com")
	require.NoError(t, err)
	assert.Equal(t, "", rel)

	_, err = relativeOwnerName("www.other.com.", "example.com")
	require.Error(t, err)
}
