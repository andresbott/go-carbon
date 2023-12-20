package env_test

import (
	"git.andresbott.com/Golang/carbon/libs/env"
	"os"
	"testing"
)

func TestStringVar(t *testing.T) {
	setEnv(t, "FOO", "BAR")

	var value string

	env.StringVar(&value, "FOO", "")

	if err := env.Parse(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if value != "BAR" {
		t.Fatalf("invalid value: %v", value)
	}
}

func TestStringVarEmpty(t *testing.T) {
	setEnv(t, "FOO", "")

	var value string

	env.StringVar(&value, "FOO", "BAR")

	if err := env.Parse(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if value != "BAR" {
		t.Fatalf("invalid value: %v", value)
	}
}

func TestRequiredStringVar(t *testing.T) {
	setEnv(t, "FOO", "BAR")

	var value string

	env.RequiredStringVar(&value, "FOO")

	if err := env.Parse(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if value != "BAR" {
		t.Fatalf("invalid value: %v", value)
	}
}

func TestRequiredStringVarEmpty(t *testing.T) {
	setEnv(t, "FOO", "")

	value := "BAR"

	env.RequiredStringVar(&value, "FOO")

	if err := env.Parse(); err == nil {
		t.Fatalf("no error returned")
	}
	if value != "BAR" {
		t.Fatalf("invalid value: %v", value)
	}
}

func TestBoolVar(t *testing.T) {
	setEnv(t, "FOO", "true")

	var value bool

	env.BoolVar(&value, "FOO", false)

	if err := env.Parse(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if value != true {
		t.Fatalf("invalid value")
	}
}

func TestBoolVarEmpty(t *testing.T) {
	setEnv(t, "FOO", "")

	var value bool

	env.BoolVar(&value, "FOO", true)

	if err := env.Parse(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if value != true {
		t.Fatalf("invalid value")
	}
}

func TestBoolVarInvalid(t *testing.T) {
	setEnv(t, "FOO", "BAR")

	value := true

	env.BoolVar(&value, "FOO", false)

	if err := env.Parse(); err == nil {
		t.Fatalf("no error returned")
	}
	if value != true {
		t.Fatalf("invalid value")
	}
}

func setEnv(t *testing.T, key, value string) {
	prevValue, ok := os.LookupEnv(key)

	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("cannot set environment variable: %v", err)
	}

	if ok {
		t.Cleanup(func() {
			if err := os.Setenv(key, prevValue); err != nil {
				t.Fatalf("restore previous environment variable value: %v", err)
			}
		})
	} else {
		t.Cleanup(func() {
			if err := os.Unsetenv(key); err != nil {
				t.Fatalf("remove environment variable: %v", err)
			}
		})
	}
}
