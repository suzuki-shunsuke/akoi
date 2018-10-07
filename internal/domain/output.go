package domain

import (
	"encoding/json"
)

// String converts result into a string.
func (result *Result) String(format string) string {
	switch format {
	case "ansible":
		b, err := json.Marshal(result)
		if err != nil {
			return `{"changed": true}`
		}
		return string(b)
	default:
		return ""
	}
}
