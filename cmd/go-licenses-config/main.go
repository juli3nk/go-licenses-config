package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/juli3nk/go-licenses-config/pkg/config"
)

func main() {
	pkgPath := flag.String("packages", "./...", "Go package path(s) to scan")
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Erreur config: %v", err)
	}

	args, err := checkArgsBase64(cfg, *pkgPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(args)
}

// CheckArgs returns the CLI args for: go-licenses check [packages...]
// Example output: ["check", "./...", "--ignore", "github.com/user1/repo1", "--allowed_licenses", "MIT"]
func checkArgs(cfg *config.Licenses, packages ...string) []string {
	args := []string{"check"}

	// --ignore <prefix>  (repeatable flag)
	for _, ign := range cfg.Ignore {
		args = append(args, "--ignore", ign)
	}

	// --allowed_licenses <spdx>  (repeatable flag)
	for _, lic := range cfg.Check.AllowedLicenses {
		args = append(args, "--allowed_licenses", lic)
	}

	// --disallowed_types <type>  (repeatable flag)
	for _, t := range cfg.Check.DisallowedTypes {
		args = append(args, "--disallowed_types", t)
	}

	// packages to scan (default: "./..." if none provided)
	if len(packages) == 0 {
		packages = []string{"./..."}
	}
	args = append(args, packages...)

	return args
}

// CheckArgsBase64 returns a base64-encoded string of the check CLI args.
// Useful for CI cache keys, job IDs, or artifact naming.
func checkArgsBase64(cfg *config.Licenses, packages ...string) (string, error) {
	args := checkArgs(cfg, packages...)

	// Use JSON as a stable, deterministic serialization format.
	data, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("marshal args: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
