package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
)

func pamAuthenticate(l logger, username, token string, conf *config) error {
	// The whole procedure can't take more than 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Create a custom client context for custom root ca or insecure tls verification
	if cctx, err := createClientContext(l, ctx, conf); err != nil {
		return err
	} else {
		ctx = cctx
	}

	// Create a custom insecure context as issuer and discovery URLs don't match
	if conf.Issuer != conf.DiscoveryURL {
		l.Warnf("discovery url and issuer don't match, creating an insecure issuer url context")

		ctx = oidc.InsecureIssuerURLContext(ctx, conf.Issuer)
	}

	// Create a provider for OIDC
	provider, err := oidc.NewProvider(ctx, conf.DiscoveryURL)
	if err != nil {
		return fmt.Errorf("unable to discover OIDC endpoint: %s", err)
	}

	// Create a verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: conf.Audience,
	})

	// Verify id token
	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		return fmt.Errorf("unable to verify id token: %s", err)
	}

	// Verify username with subject
	if err := matchUserSubject(username, idToken.Subject); err != nil {
		return err
	}

	return nil
}

func createClientContext(l logger, ctx context.Context, conf *config) (context.Context, error) {
	tlsConf := &tls.Config{}

	if conf.VerifyTLS {
		certData, err := os.ReadFile(conf.CAFile)
		if err != nil {
			return nil, fmt.Errorf("could not read ca file: %s", err)
		}

		cp := x509.NewCertPool()
		if ok := cp.AppendCertsFromPEM(certData); !ok {
			return nil, fmt.Errorf("could not parse ca file data")
		}

		tlsConf.RootCAs = cp
	} else {
		l.Warnf("skipping tls verification")

		tlsConf.InsecureSkipVerify = true
	}

	return oidc.ClientContext(ctx, &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConf}}), nil
}

func matchUserSubject(username, subject string) error {
	parts := strings.SplitN(username, "$", 2)

	// Check that split was correct
	if len(parts) != 2 {
		return fmt.Errorf("username does not fit the pattern '{{service_account}}${{namespace}}'")
	}

	// Compare username data with subject
	expected := fmt.Sprintf("system:serviceaccount:%s:%s", parts[1], parts[0])
	if subject != expected {
		return fmt.Errorf("token subject '%s' did not match the expected '%s'", subject, expected)
	}

	return nil
}
