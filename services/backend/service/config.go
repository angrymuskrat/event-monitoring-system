package service

import (
	"github.com/BurntSushi/toml"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

type Config struct {
	Address         string
	LogPath         string
	AuthLogPath     string
	SessionKey      string
	TimerLogPath    string
	ProxyDataPath   string
	User            string
	Password        string
	TestMod         bool
	Connector       string
	ConnectorParams map[string]string `toml:"conn-params"`
	CORSOrigin      string            `toml:"cors-origin"`
}

func readConfig(path string) (cfg Config, err error) {
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		unilog.Logger().Error("unable to read config file", zap.String("path", path), zap.Error(err))
	}
	if cfg.ProxyDataPath == "" {
		cfg.ProxyDataPath = "proxy_data"
	}
	return
}
