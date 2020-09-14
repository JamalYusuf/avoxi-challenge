package server

import (
	"context"
	"errors"
	"fmt"
	pbService "github.com/jamalyusuf/avoxi-challenge/proto"
	"github.com/jpillora/ipfilter"
	"log"
)

// Backend implements the protobuf interface
type Backend struct {
	Infolog    *log.Logger
	Errorlog 	*log.Logger
	Filter *ipfilter.IPFilter
}

var (
	invalidIPFormat    error = errors.New("invalid ip format, ip must be string in ip4 or ip6 format")
	invalidCountryCode error = errors.New("invalid country code format, ISO3166-1 Alpha-2 format required")
	unknownLocation error = errors.New("Unknown")
)

// New initializes a new Backend struct.
func New(infoLog *log.Logger, errorLog *log.Logger) *Backend {

	return &Backend{
		Infolog: infoLog,
		Errorlog: errorLog,
		Filter: ipfilter.New(ipfilter.Options{}),
	}
}
// GeoIPCheck given an IP and Locations, will check if IP is connecting from a whitelisted country
func (b Backend) GeoIPCheck(_ context.Context, req *pbService.GeoIPCheckRequest) (*pbService.GeoIPResponse, error) {

	var response pbService.GeoIPResponse

	// validate IP address
	if ok := isValidIP(req.GetIP()); !ok {
		response.Result = false
		response.Status = invalidIPFormat.Error()
		b.Errorlog.Printf("GeoIPCheck %s ip:%s", invalidIPFormat, req.GetIP())
		return &response, nil
	}

	// covert  to []string
	var countries []string
	for _, country := range req.GetAllowedCountries() {
		countries = append(countries, country.GetCountry())
	}

	// validate that all countries are aplha-2 country codes
	if ok := isValidCountries(req.GetAllowedCountries()); !ok {
		response.Result = false
		response.Status = invalidCountryCode.Error()
		b.Errorlog.Printf("GeoIPCheck %s countries: %s", invalidCountryCode, countries)
		return &response, nil
	}

	b.Filter = ipfilter.New(ipfilter.Options{
		AllowedCountries: countries,
	})

	response.Result = b.Filter.Allowed(req.GetIP())
	response.Status = fmt.Sprintf("ip:%s from:%s whitelisted:%t", req.GetIP(), b.Filter.IPToCountry(req.GetIP()), response.Result)
	return &response, nil
}

// IPLocation give a IP will return the location
func (b Backend) IPLocation(_ context.Context, req *pbService.IPLocationRequest) (*pbService.IPLocationResponse, error) {
	var response pbService.IPLocationResponse

	// validate IP address
	if ok := isValidIP(req.GetIP()); !ok {
		response.IP = req.GetIP()
		response.Location = unknownLocation.Error()
		b.Errorlog.Printf("GeoIPCheck %s ip:%s", invalidIPFormat, req.GetIP())
		return &response, nil
	}

	response.IP = req.GetIP()
	response.Location = getCountryFromIP(req.GetIP())
	return &response, nil
}

// Health will return a status ok
func (b Backend) Health(_ context.Context, req *pbService.HealthRequest) (*pbService.HealthResponse, error) {

	return &pbService.HealthResponse{
		Status: "OK",
	}, nil
}

