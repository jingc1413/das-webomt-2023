package layout

import (
	"regexp"
	"strings"
)

func NameFormat(in string) string {
	out := in
	out = strings.ReplaceAll(out, "   ", " ")
	out = strings.ReplaceAll(out, "  ", " ")
	if re := regexp.MustCompile(`^Module\s*\d+\s*-C\d+\s*(.*)$`); re != nil {
		out = re.ReplaceAllString(out, "Carrier $1")
		out = strings.ReplaceAll(out, "Carrier Carrier", "Carrier")
	}
	if re := regexp.MustCompile(`^Port\s*\d+\s*(.*)$`); re != nil {
		out = re.ReplaceAllString(out, "Port $1")
		out = strings.ReplaceAll(out, "Port Port", "Port")
	}
	if re := regexp.MustCompile(`^Carrier\s*(UL|DL)\s*Attenuation$`); re != nil {
		out = re.ReplaceAllString(out, "Carrier $1 Attenuation")
		out = strings.ReplaceAll(out, "Carrier Carrier", "Carrier")
	}
	if re := regexp.MustCompile(`^(UL|DL)\s*Carrier\s*Power$`); re != nil {
		out = re.ReplaceAllString(out, "Carrier $1 Power")
		out = strings.ReplaceAll(out, "Carrier Carrier", "Carrier")
	}
	if re := regexp.MustCompile(`Carrier\s*Carrier`); re != nil {
		out = re.ReplaceAllString(out, "Carrier")
	}
	out = strings.ReplaceAll(out, "Frequency Start - End", "Frequency")
	out = strings.ReplaceAll(out, "Frequency Start - High", "Frequency")

	out = strings.ReplaceAll(out, "Label (Operator/Service)", "Operator/Service")

	return out
}
