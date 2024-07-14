package config

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

// todo set own struct annotations like `config:fieldName, required`

func TestLoad2(t *testing.T) {
	tcs := []struct {
		name      string
		opts      []any
		envs      map[string]string
		expectVal map[string]string
	}{
		{
			name: "load from file",
			opts: []any{CfgFile{"sampledata/testSingleFile.yaml"}},

			// intentionally setting envs that do NOT apply because we did not set the Option
			envs: map[string]string{
				"IGNORE_FLOATNUM": "44.4",
				"IGNORE_NUMBER":   "9090",
			},
			expectVal: map[string]string{
				"floatNum":             "3.14",
				"nested.child.renamed": "renamedString",
				"bol":                  "true",
			},
		},
		{
			name: "12 factor only envs no prefix",
			opts: []any{EnvVar{}},

			// intentionally setting envs that do NOT apply because we did not set the Option
			envs: map[string]string{
				"NUMBER":               "60",
				"FLOATNUM":             "6.65",
				"TEXT":                 "this is a string",
				"BOL":                  "true",
				"LISTSTRING.0":         "string 1",
				"LISTSTRING.1":         "string 2",
				"NESTED.CHILD.RENAMED": "envValue",
			},
			expectVal: map[string]string{
				"floatNum":             "6.65",
				"nested.child.renamed": "envValue",
				"bol":                  "true",
			},
		},
		{
			name: "12 factor only envs with prefix",
			opts: []any{EnvVar{"TEST"}},

			// intentionally setting envs that do NOT apply because we did not set the Option
			envs: map[string]string{
				"TEST_NUMBER":               "60",
				"TEST_FLOATNUM":             "6.65",
				"TEST_TEXT":                 "this is a string",
				"TEST_BOL":                  "true",
				"TEST_LISTSTRING.0":         "string 1",
				"TEST_LISTSTRING.1":         "string 2",
				"TEST_NESTED.CHILD.RENAMED": "envValue",
			},
			expectVal: map[string]string{
				"floatNum":             "6.65",
				"nested.child.renamed": "envValue",
				"bol":                  "true",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}
			cfg, err := Load(tc.opts...)
			if err != nil {
				t.Fatal(err)
			}

			t.Run("get values", func(t *testing.T) {
				for k, v := range tc.expectVal {
					got := cfg.GetString(k)
					if diff := cmp.Diff(got, v); diff != "" {
						t.Errorf("unexpected value (-got +want)\n%s", diff)
					}
				}
			})

		})
	}
}

func TestFlattenMap(t *testing.T) {
	byt, err := os.ReadFile("sampledata/testFlatten.yaml")
	if err != nil {
		t.Fatal(err)
	}
	in, err := readCfgBytes(byt, ExtYaml)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]interface{}{}
	flatten("", in, got)
	want := map[string]interface{}{
		"general.child1.list.0.name":         "item1",
		"general.child1.list.0.value":        true,
		"general.child1.list.1.name":         "item2",
		"general.child1.list.1.value":        "value2",
		"general.child1.list.1.data.subData": "my string",
		"general.child1.list.2.name":         "float",
		"general.child1.list.2.val":          2.5,
		"general.child2.sub1.sub2":           1,
		"top":                                "level",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected value (-got +want)\n%s", diff)
	}

}
func TestFlattenStruct(t *testing.T) {
	in := DefaultCfg

	got := map[string]interface{}{}
	err := flattenStruct(in, got)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{
		"number":              100,
		"floatNum":            100.1,
		"text":                "default text",
		"listString.0":        "default 1",
		"listString.1":        "default 2",
		"nested.child.number": 101,
		"nested.child.text":   "child text default",
		"userList.0.name":     "a1",
		"userList.0.pass":     "b1",
		"userList.1.name":     "a2",
		"userList.1.pass":     "b2",
		"userList.2.pass":     "b3",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected value (-got +want)\n%s", diff)
	}

}
