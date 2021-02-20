package response

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// JSON response as json
func JSON(w http.ResponseWriter, object interface{}) {

	b, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

// Error response when error from internal server
func Error(w http.ResponseWriter, err error, status ...int) {
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		log.Printf("ERR %+v", err)
	} else {
		log.Printf("ERR %s", err)
	}
	if len(status) > 0 {
		w.WriteHeader(status[0])
	}
	w.WriteHeader(http.StatusInternalServerError)
}
