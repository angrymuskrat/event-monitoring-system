package crawler

import (
	"github.com/BurntSushi/toml"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

type Configuration struct {
	RootDir                 string
	CheckpointUpdateTimeout string
	WorkersNumber           int
	DataStorageURL          string
	UseDataStorage          bool
}

func readConfig(path string) (cfg Configuration, err error) {
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		unilog.Logger().Error("unable to read config file", zap.String("path", path), zap.Error(err))
	}
	return
}
