package lib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func DateTimeParsing(textContent string, iataMap map[string]AirportInfo, icaoMap map[string]AirportInfo) (string, string) {

	plainText := textContent
	colorText := textContent

	for code, info := range iataMap {
		match := "*#" + code
		plainText = strings.ReplaceAll(plainText, match, info.City)
		colorText = strings.ReplaceAll(colorText, match, ColorGreen+info.City+ColorReset)
	}

	for code, info := range icaoMap {
		match := "*##" + code
		plainText = strings.ReplaceAll(plainText, match, info.Name)
		colorText = strings.ReplaceAll(colorText, match, ColorGreen+info.Name+ColorReset)
	}

	for code, info := range iataMap {
		match := "#" + code
		plainText = strings.ReplaceAll(plainText, match, info.Name)
		colorText = strings.ReplaceAll(colorText, match, ColorGreen+info.Name+ColorReset)
	}

	for code, info := range icaoMap {
		match := "##" + code
		plainText = strings.ReplaceAll(plainText, match, info.Name)
		colorText = strings.ReplaceAll(colorText, match, ColorGreen+info.Name+ColorReset)
	}

	reg := regexp.MustCompile(`(D|T12|T24)\((.*?)\)`)
	matches := reg.FindAllStringSubmatch(plainText, -1)

	monthMap := map[string]string{
		"01": "Jan", "02": "Feb", "03": "Mar", "04": "Apr", "05": "May", "06": "Jun",
		"07": "Jul", "08": "Aug", "09": "Sep", "10": "Oct", "11": "Nov", "12": "Dec",
	}

	for _, match := range matches {

		defaultDate := match[0]
		dateTag := match[1]
		isoDate := match[2]

		if !strings.Contains(isoDate, "T") || !strings.Contains(isoDate, "+") && !strings.Contains(isoDate, "-") && !strings.Contains(isoDate, "Z") {
			continue
		}

		parts := strings.Split(isoDate, "T")

		if len(parts) != 2 {
			continue
		}

		date := parts[0]
		timeWoffset := parts[1]

		dateSplit := strings.Split(date, "-")
		year := dateSplit[0]
		month := dateSplit[1]
		day := dateSplit[2]

		var time string
		var offset string

		if strings.HasSuffix(timeWoffset, "Z") {
			time = strings.TrimSuffix(timeWoffset, "Z")
			offset = "+00:00"
		} else if strings.Contains(timeWoffset, "+") {
			toffSplit := strings.Split(timeWoffset, "+")
			time = toffSplit[0]
			offset = "+" + toffSplit[1]
		} else if strings.Contains(timeWoffset, "-") {
			toffSplit := strings.Split(timeWoffset, "-")
			time = toffSplit[0]
			offset = "-" + toffSplit[1]
		} else {
			continue
		}

		timeSplit := strings.Split(time, ":")
		hoursStr := timeSplit[0]
		minutes := timeSplit[1]

		var formatResult string

		switch dateTag {
		case "D":
			monthName, _ := monthMap[month]
			formatResult = fmt.Sprintf("%s-%s-%s", day, monthName, year)

		case "T12":
			hours, err := strconv.Atoi(hoursStr)
			if err != nil {
				continue
			}
			AMPM := "AM"
			if hours >= 12 {
				AMPM = "PM"
			}
			if hours > 12 {
				hours = hours - 12
			}
			if hours == 0 {
				hours = 12
			}
			formatResult = fmt.Sprintf("%02d:%s:%s (%s)", hours, minutes, AMPM, offset)

		case "T24":
			formatResult = fmt.Sprintf("%s:%s (%s)", hoursStr, minutes, offset)
		}

		if formatResult != "" {
			plainText = strings.Replace(plainText, defaultDate, formatResult, 1)
			colorText = strings.Replace(colorText, defaultDate, ColorYellow+formatResult+ColorReset, 1)
		}
	}

	return plainText, colorText
}
