package main

import (
	"flag"
	"fmt"
	"github.com/Ericwyn/EzeShare/conf"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"github.com/Ericwyn/EzeShare/ui"
	"github.com/Ericwyn/EzeShare/ui/terminalui"
	"github.com/Ericwyn/EzeShare/utils"
	"github.com/gin-gonic/gin"
	"os"
)

var version = "v1.0-beta2"

var showVersion = flag.Bool("v", false, "show version")
var uiMode = flag.String("ui", "terminal", "set ui mode")
var sender = flag.Bool("sender", false, "run as sender")
var senderShort = flag.Bool("s", false, "run as sender")
var receiver = flag.Bool("receiver", false, "run as receiver")
var receiverShort = flag.Bool("r", false, "run as receiver")

var permCheck = flag.Bool("perm", false, "perm check")

var ipAddr = flag.String("ip", "", "set ip address or this device")
var sendFilePath = flag.String("f", "", "[sender only] file path of the file which will send to others")

var debug = flag.Bool("debug", false, "print debug log and sql")

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(
			"  _____         ____  _                    \n" +
				" | ____|_______/ ___|| |__   __ _ _ __ ___ \n" +
				" |  _| |_  / _ \\___ \\| '_ \\ / _` | '__/ _ \\\n" +
				" | |___ / /  __/___) | | | | (_| | | |  __/\n" +
				" |_____/___\\___|____/|_| |_|\\__,_|_|  \\___|\n" +
				"                                           \n" +
				"          versioon: " + version + "\n")
		return
	}

	log.SetPrint(log.LogPrintSet{
		PrintI: true,
		PrintD: *debug,
		PrintE: true,
	})

	storage.InitDb(*debug)
	gin.SetMode(utils.Check(*debug, gin.DebugMode, gin.ReleaseMode).(string))

	uiArgs := checkParam()

	conf.RunInPermCheckMode = *permCheck

	RunMainUi(*uiArgs)
}

func RunMainUi(args ui.MainUiArgs) {
	if args.UiMode == "terminal" {
		terminalui.TerminalUi.RunMainUI(args)
	}
}

func checkParam() *ui.MainUiArgs {
	if *uiMode == "terminal" {
		//if (!*sender && !*receiver) || (*sender && *receiver) {
		//	log.E("please set one arg in '-sender' or '-receiver'")
		//	os.Exit(1)
		//}
		var runMode ui.MainUiRunMode
		if *sender || *senderShort {
			runMode = ui.MainUiRunModeSender
			if *sendFilePath == "" {
				log.E("the file path for send is empty")
				os.Exit(-1)
			}
		} else if *receiver || *receiverShort {
			runMode = ui.MainUiRunModeReceiver
		} else {
			log.E("please set one arg in '-sender' /  '-s' or '-receiver' / '-r'")
			os.Exit(1)
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
