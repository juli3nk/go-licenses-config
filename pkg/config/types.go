package config

type LicensesCheck struct {
	// list of allowed license names, can't be used in combination with disallowed_types
    AllowedLicenses []string `yaml:"allowed_licenses,omitempty"`
	// list of disallowed license types, can't be used in combination with allowed_licenses (default: forbidden, unknown)
    DisallowedTypes []string `yaml:"disallowed_types,omitempty"`
}

type Licenses struct {
	// Package path prefixes to be ignored.
	// Dependencies from the ignored packages are still checked.
    Ignore []string `yaml:"ignore,omitempty"`

    Check  LicensesCheck `yaml:"check,omitempty"`
}
