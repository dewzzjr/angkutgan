package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON response as json
func JSON(w http.ResponseWriter, object interface{}) {

	b, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// Error response when error from internal server
func Error(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}
