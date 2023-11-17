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

func TestMatchUserSubjectCorrect(t *testing.T) {
	tpl := `{{.Name | replace "-" "_"}}${{.Namespace | replace "-" "_"}}`
	username := "database_user$my_ns"
	subject := "system:serviceaccount:my-ns:database-user"

	if err := matchUserSubject(username, subject, tpl); err != nil {
		t.Fatalf("matchUserSubject unexpectedly failed: %s", err)
	}
}

func TestMatchUserSubjectWrong(t *testing.T) {
	tpl := `{{.Name | replace "-" "_"}}${{.Namespace | replace "-" "_"}}`
	username := "database-user@my-ns"
	subject := "system:serviceaccount:my-ns:database-user"

	if err := matchUserSubject(username, subject, tpl); err == nil {
		t.Fatalf("matchUserSubject did not fail when suppose to")
	}
}
