package main

import "testing"

func TestTemplateUsername(t *testing.T) {
	tpl := `{{.Name | replace "-" "_"}}${{.Namespace | replace "-" "_"}}`
	name := "database-user"
	namespace := "my-app"
	expected := "database_user$my_app"

	user, err := templateUsername(name, namespace, tpl)
	if err != nil {
		t.Fatalf("templateUsername failed with error: %s", err)
	}

	if user != expected {
		t.Fatalf("templated username does not match expected: '%s' != '%s'", expected, user)
	}
}
