package vendors

type RtspVendor struct {
	GenericVendor
}

func CreateRtspVendor(url string, surl string, cType string, name string) *RtspVendor {
	return &RtspVendor{
		GenericVendor: GenericVendor{
			camType: cType,
			url:     url,
			surl:    surl,
			camName: name,
		},
	}
}

func (rv *RtspVendor) URL() string {
	return rv.url
}

func (rv *RtspVendor) SURL() string {
	return rv.surl
}

func (rv *RtspVendor) Type() string {
	return rv.camType
}

func (rv *RtspVendor) CamName() string {
	return rv.camName
}
