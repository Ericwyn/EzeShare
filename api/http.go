package api

import (
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/utils/netutils"
	"net/http"
	"strconv"
	"time"
)

func StartReceiverHttpServer() {
	addr := netutils.GetIPv4().String() + ":" + strconv.Itoa(apidef.HttpApiServerPort)
	s := &http.Server{
		Addr:           addr,
		Handler:        NewMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.I("start http server in Addr ", addr)
	_ = s.ListenAndServe()
}
