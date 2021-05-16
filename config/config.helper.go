package config

import (
	"fmt"
	"time"
)

// AllowDomain const
const AllowDomain = "www.tipranks.com"

// DomainGlob const
const DomainGlob = "*tipranks.*"

// GetFundOverviewURL get fund overview url
func GetDividendStockByDateURL(countryCode string, date time.Time) string {
	return fmt.Sprintf("https://www.tipranks.com/api/dividends/getByDate/?name=%s&country=%s", date.Format("2006-01-02"), countryCode)
}
