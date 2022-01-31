package service

import (
	"github.com/BurntSushi/toml"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

type Config struct {
	LogPath            string
	WorkersNumber      int
	MaxPoints          float64
	MinUsers           int
	AnomalyMode        bool
	DataStorageAddress string
	Address            string
}

func readConfig(path string) (cfg Config, err error) {
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		unilog.Logger().Error("unable to read config file", zap.String("path", path), zap.Error(err))
	}
	return
}
