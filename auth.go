package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
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
	if conf.ServerURL != conf.Issuer {
		l.Warnf("server url and issuer don't match, creating an insecure issuer url context")

		ctx = oidc.InsecureIssuerURLContext(ctx, conf.Issuer)
	}

	// Create a provider for OIDC
	provider, err := oidc.NewProvider(ctx, conf.ServerURL)
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
	if err := matchUserSubject(username, idToken.Subject, conf.UsernameTemplate); err != nil {
		return err
	}

	return nil
}

func createClientContext(l logger, ctx context.Context, conf *config) (context.Context, error) {
	client, err := newHttpClient(l, conf)
	if err != nil {
		return nil, err
	}

	return oidc.ClientContext(ctx, client), nil
}

func matchUserSubject(username, subject, userTpl string) error {
	// Subject should be in the form "system:serviceaccount:namespace:name"
	parts := strings.Split(subject, ":")

	// Check that subject parsing is correct
	if len(parts) != 4 {
		return fmt.Errorf("subject format is unknown: '%s'", subject)
	}

	// Render template into expected username
	expected, err := templateUsername(parts[3], parts[2], userTpl)
	if err != nil {
		return err
	}

	// Compare provided username with expected username
	if username != expected {
		return fmt.Errorf("username did not match expected username for token")
	}

	return nil
}

func templateUsername(name, namespace, userTpl string) (string, error) {
	tpl, err := template.New("").Funcs(template.FuncMap{
		"replace": func(old, new, input string) string { return strings.ReplaceAll(input, old, new) },
	}).Parse(userTpl)
	if err != nil {
		return "", fmt.Errorf("could not parse username template: %s", err)
	}

	buf := &bytes.Buffer{}

	// Render template
	data := map[string]string{
		"Name":      name,
		"Namespace": namespace,
	}
	if err := tpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("could not render username template:%s", err)
	}

	return buf.String(), nil
}
