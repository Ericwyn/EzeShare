package main

import (
	"flag"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"github.com/Ericwyn/EzeShare/ui"
	"github.com/Ericwyn/EzeShare/ui/terminalui"
	"github.com/Ericwyn/EzeShare/utils"
	"github.com/gin-gonic/gin"
	"os"
)

var uiMode = flag.String("ui", "terminal", "set ui mode")
var sender = flag.Bool("sender", false, "run as sender")
var receiver = flag.Bool("receiver", false, "run as receiver")

var ipAddr = flag.String("ip", "", "[sender][receiver] set ip address or this device")
var sendFilePath = flag.String("f", "", "[sender] file path of the file which will send to others")

var debug = flag.Bool("debug", false, "print debug log and sql")

func main() {
	flag.Parse()

	log.SetPrint(log.LogPrintSet{
		PrintI: true,
		PrintD: *debug,
		PrintE: true,
	})

	storage.InitDb(*debug)
	gin.SetMode(utils.Check(*debug, gin.DebugMode, gin.ReleaseMode).(string))

	uiArgs := checkParam()

	RunMainUi(*uiArgs)
}

func RunMainUi(args ui.MainUiArgs) {
	if args.UiMode == "terminal" {
		terminalui.TerminalUi.RunMainUI(args)
	}
}

func checkParam() *ui.MainUiArgs {
	if *uiMode == "terminal" {
		if (!*sender && !*receiver) || (*sender && *receiver) {
			log.E("please set one arg in '-sender' or '-receiver'")
			os.Exit(1)
		}
		var runMode ui.MainUiRunMode
		if *sender {
			runMode = ui.MainUiRunModeSender
			if *sendFilePath == "" {
				log.E("the file path for send is empty")
				os.Exit(-1)
			}
		} else {
			runMode = ui.MainUiRunModeReceiver
		}
		return &ui.MainUiArgs{
			UiMode:     "terminal",
			RunMode:    runMode,
			SenderFile: *sendFilePath,
			IpAddr:     *ipAddr,
		}

	}
	// TODO other ui mode
	return nil
}