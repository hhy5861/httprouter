package httprouter

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	Debug          bool
	LogruserSystem = logrus.New()
)

func init() {
	fileName := "httprouter.log"
	path := "/var/log/httprouter"

	_, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			params := make(map[string]interface{})
			params["file_path"] = path
			params["file_name"] = fileName
			LogruserSystem.WithFields(logrus.Fields{
				"error":  err,
				"params": params,
			}).Fatal("System log module open file errer")
		}
	}

	fileNames := fmt.Sprintf("%s/%s", path, fileName)
	file, err := os.OpenFile(fileNames, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		LogruserSystem.Out = file
	} else {
		params := make(map[string]interface{})
		params["file_path"] = path
		params["file_name"] = fileName
		LogruserSystem.WithFields(logrus.Fields{
			"error":  err,
			"params": params,
		}).Fatal("System log module open file errer")
	}

	LogruserSystem.Level = logrus.InfoLevel
	LogruserSystem.Formatter = &logrus.JSONFormatter{}
}
