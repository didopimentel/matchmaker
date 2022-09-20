package handlers

import "net/http"

func AddHeaders(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")
}
