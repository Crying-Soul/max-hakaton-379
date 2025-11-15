package handler

import (
	"strings"
	"unicode"
)

var verificationStatusTranslations = map[string]string{
	"pending":  "На проверке",
	"approved": "Одобрена",
	"rejected": "Отклонена",
}

func translateVerificationStatus(status string) string {
	if status == "" {
		return "Неизвестно"
	}
	normalized := strings.ToLower(status)
	if translated, ok := verificationStatusTranslations[normalized]; ok {
		return translated
	}
	return capitalize(status)
}

func translateVerificationStatusPtr(status *string) string {
	if status == nil {
		return "Неизвестно"
	}
	return translateVerificationStatus(*status)
}

func capitalize(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return value
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
