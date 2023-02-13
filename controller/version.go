package controller

var buildVersion string

func Version() string {
	return buildVersion
}

type VersionResponse struct {
	BuildVersion string `json:"buildVersion"`
}
