package errors

import "net/http"

func HttpError(w http.ResponseWriter, err string, code int) {
	http.Error(w, err, code)
}
