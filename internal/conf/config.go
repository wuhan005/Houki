package conf

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

const configName = "config.json"

type config struct {
	Version   string   `json:"version"`
	ProxyAddr string   `json:"proxy_addr"`
	Modules   []string `json:"modules"`
}

var conf config = config{
	Version:   "",
	ProxyAddr: ":8080",
	Modules:   nil,
}

func Initialize() error {
	_, err := os.Stat(configName)
	if os.IsNotExist(err) {
		err = saveConfig()
		if err != nil {
			return err
		}

		log.Info("Create config file: config.json")
		return nil
	} else if err != nil {
		return errors.Wrap(err, "unexpect error")
	}

	configBytes, err := os.ReadFile(configName)
	if err != nil {
		return errors.Wrap(err, "read file")
	}

	log.Info("Get config from config.json")
	return json.Unmarshal(configBytes, &conf)
}

func Get() *config {
	return &conf
}

func saveConfig() error {
	configBytes, err := json.Marshal(conf)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	err = os.WriteFile(configName, configBytes, 0666)
	if err != nil {
		return errors.Wrap(err, "write file")
	}
	return nil
}
