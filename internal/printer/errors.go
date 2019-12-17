package printer

import (
	"encoding/json"
)

func standardErrorJSON(data interface{}) (json.RawMessage, bool) {
	dataErr, ok := data.(interface {
		GetRawBody() json.RawMessage
	})
	if !ok {
		return nil, false
	}

	return dataErr.GetRawBody(), true
}
