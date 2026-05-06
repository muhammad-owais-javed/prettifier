package lib

import (
	"regexp"
	"strings"
	"time"
)

func AirportCodesAndDateTimeParsing(textContent string, iataMap map[string]AirportInfo, icaoMap map[string]AirportInfo) (string, string) {

	plainText := textContent
	colorText := textContent

	for code, info := range icaoMap {
		match := "*##" + code
		plainText = strings.ReplaceAll(plainText, match, info.City)
		colorText = strings.ReplaceAll(colorText, match, ColorGreen+info.City+ColorReset)
	}

	for code, info := range iataMap {
		match := "*#" + code
		plainText = strings.ReplaceAll(plainText, match, info.City)
		colorText = strings.ReplaceAll(colorText, match, ColorGreen+info.City+ColorReset)
	}

	for code, info := range icaoMap {
		match := "##" + code
		plainText = strings.ReplaceAll(plainText, match, info.Name)
		colorText = strings.ReplaceAll(colorText, match, ColorGreen+info.Name+ColorReset)
	}

	for code, info := range iataMap {
		pattern := "([^#]|^)#" + regexp.QuoteMeta(code) + "\\b"
		reg := regexp.MustCompile(pattern)
		plainText = reg.ReplaceAllString(plainText, "$1"+info.Name)
		colorText = reg.ReplaceAllString(colorText, "$1"+ColorGreen+info.Name+ColorReset)
	}

	reg := regexp.MustCompile(`\b(D|T12|T24)\((.*?)\)`)
	matches := reg.FindAllStringSubmatch(plainText, -1)

	for _, match := range matches {

		defaultDate := match[0]
		dateTag := match[1]
		isoDate := match[2]

		idx := strings.Index(plainText, defaultDate)
		if idx != -1 && idx+len(defaultDate) < len(plainText) {
			nextChar := plainText[idx+len(defaultDate)]
			if nextChar == 'R' || nextChar == 'S' {
				continue
			}
		}

		if !strings.Contains(isoDate, "T") || !strings.Contains(isoDate, "+") && !strings.Contains(isoDate, "-") && !strings.Contains(isoDate, "Z") {
			continue
		}

		t, err := time.Parse("2006-01-02T15:04Z07:00", isoDate)
		if err != nil {
			continue
		}

		var formatResult string

		switch dateTag {
		case "D":
			formatResult = t.Format("02 Jan 2006")
		case "T12":
			formatResult = t.Format("03:04PM (-07:00)")
		case "T24":
			formatResult = t.Format("15:04 (-07:00)")
		}

		if formatResult != "" {
			plainText = strings.Replace(plainText, defaultDate, formatResult, 1)
			colorText = strings.Replace(colorText, defaultDate, ColorYellow+formatResult+ColorReset, 1)
		}
	}
	return plainText, colorText
}
