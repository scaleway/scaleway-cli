package secret

// ParsedSecretRef is an exported view of parsedSecretRef for use in tests.
type ParsedSecretRef struct {
	SecretID   string
	SecretName string
	SecretPath string
	Revision   string
	Field      string
}

// ParseSecretRef wraps the unexported parseSecretRef for external test access.
func ParseSecretRef(raw string) (*ParsedSecretRef, error) {
	ref, err := parseSecretRef(raw)
	if err != nil {
		return nil, err
	}

	return &ParsedSecretRef{
		SecretID:   ref.secretID,
		SecretName: ref.secretName,
		SecretPath: ref.secretPath,
		Revision:   ref.revision,
		Field:      ref.field,
	}, nil
}

var (
	RenderTemplate  = renderTemplate
	IsSecretUUID    = isSecretUUID
	WriteOutputFile = writeOutputFile
)
