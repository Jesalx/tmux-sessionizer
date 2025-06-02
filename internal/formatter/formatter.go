package formatter

import "strings"

func FormatSessionName(name string) string {
	sessionName := strings.TrimSpace(name)
	sessionName = strings.Trim(sessionName, "-_.")
	sessionName = strings.ReplaceAll(sessionName, ".", "_")
	return sessionName
}
