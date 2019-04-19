package errors

import "net/http"

// HTTPError replies as response with specific message & HTTP code.
func HTTPError(w http.ResponseWriter, err string, code int) {
	http.Error(w, err, code)
}
