package formatter

import "strings"

func FormatSessionName(name string) string {
	sessionName := strings.ReplaceAll(name, ".", "_")
	return sessionName
}
