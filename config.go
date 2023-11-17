package main

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	defaultServerURL        = "https://kubernetes.default.svc.cluster.local"
	defaultIssuer           = "https://kubernetes.default.svc.cluster.local"
	defaultAudience         = "https://kubernetes.default.svc.cluster.local"
	defaultCAFile           = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	defaultTokenFile        = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	defaultUsernameTemplate = "{{.Name | replace \"-\" \"_\"}}${{.Namespace | replace \"-\" \"_\"}}"
)

type config struct {
	ServerURL        string
	Issuer           string
	Audience         string
	UsernameTemplate string
	TokenFile        string
	CAFile           string
	VerifyTLS        bool
}

func defaultConfig() *config {
	return &config{
		ServerURL:        defaultServerURL,
		Issuer:           defaultIssuer,
		Audience:         defaultAudience,
		UsernameTemplate: defaultUsernameTemplate,
		TokenFile:        defaultTokenFile,
		CAFile:           defaultCAFile,
		VerifyTLS:        true,
	}
}

func parseConfig(args []string) (*config, error) {
	conf := defaultConfig()

	for _, arg := range args {
		opt := strings.Split(arg, "=")

		switch opt[0] {
		case "server_url":
			conf.ServerURL = opt[1]
		case "issuer":
			conf.Issuer = opt[1]
		case "audience":
			conf.Audience = opt[1]
		case "username_template":
			conf.UsernameTemplate = opt[1]
		case "token_file":
			conf.TokenFile = opt[1]
		case "ca_file":
			conf.CAFile = opt[1]
		case "verify_tls":
			val, err := strconv.ParseBool(opt[1])
			if err != nil {
				return nil, fmt.Errorf("unable to parse 'verify_tls' as bool: %s", err)
			}
			conf.VerifyTLS = val
		}
	}

	return conf, nil
}
