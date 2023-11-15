package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type tokenRoundTripper struct {
	token     string
	transport *http.Transport
}

func (t *tokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Inject the token into the request
	if t.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(t.token)))
	}

	return t.transport.RoundTrip(req)
}

func newHttpClient(l logger, conf *config) (*http.Client, error) {
	tlsConf := &tls.Config{}

	if conf.VerifyTLS {
		// Read CA file
		certData, err := os.ReadFile(conf.CAFile)
		if err != nil {
			return nil, fmt.Errorf("could not read ca file: %s", err)
		}

		// Create cert pool with CA file data
		cp := x509.NewCertPool()
		if ok := cp.AppendCertsFromPEM(certData); !ok {
			return nil, fmt.Errorf("could not parse ca file data")
		}

		tlsConf.RootCAs = cp
	} else {
		l.Warnf("skipping tls verification")

		tlsConf.InsecureSkipVerify = true
	}

	// Create a variable for final round tripper
	var rt http.RoundTripper
	transport := &http.Transport{TLSClientConfig: tlsConf}

	// Check if the configured token file exists
	if _, err := os.Stat(conf.TokenFile); err == nil {
		// Read the token file
		if token, err := os.ReadFile(conf.TokenFile); err == nil {
			// Wrap the transport in our tokenRoundTripper
			rt = &tokenRoundTripper{
				token:     string(token),
				transport: transport,
			}
		}
	} else {
		// Otherwise we just use the plain transport
		rt = transport
	}

	return &http.Client{Transport: rt}, nil
}
