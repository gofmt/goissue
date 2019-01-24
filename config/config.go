package config

import (
	"io/ioutil"
	"os"

	"goissue/pkgs/utils"

	"gopkg.in/yaml.v2"
)

var C = struct {
	Debug  bool
	Addr   string
	DBAddr string
}{
	Debug:  true,
	Addr:   ":6868",
	DBAddr: "postgres://postgres:12345678@localhost:5432/goissue?sslmode=disable",
}

func Load(cfgPath string) error {
	if !utils.PathExist(cfgPath) {
		b, err := yaml.Marshal(&C)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(cfgPath, b, os.ModePerm); err != nil {
			return err
		}
	}

	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, &C)
}
