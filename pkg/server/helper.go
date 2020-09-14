package server

import (
	"net"

	"github.com/jamalyusuf/avoxi-challenge/pkg/countrycodes"
	pbService "github.com/jamalyusuf/avoxi-challenge/proto"
	"github.com/jpillora/ipfilter"
)

// isValidIP checks if a string is a valid ip address
func isValidIP(address string) bool {

	if IP := net.ParseIP(address); IP == nil {
		return false
	}
	return true
}

// isValidCountries checks if passed in countries are valid countryCodes
func isValidCountries(countries []*pbService.Country) bool {
	// ensure all countries sent are valid aplha-2 country codes
	for _, country := range countries {
		if _, ok := countrycodes.GetByAlpha2(country.GetCountry()); !ok {
			return false
		}
	}

	return true
}

// getCountryFromIP given a IP address will return the country
func getCountryFromIP(address string) string{
	country := ipfilter.IPToCountry(address)
	if country == ""{
		return unknownLocation.Error()
	}
	return country
}