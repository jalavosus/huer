package serverutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SimpleMessageHandler(msg string, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendSimpleResponseWithStatusCode(msg, statusCode, w)
		return
	}
}

func SendSimpleResponseWithStatusCode(msg string, statusCode int, w http.ResponseWriter) {
	d, _ := json.MarshalIndent(
		map[string]interface{}{"message": msg, "code": statusCode},
		"",
		"  ",
	)

	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, "%s", string(d))

	return
}