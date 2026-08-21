package platform

import "strings"

var sensitive = []string{"patientid", "patientname", "accessionnumber", "authorization", "api-key"}

func RedactMap(fields map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range fields {
		redacted := false
		for _, s := range sensitive {
			if strings.EqualFold(k, s) {
				redacted = true
				break
			}
		}
		if redacted {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
