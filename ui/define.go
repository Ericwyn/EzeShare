package ui

import "github.com/Ericwyn/EzeShare/api"

type PermReqUiCallback func(permType api.PermReqRespType)

type UI struct {
	Name               string
	ShowPermReqUiAsync func(permReq api.ApiPermReq, callback PermReqUiCallback)
}
