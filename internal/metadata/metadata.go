package metadata

import (
	"net"
	"net/http"
	"time"

	"github.com/mileusna/useragent"
	"github.com/oschwald/maxminddb-golang"
)



type Click struct {
	Timestmap time.Time `json:"timestamp"`
	IPAddress string `json:"ip_address"`
	Referrer string `json:"referrer"`
	UserAgent string `json:"user_agent"`
	Device string `json:"device"`
	OS string `json:"os"`
	Browser string `json:"browser"`
	Country string `json:"country"`
	City string `json:"city"`
}

type GeoIPCity struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`

	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
}


func MetadataMiddleware(db *maxminddb.Reader) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ipStr, _, _ :=  net.SplitHostPort(r.RemoteAddr)
			ip := net.ParseIP(ipStr)

			ua := useragent.Parse(r.Header.Get("User-Agent"))
			var deviceType string

			if ua.Mobile {
				deviceType = "Mobile"
			} else if ua.Tablet {
				deviceType = "Tablet"
			} else if ua.Desktop {
				deviceType = "Desktop"
			} else {
				deviceType = "Unknown"
			}

			var geoData GeoIPCity
			country, city := "Unknown", "Unknown"
			if err := db.Lookup(ip, &geoData); err == nil {
				country = geoData.Country.ISOCode
				if cityName, ok := geoData.City.Names["en"]; ok {
					city = cityName
				}
			} 

			clickData := Click{
				Timestmap: time.Now().UTC(),
				IPAddress: ipStr,
				Referrer: r.Header.Get("Referer"),
				UserAgent: r.Header.Get("User-Agent"),
				Device: deviceType,
				OS: ua.OS,
				Browser: ua.Name,
				Country: country,
				City: city,
			}


		})
	}
} 



