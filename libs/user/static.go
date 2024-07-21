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
	Users []User `yaml:"users"`
}

func (stu *StaticUsers) AllowLogin(user string, pw string) bool {
	for _, u := range stu.Users {
		if user == u.Name {
			if !u.Enabled {
				return false
			}
			access, err := checkPass(pw, u.Pw)
			if err != nil {
				return false
			}
			return access
		}
	}
	return false
}
func (stu *StaticUsers) Add(user string, pw string) {
	stu.Users = append(stu.Users, User{
		Name:    user,
		Pw:      pw,
		Enabled: true,
	})
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
		if len(record) != 2 {
			return nil, fmt.Errorf("the file does not seem to be a valid htpasswd file")
		}
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
