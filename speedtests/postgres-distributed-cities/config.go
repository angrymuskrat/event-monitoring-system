package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

type Configuration struct {
	User     string
	Password string
	Host     string
	Port     string
}

func readConfig(path string) (cfg Configuration, err error) {
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		unilog.Logger().Error("unable to read config file", zap.String("path", path), zap.Error(err))
	}
	return
}

func (c *Configuration) makeAuthToken(dbname string) string {
	return fmt.Sprintf("database=%v user=%v password=%v sslmode=disable host=%v port=%v", dbname, c.User, c.Password, c.Host, c.Port)
}
