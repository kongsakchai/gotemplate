package validator

import "strings"

type errorMap map[string]string

func (em errorMap) Error() string {
	var sb strings.Builder
	for field, msg := range em {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(field)
		sb.WriteString(": ")
		sb.WriteString(msg)
	}
	return sb.String()
}
