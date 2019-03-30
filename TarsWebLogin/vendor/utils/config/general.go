package config

import (
	"github.com/Andrew-M-C/go-tools/jsonconv"
	"github.com/Andrew-M-C/go-tools/log"
	"io/ioutil"
	// "os"
)

const (
	CONFIG_FILE_PATH	= "/etc/TarsWebLogin/config.json"
	PID_FILE_PATH		= "/etc/TarsWebLogin/pid"
)

var (
	appConf *jsonconv.JsonValue
)

func init() {
	// parse config
	// read config file
	conf_bytes, err := ioutil.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		log.Error("Failed to read configuration file: %s", err.Error())
		appConf = jsonconv.NewObject()
		return
	}
	conf_str := string(conf_bytes)
	log.Info("Configuration: >> %s <<", conf_str)
	appConf, err = jsonconv.NewFromString(conf_str)
	if err != nil {
		log.Error("Unmarshal config file failed: %s", err.Error())
		appConf = jsonconv.NewObject()
		return
	}
	return
}


func GetConfig() *jsonconv.JsonValue {
	return appConf
}


func GetHomeDir() string {
	dir, err := appConf.GetString("homeDir")
	if err != nil {
		return "./"
	} else {
		return dir
	}
}


// func SaveConfig() error {
// 	str, _ := appConf.Unmarshal()
// 	f, err := os.OpenFile(CONFIG_FILE_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
// 	defer f.Close()
// 	if err != nil {
// 		log.Error("Failed to open config file: %s", err.Error())
// 		return err
// 	}
// 	err = f.Write([]byte(str))
// 	return err
// }
