package config_viper

// Config is a thin opinionated wrapper around viper (https://github.com/spf13/viper) as means
// to keep behaviour consistent, the overall merit goes to the viper project
// it only uses a subset of vipers potential, but you can get a reference of the underlying viper instance
import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// PathDetail represents a specific way of loading config files.
// note to myself: I'm not sure if abstracting this into an interface is overcomplicating something the would be
// also achievable with N factory functions, one for each config file mechanism we want to support.
type PathDetail interface {
	getViperDetails() (viperDetails, error)
}

type viperDetails struct {
	filename   string
	extension  string
	searchPath []string
}

type SingleConfigFile struct {
	Path string
}

func (c SingleConfigFile) getViperDetails() (viperDetails, error) {
	// we need to load a file
	filename := filepath.Base(c.Path)
	var extension = strings.TrimPrefix(filepath.Ext(filename), ".")
	var name = filename[0 : len(filename)-len(extension)]

	name = strings.TrimSuffix(name, ".")
	err := validExtension(extension)
	if err != nil {
		return viperDetails{}, err
	}

	dirs := []string{}
	dirs = append(dirs, filepath.Dir(c.Path))

	return viperDetails{
		filename:   name,
		extension:  extension,
		searchPath: dirs,
	}, nil
}

type TwelveFactor struct {
}

func (c TwelveFactor) getViperDetails() (viperDetails, error) {
	return viperDetails{}, nil
}

type OptUseEnvVar struct {
	Prefix string
}

func Load[T OptUseEnvVar](viperDetails PathDetail, opts ...T) (*Config, error) {
	details, err := viperDetails.getViperDetails()
	if err != nil {
		return nil, err
	}
	vp := viper.New()

	if len(opts) == 0 {
		return nil, fmt.Errorf("you need to provide at least one config option")
	}

	// additional options
	for _, opt := range opts {
		switch i := any(opt).(type) {
		case OptUseEnvVar:
			if i.Prefix != "" {
				vp.SetEnvPrefix(i.Prefix)
			}
			vp.AutomaticEnv()
		}
	}

	// set file path details
	if details.filename != "" {
		vp.SetConfigName(details.filename)
		vp.SetConfigType(details.extension)
		for _, p := range details.searchPath {
			vp.AddConfigPath(p)
		}
		err = vp.ReadInConfig()
	} else {
		// load an empty config
		viper.SetConfigType("yaml")
		var yamlExample = []byte(``)

		err = viper.ReadConfig(bytes.NewBuffer(yamlExample))
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	c := Config{
		viper: vp,
	}
	return &c, nil
}

// Config is a thin opinionated wrapper around viper (https://github.com/spf13/viper) as means
// to keep behaviour consistent, the overall merit goes to the viper project
// it only uses a subset of vipers potential, but you can get a reference of the underlying viper instance
type Config struct {
	viper *viper.Viper
}

func Read(path, envVarPrefix string) (*Config, error) {
	vp := viper.New()

	// check if we load a directory of a file
	dir, err := isDir(path)
	if err != nil {
		return nil, err
	}
	if !dir {
		// we need to load a file
		filename := filepath.Base(path)

		var extension = strings.TrimPrefix(filepath.Ext(filename), ".")
		var name = filename[0 : len(filename)-len(extension)]
		name = strings.TrimSuffix(name, ".")
		err = validExtension(extension)
		if err != nil {
			return nil, err
		}

		vp.SetConfigName(name)
		vp.SetConfigType(extension)
		p := filepath.Dir(path)
		vp.AddConfigPath(p)
	}

	vp.SetEnvPrefix(envVarPrefix)
	vp.AutomaticEnv()

	err = vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c := Config{
		viper: vp,
	}
	return &c, nil
}

func validExtension(in string) error {
	s := strings.TrimPrefix(in, ".")

	AllowedExtansions := map[string]bool{
		"YAML": true,
		"JSON": true,
		"TOML": true,
		"INI":  true,
	}
	if AllowedExtansions[strings.ToUpper(s)] {
		return nil
	}
	return fmt.Errorf("extension: %s is not in allowed list", s)

}

func loadCfg(path string, opts ...Option) *Config {

	// if path is a file load the file, else load all files in that dir
	// use default paths  (?) if sting is empty

	vp := viper.New()

	for _, opt := range opts {
		switch opt.Typ {
		case EnvVarPrefix:
			vp.SetEnvPrefix(opt.Value.(string))
			vp.AutomaticEnv()
			//case SingleCfgFile:
			//	f := filepath.Base(path)
			//	vp.SetConfigName(f)
			//	p := filepath.Dir(path)
			//	vp.AddConfigPath(p)
		}

	}

	c := Config{
		viper: vp,
	}
	return &c
}

func isDir(path string) (bool, error) {

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// ReadInDir ignores the file name but loads all files alphabetically in a directory
// this allows to split config into multiple files
func ReadInDir(path string) (*viper.Viper, error) {
	//viper.SetConfigName("default")
	//viper.AddConfigPath(path)
	//viper.ReadInConfig()
	//
	//if context != "" {
	//	viper.SetConfigName(context)
	//	viper.AddConfigPath(path)
	//	viper.MergeInConfig()
	//}
	//
	//viper.SetConfigName("config")
	//viper.AddConfigPath(".")
	//viper.MergeInConfig()
	return nil, nil
}

func (c *Config) Viper() *viper.Viper {
	return c.viper
}

func (c *Config) Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error {
	err := c.viper.Unmarshal(rawVal, opts...)
	if err != nil {
		return err
	}

	//validate := validator.New()
	//if err := validate.Struct(&config); err != nil {
	//	log.Fatalf("Missing required attributes %v\n", err)
	//}

	return nil
}
