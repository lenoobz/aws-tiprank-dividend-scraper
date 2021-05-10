package ticker

import (
	"fmt"
	"strings"
)

// GetYahooTicker gets yahoo ticker
func GetYahooTicker(tipRankTicker string) string {
	r := strings.Replace(tipRankTicker, ".", "-", 1)
	s := strings.Split(r, ":")

	if len(s) > 1 {
		if strings.EqualFold(s[0], "TSE") {
			return fmt.Sprintf("%s.TO", s[1])
		}
	}

	return r
}
