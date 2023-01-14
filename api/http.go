package api

import (
	"github.com/Ericwyn/EzeShare/log"
	"net/http"
	"strconv"
	"time"
)

const HttpServerPort = 23019

func StartReceiverHttpServer() {
	addr := ":" + strconv.Itoa(HttpServerPort)
	s := &http.Server{
		Addr:           addr,
		Handler:        NewMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.I("start http server in Addr : ", addr)
	_ = s.ListenAndServe()
}
