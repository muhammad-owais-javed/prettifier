package lib

import "strings"

func TrimSpaces(processedContent string) string {

	processedContent = strings.ReplaceAll(processedContent, "\r\n", "\n") // Handle Windows line endings first
	processedContent = strings.ReplaceAll(processedContent, "\r", "\n")
	processedContent = strings.ReplaceAll(processedContent, "\v", "\n")
	processedContent = strings.ReplaceAll(processedContent, "\f", "\n")

	for strings.Contains(processedContent, "\n\n\n") {
		processedContent = strings.ReplaceAll(processedContent, "\n\n\n", "\n\n")
	}

	if strings.HasPrefix(processedContent, "\n\n") {
		processedContent = strings.TrimPrefix(processedContent, "\n")
	}

	if strings.HasSuffix(processedContent, "\n\n") {
		processedContent = strings.TrimSuffix(processedContent, "\n")
	}

	processedContent = strings.Trim(processedContent, " \t")

	return processedContent
}
