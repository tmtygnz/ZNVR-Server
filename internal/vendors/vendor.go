package vendors

import (
	"corvette/internal/config"
	"log/slog"
)

type GenericVendor struct {
	url     string
	surl    string
	camType string
	camName string
}

type Vendor interface {
	URL() string
	SURL() string
	Type() string
	CamName() string
}

func VendorMapper(configInfo []config.CameraInfo) []Vendor {
	var vendors []Vendor
	for _, vendorInfo := range configInfo {
		switch vendorInfo.Type {
		case "Generic":
			newGenericCamera := CreateRtspVendor(vendorInfo.URL, vendorInfo.SURL, vendorInfo.Type, vendorInfo.CamName)
			vendors = append(vendors, newGenericCamera)
		}
	}

	slog.Info("Vendors mapped.")

	return vendors
}
