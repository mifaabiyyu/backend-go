package main

import (
	"net/http"
)

func (app AppWrapper) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("oke1"))
}
