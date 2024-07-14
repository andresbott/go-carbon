package config

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CfgFile enables config to be loaded from a single file
type CfgFile struct {
	path string
}

// CfgDir enables config to be loaded from a conf.d directory,
// note that directory values will take precedence over single file
type CfgDir struct {
	path string
}

// EnvVar enables to load config using an env vars
// note that Envs will take precedence over file persisted values
type EnvVar struct {
	Prefix string
}

type Config struct {
	loadEnvs  bool
	envPrefix string
	flatData  map[string]any
	subset    string
}

func Load(opts ...any) (*Config, error) {

	c := Config{
		//data:     map[string]interface{}{},
		flatData: map[string]any{},
	}
	// add cfg options into struct to control the order of precedence
	type cfgLoader struct {
		file *CfgFile
		dir  *CfgDir
		env  *EnvVar
	}
	cl := cfgLoader{}

	for _, opt := range opts {

		switch item := opt.(type) {
		case CfgFile:
			cl.file = &item
		case CfgDir:
			spew.Dump("CfgDir")
		case EnvVar:
			c.loadEnvs = true
			c.envPrefix = item.Prefix
		case []any:
			return nil, fmt.Errorf("wrong options payload: [][]any, only pass an array of options")

		}
	}

	var data map[string]any
	if cl.file != nil {
		extType := fileType(cl.file.path)
		if extType == ExtUnsupported {
			return nil, fmt.Errorf("file %s is of unsuporeted type", cl.file.path)
		}
		byt, err := os.ReadFile(cl.file.path)
		if err != nil {
			return nil, err
		}

		data, err = readCfgBytes(byt, extType)
		if err != nil {
			return nil, err
		}
	}
	flatten("", data, c.flatData)

	return &c, nil
}

func (c *Config) GetString(fieldName string) string {
	// check ENV firs
	envName := fieldName
	if c.envPrefix != "" {
		envName = c.envPrefix + "_" + fieldName
	}
	envName = strings.ToUpper(envName)
	envVal := os.Getenv(envName)
	if c.loadEnvs && envVal != "" {
		return envVal
	}

	val, ok := c.flatData[fieldName]
	if ok {
		switch val.(type) {
		case map[string]interface{}:
			return ""

		case string:
			return val.(string)
		default:
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

// Subset returns a config that only handles a subset of the overall config
func (c *Config) Subset(key string) *Config {
	newC := Config{
		loadEnvs:  c.loadEnvs,
		envPrefix: c.envPrefix,
		subset:    c.subset + sep + key,
		//data:      c.data,
		flatData: c.flatData,
	}

	return &newC
}

func flatten(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	// got from: https://stackoverflow.com/questions/64419565/how-to-efficiently-flatten-a-map
	if len(prefix) > 0 {
		prefix += sep
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flatten(prefix+k, child, dest)
		case []interface{}:
			for i := 0; i < len(child); i++ {
				switch child[i].(type) {
				case map[string]interface{}:
					flatten(prefix+k+sep+strconv.Itoa(i), child[i].(map[string]interface{}), dest)
				default:
					dest[prefix+k+sep+strconv.Itoa(i)] = child[i]
				}
			}
		default:
			dest[prefix+k] = v
		}
	}
}

// readCfgBytes takes a []byte, normally from reading a file, and will parse it's content depending
// on the extension passed in ext
// it returns a map[string]any
func readCfgBytes(bytes []byte, ext string) (map[string]any, error) {
	var data map[string]any

	if ext == ExtYaml {
		err := yaml.Unmarshal(bytes, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if ext == ExtJson {
		err := json.Unmarshal(bytes, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("unsuported file type")

}

const (
	ExtYaml        = "YAML"
	ExtJson        = "JSON"
	ExtUnsupported = "unsupported"
)

func fileType(fpath string) string {
	filename := filepath.Base(fpath)
	extension := strings.TrimPrefix(filepath.Ext(filename), ".")
	extension = strings.ToUpper(extension)
	switch extension {
	case ExtYaml:
		return ExtYaml
	case ExtJson:
		return ExtJson
	default:
		return ExtUnsupported
	}
}
