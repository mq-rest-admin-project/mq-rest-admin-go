//go:build integration

package devtls

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var certBlockRE = regexp.MustCompile(`(?s)-----BEGIN CERTIFICATE-----.*?-----END CERTIFICATE-----`)

// CAFileFor returns the path to a PEM CA bundle that trusts the given MQ REST
// endpoints. If MQ_REST_TLS_CA_FILE is set it is returned unchanged; otherwise
// each server's certificate chain is extracted at runtime via the openssl CLI
// into a temporary bundle. TLS verification is never disabled.
func CAFileFor(restBaseURLs []string) (string, error) {
	if override := os.Getenv("MQ_REST_TLS_CA_FILE"); override != "" {
		return override, nil
	}

	blocks, err := extractServerCertChains(restBaseURLs)
	if err != nil {
		return "", err
	}
	if len(blocks) == 0 {
		return "", fmt.Errorf("could not extract dev MQ server certificates; set MQ_REST_TLS_CA_FILE")
	}

	file, err := os.CreateTemp("", "mq-dev-ca-*.pem")
	if err != nil {
		return "", fmt.Errorf("create temp CA bundle: %w", err)
	}
	defer func() { _ = file.Close() }()
	if _, err := file.WriteString(strings.Join(blocks, "\n") + "\n"); err != nil {
		return "", fmt.Errorf("write CA bundle: %w", err)
	}
	return file.Name(), nil
}

// extractServerCertChains reads each server's certificate chain with the
// openssl CLI. It only reads the certificate the server presents (it does not
// establish a trusted, application-level connection), so it stays free of any
// verification-disabling code in the library or test suite.
func extractServerCertChains(restBaseURLs []string) ([]string, error) {
	seen := make(map[string]bool)
	var blocks []string
	for _, raw := range restBaseURLs {
		parsed, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("parse url %q: %w", raw, err)
		}
		// #nosec G204 -- host:port comes from trusted dev configuration and is
		// only used to read the presented certificate.
		cmd := exec.Command("openssl", "s_client", "-connect", parsed.Host,
			"-servername", parsed.Hostname(), "-showcerts")
		cmd.Stdin = strings.NewReader("")
		out, _ := cmd.Output()
		for _, block := range certBlockRE.FindAllString(string(out), -1) {
			if !seen[block] {
				seen[block] = true
				blocks = append(blocks, block)
			}
		}
	}
	return blocks, nil
}
