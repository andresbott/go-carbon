package config

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	c, err := Load(CfgFile{
		path: "sampledata/server.yaml",
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = c

	spew.Dump(c.data)
}

type serverCfg struct {
	General serverGeneral
	DevMode bool `mapstructure:"isDevMode" `
}
type serverGeneral struct {
	Port     int    `mapstructure:"port" validate:"required"`
	LogLevel string `mapstructure:"log_level" `
}

func TestLoad2(t *testing.T) {
	tcs := []struct {
		name         string
		opts         []any
		envs         map[string]string
		expectVal    map[string]string
		expectParams serverCfg
	}{
		{
			name: "load from file",
			opts: []any{CfgFile{"sampledata/server.yaml"}},

			// intentionally setting envs that do NOT apply because we did not set the Option
			envs: map[string]string{
				"TEST_ISDEVMODE":    "false",
				"TEST_GENERAL.PORT": "9090",
			},
			expectVal: map[string]string{
				"isDevMode":    "true",
				"general.port": "8080",
			},
			expectParams: serverCfg{
				General: serverGeneral{
					Port:     8080,
					LogLevel: "info",
				},
				DevMode: true,
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

			//t.Run("unmarshal", func(t *testing.T) {
			//
			//	got := serverCfg{}
			//	err = cfg.Unmarshal(&got)
			//	if err != nil {
			//		t.Fatal(err)
			//	}
			//	if diff := cmp.Diff(got, tc.expectParams); diff != "" {
			//		t.Errorf("unexpected value (-got +want)\n%s", diff)
			//	}
			//})

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
