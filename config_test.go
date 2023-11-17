package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseConfigEmptyList(t *testing.T) {
	expected := defaultConfig()

	conf, err := parseConfig([]string{})
	if err != nil {
		t.Fatalf("unexpected error while parsing empty list")
	}

	if !reflect.DeepEqual(conf, expected) {
		t.Fatalf("parsed config does not match expected output: %v != %v", conf, expected)
	}
}

func TestParseConfigNonEmptyList(t *testing.T) {
	expected := &config{
		ServerURL:        "https://127.0.0.1:6443",
		Issuer:           "https://127.0.0.1:6443",
		Audience:         "k3s",
		UsernameTemplate: "{{.Name}}@{{.Namespace}}",
		TokenFile:        "/tmp/k8s_token",
		CAFile:           "/tmp/k8s_ca.crt",
		VerifyTLS:        false,
	}
	args := []string{
		fmt.Sprintf("server_url=%s", expected.ServerURL),
		fmt.Sprintf("issuer=%s", expected.Issuer),
		fmt.Sprintf("audience=%s", expected.Audience),
		fmt.Sprintf("username_template=%s", expected.UsernameTemplate),
		fmt.Sprintf("token_file=%s", expected.TokenFile),
		fmt.Sprintf("ca_file=%s", expected.CAFile),
		"verify_tls=false",
	}

	conf, err := parseConfig(args)
	if err != nil {
		t.Fatalf("unexpected error when parsing valid parameters")
	}

	if !reflect.DeepEqual(conf, expected) {
		t.Fatalf("parsed config does not match expected output: %v != %v", conf, expected)
	}
}

func TestParseConfigVerifyTLSWrongType(t *testing.T) {
	if _, err := parseConfig([]string{"verify_tls=no"}); err == nil {
		t.Fatalf("parseConfig did not error when passing invalid parameter")
	}
}
