package log

import (
	"log"
	"os"
)

// Logger logger of application.
var Logger = log.New(os.Stdout, "[MonitoringAzureAD]\t", log.Ldate|log.Ltime|log.LUTC)
