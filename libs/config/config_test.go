package config_test

import (
	"git.andresbott.com/Golang/carbon/libs/config"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type serverCfg struct {
	General serverGeneral
	DevMode bool `mapstructure:"isDevMode" `
}
type serverGeneral struct {
	Port     int    `mapstructure:"port" validate:"required"`
	LogLevel string `mapstructure:"log_level" `
}

func TestConfig2(t *testing.T) {
	tcs := []struct {
		name         string
		path         config.PathDetail
		options      []config.Option
		envs         map[string]string
		expect       map[string]string
		expectParams serverCfg
	}{
		{
			name: "load from file",
			path: config.SingleConfigFile{
				Path: "./sampledata/server.yaml",
			},
			// intentionally setting envs that do NOT apply because we did not set the Option
			envs: map[string]string{
				"TEST_ISDEVMODE":    "false",
				"TEST_GENERAL.PORT": "9090",
			},
			expect: map[string]string{
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
		{
			name: "load from file and Override Envs",
			path: config.SingleConfigFile{
				Path: "./sampledata/server.yaml",
			},
			options: []config.Option{{Typ: config.EnvVarPrefix, Value: "TEST"}},
			envs: map[string]string{
				"TEST_ISDEVMODE":    "false",
				"TEST_GENERAL.PORT": "9090",
			},
			expect: map[string]string{
				"isDevMode":    "false",
				"general.port": "9090",
			},
			expectParams: serverCfg{
				General: serverGeneral{
					Port:     9090,
					LogLevel: "info",
				},
				DevMode: false,
			},
		},
		{
			name:    "12 factor without env prefix",
			path:    config.TwelveFactor{},
			options: []config.Option{{Typ: config.EnvVarPrefix}},
			envs: map[string]string{
				"ISDEVMODE":    "false",
				"GENERAL.PORT": "9090",
			},
			expect: map[string]string{
				"isDevMode":    "false",
				"general.port": "9090",
			},
			expectParams: serverCfg{
				General: serverGeneral{
					Port:     9090,
					LogLevel: "info",
				},
				DevMode: false,
			},
		},
		{
			name:    "12 factor use env prefix",
			path:    config.TwelveFactor{},
			options: []config.Option{{Typ: config.EnvVarPrefix, Value: "TEST"}},
			envs: map[string]string{
				"TEST_ISDEVMODE":    "false",
				"TEST_GENERAL.PORT": "9090",
			},
			expect: map[string]string{
				"isDevMode":    "false",
				"general.port": "9090",
			},
			expectParams: serverCfg{
				General: serverGeneral{
					Port:     9090,
					LogLevel: "info",
				},
				DevMode: false,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}
			cfg, err := config.Load(tc.path, tc.options...)
			if err != nil {
				t.Fatal(err)
			}

			t.Run("get values", func(t *testing.T) {
				for k, v := range tc.expect {
					got := cfg.Viper().GetString(k)
					if diff := cmp.Diff(got, v); diff != "" {
						t.Errorf("unexpected value (-got +want)\n%s", diff)
					}
				}
			})

			t.Run("unmarshal", func(t *testing.T) {

				got := serverCfg{}
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

type params struct {
	General struct {
		MyKey string `mapstructure:"my_key" validate:"required"`
	} `mapstructure:"general"`
	TopLevel string `mapstructure:"Toplevel" validate:"required"`
}

func TestConfig(t *testing.T) {
	envPrefix := "TEST"
	tcs := []struct {
		name         string
		path         string
		envs         map[string]string
		expect       map[string]string
		expectParams params
	}{
		{
			name: "load from file",
			path: "./sampledata/sample.yaml",
			expect: map[string]string{
				"Toplevel":       "banana",
				"general.my_key": "my value",
				"SomeNumber":     "123",
			},
			expectParams: params{
				General: struct {
					MyKey string `mapstructure:"my_key" validate:"required"`
				}(struct {
					MyKey string
				}{MyKey: "my value"}),
				TopLevel: "banana",
			},
		},
		{
			name: "load from file and Override Envs",
			path: "./sampledata/sample.yaml",
			envs: map[string]string{
				envPrefix + "_TOPLEVEL":       "apple",
				envPrefix + "_GENERAL.MY_KEY": "some Value",
			},
			expect: map[string]string{
				"Toplevel":       "apple",
				"general.my_key": "some Value",
				"SomeNumber":     "123",
			},
			expectParams: params{
				General: struct {
					MyKey string `mapstructure:"my_key" validate:"required"`
				}(struct {
					MyKey string
				}{MyKey: "some Value"}),
				TopLevel: "apple",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			cfg, err := config.Read(tc.path, envPrefix)
			if err != nil {
				t.Fatal(err)
			}

			t.Run("get values", func(t *testing.T) {
				for k, v := range tc.expect {
					got := cfg.Viper().GetString(k)
					if diff := cmp.Diff(got, v); diff != "" {
						t.Errorf("unexpected value (-got +want)\n%s", diff)
					}
				}
			})

			t.Run("unmarshal", func(t *testing.T) {

				got := params{}
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

func TestConfigErrors(t *testing.T) {

	tcs := []struct {
		name      string
		path      string
		envPrefix string
		expect    string
	}{
		{
			name:   "unsupported config file",
			path:   "./sampledata/style.css",
			expect: "extension: css is not in allowed list",
		},
		{
			name:   "file does not exist",
			path:   "./sampledata/no",
			expect: "open ./sampledata/no: no such file or directory",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := config.Read(tc.path, tc.envPrefix)
			if err == nil {
				t.Errorf("expecting an error but got none")
			}
			if err.Error() != tc.expect {
				if diff := cmp.Diff(err.Error(), tc.expect); diff != "" {
					t.Errorf("unexpected error message (-got +want)\n%s", diff)
				}
			}

		})
	}
}
