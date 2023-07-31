package utils

import "strings"

func PureText(oText string) string {
	//排版美观，只用第一段作为content
	lines := strings.Split(oText, "\n")
	if len(lines) > 0 {
		contentBeforeNewline := lines[0]
		return contentBeforeNewline
	}
	return ""
}
