package handlers

import (
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "service running")
}
