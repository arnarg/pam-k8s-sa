package main

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	defaultIssuer   = "https://kubernetes.default.svc.cluster.local"
	defaultAudience = "https://kubernetes.default.svc.cluster.local"
	defaultCAFile   = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

type config struct {
	DiscoveryURL string
	Issuer       string
	Audience     string
	CAFile       string
	VerifyTLS    bool
}

func parseConfig(args []string) (*config, error) {
	conf := &config{
		Issuer:    defaultIssuer,
		Audience:  defaultAudience,
		CAFile:    defaultCAFile,
		VerifyTLS: true,
	}

	for _, arg := range args {
		opt := strings.Split(arg, "=")

		switch opt[0] {
		case "discovery_url":
			conf.DiscoveryURL = opt[1]
		case "issuer":
			conf.Issuer = opt[1]
		case "audience":
			conf.Audience = opt[1]
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

	if conf.DiscoveryURL == "" {
		conf.DiscoveryURL = conf.Issuer
	}

	return conf, nil
}
