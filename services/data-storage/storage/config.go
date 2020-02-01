package storage

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"log"
	"os"
)

type Configuration struct {
	AuthDB            string
	AggrPostsGRIDSize float64
}

func readConfig(path string) (cfg Configuration, err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		unilog.Logger().Error("unable to read config file", zap.String("path", path), zap.Error(err))
	}
	return
}
