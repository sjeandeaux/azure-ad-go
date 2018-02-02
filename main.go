package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sjeandeaux/azure-ad-go/information"
	"github.com/sjeandeaux/azure-ad-go/log"
	"github.com/sjeandeaux/azure-ad-go/ugly"
)

//commandLineArgs all parameters in command line
type commandLineArgs struct {
	user     string
	password string
}

var commandLine = &commandLineArgs{}

//init configuration
func init() {
	log.Logger.SetOutput(os.Stdout)
	flag.StringVar(&commandLine.user, "user", os.Getenv("AZURE_AD_USER"), "user for Azure AD")
	flag.StringVar(&commandLine.password, "password", os.Getenv("AZURE_AD_PASSWORD"), "user for Azure AD")
	flag.Parse()
}

func main() {
	log.Logger.Println(information.Print())
	monitoringAzureAD := ugly.NewMonitoringAzureAD(commandLine.user, commandLine.password)
	if dir, err := monitoringAzureAD.GetDirSyncManagement(); err != nil {
		log.Logger.Panicln(err)
	} else {
		log.Logger.Println(fmt.Sprintf("%+v", dir))
	}
}
