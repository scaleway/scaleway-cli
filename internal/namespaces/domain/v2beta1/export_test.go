package domain

// Test-only aliases so record_import_parse_test.go can live in package domain_test.
var (
	ParseImportBind   = parseImportBind
	ParseImportJSON   = parseImportJSON
	RelativeOwnerName = relativeOwnerName
)

const DNSImportDefaultTTL = dnsImportDefaultTTL
