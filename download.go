package geodata

import (
	"io"
	"net/http"
)

const (
	GeoIPLatestURL   = "https://github.com/v2fly/geoip/releases/latest/download/geoip.dat"
	GeoSiteLatestURL = "https://github.com/v2fly/domain-list-community/releases/latest/download/dlc.dat"
)

func DownloadGeoIPLatest(client *http.Client, from io.ReaderFrom) error {
	return download(client, GeoIPLatestURL, from)
}

func DownloadGeoSiteLatest(client *http.Client, from io.ReaderFrom) error {
	return download(client, GeoSiteLatestURL, from)
}

func download(client *http.Client, url string, from io.ReaderFrom) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return http.ErrNotSupported
	}
	if _, err = from.ReadFrom(resp.Body); err != nil {
		return err
	}
	return resp.Body.Close()
}
