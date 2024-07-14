package config

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

type testConfig struct {
	Number     int          `config:"number"`
	FloatNum   float64      `config:"floatNum"`
	Text       string       `config:"text"`
	Bol        bool         `config:"bol"`
	StringList []string     `config:"listString"`
	StructList []userData   `config:"userList"`
	Nested     NestedConfig `config:"nested"`
}
type NestedConfig struct {
	Child  Child `config:"child"`
	Child2 struct {
		Number int `config:"number"`
	} `config:"child_2"`
}
type userData struct {
	Name string `config:"name"`
	Pass string `config:"pass"`
}

type Child struct {
	Number      int    `config:"number"`
	Text        string `config:"text"`
	AnotherName string `config:"renamed"`
}

// todo set own struct annotations like `config:fieldName, required`

func TestLoad2(t *testing.T) {
	tcs := []struct {
		name         string
		opts         []any
		envs         map[string]string
		expectVal    map[string]string
		expectParams testConfig
	}{
		{
			name: "load from file",
			opts: []any{CfgFile{"sampledata/testSingleFile.yaml"}},

			// intentionally setting envs that do NOT apply because we did not set the Option
			envs: map[string]string{
				"TEST_ISDEVMODE":    "false",
				"TEST_GENERAL.PORT": "9090",
			},
			expectVal: map[string]string{
				"floatNum":             "3.14",
				"nested.child.renamed": "renamedString",
			},
			expectParams: testConfig{
				Number:     60,
				FloatNum:   3.14,
				Text:       "this is a string",
				Bol:        true,
				StringList: []string{"sting 1", "string 2"},
				StructList: []userData{
					{Name: "u1", Pass: "p1"},
					{Name: "u2", Pass: "p2"},
					{Pass: "p3"},
				},
				Nested: NestedConfig{
					Child: Child{
						Number:      61,
						Text:        "this is a string 2",
						AnotherName: "renamedString",
					},
					Child2: struct {
						Number int `config:"number"`
					}(struct{ Number int }{
						Number: 62,
					}),
				},
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
			expectParams: testConfig{
				Number:     60,
				FloatNum:   6.65,
				Text:       "this is a string",
				Bol:        true,
				StringList: []string{"string 1", "string 2"},
				StructList: nil,
				Nested: NestedConfig{
					Child: Child{
						AnotherName: "envValue",
					},
				},
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

			t.Run("unmarshal", func(t *testing.T) {

				got := testConfig{}
				err = cfg.Unmarshal(&got)
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(got, tc.expectParams); diff != "" {
					t.Errorf("unexpected value (-got +want)\n%s", diff)
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
