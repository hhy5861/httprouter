package logs

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hhy5861/httprouter"
)

var (
	logruserSystem *logrus.Logger
)

type (
	statusWriter struct {
		http.ResponseWriter
		status int
		length int
	}

	Params map[string]interface{}
)

func instantiation(path string) {
	fileName := "access.log"
	if path == "" {
		path = "/data/logs/httprouter"
	}

	_, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			params := make(Params)
			params["file_path"] = path
			params["file_name"] = fileName

			Fatal(params, "System log module open file errer")
		}
	}

	logruserSystem = logrus.New();
	fileNames := fmt.Sprintf("%s/%s", path, fileName)
	file, err := os.OpenFile(fileNames, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		logruserSystem.Out = file
	} else {
		params := make(Params)
		params["file_path"] = path
		params["file_name"] = fileName
		Fatal(params, "System log module open file errer")
	}

	logruserSystem.Level = logrus.InfoLevel
	logruserSystem.Formatter = &logrus.JSONFormatter{}
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}

	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

func WriteLog(r *httprouter.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if logruserSystem ==nil{
			instantiation(r.Path)
		}
		start := time.Now()
		writer := statusWriter{w, 0, 0}
		r.ServeHTTP(&writer, request)
		if r.Debug {
			end := time.Now()
			latency := end.Sub(start)
			statusCode := writer.status
			length := writer.length

			params := make(Params)

			params["code"] = statusCode
			params["length"] = length
			params["method"] = request.Method
			params["path"] = request.URL
			params["time"] = fmt.Sprintf("%s", latency)
			params["ip"] = request.RemoteAddr
			params["user_agent"] = request.UserAgent()

			Info(params, "System log info")
		}
	}
}

func Info(ps Params, message string) {
	logruserSystem.WithFields(logrus.Fields{
		"params": ps,
	}).Info(message)
}

func Warn(ps Params, message string) {
	logruserSystem.WithFields(logrus.Fields{
		"params": ps,
	}).Warn(message)
}

func Fatal(ps Params, message string) {
	logruserSystem.WithFields(logrus.Fields{
		"params": ps,
	}).Fatal(message)
}
