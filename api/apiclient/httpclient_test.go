package apiclient

import (
	"encoding/json"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/GoTools/file"
	"testing"
)

func TestDoPermRequest(t *testing.T) {
	filePath := "D:/Chaos/go/EzeShare/README.md"
	file := file.OpenFile(filePath)

	file.Size()

	DoPermRequest("http://127.0.0.1:23019", file, apidef.PermTypeAlways, func(per int) {

	})
}

func TestDo(t *testing.T) {
	var permReqJson = `{"Code":2000,"Msg":"success","Data":{"SecToken":"JnhJiKjeD6x/7waGTh13zIq51Q8897qjbLRg34FF+L728ZODokb5tCU+CXLKH8sy6QBkf4ZfznGuEPgUrQuf++HIgFSGDzhCJmuEryMSR5Kwn9zzi3QrQLxWYerBdJzzxVKK6WPLuzr7YJzrxIoIvMCuINhj36PxlE4H8Bua3ccN1Vm2MHA9rkYIeamN56SSS+H3SmakyVXUe1F4KXiDMO/O9A+O33/N0TJFAQzKLu26OCA8m7zuGoI3BDgP31lvOfD5ZPfxEdhzwAf+gKps91Briifr9a2LZNZGTgjqRCnIQG78UnahGLNRgpw77sF6dyW65IrcAtH3psgi5jnorA==","PermType":"AllowOnce","TransferId":"c03fd63e-486d-422a-8e33-45daf65a91fd"}}`

	var apiResp apidef.PubResp
	err := json.Unmarshal([]byte(permReqJson), &apiResp)
	if err != nil {
		log.I("parse respBody error: ", permReqJson)
	}
	if apiResp.Code == apidef.RespCodeSuccess {
		data := apiResp.Data.(map[string]interface{})
		secToken := data["SecToken"].(string)
		//permType := data["PermType"].(apidef.PermReqRespType)
		//transferId := data["TransferId"].(string)

		decryptToken, err := auth.DecryptRSA(secToken, auth.GetRsaPrivateKeyPath())
		if err != nil {
			return
		}
		log.I("decrypt token: ", decryptToken)
	}
}
