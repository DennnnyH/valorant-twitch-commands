package httpapi

import (
	_ "embed"
	"net/http"
	"strings"
)

//go:embed riot.txt
var riotTxt string

func riotVerification(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(strings.TrimSpace(riotTxt)))
}
