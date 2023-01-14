package ui

import (
	"github.com/Ericwyn/EzeShare/api/apidef"
)

type PermReqUiCallback func(permType apidef.PermReqRespType)

type UI struct {
	Name               string
	ShowPermReqUiAsync func(permReq apidef.ApiPermReq, callback PermReqUiCallback)
}
