package http

import (
	"net/http"
	"strconv"
	"time"
)

const HttpServerPort = 23019

func StartReceiverHttpServer() {
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(HttpServerPort),
		Handler:        NewMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}
