package geodata

import (
	"fmt"
	"net/http"
	"testing"
)

func TestDownloadGeoIPLatest(t *testing.T) {
	var seeker GeoIPSeeker
	if err := DownloadGeoIPLatest(http.DefaultClient, &seeker); err != nil {
		t.Error(err)
	}
	if seeker.IsZero() {
		t.Error("initialize failed")
	}
	ipset := seeker.Query("hk")
	fmt.Println(ipset)
}

func TestDownloadGeoSiteLatest(t *testing.T) {
	var seeker GeoSiteSeeker
	if err := DownloadGeoSiteLatest(http.DefaultClient, &seeker); err != nil {
		t.Error(err)
	}
	if seeker.IsZero() {
		t.Error("initialize failed")
	}
	domains := seeker.Query("geosite:google@ads")
	fmt.Println(domains)
}

func ExampleDownloadGeoSiteLatest() {
	var seeker GeoSiteSeeker
	if err := DownloadGeoSiteLatest(http.DefaultClient, &seeker); err != nil {
		panic(err)
	}
	domains := seeker.Query("geosite:google@ads")
	fmt.Println(domains)
}

func ExampleDownloadGeoIPLatest() {
	var seeker GeoSiteSeeker
	if err := DownloadGeoIPLatest(http.DefaultClient, &seeker); err != nil {
		panic(err)
	}
	ipset := seeker.Query("cn")
	fmt.Println(ipset)
}
