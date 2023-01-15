package api

import (
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/log"
	"net/http"
	"strconv"
	"time"
)

func StartReceiverHttpServer() {
	addr := ":" + strconv.Itoa(apidef.HttpApiServerPort)
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
