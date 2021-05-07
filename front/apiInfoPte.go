package front

import "net/http"

func apiInfoPte(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}
