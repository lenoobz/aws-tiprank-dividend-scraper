package currency

import (
	"fmt"

	"github.com/lenoobz/aws-tiprank-dividend-scraper/consts"
)

// GetCountryCurrency gets currency code of a given country name
func GetCountryCurrency(countryName string) (string, error) {
	currency, found := consts.Currencies[countryName]

	if !found {
		return "", fmt.Errorf("not found %s's currency", countryName)
	}

	return currency.Code, nil
}
