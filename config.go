package main

import (
	"os"
	"strings"
)

type Config struct {
	Regions      []string
	Roles        []string
	OutputFormat string
	OutputFile   string
	ScanInterval string
	Tags         TagMapping
}

type TagMapping struct {
	netbios                   string
	assetType                 string
	baselineConfig            string
	authenticatedScanPlanned  string
	lastScanned               string
	applicationVendor         string
	applicationNameAndVersion string
	applicationPatchLevel     string
	applicationFunction       string
	comments                  string
	serialNumber              string
	sysadmin                  string
	appadmin                  string
}

func NewConfig() Config {
	config := configDefaults
	envMap := map[string]*string{
		"OUTPUT_FORMAT":                  &config.OutputFormat,
		"OUTPUT_FILE":                    &config.OutputFile,
		"SCAN_INTERVAL":                  &config.ScanInterval,
		"TAG_NETBIOS":                    &config.Tags.netbios,
		"TAG_ASSET_TYPE":                 &config.Tags.assetType,
		"TAG_BASELINE_CONFIG":            &config.Tags.baselineConfig,
		"TAG_AUTHENTICATED_SCAN_PLANNED": &config.Tags.authenticatedScanPlanned,
		"TAG_LAST_SCANNED":               &config.Tags.lastScanned,
		"TAG_APPLICATION_VENDOR":         &config.Tags.applicationVendor,
		"TAG_APPLICATION_NAME_VERSION":   &config.Tags.applicationNameAndVersion,
		"TAG_APPLICATION_PATCH_LEVEL":    &config.Tags.applicationPatchLevel,
		"TAG_APPLICATION_FUNCTION":       &config.Tags.applicationFunction,
		"TAG_COMMENTS":                   &config.Tags.comments,
		"TAG_SERIAL_NUMBER":              &config.Tags.serialNumber,
		"TAG_SYSADMIN":                   &config.Tags.sysadmin,
		"TAG_APPADMIN":                   &config.Tags.appadmin,
	}

	if len(os.Getenv("REGIONS")) > 0 {
		config.Regions = strings.Split(os.Getenv("REGIONS"), ",")
	}
	if len(os.Getenv("ROLES")) > 0 {
		config.Roles = strings.Split(os.Getenv("ROLES"), ",")
	}
	for k, v := range envMap {
		SetStringFromEnv(k, v)
	}
	return config
}

func SetStringFromEnv(key string, pointer *string) {
	if len(os.Getenv(key)) > 0 {
		*pointer = os.Getenv(key)
	}
}

var configDefaults = Config{
	Regions:      []string{"us-gov-west-1"},
	Roles:        []string{},
	OutputFormat: "csv",
	OutputFile:   "output.csv",
	ScanInterval: "1d",
	Tags: TagMapping{
		netbios:                   "NetBIOS",
		assetType:                 "AssetType",
		baselineConfig:            "BaselineConfiguration",
		authenticatedScanPlanned:  "AuthenticatedScanPlanned",
		lastScanned:               "LastScanned",
		applicationVendor:         "ApplicationVendor",
		applicationNameAndVersion: "ApplicationNameAndVersion",
		applicationPatchLevel:     "ApplicationPatchLevel",
		applicationFunction:       "ApplicationFunction",
		comments:                  "Comments",
		serialNumber:              "SerialNumber",
		sysadmin:                  "SysAdmin",
		appadmin:                  "AppAdmin",
	},
}
