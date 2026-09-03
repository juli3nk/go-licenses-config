package config

import (
	"errors"
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

const (
	configYML  = ".go-licenses.yml"
	configYAML = ".go-licenses.yaml"
)

// New loads the first existing config file:
//  1. .go-licenses.yml
//  2. .go-licenses.yaml
func New() (*Licenses, error) {
	cfg := &Licenses{
		Check: LicensesCheck{
			DisallowedTypes: []string{"forbidden", "unknown"},
		},
	}

	var path string
	for _, candidate := range []string{configYML, configYAML} {
		if _, err := os.Stat(candidate); err == nil {
			path = candidate
			break
		}
	}

	if path == "" {
		// Aucun fichier trouvé -> on retourne la config par défaut.
		return cfg, nil
	}

	//nolint:gosec
	// false positive: reading config file from trusted user path
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config %s: %w", path, err)
	}

	return cfg, nil
}

func (c *Licenses) Validate() error {
	if len(c.Check.AllowedLicenses) > 0 && len(c.Check.DisallowedTypes) > 0 {
		return errors.New("allowed_licenses and disallowed_types cannot be used together")
	}

	seen := make(map[string]struct{})
	for _, lic := range c.Check.AllowedLicenses {
		if _, ok := seen[lic]; ok {
			return fmt.Errorf("duplicate allowed license: %s", lic)
		}
		seen[lic] = struct{}{}
	}

	seen = make(map[string]struct{})
	for _, t := range c.Check.DisallowedTypes {
		if _, ok := seen[t]; ok {
			return fmt.Errorf("duplicate disallowed type: %s", t)
		}
		seen[t] = struct{}{}
	}

	return nil
}
