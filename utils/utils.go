package utils

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func InfoHandleFunc(log *logrus.Logger, r *http.Request) {
	log.Infof("recieve %s %s", r.Method, r.URL.Path)
}
