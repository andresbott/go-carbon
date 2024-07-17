package user

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type StaticUsers struct {
	//Users map[string]string
	Users []User `yaml:"users"`
}

func (fu StaticUsers) AllowLogin(user string, hash string) bool {
	for _, u := range fu.Users {
		if user == u.Name {
			if !u.Enabled {
				return false
			}
			if u.Pw == hash {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

// todo implement Reload functionality
// maybe add filewatcher to auto reload?
func (u User) Reload() error {
	return fmt.Errorf("not implemented")
}

func FromFile(file string) (*StaticUsers, error) {
	fType := fileType(file)

	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if fType == ExtYaml || fType == ExtYml {
		return yamlBytes(b)
	}
	if fType == ExtJson {
		return jsonBytes(b)
	}
	return htpasswdBytes(b)

}

func yamlBytes(in []byte) (*StaticUsers, error) {
	data := StaticUsers{}
	err := yaml.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}
	return &data, err
}
func jsonBytes(in []byte) (*StaticUsers, error) {
	data := StaticUsers{}
	err := json.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}
	return &data, err
}

func htpasswdBytes(in []byte) (*StaticUsers, error) {
	reader := csv.NewReader(bytes.NewBuffer(in))
	reader.Comma = ':'
	reader.Comment = '#'
	reader.TrimLeadingSpace = true

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	data := StaticUsers{}

	for _, record := range lines {
		u := User{
			Name:    record[0],
			Pw:      record[1],
			Enabled: true,
		}
		data.Users = append(data.Users, u)
	}
	return &data, err
}

const (
	ExtYaml     = "YAML"
	ExtYml      = "YML"
	ExtJson     = "JSON"
	ExtHtpasswd = "htpasswd"
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
		return ExtHtpasswd
	}
}
