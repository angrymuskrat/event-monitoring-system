package crawler

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

type Configuration struct {
	RootDir        string
	DataStorageURL string
	Groups         []Group
}

type Group struct {
	TorPorts  []int
	Token     string
	SessionID string
}

func readConfig(path string) (cfg Configuration, err error) {
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		unilog.Logger().Error("unable to read config file", zap.String("path", path), zap.Error(err))
	}
	fmt.Println(cfg)
	return
}
